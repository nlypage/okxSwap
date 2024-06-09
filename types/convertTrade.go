package types

type ConvertTradeResponse struct {
	Code string             `json:"code"`
	Data []ConvertTradeData `json:"data"`
	Msg  string             `json:"msg"`
}

type ConvertTradeData struct {
	BaseCcy     string `json:"baseCcy"`
	ClTReqId    string `json:"clTReqId"`
	FillBaseSz  string `json:"fillBaseSz"`
	FillPx      string `json:"fillPx"`
	FillQuoteSz string `json:"fillQuoteSz"`
	InstId      string `json:"instId"`
	QuoteCcy    string `json:"quoteCcy"`
	QuoteId     string `json:"quoteId"`
	Side        string `json:"side"`
	State       string `json:"state"`
	TradeId     string `json:"tradeId"`
	Ts          string `json:"ts"`
}
