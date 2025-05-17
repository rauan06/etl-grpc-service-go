module github.com/rauan06/etl-grpc-service-go

go 1.24.2

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3
	github.com/joho/godotenv v1.5.1
	google.golang.org/genproto/googleapis/api v0.0.0-20250303144028-a0af3efb3deb
	google.golang.org/grpc v1.72.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/kr/text v0.2.0 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.37.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/lint v0.0.0-20241112194109-818c5a804067 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.5.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

tool (
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)
