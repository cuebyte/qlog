package qlog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type MutexWriter struct {
	sync.Mutex
	fd *os.File
}

func (m *MutexWriter) Write(b []byte) (int, error) {
	m.Lock()
	defer m.Unlock()
	return m.fd.Write(b)
}

// implament
type FileLog struct {
	*log.Logger
	writer *MutexWriter

	Filename string
	Suffix   string

	Daily  bool
	Maxage int64

	Rotate  bool
	Maxsize int64
}

func NewFileLog() *FileLog {
	fl := &FileLog{
		Filename: "C:\\go\\workspace\\src\\qlog\\badnews",
		Suffix:   ".log",
		Daily:    false,
		Maxage:   0,
		Rotate:   false,
		Maxsize:  0,
	}
	fl.writer = new(MutexWriter)
	fl.Logger = log.New(fl.writer, "", 0)
	return fl
}
func (this *FileLog) Init() error {
	this.writer = new(MutexWriter)
	var err error
	f, _ := this.getRealFilename()
	this.writer.fd, err = os.OpenFile(f, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		return err
	}
	this.Logger = log.New(this.writer, "", 0)
	return nil
}
func (this *FileLog) Write(msg *Message) {
	if !this.Daily && !this.Rotate {
		this.insert(msg.ToString())
		return
	}
	var err error
	n, isNew := this.getRealFilename()
	if isNew {
		this.Close()
		this.writer.fd, err = os.OpenFile(n, os.O_WRONLY|os.O_APPEND, 0)
	}
	if err != nil {
		log.Fatal(err)
	}
	this.insert(msg.ToString())
}
func (this *FileLog) Flush() {
	this.writer.fd.Sync()
}
func (this *FileLog) Close() {
	this.writer.fd.Close()
}

func isFileExist(Filename string) bool {
	_, err := os.Lstat(Filename)
	return err == nil
}

// Filename must exist
func isOverSize(Filename string, maxSize int64) bool {
	if maxSize != 0 {
		fi, err := os.Stat(Filename)
		if err != nil {
			panic("1." + err.Error())
		}
		return fi.Size() >= maxSize
	}
	return false
}

func (this *FileLog) getLastFile() (string, string) {
	fn := this.Filename
	lastFile := fn
	newFile := fn
	if this.Daily {
		fn += "_" + time.Now().Format("20060102")
		lastFile = fn
		newFile = lastFile
	}
	if this.Rotate {
		lastFile = fn + fmt.Sprintf("_%03d", 0)
		newFile = lastFile
		for i := 1; isFileExist(newFile + this.Suffix); i++ {
			lastFile = newFile
			newFile = fn + fmt.Sprintf("_%03d", i)
		}
	}
	return lastFile + this.Suffix, newFile + this.Suffix
}

// create new log file
func (this *FileLog) getRealFilename() (string, bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	var lastFile, newFile string
	lastFile, newFile = this.getLastFile()
	if isFileExist(lastFile) && !isOverSize(lastFile, this.Maxsize) {
		return lastFile, false
	} else {
		path, _ := filepath.Split(newFile)
		if err := os.MkdirAll(path, 0770); err != nil {
			panic(path + "---" + err.Error())
		}
		f, err := os.OpenFile(newFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			panic(err.Error())
		}
		return f.Name(), true
	}
}

// delete old log file
func (this *FileLog) delete() { // 应该用go-crontab独立出一个进程进行expire操作

}

// add one record
func (this *FileLog) insert(input string) {
	this.Logger.Println(input)
}
