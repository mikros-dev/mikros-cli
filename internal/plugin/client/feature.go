package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/structpb"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mikros-dev/mikros-cli/internal/survey"
	featurepb "github.com/mikros-dev/mikros-cli/pkg/plugin/feature"
)

type Feature struct {
	cmd    *exec.Cmd
	conn   *grpc.ClientConn
	client featurepb.PluginClient
}

func NewFeature(path, name string) *Feature {
	return &Feature{
		cmd: exec.Command(filepath.Join(path, name)),
	}
}

func (f *Feature) Start() error {
	println("1")
	if err := f.cmd.Start(); err != nil {
		return err
	}

	//println("2")
	//// Wait for the plugin to signal it's ready
	//for {
	//	if _, err := os.Stat("plugin_ready.txt"); err == nil {
	//		break
	//	}
	//
	//	time.Sleep(100 * time.Millisecond)
	//}

	//println("3")
	//// Remove the ready signal for the next plugin
	//_ = os.Remove("plugin_ready.txt")

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	healthClient := grpc_health_v1.NewHealthClient(conn)

	// Retry checking health until the server is ready
	for {
		resp, err := healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
		fmt.Println(resp.GetStatus())
		if err == nil && resp.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	println("vai tentar conectar")
	f.client = featurepb.NewPluginClient(conn)
	f.conn = conn

	return nil
}

func (f *Feature) Stop(ctx context.Context) error {
	if _, err := f.client.Stop(ctx, &featurepb.Empty{}); err != nil {
		return err
	}

	_ = f.conn.Close()

	if err := f.cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (f *Feature) GetUIName(ctx context.Context) (string, error) {
	res, err := f.client.GetUIName(ctx, &featurepb.Empty{})
	if err != nil {
		return "", err
	}

	return res.GetName(), nil
}

func (f *Feature) GetSurvey(ctx context.Context) (*survey.Survey, error) {
	res, err := f.client.GetSurvey(ctx, &featurepb.Empty{})
	if err != nil {
		return nil, err
	}

	return survey.FromProtoSurvey(res.GetSurvey()), nil
}

func (f *Feature) ValidateAnswers(ctx context.Context, answers map[string]interface{}) (map[string]interface{}, bool, error) {
	st, err := structpb.NewStruct(answers)
	if err != nil {
		return nil, false, err
	}

	res, err := f.client.ValidateAnswers(ctx, &featurepb.ValidateAnswersRequest{
		Answers: st,
	})
	if err != nil {
		return nil, false, err
	}

	return res.GetValues().AsMap(), res.GetShouldSave(), nil
}
