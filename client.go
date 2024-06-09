package okxSwap

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qushedo/okxSwap/types"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	secretKey  []byte
	passphrase string
	restUrl    string
	client     *http.Client
}

func NewClient(apiKey, secretKey, passPhrase string) *Client {
	return &Client{
		apiKey:     apiKey,
		secretKey:  []byte(secretKey),
		passphrase: passPhrase,
		restUrl:    "https://www.okx.com",
		client:     http.DefaultClient,
	}
}

func (c *Client) sign(method, path, body string) (string, string) {
	format := "2006-01-02T15:04:05.999Z07:00"
	t := time.Now().UTC().Format(format)
	ts := fmt.Sprint(t)
	s := ts + method + path + body
	p := []byte(s)
	h := hmac.New(sha256.New, c.secretKey)
	h.Write(p)
	return ts, base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (c *Client) Do(method, endpoint string, body map[string]interface{}) ([]byte, error) {
	url := c.restUrl + endpoint

	orderBody, _ := json.Marshal(body)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(orderBody))
	if err != nil {
		return nil, err
	}

	timestamp, signature := c.sign(method, endpoint, string(orderBody))

	req.Header.Set("OK-ACCESS-KEY", c.apiKey)
	req.Header.Set("OK-ACCESS-SIGN", signature)
	req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("OK-ACCESS-PASSPHRASE", c.passphrase)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error: %s", responseBody)
	}

	return responseBody, nil
}

func (c *Client) EstimateQuote(fromCurrency, toCurrency, amountCurrency string, toAmount float64) (*types.QuoteResponse, error) {
	order := map[string]interface{}{
		"baseCcy":  toCurrency,
		"quoteCcy": fromCurrency,
		"side":     "buy",
		"rfqSz":    toAmount,
		"rfqSzCcy": amountCurrency,
	}
	resp, err := c.Do("POST", "/api/v5/asset/convert/estimate-quote", order)

	var quoteResp types.QuoteResponse
	err = json.Unmarshal(resp, &quoteResp)
	if err != nil {
		return nil, err
	}

	if quoteResp.Code != "0" {
		return nil, errors.New(quoteResp.Msg)
	}

	return &quoteResp, err
}

func (c *Client) Convert(fromCurrency, toCurrency, amountCurrency string, amount float64) (*types.ConvertTradeResponse, error) {
	estimateQuote, err := c.EstimateQuote(fromCurrency, toCurrency, amountCurrency, amount)
	if err != nil {
		return nil, err
	}

	order := map[string]interface{}{
		"quoteId":  estimateQuote.Data[0].QuoteId,
		"baseCcy":  toCurrency,
		"quoteCcy": fromCurrency,
		"side":     "buy",
		"szCcy":    amountCurrency,
		"sz":       amount,
	}
	resp, err := c.Do("POST", "/api/v5/asset/convert/trade", order)

	var convertTradeResp types.ConvertTradeResponse
	err = json.Unmarshal(resp, &convertTradeResp)
	if err != nil {
		return nil, err
	}

	switch convertTradeResp.Msg {
	case "":
		return &convertTradeResp, nil
	case types.NotEnoughBalance.Error():
		return &convertTradeResp, types.NotEnoughBalance
	default:
		return &convertTradeResp, errors.New(convertTradeResp.Msg)
	}
}

func (c *Client) GetCurrency(currency string) (*types.CurrencyInfoResponse, error) {
	resp, err := c.Do("GET", fmt.Sprintf("/api/v5/asset/currencies?ccy=%s", currency), nil)

	var currencyResp types.CurrencyInfoResponse
	err = json.Unmarshal(resp, &currencyResp)
	if err != nil {
		return nil, err
	}

	if currencyResp.Code != "0" {
		return nil, errors.New(currencyResp.Msg)
	}

	return &currencyResp, err
}

func (c *Client) Withdrawal(currency, toAddr string, amount float64) (*types.WithdrawalResponse, error) {
	currencyData, err := c.GetCurrency(currency)
	if err != nil {
		return nil, err
	}

	order := map[string]interface{}{
		"ccy":        currency,
		"amt":        amount,
		"dest":       "4",
		"toAddr":     toAddr,
		"fee":        currencyData.Data[0].MinFee,
		"walletType": "private",
	}
	resp, err := c.Do("POST", "/api/v5/asset/withdrawal", order)

	var withdrawalResp types.WithdrawalResponse
	err = json.Unmarshal(resp, &withdrawalResp)
	if err != nil {
		return nil, err
	}

	if withdrawalResp.Code != "0" {
		return nil, errors.New(withdrawalResp.Msg)
	}

	return &withdrawalResp, err
}

func (c *Client) WaitForWithdrawal(wdID string) (bool, error) {
	for {
		response, err := c.Do("GET", fmt.Sprintf("/api/v5/asset/withdrawal-history?wdId=%s", wdID), nil)
		if err != nil {
			log.Println(err)
			continue
		}

		var withdrawalStatus types.WithdrawalStatusResponse
		if err := json.Unmarshal(response, &withdrawalStatus); err != nil {
			return false, err
		}

		if withdrawalStatus.Code != "0" {
			log.Println(withdrawalStatus.Msg)
			continue
		}

		if len(withdrawalStatus.Data) > 0 && withdrawalStatus.Data[0].State == "2" {
			return true, nil
		} else if len(withdrawalStatus.Data) > 0 && withdrawalStatus.Data[0].State == "-2" {
			return false, fmt.Errorf("withdrawal %s canceled", wdID)
		} else if len(withdrawalStatus.Data) > 0 && withdrawalStatus.Data[0].State == "-1" {
			return false, fmt.Errorf("withdrawal %s failed", wdID)
		} else {
			time.Sleep(3 * time.Second)
			continue
		}
	}
}
func (c *Client) CheckForExistsWithdraw() (bool, error) {
	response, err := c.Do("GET", "/api/v5/asset/withdrawal-history", nil)
	if err != nil {
		return false, err
	}

	var withdrawalStatus types.WithdrawalStatusResponse
	if err := json.Unmarshal(response, &withdrawalStatus); err != nil {
		return false, err
	}

	if withdrawalStatus.Code != "0" {
		return false, errors.New(withdrawalStatus.Msg)
	}

	return withdrawalStatus.Data[0].State != "2" && withdrawalStatus.Data[0].State != "-1" && withdrawalStatus.Data[0].State != "-2", err
}

func (c *Client) GetBalance(currency string) (*types.AccountBalanceResponseFloat, error) {
	response, err := c.Do("GET", fmt.Sprintf("/api/v5/asset/balances?ccy=%s", currency), nil)
	if err != nil {
		return nil, err
	}

	var balanceResp types.AccountBalanceResponse
	if err := json.Unmarshal(response, &balanceResp); err != nil {
		return nil, err
	}

	floatResp, err := balanceResp.Float()
	if err != nil {
		return nil, err
	}

	if balanceResp.Code != "0" {
		return nil, errors.New(balanceResp.Msg)
	}

	return &floatResp, nil
}
