package protobuf

// Answers represents the user-provided configuration for generating a protobuf
// module.
type Answers struct {
	ServiceName string
	Kind        string
	Grpc        *GrpcAnswers
	HTTP        *HTTPAnswers
}

// GrpcAnswers defines the properties for configuring gRPC service generation.
type GrpcAnswers struct {
	EntityName     string
	UseDefaultRPCs bool
	CustomRPCs     []string
}

// HTTPAnswers defines the properties for configuring HTTP service generation.
type HTTPAnswers struct {
	IsAuthenticated bool
	RPCs            []*RPC
}
