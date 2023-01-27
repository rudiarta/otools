package outils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	Byte          uint64 = 1
	maxPacketSize        = int(65000 * Byte)
)

func ToString(v interface{}) (s string) {
	switch val := v.(type) {
	case error:
		if val != nil {
			s = val.Error()
		}
	case string:
		s = val
	case int:
		s = strconv.Itoa(val)
	default:
		b, _ := json.Marshal(val)
		s = string(b)
	}

	if len(s) >= maxPacketSize {
		s = fmt.Sprintf("overflow, size is: %d, max: %d", len(s), maxPacketSize)
	}
	return
}
