package bittrex

type Fill struct {
	OrderType string
	Price     float64 `json:"Rate"`
	Quantity  float64
	TimeStamp string
}
