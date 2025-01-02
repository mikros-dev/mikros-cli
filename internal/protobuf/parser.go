package protobuf

import (
	"os"

	protofile "github.com/emicklei/proto"
)

type Proto struct {
	ModuleName  string
	ServiceName string
	Methods     []*Method
}

func (p *Proto) parse(definitions *protofile.Proto) {
	protofile.Walk(definitions,
		protofile.WithRPC(p.parseMethods))

	protofile.Walk(definitions,
		protofile.WithService(func(service *protofile.Service) {
			p.ServiceName = service.Name
		}))

	protofile.Walk(definitions,
		protofile.WithPackage(func(pkg *protofile.Package) {
			p.ModuleName = pkg.Name
		}))
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
	defer func(reader *os.File) {
		_ = reader.Close()
	}(reader)

	parser := protofile.NewParser(reader)
	definitions, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	p := &Proto{}
	p.parse(definitions)

	return p, nil
}
