package maplib

func Choice(m map[string]string) string {
	for _, v := range m {
		return v
	}

	return ""
}
