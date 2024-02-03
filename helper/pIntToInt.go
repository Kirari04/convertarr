package helper

func PIntToInt(str *int) int {
	if str == nil {
		return 0
	}

	return *str
}
