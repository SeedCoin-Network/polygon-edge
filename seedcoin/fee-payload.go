package seedcoin

import "fmt"

// {"status":"active","updated":"2022-10-13 19:00:00","price":0.1}
type FeePayload struct {
	Status  string  `json:"status"`
	Updated string  `json:"updated"`
	Price   float64 `json:"price"`
}

func (f FeePayload) Description() string {
	return fmt.Sprintf("Status: %s, ", f.Status) +
		fmt.Sprintf("Updated: %s, ", f.Updated) +
		fmt.Sprintf("Price: %f", f.Price)
}
