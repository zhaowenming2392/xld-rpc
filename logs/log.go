package logs

import (
	"io"
	"log"
	"os"
)

//*******************************日志者******************
type Loger struct {
	log *log.Logger
}

//创建自定义日志记录
func NewLoger(out io.Writer) *Loger {
	return &Loger{
		log: log.New(out, "[info]", log.Ldate|log.Lmicroseconds),
	}
}

//信息
func (l *Loger) Info(v ...interface{}) {
	l.log.SetPrefix("[info]")
	l.log.SetFlags(log.Ldate | log.Lmicroseconds)
	l.log.Println(v...)
}

//信息格式化
func (l *Loger) Infof(format string, v ...interface{}) {
	l.log.SetPrefix("[info]")
	l.log.SetFlags(log.Ldate | log.Lmicroseconds)
	l.log.Printf(format+"\n", v...)
}

//警告
func (l *Loger) Warning(v ...interface{}) {
	l.log.SetPrefix("[warn]")
	l.log.SetFlags(log.Ldate | log.Ltime)
	l.log.Println(v...)
}

//恐慌
func (l *Loger) Panic(v ...interface{}) {
	l.log.SetPrefix("[panic]")
	l.log.SetFlags(log.Ldate | log.Ltime)
	l.log.Panicln(v...)
}

//恐慌格式化
func (l *Loger) Panicf(format string, v ...interface{}) {
	l.log.SetPrefix("[panic]")
	l.log.SetFlags(log.Ldate | log.Ltime)
	l.log.Panicf(format+"\n", v...)
}

//致命错误
func (l *Loger) Fatal(v ...interface{}) {
	l.log.SetPrefix("[fatal]")
	l.log.SetFlags(log.Ldate | log.Ltime)
	l.log.Fatalln(v...)
}

//致命错误格式化
func (l *Loger) Fatalf(format string, v ...interface{}) {
	l.log.SetPrefix("[fatal]")
	l.log.SetFlags(log.Ldate | log.Ltime)
	l.log.Fatalf(format+"\n", v...)
}

//**********************默认日志*******************************
var defLoger *Loger = NewLoger(os.Stdout)

//默认信息
func Info(v ...interface{}) {
	defLoger.log.SetPrefix("[info]")
	defLoger.log.SetFlags(log.Ldate | log.Lmicroseconds)
	defLoger.log.Println(v...)
}

//默认信息格式化
func Infof(format string, v ...interface{}) {
	defLoger.log.SetPrefix("[info]")
	defLoger.log.SetFlags(log.Ldate | log.Lmicroseconds)
	defLoger.log.Printf(format+"\n", v...)
}

//默认警告
func Warning(v ...interface{}) {
	defLoger.log.SetPrefix("[warn]")
	defLoger.log.SetFlags(log.Ldate | log.Ltime)
	defLoger.log.Println(v...)
}

//默认恐慌
func Panic(v ...interface{}) {
	defLoger.log.SetPrefix("[panic]")
	defLoger.log.SetFlags(log.Ldate | log.Ltime)
	defLoger.log.Panicln(v...)
}

//默认恐慌格式化
func Panicf(format string, v ...interface{}) {
	defLoger.log.SetPrefix("[panic]")
	defLoger.log.SetFlags(log.Ldate | log.Ltime)
	defLoger.log.Panicf(format+"\n", v...)
}

//默认致命错误
func Fatal(v ...interface{}) {
	defLoger.log.SetPrefix("[fatal]")
	defLoger.log.SetFlags(log.Ldate | log.Ltime)
	defLoger.log.Fatalln(v...)
}

//默认致命错误格式化
func Fatalf(format string, v ...interface{}) {
	defLoger.log.SetPrefix("[fatal]")
	defLoger.log.SetFlags(log.Ldate | log.Ltime)
	defLoger.log.Fatalf(format+"\n", v...)
}
