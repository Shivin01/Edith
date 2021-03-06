package jenkins

import (
	"context"
	"net/http"
	"strings"

	"github.com/Shivin01/Edith/slack-bot/pkg/client"
	"github.com/Shivin01/Edith/slack-bot/pkg/config"
	"github.com/bndr/gojenkins"
)

// Client is an interface representing used jenkins functions of gojenkins.
type Client interface {
	GetJob(ctx context.Context, id string) (*gojenkins.Job, error)
	BuildJob(ctx context.Context, name string, params map[string]string) (int64, error)
	GetAllNodes(ctx context.Context) ([]*gojenkins.Node, error)
}

// GetClient created Jenkins client with given options/credentials
func GetClient(cfg config.Jenkins) (Client, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	return createJenkinsClient(context.TODO(), client.GetHTTPClient(), cfg)
}

// implementation of Client interface. proxies to gojenkins with additional handling for inner jenkins jobs.
type jenkinsClientImpl struct {
	client *gojenkins.Jenkins
}

func createJenkinsClient(ctx context.Context, httpClient *http.Client, cfg config.Jenkins) (*jenkinsClientImpl, error) {
	jenkins := gojenkins.CreateJenkins(
		httpClient,
		cfg.Host,
		cfg.Username,
		cfg.Password,
	)

	jenkinsClient, err := jenkins.Init(ctx)
	if err != nil {
		return nil, err
	}

	return &jenkinsClientImpl{
		client: jenkinsClient,
	}, nil
}

func (c *jenkinsClientImpl) GetJob(ctx context.Context, id string) (*gojenkins.Job, error) {
	// split jobs id by "/"" to be able to access inner job
	jobs := strings.Split(id, "/")

	jobsCount := len(jobs)
	if jobsCount > 1 {
		return c.client.GetJob(ctx, jobs[jobsCount-1], jobs[:jobsCount-1]...)
	}

	return c.client.GetJob(ctx, id)
}

func (c *jenkinsClientImpl) BuildJob(ctx context.Context, name string, params map[string]string) (int64, error) {
	return c.client.BuildJob(ctx, name, params)
}

func (c *jenkinsClientImpl) GetAllNodes(ctx context.Context) ([]*gojenkins.Node, error) {
	return c.client.GetAllNodes(ctx)
}
