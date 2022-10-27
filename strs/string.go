package strs

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

//StrToInt 字符串转成数字，注意这里排除了错误，用于明确知道是数字情况下
func StrToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func HexToStr(bytes []byte) string {
	return hex.EncodeToString(bytes) //将数组转换成切片，转换成16进制，返回字符串
}

func IsPhone(p string) bool {
	if p == "" || len(p) != 11 || string(p[0]) != "1" {
		return false
	}

	//通过转换为10进制的int64是否有误来判断是否为int64数字
	_, err := strconv.ParseInt(p, 10, 64)
	return err == nil
}

//是不是需要分割的字符串
func isSplit(r rune) bool {
	return r == ' ' || r == '-' || r == '_'
}

//格式化不同大小写格式的英文名称
//
//format=0 原样输出，1=首字母小写(后面字母不处理)，2=全小写，3=全大写，4=首字母大写(后面字母不处理)，5首字母小写的驼峰写法，6首字母大写的驼峰写法
func FormatName(name string, format int) string {
	switch format {
	case 1:
		return strings.ToLower(name[:1]) + name[1:]
	case 2:
		return strings.ToLower(name)
	case 3:
		return strings.ToUpper(name)
	case 4:
		return strings.ToUpper(name[:1]) + name[1:]
	case 5:
		//如果由空格、-、_字符，连接而成的，如a—b，a_b，a b，应该组织成aB
		names := strings.FieldsFunc(name, isSplit) //如果存在满足函数的字符串则切割

		//被切割了
		newName := ""
		for i := 0; i < len(names); i++ {
			fmt.Println("正在处理切割后的", names[i])
			if i == 0 {
				//首单词首字母小写
				newName = strings.ToLower(names[0][:1]) + names[0][1:]
				continue
			}

			//后面的单词，首字母都大写
			newName += strings.ToUpper(names[i][:1]) + names[i][1:]
		}
		return newName

	case 6:
		//如果由空格、-、_字符，连接而成的，如a—b，a_b，a b，应该组织成aB
		names := strings.FieldsFunc(name, isSplit) //如果存在满足函数的字符串则切割

		//被切割了
		newName := ""
		for i := 0; i < len(names); i++ {
			//单词，首字母都大写
			newName += strings.ToUpper(names[i][:1]) + names[i][1:]
		}

		return newName
	default:
		return name
	}
}
