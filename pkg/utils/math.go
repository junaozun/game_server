package utils

func MinInt(x,y int) int  {
	if x > y {
		return y
	}
	return x
}

func MaxInt(x,y int) int  {
	if x > y {
		return x
	}
	return y
}

func AbsInt(x int) int{
	if x >= 0 {
		return x
	}
	return -x
}