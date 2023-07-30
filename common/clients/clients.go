package clients

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients/dead_nation"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/files"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/payments"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/receipts"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/scoreboard"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/spreadsheets"
)

type Clients struct {
	Files        files.ClientWithResponsesInterface
	Payments     payments.ClientWithResponsesInterface
	Receipts     receipts.ClientWithResponsesInterface
	Scoreboard   scoreboard.ClientWithResponsesInterface
	Spreadsheets spreadsheets.ClientWithResponsesInterface
	DeadNation   dead_nation.ClientWithResponsesInterface
}

func NewClients(
	gatewayAddress string,
	requestEditorFn RequestEditorFn,
) (*Clients, error) {
	return NewClientsWithHttpClient(gatewayAddress, requestEditorFn, http.DefaultClient)
}

func NewClientsWithHttpClient(
	gatewayAddress string,
	requestEditorFn RequestEditorFn,
	httpDoer HttpDoer,
) (*Clients, error) {
	if gatewayAddress == "" {
		return nil, fmt.Errorf("gateway address is required")
	}

	if requestEditorFn == nil {
		requestEditorFn = func(_ context.Context, _ *http.Request) error {
			return nil
		}
	}

	filesClient, err := newClient(
		gatewayAddress,
		"files-api",
		files.NewClientWithResponses,
		files.WithRequestEditorFn(files.RequestEditorFn(requestEditorFn)),
		files.WithRequestEditorFn(files.RequestEditorFn(requestEditorFn)),
		files.WithHTTPClient(httpDoer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create files client: %w", err)
	}

	paymentsClient, err := newClient(
		gatewayAddress,
		"payments-api",
		payments.NewClientWithResponses,
		payments.WithRequestEditorFn(payments.RequestEditorFn(requestEditorFn)),
		payments.WithHTTPClient(httpDoer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create payments client: %w", err)
	}

	receiptsClient, err := newClient(
		gatewayAddress,
		"receipts-api",
		receipts.NewClientWithResponses,
		receipts.WithRequestEditorFn(receipts.RequestEditorFn(requestEditorFn)),
		receipts.WithHTTPClient(httpDoer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create receipts client: %w", err)
	}

	scoreboardClient, err := newClient(
		gatewayAddress,
		"scoreboard-api",
		scoreboard.NewClientWithResponses,
		scoreboard.WithRequestEditorFn(scoreboard.RequestEditorFn(requestEditorFn)),
		scoreboard.WithHTTPClient(httpDoer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create scoreboard client: %w", err)
	}

	spreadsheetsClient, err := newClient(
		gatewayAddress,
		"spreadsheets-api",
		spreadsheets.NewClientWithResponses,
		spreadsheets.WithRequestEditorFn(spreadsheets.RequestEditorFn(requestEditorFn)),
		spreadsheets.WithHTTPClient(httpDoer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create spreadsheets client: %w", err)
	}

	deadNationClient, err := newClient(
		gatewayAddress,
		"dead-nation-api",
		dead_nation.NewClientWithResponses,
		dead_nation.WithRequestEditorFn(dead_nation.RequestEditorFn(requestEditorFn)),
		dead_nation.WithHTTPClient(httpDoer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create dead-nation client: %w", err)
	}

	return &Clients{
		Files:        filesClient,
		Payments:     paymentsClient,
		Receipts:     receiptsClient,
		Scoreboard:   scoreboardClient,
		Spreadsheets: spreadsheetsClient,
		DeadNation:   deadNationClient,
	}, nil
}

type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestEditorFn func(ctx context.Context, req *http.Request) error

func newClient[Client any, ClientOption any](
	gatewayAddress string,
	serviceName string,
	clientConstructor func(server string, opts ...ClientOption) (*Client, error),
	requestEditorFn ...ClientOption,
) (*Client, error) {
	apiServerAddr, err := url.JoinPath(gatewayAddress, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create files api server address of %s: %w", serviceName, err)
	}

	apiClient, err := clientConstructor(apiServerAddr, requestEditorFn...)
	if err != nil {
		return nil, fmt.Errorf("failed to create files client: %w", err)
	}

	return apiClient, nil
}
