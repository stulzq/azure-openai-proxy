package apis

var (
	ignoreHeaders = map[string]int{
		"Content-Type":      1,
		"Transfer-Encoding": 1,
		"Date":              1,
	}
)
