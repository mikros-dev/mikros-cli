package protobuf

import (
	protofile "github.com/emicklei/proto"
)

// Method represents a method of a protobuf service.
type Method struct {
	Name       string
	InputName  string
	OutputName string
}

func loadMethod(r *protofile.RPC) *Method {
	return &Method{
		Name:       r.Name,
		InputName:  r.RequestType,
		OutputName: r.ReturnsType,
	}
}
