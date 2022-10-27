package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"helpers.zhaowenming.cn/errs"
	"helpers.zhaowenming.cn/files"
)

/*
FastLog 快速记录日志到文件

name 文件名称会自动加上.log

tag 消息标记[tag]

msg interface{} 消息内容

日志文件：./runtime/log/年月/name-日.log

内容为追加的：2006-01-02 11:06:39.1234 [tag] msg的json编码\n
*/
type Fastlog struct {
	//路径
	Path string

	//文件名称
	Name string
	//文件名称的前后，防止重复名称而使用的
	NameSuffix string
	NamePrefix string

	//消息内容的前后
	MsgSuffix string
	MsgPrefix string

	//是否同时打印到控制台
	Display bool
}

var defaultFastLog = Fastlog{
	Name:       "test",
	NameSuffix: time.Now().Format("02"),
	MsgSuffix:  "\n",
}
var _ = SetTempLogDir("")

//创建新的快速日志助手
func NewFastlog(path, name string, display bool) (*Fastlog, error) {
	if path == "" || name == "" {
		return nil, errs.EmptyErr("path 和 name")
	}

	if !files.IsExist(path) {
		return nil, errs.NotExistErr("path")
	}

	return &Fastlog{
		Path:    strings.TrimRight(path, "/") + "/",
		Name:    name,
		Display: display,
	}, nil
}

//设置默认助手的日志路径
func SetTempLogDir(path string) error {
	if path == "" {
		m := time.Now().Format("200601")
		path = "./runtime/log/" + m + "/"
	}

	err := os.MkdirAll(path, 0644)
	if err != nil {
		return err
	}

	defaultFastLog.Path = strings.TrimRight(path, "/") + "/"
	return nil
}

//设置默认助手的日志文件名称
func SetLogName(name string) {
	defaultFastLog.Name = name
}

//设置默认助手的日志名称的前后缀
func SetNameMark(prefix, suffix string) {
	defaultFastLog.NamePrefix = prefix
	defaultFastLog.NameSuffix = suffix
}

//设置默认助手的日志消息内容的前后缀
func SetMsgMark(prefix, suffix string) {
	defaultFastLog.MsgPrefix = prefix
	defaultFastLog.MsgSuffix = suffix
}

//设置默认助手的控制台打印
func SetDisplay(dis bool) {
	defaultFastLog.Display = dis
}

func (f Fastlog) Log(tag string, msg interface{}) error {
	//非字符串进行json编码
	msgStr, ok := msg.(string)
	if !ok {
		body, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		msgStr = string(body)
	}

	//组织tag
	if tag != "" {
		tag = "[" + tag + "] "
	}

	//写入文件
	name := f.Path + f.NamePrefix + f.Name + f.NameSuffix

	t := time.Now().Format("2006-01-02 15:04:05.000 ")
	str := t + tag + f.MsgPrefix + msgStr + f.MsgSuffix

	//TODO 如果想快，减少写入次数，可以先存到缓冲，再写入
	_, err := files.WriteFile(name+".log", str, true)
	if err != nil {
		return err
	}

	if f.Display {
		fmt.Println(str)
	}

	return nil
}

//默认的快速日志记录，所有选项才有默认值
func FastLog(tag string, msg interface{}) error {
	return defaultFastLog.Log(tag, msg)
}

//提供自定义名称的日志记录，名称临时生效
func FastNameLog(name, tag string, msg interface{}) error {
	old := defaultFastLog.Name
	SetLogName(name)
	if err := defaultFastLog.Log(tag, msg); err != nil {
		return err
	}
	SetLogName(old)
	return nil
}

//提供自定义路径和名称的日志记录，名称、路径临时生效
func FastNamePathLog(path, name, tag string, msg interface{}) error {
	old := defaultFastLog.Path
	if defaultFastLog.Path != path {
		if err := SetTempLogDir(path); err != nil {
			return err
		}
	}
	if err := FastNameLog(name, tag, msg); err != nil {
		return err
	}
	if err := SetTempLogDir(old); err != nil {
		return err
	}
	return nil
}
