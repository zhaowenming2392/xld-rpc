package files

import (
	"bufio"
	"io"
	"os"
	"strings"
)

//ReadFile 读取文件
func ReadFile(filePath string) (string, error) {
	f, err := os.OpenFile(filePath, os.O_RDWR, 066)
	if err != nil {
		return "", err
	}

	defer f.Close()

	r := bufio.NewReader(f)
	strs := ""

	for {
		str, err := r.ReadString('\n') //读到一个换行就结束
		strs += str
		//读到末尾又没有字符串了
		if err == io.EOF { //io.EOF 表示文件的末尾
			break
		}
	}

	return strs, nil
}

//ReadFileToArray 读取文件到数组中，每行都是数组的一个元素
func ReadFileToArray(filePath string) ([]string, error) {
	f, err := os.OpenFile(filePath, os.O_RDWR, 066)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	r := bufio.NewReader(f)
	strs := []string{}

	for {
		str, err := r.ReadString('\n') //读到一个换行就结束
		str = strings.TrimSpace(str)   //去除首尾空格，特别是换行符
		strs = append(strs, str)
		//读到末尾又没有字符串了
		if err == io.EOF { //io.EOF 表示文件的末尾
			break
		}
	}

	return strs, nil
}

//WriteFile 写入文件
func WriteFile(filePath, str string, isAppend bool) (int, error) {
	perm := 0
	if isAppend {
		perm = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	} else {
		perm = os.O_WRONLY | os.O_CREATE
	}

	f, err := os.OpenFile(filePath, perm, 0775)

	if err != nil {
		return 0, err
	}

	defer f.Close()

	n, err := f.WriteString(str)
	if err != nil {
		return 0, err
	}

	return n, nil
}

//IsExist 判断文件或者目录是否存在
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

//IsFile 判断是否是文件，不存在也是假
func IsFile(f string) bool {
	return !IsDir(f)
}

//IsDir 判断是否是目录，不存在也是假
func IsDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

//DirName 按照惯例返回文件的目录名称
//
//不检查是否真实和存在
//
//只是去除最后的文件名来得到目录名称
func DirName(f string) string {
	if strings.ContainsRune(f, '/') && !strings.HasSuffix(f, "/") {
		//存在/且不是最后一个字符（是最后一个字符则本身就是目录）
		fs := strings.Split(f, "/")
		return strings.Join(fs[:len(fs)-1], "/")
	} else if strings.ContainsRune(f, '\\') && !strings.HasSuffix(f, "\\") {
		//存在\且不是最后一个字符（是最后一个字符则本身就是目录）
		fs := strings.Split(f, "\\")
		return strings.Join(fs[:len(fs)-1], "/")
	}

	return f
}
