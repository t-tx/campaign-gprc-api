FROM ubuntu AS swagger-ui
RUN apt-get update -y 
RUN apt install -y git
WORKDIR /
RUN git clone https://github.com/swagger-api/swagger-ui.git

RUN cat <<EOF > /swagger-ui/dist/swagger-initializer.js
window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    url: "swagger.json",
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};
EOF

FROM golang:1.23.4-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin

FROM alpine
COPY --from=swagger-ui /swagger-ui/dist /public
COPY --from=builder /app/bin /app/bin
COPY swagger.json /app/swagger.json
WORKDIR /app
CMD ["./bin", "serve"]