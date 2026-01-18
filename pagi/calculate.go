package pagi

func CalculateLimit(limit, def, max uint) uint {
	if limit <= 0 {
		return def
	}
	if limit > max {
		return max
	}

	return limit
}
