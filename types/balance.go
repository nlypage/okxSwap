package types

import "strconv"

type AccountBalanceResponse struct {
	Code string        `json:"code"`
	Data []BalanceData `json:"data"`
	Msg  string        `json:"msg"`
}

type BalanceData struct {
	AvailBal  string `json:"availBal"`
	Bal       string `json:"bal"`
	Ccy       string `json:"ccy"`
	FrozenBal string `json:"frozenBal"`
}

type AccountBalanceResponseFloat struct {
	Code string             `json:"code"`
	Data []BalanceDataFloat `json:"data"`
	Msg  string             `json:"msg"`
}

type BalanceDataFloat struct {
	AvailBal  float64 `json:"availBal"`
	Bal       float64 `json:"bal"`
	Ccy       string  `json:"ccy"`
	FrozenBal float64 `json:"frozenBal"`
}

func (r AccountBalanceResponse) Float() (AccountBalanceResponseFloat, error) {
	var result AccountBalanceResponseFloat
	result.Code = r.Code
	result.Msg = r.Msg

	for _, data := range r.Data {
		availBal, err := strconv.ParseFloat(data.AvailBal, 64)
		if err != nil {
			return AccountBalanceResponseFloat{}, err
		}

		bal, err := strconv.ParseFloat(data.Bal, 64)
		if err != nil {
			return AccountBalanceResponseFloat{}, err
		}

		frozenBal, err := strconv.ParseFloat(data.FrozenBal, 64)
		if err != nil {
			return AccountBalanceResponseFloat{}, err
		}

		newData := BalanceDataFloat{
			AvailBal:  availBal,
			Bal:       bal,
			Ccy:       data.Ccy,
			FrozenBal: frozenBal,
		}

		result.Data = append(result.Data, newData)
	}

	return result, nil
}
