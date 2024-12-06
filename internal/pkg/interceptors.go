package pkg

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthInterceptor(headerKey string, jwtSecret []byte) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}
		tokens := md.Get("token")
		if len(tokens) == 0 {
			return handler(ctx, req)
		}
		token := tokens[0]
		claims, err := VerifyJWT(token, jwtSecret)
		if err != nil {
			return handler(ctx, req)
		}
		ctx = context.WithValue(ctx, contextKey, &AuthContext{Username: claims.Username, Payload: claims.Data})

		resp, err = handler(ctx, req)

		return resp, err
	}
}

func GetAuthContext(ctx context.Context) (*AuthContext, bool) {
	v, ok := ctx.Value(contextKey).(*AuthContext)
	return v, ok

}

type AuthContext struct {
	Username string
	Payload  map[string]string
}
type keyType struct{}

var contextKey = keyType{}
