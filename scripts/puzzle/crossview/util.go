package crossview

import "math/rand"

func chance(perc float64) bool {
	return rand.Float64() <= perc
}
