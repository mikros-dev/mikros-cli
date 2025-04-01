package template

// Language represents a supported programming language for services generated
// by us.
type Language int

const (
	LanguageGolang Language = iota
	LanguageRust
)

func (k Language) String() string {
	switch k {
	case LanguageGolang:
		return "go"
	case LanguageRust:
		return "rust"
	}

	return "unknown"
}
