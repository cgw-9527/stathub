package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	sys   = ""
	file  = "/var/log"
	level = "debug"
)

func init() {
	if runtime.GOOS == "windows" {
		file = "d:/log"
	}

	if Exists(file) == false { //不存在路径则创建
		os.Mkdir(file, 0777)
	}

	//file = file + "/node_main.log"
	file = file + "/stathub.log"

	InitLogger()
	// file, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE, 666)
	// if err != nil {
	// 	log.Fatalln("fail to create test.log file!")
	// }
	// defer file.Close()
	// logger = log.New(file, "", log.LstdFlags|log.Lshortfile) // 日志文件格式:log包含时间及文件行数
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func rotateCreaten() (*rotatelogs.RotateLogs, error) {
	baseLogPath := file
	if runtime.GOOS == "windows" {
		return rotatelogs.New(
			baseLogPath+".%Y%m%d%H%M",
			rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
			rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
		)
	}

	return rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
}

//初始化日志
func InitLogger() {
	writer, err := rotateCreaten()
	if err != nil {
		log.Errorf("config local file system logger error. %v", errors.WithStack(err))
	}

	//log.SetFormatter(&log.TextFormatter{})
	switch level := level; level {
	/*
		如果日志级别不是debug就不要打印日志到控制台了
	*/
	case "debug":
		log.SetLevel(log.DebugLevel)
		log.SetOutput(os.Stderr)
	case "info":
		setNull()
		log.SetLevel(log.InfoLevel)
	case "warn":
		setNull()
		log.SetLevel(log.WarnLevel)
	case "error":
		setNull()
		log.SetLevel(log.ErrorLevel)
	default:
		setNull()
		log.SetLevel(log.InfoLevel)
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, nil)
	log.AddHook(lfHook)

}

func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	log.SetOutput(writer)
}

func Nlog1(tag string, msg interface{}) {

	log.WithFields(logrus.Fields{
		"tag": tag,
	}).Info(msg)
}

//默认日志调用
// 	xxx := [...]int{1, 2, 3, 4}
func Nlog(msg ...interface{}) {

	if len(msg) == 1 {
		Nlog1("default", msg)
	} else {

		fmtstr := msg[0].(string)
		// args := msg[1:]
		// s := fmt.Sprintf(fmtstr, args)
		s := ""
		if len(msg) == 2 {
			s = fmt.Sprintf(fmtstr, msg[1])
		} else if len(msg) == 3 {
			s = fmt.Sprintf(fmtstr, msg[1], msg[2])
		} else if len(msg) == 4 {
			s = fmt.Sprintf(fmtstr, msg[1], msg[2], msg[3])
		} else if len(msg) == 5 {
			s = fmt.Sprintf(fmtstr, msg[1], msg[2], msg[3], msg[4])
		} else if len(msg) == 6 {
			s = fmt.Sprintf(fmtstr, msg[1], msg[2], msg[3], msg[4], msg[5])
		} else if len(msg) == 7 {
			s = fmt.Sprintf(fmtstr, msg[1], msg[2], msg[3], msg[4], msg[5], msg[6])
		} else { //默认
			s = fmt.Sprintf(fmtstr, msg[1:])
		}

		Nlog1("default", s)
	}
	//	InitLogger()

}

/*

  log.WithFields(logrus.Fields{
    "animal": "walrus",
    "size":   10,
  }).Info("A group of walrus emerges from the ocean")

*/
