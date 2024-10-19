package types

type TransferAccount int

const (
	Funding TransferAccount = 6
	Trading TransferAccount = 18
)

type TransferResponse struct {
	Code string         `json:"code"`
	Msg  string         `json:"msg"`
	Data []TransferData `json:"data"`
}

type TransferData struct {
	TransferID  string  `json:"transId"`
	Currency    string  `json:"ccy"`
	ClientID    string  `json:"clientId"`
	FromAccount int     `json:"from,string"`
	ToAccount   int     `json:"to,string"`
	Amount      float64 `json:"amt,string"`
}
