package protobuf

import (
	"os"
	"strings"

	protofile "github.com/emicklei/proto"
)

type Proto struct {
	ServiceName string
	Methods     []*Method
}

func (p *Proto) parse(definitions *protofile.Proto) {
	protofile.Walk(definitions,
		protofile.WithPackage(p.parsePackage),
		protofile.WithRPC(p.parseMethods))
}

func (p *Proto) parseMethods(r *protofile.RPC) {
	m := loadMethod(r)
	p.Methods = append(p.Methods, m)
}

func (p *Proto) parsePackage(pkg *protofile.Package) {
	name := pkg.Name
	if strings.Contains(name, ".") {
		parts := strings.Split(name, ".")
		name = parts[len(parts)-1]
	}
	p.ServiceName = name
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
