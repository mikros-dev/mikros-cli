package client

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/mikros-dev/mikros-cli/internal/survey"
	servicepb "github.com/mikros-dev/mikros-cli/pkg/plugin/service"
)

type Service struct {
	cmd    *exec.Cmd
	conn   *grpc.ClientConn
	client servicepb.PluginClient
}

func NewService(path, name string) *Service {
	return &Service{
		cmd: exec.Command(filepath.Join(path, name)),
	}
}

func (s *Service) Start() error {
	if err := s.cmd.Start(); err != nil {
		return err
	}

	// Wait for the plugin to signal it's ready
	for {
		if _, err := os.Stat("plugin_ready.txt"); err == nil {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	// Remove the ready signal for the next plugin
	_ = os.Remove("plugin_ready.txt")

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	s.client = servicepb.NewPluginClient(conn)
	s.conn = conn

	return nil
}

func (s *Service) GetKind(ctx context.Context) (string, error) {
	res, err := s.client.GetKind(ctx, &servicepb.Empty{})
	if err != nil {
		return "", err
	}

	return res.GetKind(), nil
}

func (s *Service) GetName(ctx context.Context) (string, error) {
	res, err := s.client.GetName(ctx, &servicepb.Empty{})
	if err != nil {
		return "", err
	}

	return res.GetName(), nil
}

func (s *Service) GetSurvey(ctx context.Context) (*survey.Survey, error) {
	res, err := s.client.GetSurvey(ctx, &servicepb.Empty{})
	if err != nil {
		return nil, err
	}

	return survey.FromProtoSurvey(res.GetSurvey()), nil
}

func (s *Service) ValidateAnswers(ctx context.Context, answers map[string]interface{}) (map[string]interface{}, bool, error) {
	st, err := structpb.NewStruct(answers)
	if err != nil {
		return nil, false, err
	}

	res, err := s.client.ValidateAnswers(ctx, &servicepb.ValidateAnswersRequest{
		Answers: st,
	})
	if err != nil {
		return nil, false, err
	}

	return res.GetValues().AsMap(), res.GetShouldSave(), nil
}

func (s *Service) GetTemplates(ctx context.Context) (*servicepb.GetTemplateResponse, error) {
	return s.client.GetTemplate(ctx, &servicepb.Empty{})
}

func (s *Service) Stop(ctx context.Context) error {
	if _, err := s.client.Stop(ctx, &servicepb.Empty{}); err != nil {
		return err
	}

	_ = s.conn.Close()

	if err := s.cmd.Wait(); err != nil {
		return err
	}

	return nil
}
