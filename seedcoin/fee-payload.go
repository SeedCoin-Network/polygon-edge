package seedcoin

// {"status":"active","updated":"2022-10-13 19:00:00","price":0.1}
type FeePayload struct {
	Status  string  `json:"status"`
	Updated string  `json:"updated"`
	Price   float64 `json:"price"`
}
