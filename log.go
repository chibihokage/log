package log

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type LogLevel struct {
	Debuglog, Infolog, Errorlog Log
}

type logPattern struct {
	Date, Level, Refnum, Subrnumb, Service, Proc, Msg string
}
type logTrnsPattern struct {
	Date, Service, Subrnumb, Status, Errcode string
}
type Log struct {
	Debug       *log.Logger
	patternInit string
	patternRec  string
	patternTrns string
}
type mytime struct {
	time.Time
	format string
}

func (t mytime) String() string {
	return t.Format(t.format)
}

func newTimeFmt(t time.Time, f string) mytime {
	return mytime{t, f}
}

func NewDebugLog(w io.Writer, level string) Log {
	format := `{{.Date}}|{{.Level}}|{{.Refnum}}|{{.Subrnumb}}|{{.Service}}|{{.Proc}}|{{.Msg}}`

	logDebug := logPattern{
		Date:     "$date",
		Level:    level,
		Refnum:   "$ref",
		Subrnumb: "$sub",
		Service:  "$ser",
		Proc:     "$proc",
		Msg:      "%s",
	}

	bufDebug := bytes.NewBuffer([]byte{})
	t := template.Must(template.New(format).Parse(format))
	err := t.Execute(bufDebug, logDebug)

	if err != nil {
		log.Println("executing template:", err)
	}

	return Log{
		Debug:       log.New(w, "", 0),
		patternInit: bufDebug.String(),
	}
}
func NewTrnsLog(w io.Writer) Log {
	format := `{{.Date}}|{{.Service}}|{{.Subrnumb}}|{{.Status}}|{{.Errcode}}`
	return Log{
		Debug:       log.New(w, "", 0),
		patternTrns: format,
	}
}

func (log *Log) SetRecordDetail(refnumb, subrnumb string) {
	r := strings.NewReplacer(
		"$ref", refnumb,
		"$sub", subrnumb,
	)
	log.patternRec = r.Replace(log.patternInit)
}
func (log *Log) SetInitlogDetail(service, proc string) {
	r := strings.NewReplacer(
		"$ser", service,
		"$proc", proc,
	)
	log.patternInit = r.Replace(log.patternInit)
}
func (log Log) PrintTrns(service, subrnumb, status, errcode string) {
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	dateNow := newTimeFmt(time.Now().In(bangkok), "2006-01-02T15:04:05")
	logTrns := logTrnsPattern{
		Date:     dateNow.String(),
		Service:  service,
		Subrnumb: subrnumb,
		Status:   status,
		Errcode:  errcode,
	}

	bufTrns := bytes.NewBuffer([]byte{})
	t := template.Must(template.New(log.patternTrns).Parse(log.patternTrns))
	err = t.Execute(bufTrns, logTrns)
	if err != nil {
		log.Printf("executing template: %v", err)
	}
	log.Debug.Println(bufTrns.String())
}

func (log Log) Println(msg string) {
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatalf("Failed to LoadLocation %v", err)
	}

	dateNow := newTimeFmt(time.Now().In(bangkok), "2006-01-02T15:04:05")
	var patternDebug string
	if log.patternRec == "" {
		r := strings.NewReplacer(
			"$ref", "",
			"$sub", "",
			"$date", dateNow.String(),
		)
		patternDebug = r.Replace(log.patternInit)
	} else {
		patternDebug = strings.Replace(log.patternRec, "$date", dateNow.String(), -1)
	}

	log.Debug.Println(fmt.Sprintf(patternDebug, msg))
}

func (log Log) Printf(msg string, a ...interface{}) {
	log.Println(fmt.Sprintf(msg, a...))
}

func (log Log) Fatalln(msg string) {
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatalf("Failed to LoadLocation %v", err)
	}

	dateNow := newTimeFmt(time.Now().In(bangkok), "2006-01-02T15:04:05")
	var patternDebug string
	if log.patternRec == "" {
		r := strings.NewReplacer(
			"$ref", "",
			"$sub", "",
			"$date", dateNow.String(),
		)
		patternDebug = r.Replace(log.patternInit)
	} else {
		patternDebug = strings.Replace(log.patternRec, "$date", dateNow.String(), -1)
	}

	log.Debug.Fatalln(fmt.Sprintf(patternDebug, msg))
}

func (log Log) Fatalf(msg string, a ...interface{}) {
	log.Fatalln(fmt.Sprintf(msg, a...))
}

func CreateLogFile(filename string) io.Writer {
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	info, err := os.Stat(filename)
	if err == nil {
		modifiedtime := info.ModTime()
		dateNow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, bangkok)
		modifiedDate := time.Date(modifiedtime.Year(), modifiedtime.Month(), modifiedtime.Day(), 0, 0, 0, 0, bangkok)
		modifiedDateFm := newTimeFmt(modifiedDate, "2006-01-02")
		if !dateNow.Equal(modifiedDate) {
			backupName := fmt.Sprintf("%s.%s", filename, modifiedDateFm)
			os.Rename(filename, backupName)
		}

	}
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln("Failed to open log file", logFile, ":", err)
	}

	return logFile
}
func InitDebuglog(filename, service, proc string) LogLevel {
	var logs LogLevel
	fileDebug := CreateLogFile(filename)
	logs.Debuglog = NewDebugLog(fileDebug, "DEBUG")
	logs.Infolog = NewDebugLog(fileDebug, "INFO ")
	logs.Errorlog = NewDebugLog(fileDebug, "ERROR")
	logs.Debuglog.SetInitlogDetail(service, proc)
	logs.Infolog.SetInitlogDetail(service, proc)
	logs.Errorlog.SetInitlogDetail(service, proc)
	return logs
}
func InitTrnslog(filename string) Log {
	fileTrns := CreateLogFile(filename)
	trnslog := NewTrnsLog(fileTrns)
	return trnslog
}
