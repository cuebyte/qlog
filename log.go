package qlog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"qlog/m2s"
	"strings"
	"sync"
	"time"

	"github.com/Unknwon/goconfig"
)

const (
	EMER = iota
	ALER
	CRIT
	ERRO
	WARN
	NOTI
	INFO
	DBUG
	TRAC
)

var lv = []string{"EMER", "ALER", "CRIT", "ERRO", "WARN", "NOTI", "INFO", "DBUG", "TRAC"}
var conf_type = []string{"file", "mysql"}
var conf_type_split string = "-"

//////// Message Start ////////

type Message struct {
	Time    string
	Level   int
	Content interface{}
}

func NewMessage(method int, input interface{}) *Message {
	pl := &Message{Level: method, Content: input}
	pl.Time = time.Now().Format("2006/01/02 15:04:05")
	return pl
}

func (this *Message) GetContentString() string {
	switch this.Content.(type) {
	case string:
		return this.Content.(string)
	case error:
		return this.Content.(error).Error()
	default:
		return ""
		// TODO: User Action(AOP)
	}
}
func (this *Message) ToString() string {
	return fmt.Sprintf("%s [%s] %s", this.Time, lv[this.Level], this.GetContentString())
}

//////// Message End ////////

type LogWriter interface {
	Init() error
	Write(*Message)
	Flush()
	Close()
}

//////// Logger start ////////

type Logger struct {
	Level   int
	Writers []LogWriter
}

func InitLogger(level int, paths ...string) (Logger, error) {
	log := &Logger{Level: level}
	log.Writers = make([]LogWriter, 0)
	// read config
	// level, output, and so on, should be configured in the conf/ini file
	if err := log.Configure(paths); err != nil {
		return Logger{}, err
	}
	// 并发控制
	return *log, nil
}
func (this *Logger) Close() {
	for _, x := range this.Writers {
		x.Flush()
		x.Close()
	}
}

// input: path of .ini files or .ini file route
func (this *Logger) Configure(paths []string) error {
	if len(paths) == 0 {
		paths = make([]string, 1)
		paths[0] = "."
	}
	for _, s := range paths {
		if !filter(s) {
			return errors.New(fmt.Sprintf("\"%s\" is not available", s))
		}
		if isIni(s) {
			this.readIni(s)
		} else {
			// provide recursion
			filepath.Walk(s, func(ipath string, iinfo os.FileInfo, ierr error) error {
				if isIni(ipath) {
					if err := this.readIni(ipath); err != nil {
						return err
					}
				}
				return nil
			})
		}
	}
	return nil
}

func (this *Logger) readIni(filename string) error {
	cfg, err := goconfig.LoadConfigFile(filename)
	if err != nil {
		return err
	}
	nodes := cfg.GetSectionList()
	for _, node := range nodes {
		switch strings.Split(strings.ToLower(node), conf_type_split)[0] {
		case "file":
			m, _ := cfg.GetSection(node)
			fl := NewFileLog()
			m2s.Map2Struct(m, fl)
			if err := fl.Init(); err != nil {
				return err
			}
			//this.Writers[node] = fl
			this.Writers = append(this.Writers, fl)
		case "mysql":
		default:
			continue
		}
	}
	return nil
}

func (this *Logger) write(info *Message) {
	if this.Level < info.Level {
		return
	}
	//var wg sync.WaitGroup
	lenght := len(this.Writers)
	wgs := make([]sync.WaitGroup, 0)
	for i := 0; i < lenght; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		wgs = append(wgs, wg)
	}
	//cs := make(chan *Message, len(this.Writers))
	for i := 0; i < lenght; i++ {
		// cs <- info
		go func(info *Message, i int) {
			//fmt.Println(2)
			this.Writers[i].Write(info)
			wgs[i].Done()
		}(info, i)
	}
	for i := 0; i < lenght; i++ {
		//fmt.Println(3)
		wgs[i].Wait()
	}
}

func (this *Logger) Emerge(input interface{}) {
	this.write(NewMessage(EMER, input))

}
func (this *Logger) Alert(input interface{}) {
	this.write(NewMessage(ALER, input))

}
func (this *Logger) Critical(input interface{}) {
	this.write(NewMessage(CRIT, input))

}
func (this *Logger) Error(input interface{}) {
	this.write(NewMessage(ERRO, input))

}
func (this *Logger) Warn(input interface{}) {
	this.write(NewMessage(WARN, input))
}
func (this *Logger) Notice(input interface{}) {
	this.write(NewMessage(NOTI, input))

}
func (this *Logger) Info(input interface{}) {
	this.write(NewMessage(INFO, input))

}
func (this *Logger) Debug(input interface{}) {
	this.write(NewMessage(DBUG, input))

}
func (this *Logger) Trace(input interface{}) {
	this.write(NewMessage(TRAC, input))

}

//////// Logger end ////////
