package kmp

func FindName(set []string, find string) []string {
	var res []string
	for _, v := range set {
		pos := strStrV2(v, find)
		if pos != -1 {
			res = append(res, v)
		}
	}
	return res
}
