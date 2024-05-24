package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type WalletClient struct {
	BaseURL string
}

func NewWalletClient(baseURL string) *WalletClient {
	return &WalletClient{BaseURL: baseURL}
}

type UpdateBalanceRequest struct {
	AccountNumber string  `json:"account_number"`
	Amount        float64 `json:"amount"`
}

func (client *WalletClient) UpdateBalance(ctx context.Context, accountNumber string, amount float64) error {
	url := fmt.Sprintf("%s/update_balance", client.BaseURL)
	reqBody := UpdateBalanceRequest{
		AccountNumber: accountNumber,
		Amount:        amount,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update balance, status code: %d", resp.StatusCode)
	}

	return nil
}
