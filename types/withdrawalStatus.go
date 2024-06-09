package types

type WithdrawalStatusResponse struct {
	Code string                 `json:"code"`
	Msg  string                 `json:"msg"`
	Data []WithdrawalStatusData `json:"data"`
}

type WithdrawalStatusData struct {
	State string `json:"state"`
	WdID  string `json:"wdId"`
}
