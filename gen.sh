#!/bin/bash
protoc -I . \
        -I=/home/coder/googleapis/ \
        -I=/home/coder/grpc-gateway/ \
        --go_out .   \
        --go-grpc_out .   \
        --grpc-gateway_out . --grpc-gateway_opt logtostderr=true \
        --openapiv2_out . --openapiv2_opt logtostderr=true,use_go_templates=true \
        proto/service.proto
mv proto/service.swagger.json swagger.json