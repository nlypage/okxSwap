package types

type CurrencyInfoResponse struct {
	Code string         `json:"code"`
	Msg  string         `json:"msg"`
	Data []CurrencyInfo `json:"data"`
}

type CurrencyInfo struct {
	CanDep               bool   `json:"canDep"`
	CanInternal          bool   `json:"canInternal"`
	CanWd                bool   `json:"canWd"`
	Ccy                  string `json:"ccy"`
	Chain                string `json:"chain"`
	DepQuotaFixed        string `json:"depQuotaFixed"`
	DepQuoteDailyLayer2  string `json:"depQuoteDailyLayer2"`
	LogoLink             string `json:"logoLink"`
	MainNet              bool   `json:"mainNet"`
	MaxFee               string `json:"maxFee"`
	MaxFeeForCtAddr      string `json:"maxFeeForCtAddr"`
	MaxWd                string `json:"maxWd"`
	MinDep               string `json:"minDep"`
	MinDepArrivalConfirm string `json:"minDepArrivalConfirm"`
	MinFee               string `json:"minFee"`
	MinFeeForCtAddr      string `json:"minFeeForCtAddr"`
	MinWd                string `json:"minWd"`
	MinWdUnlockConfirm   string `json:"minWdUnlockConfirm"`
	Name                 string `json:"name"`
	NeedTag              bool   `json:"needTag"`
	UsedDepQuotaFixed    string `json:"usedDepQuotaFixed"`
	UsedWdQuota          string `json:"usedWdQuota"`
	WdQuota              string `json:"wdQuota"`
	WdTickSz             string `json:"wdTickSz"`
}
