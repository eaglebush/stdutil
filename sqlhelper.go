package stdutil

// BoolToInt returns 1 when value is true, 0 when value is false
func BoolToInt(value bool) int {
	if value {
		return 1
	} else {
		return 0
	}
}

// BoolToUInt returns 1 when value is true, 0 when value is false
func BoolToUInt(value bool) uint {
	if value {
		return 1
	} else {
		return 0
	}
}
