package protobuf

import (
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/ui"
)

func chooseService(cfg *settings.Settings) (string, string, error) {
	var (
		serviceName string
		serviceKind string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Service name. Enter the service name").
				Value(&serviceName).
				Validate(ui.IsEmpty("service name cannot be empty")),

			huh.NewSelect[string]().
				Title("Select the service type").
				Options(
					huh.NewOption("grpc", "grpc"),
					huh.NewOption("http", "http"),
				).
				Value(&serviceKind),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return "", "", err
	}

	return serviceName, serviceKind, nil
}

type gRPCForm struct {
	EntityName  string
	DefaultRPCs bool
	CustomRPCs  []string
}

func runGrpcForm(cfg *settings.Settings) (*gRPCForm, error) {
	var (
		entityName  string
		defaultRPCs = true
		customRPCs  []string
		text        string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Entity name. Enter the service main entity name:").
				Validate(ui.IsEmpty("entity name cannot be empty")).
				Value(&entityName),

			huh.NewConfirm().
				Title("Use default CRUD RPCs for the service?").
				Value(&defaultRPCs),

			huh.NewText().
				Title("Enter the custom RPCs names (one per line)").
				Value(&text),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return nil, err
	}

	if text != "" {
		customRPCs = strings.Split(text, "\n")
	}

	return &gRPCForm{
		EntityName:  entityName,
		DefaultRPCs: defaultRPCs,
		CustomRPCs:  customRPCs,
	}, nil
}

func runHTTPForm(cfg *settings.Settings) (bool, []*RPC, error) {
	var isAuthenticated bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Is the service authenticated?").
				Value(&isAuthenticated),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return false, nil, nil
	}

	rpcs, err := runHTTPRPCForm(cfg, isAuthenticated)
	if err != nil {
		return false, nil, nil
	}

	return isAuthenticated, rpcs, nil
}

func runHTTPRPCForm(cfg *settings.Settings, isAuthenticated bool) ([]*RPC, error) {
	var rpcs []*RPC

	for {
		rpc, err := promptSingleHTTPRPC(cfg, isAuthenticated)
		if err != nil {
			return nil, err
		}
		rpcs = append(rpcs, rpc)

		continueAdding, err := confirmAddMoreRPC(cfg)
		if err != nil {
			return nil, err
		}
		if !continueAdding {
			break
		}
	}

	return rpcs, nil
}

func promptSingleHTTPRPC(cfg *settings.Settings, isAuthenticated bool) (*RPC, error) {
	var (
		name     string
		method   string
		endpoint string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name. Enter the RPC name:").
				Value(&name).
				Validate(ui.IsEmpty("RPC name cannot be empty")),

			huh.NewSelect[string]().
				Title("Method. Enter the RPC method:").
				Options(
					huh.NewOption("GET", "get"),
					huh.NewOption("POST", "post"),
					huh.NewOption("PUT", "put"),
					huh.NewOption("DELETE", "delete"),
					huh.NewOption("PATCH", "patch"),
				).
				Value(&method),

			huh.NewInput().
				Title("Endpoint. Enter the RPC endpoint:").
				Value(&endpoint).
				Validate(ui.IsEmpty("RPC endpoint cannot be empty")),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return nil, err
	}

	return &RPC{
		IsAuthenticated: isAuthenticated,
		Name:            name,
		HTTPMethod:      method,
		HTTPEndpoint:    endpoint,
		AuthArgMode:     getAuthArgMode(method),
	}, nil
}

func confirmAddMoreRPC(cfg *settings.Settings) (bool, error) {
	var continueAdding bool

	confirm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to add a new RPC?").
				Value(&continueAdding),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := confirm.Run(); err != nil {
		return false, err
	}

	return continueAdding, nil
}
