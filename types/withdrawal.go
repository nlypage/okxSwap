package types

type WithdrawalResponse struct {
	Code string           `json:"code"`
	Msg  string           `json:"msg"`
	Data []WithdrawalData `json:"data"`
}

type WithdrawalData struct {
	Amount   string `json:"amt"`
	WdID     string `json:"wdId"`
	Currency string `json:"ccy"`
	ClientID string `json:"clientId"`
	Chain    string `json:"chain"`
}
