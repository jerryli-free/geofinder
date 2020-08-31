package lib

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UseLog struct {
	path         string
	buffer       *bytes.Buffer
	bufferTicker *time.Ticker
	bufferCap    int
	bufferChan   chan []byte
	logFile      io.Writer
}

const (
	ACCESS_PATH  = "/home/log/gpsServer/http_access.log"
	INFO_PATH    = "/home/log/gpsServer/info.log"
	FATAL_PATH   = "/home/log/gpsServer/fatal.log"
	PREFIX_INFO  = "INFO: "
	PREFIX_FATAL = "FATAL: "
	BUFFER_CAP = 1024 * 32
)

var (
	infLog   UseLog
	fatalLog UseLog
)

func init() {
	infLog.setPath(INFO_PATH)
	fatalLog.setPath(FATAL_PATH)
}

func (cl *UseLog) setPath(path string) *UseLog {
	path = strings.Trim(path, "")
	if path == "" {
		log.Fatalln("Path is empty")
	} else {
		cl.path = path
	}
	return cl
}

func (cl *UseLog) getPath() string {
	return cl.path
}

func (cl *UseLog) Info(args interface{}) {
	if cl.getPath() == "" {
		cl.setPath(INFO_PATH)
	}
	file := cl.fileInit()
	if file == nil {
		return
	}
	cl.write(file, PREFIX_INFO, args)
}
func (cl *UseLog) Fatal(args interface{}) {
	if cl.getPath() == "" {
		cl.setPath(FATAL_PATH)
	}
	file := cl.fileInit()
	if file == nil {
		return
	}
	cl.write(file, PREFIX_FATAL, args)
}
func (cl *UseLog) fileInit() *os.File {
	file, err := os.OpenFile(cl.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
		return nil
	}
	return file
}
func (cl *UseLog) write(file *os.File, prefix string, args interface{}) {
	infoVar := log.New(io.MultiWriter(file), prefix, log.Ldate|log.Ltime)
	infoVar.Println(args)
}

func Info(args ...interface{}) {
	infLog.setPath(INFO_PATH).Info(args)
}

func Fatal(args ...interface{}) {
	fatalLog.setPath(FATAL_PATH).Fatal(args)
}
func SetFatalPath(path string) *UseLog {
	fatalLog.setPath(path)
	return &fatalLog
}
func SetInfoPath(path string) *UseLog {
	infLog.setPath(path)
	return &infLog
}

//异步写盘时 - 定时器或超长时从内存写入磁盘
func (cl *UseLog) asyncWrite() bool {
	for {
		select {
		case <-cl.bufferTicker.C:
			if cl.buffer.Len() > 0 {
				cl.flush()
			}
		case record := <-cl.bufferChan:
			cl.buffer.Write(record)
			if cl.buffer.Len() >= cl.bufferCap {
				cl.flush()
			}
		}
	}
}

//异步写盘时 - 强制刷磁盘
func (cl *UseLog) flush() {
	//刷内容从buffer到文件
	_, _ = cl.logFile.Write(cl.buffer.Bytes())
	cl.buffer.Reset()
}

//备份并删除同类型的过期的Log
func (cl *UseLog) asyncModfiyLog() {
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//以下为定时0点执行的操作
		cl.backupLog()
		cl.removeExpireLog()
	}

}
func (cl *UseLog) backupLog() {
	yesTime := time.Now().Add(-time.Hour * 24)
	newName := filepath.Dir(cl.path) + "/" + cl.getMidName() + "_" + yesTime.Format("20060102") + filepath.Ext(cl.path)
	err:=os.Rename(cl.path, newName)
	if err == nil {
		fmt.Println("访问日志文件备份成功：" + newName)
	}
	logFile, err := os.OpenFile(cl.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("open log file " + cl.path + " error: " + err.Error())
	}
	cl.logFile = logFile
}
func (cl *UseLog) removeExpireLog() {
	//删除7天前的日志
	_ = filepath.Walk(filepath.Dir(cl.path), func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		fileNamePrefix := cl.getMidName()
		if !strings.Contains(path, fileNamePrefix) {
			return nil
		}
		if time.Now().Sub(info.ModTime()) > 86400*7*time.Second {
			err := os.Remove(path)
			if err == nil {
				fmt.Println("访问日志历史文件删除成功：" + path)
			}
		}
		return nil
	})
}

func (cl *UseLog) getMidName() string {
	return strings.TrimSuffix(filepath.Base(cl.path), filepath.Ext(cl.path))
}

//异步写盘时 - 写入异步队列
func (cl *UseLog) Write(p []byte) (n int, err error) {
	cl.bufferChan <- p
	return len(p), nil
}

func NewUseLog(path string) gin.HandlerFunc {
	if path == "" {
		path = ACCESS_PATH
	}
	//打开文件
	os.MkdirAll(filepath.Dir(path), 0666)
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("open log file " + path + " error: " + err.Error())
	}
	bufferCap := BUFFER_CAP
	l := &UseLog{
		path:         path,
		logFile:      logFile,
		buffer:       bytes.NewBuffer(make([]byte, 0, bufferCap)),
		bufferTicker: time.NewTicker(time.Second * 1),
		bufferCap:    bufferCap,
		bufferChan:   make(chan []byte, bufferCap),
	}
	go l.asyncWrite()
	go l.asyncModfiyLog()
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Output: l,
	})
}
