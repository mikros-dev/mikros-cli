package main

import (
	"context"
	mplugin "github.com/mikros-dev/mikros-cli/pkg/plugin"
	featurepb "github.com/mikros-dev/mikros-cli/pkg/plugin/feature"
	surveypb "github.com/mikros-dev/mikros-cli/pkg/plugin/survey"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
)

const (
	featureName   = "database"
	uiFeatureName = "nosql database"
)

type plugin struct {
	server *mplugin.BaseServer
}

// GetName must return the feature name that is used inside by mikros, which
// will be used to save definitions into the 'service.toml' file.
func (p *plugin) GetName(_ context.Context, _ *featurepb.Empty) (*featurepb.GetNameResponse, error) {
	return &featurepb.GetNameResponse{
		Name: featureName,
	}, nil
}

// GetUIName must return the feature name that will be displayed for the user
// when the mikros CLI survey is executed.
func (p *plugin) GetUIName(_ context.Context, _ *featurepb.Empty) (*featurepb.GetUINameResponse, error) {
	return &featurepb.GetUINameResponse{
		Name: uiFeatureName,
	}, nil
}

// GetSurvey should return a survey structure of questions that will be executed
// by mikros CLI if this feature is selected for use.
//
// This is optional, only useful if the feature has something that can be
// customized when created rather than having default values in the 'service.toml'
// file.
func (p *plugin) GetSurvey(_ context.Context, _ *featurepb.Empty) (*featurepb.GetSurveyResponse, error) {
	return &featurepb.GetSurveyResponse{
		Survey: &surveypb.Survey{
			Questions: []*surveypb.Question{
				{
					Name:    "database_cache",
					Message: "Use cache to optimize the queries?",
					Prompt:  surveypb.PromptKind_PROMPT_KIND_CONFIRM,
				},
				{
					Name:    "database_kind",
					Message: "Select the database kind:",
					Default: "mongo",
					Options: []string{"mongo", "postgres", "mysql", "sqlserver", "sqlite"},
					Prompt:  surveypb.PromptKind_PROMPT_KIND_MULTI_SELECT,
				},
				{
					Name:    "database_ttl",
					Message: "Enter the TTL of the entity, if it needs to be cooled:",
					Default: "0",
					Prompt:  surveypb.PromptKind_PROMPT_KIND_INPUT,
				},
				{
					Name:    "database_collections",
					Message: "Enter the name of additional collections (one by line):",
					Prompt:  surveypb.PromptKind_PROMPT_KIND_MULTI_SELECT,
				},
			},
		},
	}, nil
}

// ValidateAnswers is where the feature plugin should validate answers received
// by the mikros CLI and return the data in the format supported by 'service.toml'
// definitions file and a flag saying that if it should be written or not.
func (p *plugin) ValidateAnswers(_ context.Context, req *featurepb.ValidateAnswersRequest) (*featurepb.ValidateAnswersResponse, error) {
	// User survey answers are available here through this map
	_ = req.GetAnswers().AsMap()

	values, _ := structpb.NewStruct(map[string]interface{}{
		"collections": []string{"name1", "name2"},
		"ttl":         0,
	})

	return &featurepb.ValidateAnswersResponse{
		ShouldSave: true,
		Values:     values,
	}, nil
}

// Stop is where the mikros CLI signals that the plugin should stop its execution.
func (p *plugin) Stop(_ context.Context, _ *featurepb.Empty) (*featurepb.Empty, error) {
	p.server.Stop()
	return &featurepb.Empty{}, nil
}

func main() {
	s, err := mplugin.NewBaseServer()
	if err != nil {
		log.Fatal(err)
	}

	p := &plugin{server: s}
	featurepb.RegisterPluginServer(s.GetServer(), p)

	if err := p.server.Run(); err != nil {
		log.Fatal(err)
	}
}
