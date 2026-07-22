package domains

import "fmt"

type QRString struct {
	T  string  `json:"t"`
	S  float64 `json:"s"`
	FN string  `json:"fn"`
	FD string  `json:"i"`
	FP string  `json:"fp"`
	N  int64   `json:"n"`
}

func (self QRString) String() string {
	return fmt.Sprintf("t=%s&s=%.2f&fn=%s&i=%s&fp=%s&n=%d", self.T, self.S, self.FN, self.FD, self.FP, self.N)
}
