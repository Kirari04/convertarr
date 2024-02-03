package helper

func PStrToStr(str *string) string {
	if str == nil {
		return ""
	}

	return *str
}
