package lint

type configContext struct {
	LineLength         uint
	Cyclomatic         uint
	Complexity         uint
	FunctionStatements uint
	FunctionLines      uint
	MaxNestingLevel    uint
}

var (
	// Available profiles that can be selected when running the lint command.
	availableProfiles = map[string]configContext{
		"default": {
			LineLength:         120,
			Cyclomatic:         12,
			Complexity:         15,
			FunctionStatements: 40,
			FunctionLines:      60,
			MaxNestingLevel:    3,
		},
	}
)
