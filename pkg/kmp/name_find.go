package kmp

type NameFinder struct {
	NamePool []string
}

func NewNameFinder() *NameFinder {
	return &NameFinder{}
}

func (n *NameFinder) FindName(find string) []string {
	var res []string
	for _, v := range n.NamePool {
		pos := findSubstring(v, find)
		if pos != -1 {
			res = append(res, v)
		}
	}
	return res
}
