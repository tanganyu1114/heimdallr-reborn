package v1

import (
	"context"
	"net/http"

	epclientv1 "gin-vue-admin/pkg/client/v1/endpoint"
	mwclientv1 "gin-vue-admin/pkg/client/v1/middleware"
	svcclientv1 "gin-vue-admin/pkg/client/v1/service"
	txpclientv1 "gin-vue-admin/pkg/client/v1/transport"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
)

var clientContext context.Context

// Client is the main Heimdallr API client
type Client struct {
	client  *http.Client
	baseURL string
	svcclientv1.Factory
}

// Close closes the client
func (c *Client) Close() error {
	return nil
}

// newClient creates a new client instance
func newClient(ctx context.Context, client *http.Client, baseURL string) (*Client, error) {
	txp, err := txpclientv1.NewTransport(client, baseURL)
	if err != nil {
		return nil, err
	}

	mws := mwclientv1.GetMiddlewares()
	for name, mw := range mws {
		logV1.Infof("Loading middleware: %q", name)
		txp = mw(txp)
	}

	ep, err := epclientv1.NewEndpoint(txp)
	if err != nil {
		return nil, err
	}

	svc, err := svcclientv1.NewService(ctx, ep)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:  client,
		baseURL: baseURL,
		Factory: svc,
	}, nil
}

// NewClient creates a new Heimdallr API client
func NewClient(baseURL string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if clientContext == nil {
		clientContext = context.Background()
	}

	return newClient(clientContext, httpClient, baseURL)
}

// SetClientContext sets the global client context
func SetClientContext(ctx context.Context) {
	clientContext = ctx
}
