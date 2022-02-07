package scraping

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

type versionResponse struct {
	Version string `json:"version"`
}

type nodeClient struct {
	cl *client.Client
}

// ATTENTION! `url` MUST BE validated for proper format before passing to this function.
func newNodeClient(url string, timeout time.Duration) *nodeClient {
	opts := client.Options{
		BaseUrl: url,
		Client:  &http.Client{Timeout: timeout},
	}
	// The error can be safely ignored because `NewClient` function only checks the number of passed `opts`
	cl, _ := client.NewClient(opts)
	return &nodeClient{cl: cl}
}

func (c *nodeClient) version(ctx context.Context) (string, error) {
	versionRequest, err := http.NewRequest("GET", c.cl.GetOptions().BaseUrl+"/node/version", nil)
	if err != nil {
		log.Printf("Version request failed: %v", err)
		return "", err
	}
	versionRequest.Close = true
	resp := new(versionResponse)
	_, err = c.cl.Do(ctx, versionRequest, resp)
	if err != nil {
		log.Printf("Version request failed: %v", err)
		return "", err
	}
	return resp.Version, nil
}

func (c *nodeClient) height(ctx context.Context) (int, error) {
	height, _, err := c.cl.Blocks.Height(ctx)
	if err != nil {
		return 0, err
	}
	return int(height.Height), nil
}

func (c *nodeClient) stateHash(ctx context.Context, height int) (*proto.StateHash, error) {
	sh, _, err := c.cl.Debug.StateHash(ctx, uint64(height))
	if err != nil {
		return nil, err
	}
	return sh, nil
}
