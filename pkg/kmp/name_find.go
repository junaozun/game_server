package kmp

type NameFinder struct {
	NamePool []string
}

func (n *NameFinder) FindName(find string) []string {
	var res []string
	for _, v := range n.NamePool {
		pos := strStrV2(v, find)
		if pos != -1 {
			res = append(res, v)
		}
	}
	return res
}
