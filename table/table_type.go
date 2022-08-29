package table

type PowerData struct {
	ID           int64 `json:"id"`
	TotalPower   int64 `json:"totalPower"`
	FreePower    int64 `json:"freePower"`
	PayPowerLay1 int64 `json:"payPowerLay1"`
	PayPowerLay2 int64 `json:"payPowerLay2"`
	PayPowerLay3 int64 `json:"payPowerLay3"`
	PayPowerLay4 int64 `json:"payPowerLay4"`
	PayPowerLay5 int64 `json:"payPowerLay5"`
	AllPower     []int64
}
