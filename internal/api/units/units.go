// Package units
package units

import (
	"fmt"
	"math"
)

func BytesToBits(bytes int64) int64 {
	return bytes * 8
}

func BitsHuman(bits int64) string {
	units := []string{"b", "kb", "Mb", "Gb", "Tb"}
	power := int(math.Floor(math.Log(float64(bits)) / math.Log(1000)))
	if bits == 0 {
		power = 0
	}
	power = min(power, len(units)-1)

	return fmt.Sprintf("%.2f %sps", float64(bits)/math.Pow(1000, float64(power)), units[power])
}
