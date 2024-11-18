package kmp

import (
	"fmt"
	"time"
)

// 基于 KMP 算法实现查找字符串子串函数
func findSubstring(haystack, needle string) int {
	// 子串长度=0
	if len(needle) == 0 {
		return 0
	}
	// 主串长度=0，或者主串长度小于子串长度
	if len(haystack) == 0 || len(haystack) < len(needle) {
		return -1
	}
	// 子串长度=1时单独判断
	if len(needle) == 1 {
		i := 0
		for ; i < len(haystack); i++ {
			if haystack[i] == needle[0] {
				return i
			}
		}
		return -1
	}

	// 其他情况走 KMP 算法
	return kmpSearch(haystack, needle)
}

// KMP 算法实现函数
func kmpSearch(s, p string) int {
	n := len(s)             // 主串长度
	m := len(p)             // 模式串长度
	next := generateNext(p) // 生成 next 数组
	i, j := 0, 0
	for i < n && j < m {
		// 如果主串字符和模式串字符不相等，
		// 更新模式串坏字符下标位置为好前缀最长可匹配前缀子串尾字符下标
		// 然后从这个位置重新开始与主串匹配
		// 相当于前面提到的把模式串向后移动 j - k 位
		if j == -1 || s[i] == p[j] {
			i++
			j++
		} else {
			j = next[j]
		}
	}
	if j == m {
		// 完全匹配，返回下标位置
		return i - j
	} else {
		return -1
	}
}

// 生成 next 数组
func generateNext(p string) []int {
	m := len(p)
	next := make([]int, m, m)
	next[0] = -1
	next[1] = 0
	i, j := 0, 1 // 前缀子串、后缀子串起始位置
	// 因为是通过最长可匹配前缀子串计算，所以 j 的值最大不会超过 m-1
	for j < m-1 {
		if i == -1 || p[i] == p[j] {
			i++
			j++
			// 设置当前最长可匹配前缀子串结尾字符下标
			next[j] = i
		} else {
			// 如果 p[i] != p[j]，回到上一个最长可匹配前缀子串结尾字符下标
			i = next[i]
		}
	}
	return next
}

func TestGO() {
	var a int
	go func(a int) {
		for {
			panic("xxx")
			time.Sleep(1 * time.Second)
			fmt.Println("111")
		}
	}(a)
}

func (hMgr *ETCDHttpMgr) CallService(c *engine.Context) error {
	log.Log.Debug().Interface("server", c.Server).Interface("args", c.Args).Msg("ETCDHttpMgr:CallService")
	xclient := hMgr.getXclient(c.Server.AppName)
	if xclient == nil {
		clog.Log.Error().Str("appName", c.Server.AppName).Interface("req", c.R).Msg("CallServicePath :get service client nil ")
		return fmt.Errorf("appName:%s,path :%s,CallService xclient nil", c.Server.AppName, c.Server.Path)
	}
	hMgr.Proxy.Director = DirectorReq(c.Ctx, c)
	xclient.SetProxy(hMgr.Proxy)
	err := xclient.Proxy(c.W, c.R)
	if errors.Is(err, selector.ErrDebugNoAvailable) {
		log.Log.Debug().Interface("args", c.Args).Interface("server", c.Server).Msg("CallService:ErrDebugNoAvailable")
		c.R.Header.Del(constants.RequestUserName)
		err = xclient.Proxy(c.W, c.R)
	}
	if err != nil {
		clog.Log.Err(err).Str("appName", c.Server.AppName).Interface("args", c.Args).Interface("req", c.R).Msg("CallServicePath:Invoke has err")
	}
	return err
}
