package nils

func String(in *string, defaultValue ...string) string {
	result := ""
	if len(defaultValue) > 0 {
		result = defaultValue[0]
	}
	if in != nil {
		result = *in
	}
	return result
}
