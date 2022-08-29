package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	template  = flag.String("template", "table_template.txt", "template file path")
	tablelist = flag.String("tablelist", "table_list.txt", "table data list file")
	out       = flag.String("out", "../../src/loadTable", "out dir")
)

func initPath() {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("get abs path error")
	}
	path += "/"
	*template = path + *template
	*tablelist = path + *tablelist
	*out = path + *out
}

func formatPath() {
	// *template = filepath.ToSlash(*template)
	// *tablelist = filepath.ToSlash(*tablelist)
	// *out = filepath.ToSlash(*out)

	*template = filepath.Clean(*template)
	*tablelist = filepath.Clean(*tablelist)
	*out = filepath.Clean(*out)
}

func main() {
	if len(os.Args) == 1 {
		initPath()
	} else if len(os.Args) == 4 {
		flag.Parse()
	} else {
		panic(`运行参数错误，格式如下
			方式1: 直接运行可执行文件
			方式2: 参数1(模板文件路径)，参数二（列表文件路径），参数三（输出目录）
		`)
	}
	formatPath()
	fmt.Println("template=", *template)
	fmt.Println("tablelist=", *tablelist)
	fmt.Println("out=", *out)
	fmt.Println()
	fmt.Println()
	data, err := ioutil.ReadFile(*template)
	if nil != err {
		fmt.Println("read file error", *template)
		return
	}
	orgin := string(data)

	data1, err1 := ioutil.ReadFile(*tablelist)
	if nil != err1 {
		fmt.Println("read file error", *tablelist)
		return
	}

	ts := strings.TrimSpace(string(data1))

	datas := strings.Split(ts, "\n")

	for _, v := range datas {
		ds := strings.Split(v, "|")
		typ := ds[0]
		idTyp := "int64"
		idV := "ID"
		ret := strings.Replace(orgin, "$1", typ, -1)
		if len(ds) > 1 {
			idTyp = ds[1]
		}
		ret = strings.Replace(ret, "$2", idTyp, -1)
		if len(ds) > 2 {
			idV = ds[2]
		}
		ret = strings.Replace(ret, "$3", idV, -1)

		fileName := fmt.Sprintf("%s/%s_table.go", *out, strings.ToLower(typ))
		fmt.Println("writefile begin:" + fileName)
		// fileName = filepath.ToSlash(fileName)
		fileName = filepath.Clean(fileName)

		err := ioutil.WriteFile(fileName, []byte(ret), 0666)
		if err != nil {
			panic(err.Error() + "	file: " + fileName)
		}
		fmt.Println("writefile success:" + fileName)
	}

	fmt.Println("")
	fmt.Println("success!!")
}
