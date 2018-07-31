package types

import "time"

//RawTime representa um time
type RawTime []byte

//Time para promover conversão
func (t RawTime) Time() (time.Time, error) {
	return time.Parse("15:04:05", string(t))
}
