package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/forgoer/openssl"
	"golang.org/x/crypto/scrypt"
)

func AesCBCEncrypt(src, key, iv []byte, padding string) ([]byte, error) {
	data, err := openssl.AesCBCEncrypt(src, key, iv, padding)
	if err != nil {
		return nil, err
	}
	return []byte(hex.EncodeToString(data)), nil
}

func AesCBCDecrypt(src, key, iv []byte, padding string) ([]byte, error) {
	data, err := hex.DecodeString(string(src))
	if err != nil {
		return nil, err
	}
	return openssl.AesCBCDecrypt(data, key, iv, padding)

}
func Md5(text string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, text)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}

func Zip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&b, 9)
	if _, err := gz.Write([]byte(data)); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func UnZip(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	unzipData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return unzipData, nil
}

// ScryptPasswd 密码加密
func ScryptPasswd(password string) string {
	const keyLen = 10
	salt := make([]byte, 8)
	salt = []byte{12, 32, 4, 6, 22, 66, 222, 111}
	hashPasswd, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, keyLen)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(hashPasswd)
}
