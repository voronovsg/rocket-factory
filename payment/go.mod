module github.com/voronovsg/rocket-factory/payment

go 1.26.1

replace github.com/voronovsg/rocket-factory/shared => ../shared

require (
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.28.0
	github.com/voronovsg/rocket-factory/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.79.3
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.3.3 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260319201613-d00831a3d3e7 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260319201613-d00831a3d3e7 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
