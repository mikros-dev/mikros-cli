package protobuf

import (
	"fmt"
	"os"
	"strings"

	protofile "github.com/emicklei/proto"
)

// Proto represents a protobuf file.
type Proto struct {
	ServiceName string
	Methods     []*Method
}

// Parse parses a protobuf file.
func Parse(filename string) (*Proto, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", filename, err)
	}
	defer func(reader *os.File) {
		_ = reader.Close()
	}(reader)

	parser := protofile.NewParser(reader)
	definitions, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	p := &Proto{}
	p.parse(definitions)

	return p, nil
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
