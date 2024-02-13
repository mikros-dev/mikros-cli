package protobuf

import (
	"os"

	protofile "github.com/emicklei/proto"
)

type Proto struct {
	Methods []*Method
}

func (p *Proto) parse(definitions *protofile.Proto) {
	protofile.Walk(definitions,
		protofile.WithRPC(p.parseMethods))
}

func (p *Proto) parseMethods(r *protofile.RPC) {
	m := loadMethod(r)
	p.Methods = append(p.Methods, m)
}

func Parse(filename string) (*Proto, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	parser := protofile.NewParser(reader)
	definitions, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	p := &Proto{}
	p.parse(definitions)

	return p, nil
}
