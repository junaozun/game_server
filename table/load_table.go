package table

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

type iTable interface {
	load(dir string) error
	reload(dir string) (bool, error)
	GetFileName() string
	MD5() string
}

var dir = "./data/"

// LoadTable 加载所有表格
func LoadTable() {
	for _, v := range tableList {
		filepath := dir + v.GetFileName()
		if e := v.load(filepath); e != nil {
			// util.LOG.Error("load csv table err", e)
			panic("file:" + filepath + "	error:" + e.Error())
		}
	}
	initSecond()
	checkTable()
	runtime.GC()
}

// Reload 重新加载所有表格
// 说明：
// 1、Reload的表不会减少条数，比如A表原来有100条，然后给改成99条，Reload完还是100条
// 2、Reload不会改变数组长度，只能改变值，[1,2,3]然后表改成[2,2]，Reload后实际是[2,2,3]
func Reload() []string {
	// 中间处理不可预料得错误一定要恢复回来
	defer func() {
		if err := recover(); nil != err {
			stack := debug.Stack()
			fmt.Println("[Table.Reload]", err, stack)
			// util.LOG.Critical("[Table.Reload] %v %s", err, stack)
		}
	}()

	var ret []string
	begin := time.Now()
	for _, v := range tableList {
		filepath := dir + v.GetFileName()
		if ok, e := v.reload(filepath); nil != e {
			log.Println(e.Error())
		} else if ok {
			ret = append(ret, v.GetFileName())
		}
	}
	cost := time.Since(begin)
	log.Printf("reload table cost time: [%d]ms", cost.Milliseconds())
	if len(ret) > 0 {
		initSecond()
	}
	return ret
}

// GetFileModTime 获取文件修改时间
func GetFileMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	} else {
		return hex.EncodeToString(h.Sum(nil)), nil
	}
}
