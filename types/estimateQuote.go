package types

type QuoteResponse struct {
	Code string      `json:"code"`
	Data []QuoteData `json:"data"`
	Msg  string      `json:"msg"`
}

type QuoteData struct {
	BaseCcy   string `json:"baseCcy"`
	BaseSz    string `json:"baseSz"`
	ClQReqId  string `json:"clQReqId"`
	CnvtPx    string `json:"cnvtPx"`
	OrigRfqSz string `json:"origRfqSz"`
	QuoteCcy  string `json:"quoteCcy"`
	QuoteId   string `json:"quoteId"`
	QuoteSz   string `json:"quoteSz"`
	QuoteTime string `json:"quoteTime"`
	RfqSz     string `json:"rfqSz"`
	RfqSzCcy  string `json:"rfqSzCcy"`
	Side      string `json:"side"`
	TtlMs     string `json:"ttlMs"`
}
