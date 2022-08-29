package csvloader

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var DataBeginIdx = 5
var HeaderSeparator = "$"
var HeaderExportTag = "s"
var HeaderExportStrLine = 2
var HeaderExportStr = []string{"both", "server", "language", "client"}
var Comma = ','

// var Comment = '#'
var FiledTag = "json"

type Index map[interface{}]interface{}

type headerInfo struct {
	col  int    // header在csv中的列
	name string // header 名字
}

type csvloader struct {
	Comma rune
	// Comment    rune
	typeRecord reflect.Type
	records    []interface{}
	FileName   string
}

func (rf *csvloader) loadFromReader(r io.Reader) ([]interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var reader io.Reader = bytes.NewReader(b)
	if !validUTF8(b) {
		reader = transform.NewReader(reader, simplifiedchinese.GB18030.NewDecoder())
	}

	err = rf.read(reader)
	if err != nil {
		return nil, err
	}

	return rf.Record(), nil
}

func LoadCSVConfig(filename string, rt reflect.Type) ([]interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cl, err := newcsvloader(filename, rt)
	if err != nil {
		return nil, err
	}
	return cl.loadFromReader(f)
}

func LoadCSVConfigFromReader(r io.Reader, rt reflect.Type) ([]interface{}, error) {
	rf, err := newcsvloader("", rt)
	if err != nil {
		return nil, err
	}
	return rf.loadFromReader(r)
}

func newcsvloader(fileName string, typeRecord reflect.Type) (*csvloader, error) {
	if typeRecord == nil || typeRecord.Kind() != reflect.Struct {
		return nil, errors.New("st must be a struct")
	}

	// for i := 0; i < typeRecord.NumField(); i++ {
	// 	f := typeRecord.Field(i)

	// 	kind := f.Type.Kind()
	// 	switch kind {
	// 	case reflect.Bool:
	// 	case reflect.Int:
	// 	case reflect.Int8:
	// 	case reflect.Int16:
	// 	case reflect.Int32:
	// 	case reflect.Int64:
	// 	case reflect.Uint:
	// 	case reflect.Uint8:
	// 	case reflect.Uint16:
	// 	case reflect.Uint32:
	// 	case reflect.Uint64:
	// 	case reflect.Float32:
	// 	case reflect.Float64:
	// 	case reflect.String:
	// 	case reflect.Struct:
	// 	case reflect.Array:
	// 	case reflect.Slice:
	// 	case reflect.Map:
	// 	default:
	// 		return nil, fmt.Errorf("invalid type: %v %s",
	// 			f.Name, kind)
	// 	}

	// 	tag := f.Tag
	// 	if tag == "index" {
	// 		switch kind {
	// 		case reflect.Struct, reflect.Slice, reflect.Map:
	// 			return nil, fmt.Errorf("could not index %s field %v %v",
	// 				kind, i, f.Name)
	// 		}
	// 	}
	// }

	rf := new(csvloader)
	rf.typeRecord = typeRecord
	rf.FileName = fileName

	return rf, nil
}

func (rf *csvloader) parseHeader(header []string, typeRecord reflect.Type, lines [][]string) []headerInfo {
	headerMap := make([]headerInfo, 0, 8)
	for i, v := range header {
		needExport := false
		vArr := strings.Split(v, HeaderSeparator)
		if len(vArr) == 2 && strings.Contains(vArr[1], HeaderExportTag) {
			needExport = true
		}

		if !needExport && len(lines) > HeaderExportStrLine {
			for _, exportStr := range HeaderExportStr {
				if lines[HeaderExportStrLine][i] == exportStr {
					needExport = true
					break
				}
			}
		}

		if !needExport {
			continue
		}

		csvColumnHeader := vArr[0]
		// 强制把第一个字段修改为id
		if i == 0 && csvColumnHeader != "id" {
			csvColumnHeader = "id"
		}
		for j := 0; j < typeRecord.NumField(); j++ {
			f := typeRecord.Field(j)
			tag := f.Tag.Get(FiledTag)
			if tag == csvColumnHeader {
				headerMap = append(headerMap, headerInfo{
					col:  i,
					name: f.Name,
				})
				break
			}
		}
	}
	return headerMap
}

func (rf *csvloader) read(r io.Reader) error {

	if rf.Comma == 0 {
		rf.Comma = Comma
	}
	// if rf.Comment == 0 {
	//	rf.Comment = Comment
	// }
	reader := csv.NewReader(r)
	reader.LazyQuotes = true
	reader.Comma = rf.Comma
	// reader.Comment = rf.Comment
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// 空数据
	if len(lines) <= DataBeginIdx {
		rf.records = make([]interface{}, 0)
		return nil
	}

	// 类型数据
	typeRecord := rf.typeRecord

	// header解析
	headerMap := rf.parseHeader(lines[0], typeRecord, lines)

	// make records
	// 重置数据到从DataBeginIdx位置
	lines = lines[DataBeginIdx:]
	records := make([]interface{}, 0, len(lines))

	for n := 0; n < len(lines); n++ {
		line := lines[n]
		// 第一个字段为空则跳过解析
		if len(line) <= 0 || line[0] == "" || line[0][0] == '#' {
			continue
		}

		value := reflect.New(typeRecord)
		records = append(records, value.Interface())
		record := value.Elem()

		var firstColValue interface{}
		for _, info := range headerMap {
			i, v := info.col, info.name
			strField := line[i]
			if strField == "" {
				continue
			}

			f, ok := typeRecord.FieldByName(v)
			if !ok {
				return errors.New("type filed not found: " + v)
			}
			field := record.FieldByName(v)
			if !field.CanSet() {
				continue
			}

			if layout, ok := f.Tag.Lookup("time_parse"); ok {
				if err := rf.parseTime(field, strField, layout); err != nil {
					return fmt.Errorf("parse field (Line=%v, ColumnName=%v) error: %v",
						n+1+5, f.Tag.Get("json"), err)
				}
				continue
			}

			var err error

			kind := f.Type.Kind()
			isPtr := false
			if kind == reflect.Ptr {
				kind = f.Type.Elem().Kind()
				isPtr = true
			}
			if mlFiledName, ok := f.Tag.Lookup("multi_lang"); ok && kind == reflect.String {
				mlf := record.FieldByName(mlFiledName)
				if mlf.CanSet() {
					tagName, _ := f.Tag.Lookup("json")
					fn := path.Base(rf.FileName)
					fn = strings.TrimRight(fn, ".csv")
					str := fmt.Sprintf("%v_%v_%v", fn, tagName, firstColValue)
					if isPtr {
						mlf.SetPointer(unsafe.Pointer(&str))
					} else {
						mlf.SetString(str)
					}
				}
			}

			if kind == reflect.Bool {
				var v bool
				v, err = strconv.ParseBool(strField)
				if err == nil {
					if isPtr {
						field.SetPointer(unsafe.Pointer(&v))
					} else {
						field.SetBool(v)
					}
				}
			} else if kind == reflect.Int ||
				kind == reflect.Int8 ||
				kind == reflect.Int16 ||
				kind == reflect.Int32 ||
				kind == reflect.Int64 {

				if kind == reflect.Int64 { // 时间戳特殊处理
					if layout, ok := f.Tag.Lookup("time_stamp"); ok {
						if len(layout) == 0 {
							layout = "2006/1/2 15/4/5"
						}
						d, err := time.ParseInLocation(layout, strField, time.Local)
						if err != nil {
							return fmt.Errorf("parse field (Line=%v, ColumnName=\"%v\") error: %v", n+1+5, f.Tag.Get("json"), err)
						}
						field.SetInt(d.Unix())
						continue
					}
				}

				var v int64
				v, err = strconv.ParseInt(strField, 0, f.Type.Bits())
				if err == nil {
					if isPtr {
						field.SetPointer(unsafe.Pointer(&v))
					} else {
						field.SetInt(v)
					}
					if i == 0 {
						firstColValue = v
					}
				}
			} else if kind == reflect.Uint ||
				kind == reflect.Uint8 ||
				kind == reflect.Uint16 ||
				kind == reflect.Uint32 ||
				kind == reflect.Uint64 {
				var v uint64
				v, err = strconv.ParseUint(strField, 0, f.Type.Bits())
				if err == nil {
					if isPtr {
						field.SetPointer(unsafe.Pointer(&v))
					} else {
						field.SetUint(v)
					}
					if i == 0 {
						firstColValue = v
					}
				}
			} else if kind == reflect.Float32 ||
				kind == reflect.Float64 {
				var v float64
				v, err = strconv.ParseFloat(strField, f.Type.Bits())
				if err == nil {
					if isPtr {
						field.SetPointer(unsafe.Pointer(&v))
					} else {
						field.SetFloat(v)
					}
				}
			} else if kind == reflect.String {
				if isPtr {
					field.SetPointer(unsafe.Pointer(&strField))
				} else {
					field.SetString(strField)
				}
				if i == 0 {
					firstColValue = strField
				}
			} else if kind == reflect.Struct ||
				kind == reflect.Array ||
				kind == reflect.Slice ||
				kind == reflect.Map {
				if f.Type.String() == "*json.RawMessage" && !isBeginJson(strField) {
					strField = "\"" + strField + "\""
				}
				err = json.Unmarshal([]byte(strField), field.Addr().Interface())
			} else if kind == reflect.Interface {
				err = json.Unmarshal([]byte(strField), field.Addr().Interface())
			}

			if err != nil {
				return fmt.Errorf("parse field (Line=%v, ColumnName=\"%v\") error: %v",
					n+1+5, f.Tag.Get("json"), err)
			}
		}
	}

	rf.records = records
	return nil
}

func (rf *csvloader) Record() []interface{} {
	return rf.records
}

func (rf *csvloader) NumRecord() int {
	return len(rf.records)
}

type Hourer interface {
	SetHour(int)
}

type Minuter interface {
	SetMinute(int)
}

func (rf *csvloader) parseTime(field reflect.Value, csvData, layout string) error {
	format := strings.Split(layout, ":")
	data := strings.Split(csvData, ":")
	if len(format) != len(data) {
		return fmt.Errorf("parseTime csvData[%v] to layout[%v]", csvData, layout)
	}
	for i, v := range format {
		switch v {
		case "H", "h":
			b, ok := field.Addr().Interface().(Hourer)
			if !ok {
				return fmt.Errorf("parseTime csvData[%v] to layout[%v]", csvData, layout)
			}
			h, err := strconv.Atoi(data[i])
			if err != nil {
				return err
			}
			b.SetHour(h)
		case "M", "m":
			b, ok := field.Addr().Interface().(Minuter)
			if !ok {
				return fmt.Errorf("parseTime csvData[%v] to layout[%v]", csvData, layout)
			}
			m, err := strconv.Atoi(data[i])
			if err != nil {
				return err
			}
			b.SetMinute(m)
		}
	}
	return nil
}

// validUTF8 验证是否是utf8编码
func validUTF8(buf []byte) bool {
	temp := make([]byte, len(buf))
	copy(temp, buf)
	nBytes := 0
	for i := 0; i < len(temp); i++ {
		if nBytes == 0 {
			if (temp[i] & 0x80) != 0 { // 与操作之后不为0，说明首位为1
				for (temp[i] & 0x80) != 0 {
					temp[i] <<= 1 // 左移一位
					nBytes++      // 记录字符共占几个字节
				}

				if nBytes < 2 || nBytes > 6 { // 因为UTF8编码单字符最多不超过6个字节
					return false
				}

				nBytes-- // 减掉首字节的一个计数
			}
		} else { // 处理多字节字符
			if temp[i]&0xc0 != 0x80 { // 判断多字节后面的字节是否是10开头
				return false
			}
			nBytes--
		}
	}
	return nBytes == 0
}

func isBeginJson(s string) bool {
	c := s[0]
	switch c {
	case '{', '[', '"', '-':
		return true
	default:
		if c >= '0' && c <= '9' {
			return true
		} else {
			if s == "true" || s == "false" || s == "null" {
				return true
			}
		}
	}
	return false
}
