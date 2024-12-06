package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"net"
	"net/http"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server interface {
	RegisterGW(httpRegister GatewayRegister)
	AddInterceptors(interceptors ...grpc.UnaryServerInterceptor)
	Run() error
}

type baseService struct {
	interceptors       []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	name               string
	port               int
	listener           net.Listener
	register           Register
	gatewayRegister    []GatewayRegister
	server             *grpc.Server
	httpServer         *http.Server
	rootPath           string
	gatewayOption      *GatewayOption
}

// Register use to register with generated proto
type Register func(server *grpc.Server)

// GatewayRegister use to register with generated gateway proto
type GatewayRegister func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

type GatewayOption struct {
	ForceLowerCaseResponse bool
}

func (s *baseService) SetGWOption(option *GatewayOption) {
	s.gatewayOption = option
}

func NewGrpcServer(register Register, name string, port int) Server {
	s := baseService{register: register, name: name, port: port}
	return &s
}
func (s *baseService) AddInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	s.interceptors = append(s.interceptors, interceptors...)
}
func (s *baseService) Run() (err error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	fmt.Println("service started, listen at: " + lis.Addr().String())
	s.listener = lis

	return s.run()

}

func (s *baseService) RegisterGW(httpRegister GatewayRegister) {
	s.gatewayRegister = append(s.gatewayRegister, httpRegister)
}

const (
	contentTypeName  = "content-type"
	contentTypeValue = "application/grpc"
)

func (s *baseService) run() error {

	if s.gatewayRegister != nil {
		m := cmux.New(s.listener)
		listener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings(contentTypeName, contentTypeValue))
		httpListener := m.Match(cmux.HTTP1Fast())

		g := new(errgroup.Group)
		g.Go(func() error { return s.serve(listener) })
		g.Go(func() error { return s.gwServe(httpListener) })
		g.Go(func() error { return m.Serve() })

		return g.Wait()
	} else {
		return s.serve(s.listener)
	}
}

const (
	maxConcurrentStreams = 64
	maxRecvMsgSize       = 1024 * 1024 * 1024 * 4
	maxSendMsgSize       = 1024 * 1024 * 1024 * 4
)

func (s *baseService) createUnaryInterceptor() grpc.UnaryServerInterceptor {
	return grpc_middleware.ChainUnaryServer(
		grpc_prometheus.UnaryServerInterceptor,
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_middleware.ChainUnaryServer(s.interceptors...),
	)
}
func (s *baseService) serve(l net.Listener) error {
	sIntOpt := grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_prometheus.StreamServerInterceptor,
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_middleware.ChainStreamServer(s.streamInterceptors...),
	))

	uIntOpt := grpc.UnaryInterceptor(s.createUnaryInterceptor())

	s.server = grpc.NewServer(grpc.MaxConcurrentStreams(uint32(maxConcurrentStreams)), sIntOpt, uIntOpt, grpc.MaxRecvMsgSize(maxRecvMsgSize), grpc.MaxSendMsgSize(maxSendMsgSize))
	s.register(s.server)

	reflection.Register(s.server)
	grpc_prometheus.Register(s.server)
	return s.server.Serve(l)
}

func (s *baseService) Shutdown() error {
	if s.listener != nil {
		err := s.listener.Close()
		s.listener = nil
		if err != nil {
			return err
		}
	}
	return nil
}

const (
	jsonFile = "swagger.json"
)

func (s *baseService) gwServe(l net.Listener) error {
	ctx := context.Background()
	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{Marshaler: &runtime.JSONPb{}}),
		runtime.WithIncomingHeaderMatcher(IncomingHeaderMatcher),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxRecvMsgSize))}
	endPoint := fmt.Sprintf("localhost:%d", s.port)

	for _, gwRegister := range s.gatewayRegister {
		err := gwRegister(ctx, gwMux, endPoint, opts)
		if err != nil {
			return err
		}
	}

	//swagger
	mux := http.NewServeMux()
	mux.Handle("/", gwMux)

	fs := http.FileServer(http.Dir("/public"))
	mux.Handle("/help/", http.StripPrefix("/help", fs))

	//interceptor_server
	handlerMiddleware := handlerHTTPRequest(mux)

	handlerMiddleware = s.handlerHTTPBaseService(handlerMiddleware)

	handlerMiddleware = handleCrossOrigin(handlerMiddleware)

	s.httpServer = &http.Server{Handler: handlerMiddleware}

	return s.httpServer.Serve(l)
}
func handleCrossOrigin(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

func handlerHTTPRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Header.Get("Content-Type")) == "application/x-www-form-urlencoded" {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			jsonMap := make(map[string]interface{}, len(r.Form))
			for k, v := range r.Form {
				if len(v) > 0 {
					jsonMap[k] = v[0]
				}
			}
			jsonBody, err := json.Marshal(jsonMap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			r.Body = io.NopCloser(bytes.NewReader(jsonBody))
			r.ContentLength = int64(len(jsonBody))
		}

		handler.ServeHTTP(w, r)
	})
}

func (s *baseService) handlerHTTPBaseService(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.URL.Path {
		case "/help/swagger.json":
			http.ServeFile(w, r, jsonFile)

		case "/metrics":
			promhttp.Handler().ServeHTTP(w, r)

		default:

			handler.ServeHTTP(w, r)
		}
	})
}

func DefaultProtoErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaller runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
	w.WriteHeader(http.StatusNotFound)
}
func IncomingHeaderMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "connection":
		return "grpc-gw-internal-" + key, true
	default:
		return key, true
	}
}
