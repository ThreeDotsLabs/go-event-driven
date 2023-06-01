package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type IssueReceiptRequest struct {
	TicketID      string `json:"ticket_id"`
	CustomerEmail string `json:"customer_email"`
	Price         Money  `json:"price"`
}

type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type ReceiptsServiceClient struct{}

func (c ReceiptsServiceClient) IssueReceipt(ctx context.Context, request IssueReceiptRequest) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(payload)

	endpoint := fmt.Sprintf("%s/receipts-api/receipts", os.Getenv("GATEWAY_ADDR"))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

type SpreadsheetRow struct {
	Columns []string `json:"columns"`
}

type SpreadsheetsAPIClient struct{}

func (c SpreadsheetsAPIClient) AppendRow(ctx context.Context, spreadsheetName string, row SpreadsheetRow) error {
	payload, err := json.Marshal(row)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(payload)
	endpoint := fmt.Sprintf("%s/spreadsheets-api/sheets/%v/rows", os.Getenv("GATEWAY_ADDR"), spreadsheetName)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
