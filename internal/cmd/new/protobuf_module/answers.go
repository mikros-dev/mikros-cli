package protobuf_module

type Answers struct {
	ServiceName string
	Kind        string
	Grpc        *GrpcAnswers
	Http        *HttpAnswers
}

type GrpcAnswers struct {
	EntityName     string
	UseDefaultRPCs bool
	CustomRPCs     []string
}

type HttpAnswers struct {
	IsAuthenticated bool
	RPCs            []*RPC
}
