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

type logDebugPattern struct {
	Date, Level, SourceSystemID, SessionID, TrnsID, Subrnumb, Msg string
}

type logTrnsEndpointPattern struct {
	Date, HostName, SourceSystemID, SessionID, TrnsID, Subrnumb, RequestIP, ServiceName, FuncName, ServiceType, EndpointServiceName, EndpointStatusType, EndpointStatusCode, EndpointErrCode, ResponseTime string
}

type logTrnsPattern struct {
	Date, HostName, SourceSystemID, SessionID, TrnsID, Subrnumb, RequestIP, ServiceName, FuncName, StatusType, ErrCode, ErrMsg, EndpointErrCode, ResponseTime string
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

func NewDebugLog(w io.Writer, level, srcSysName string) Log {
	format := `{{.Date}}|{{.Level}}|{{.SourceSystemID}}|{{.SessionID}}|{{.TrnsID}}|{{.Subrnumb}}|{{.Msg}}`

	logDebug := logDebugPattern{
		Date:           "$date",
		Level:          level,
		SourceSystemID: srcSysName,
		SessionID:      "$sessionID",
		TrnsID:         "$trnsID",
		Subrnumb:       "$sub",
		Msg:            "%s",
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
func NewTrnsLog(w io.Writer, sourceSystemID, sessionID, trnsID, subrnumb, requestIP, serviceName, funcName string) Log {
	format := `{{.Date}}|{{.HostName}}|{{.SourceSystemID}}|{{.SessionID}}|{{.TrnsID}}|{{.Subrnumb}}|{{.RequestIP}}|{{.ServiceName}}|{{.FuncName}}|{{.StatusType}}|{{.ErrCode}}|{{.ErrMsg}}|{{.EndpointErrCode}}|{{.ResponseTime}}`
	hostname, _ := os.Hostname()
	logTrns := logTrnsPattern{
		Date:            "$date",
		HostName:        hostname,
		SourceSystemID:  sourceSystemID,
		SessionID:       sessionID,
		TrnsID:          trnsID,
		Subrnumb:        subrnumb,
		RequestIP:       requestIP,
		ServiceName:     serviceName,
		FuncName:        funcName,
		StatusType:      "$statusType",
		ErrCode:         "$errCode",
		ErrMsg:          "$errMsg",
		EndpointErrCode: "$endpointErrCode",
		ResponseTime:    "$responseTime",
	}

	bufTrns := bytes.NewBuffer([]byte{})
	t := template.Must(template.New(format).Parse(format))
	err := t.Execute(bufTrns, logTrns)

	if err != nil {
		log.Println("executing template:", err)
	}

	return Log{
		Debug:       log.New(w, "", 0),
		patternTrns: bufTrns.String(),
	}
}

func NewEndpointTrnsLog(w io.Writer) Log {
	format := `{{.Date}}|{{.HostName}}|{{.SourceSystemID}}|{{.SessionID}}|{{.TrnsID}}|{{.Subrnumb}}|{{.RequestIP}}|{{.ServiceName}}|{{.FuncName}}|{{.ServiceType}}|{{.EndpointServiceName}}|{{.EndpointStatusType}}|{{.EndpointStatusCode}}|{{.EndpointErrCode}}|{{.ResponseTime}}`
	return Log{
		Debug:       log.New(w, "", 0),
		patternTrns: format,
	}
}

func (log *Log) SetRecordDetail(subrnumb string) {
	r := strings.NewReplacer(
		"$sub", subrnumb,
	)
	log.patternRec = r.Replace(log.patternInit)
}
func (log *Log) SetInitlogDetail(sessionID, trnsID string) {
	r := strings.NewReplacer(
		"$sessionID", sessionID,
		"$trnsID", trnsID,
	)
	log.patternInit = r.Replace(log.patternInit)
}

func (log Log) PrintTrns(statusType, errCode, errMsg, endpointErrCode, responseTime string) {
	var patternTrns string
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	dateNow := newTimeFmt(time.Now().In(bangkok), "2006-01-02T15:04:05")
	if err != nil {
		log.Fatalf("Failed to LoadLocation %v", err)
	}

	r := strings.NewReplacer(
		"$date", dateNow.String(),
		"$statusType", statusType,
		"$errCode", errCode,
		"$errMsg", errMsg,
		"$endpointErrCode", endpointErrCode,
		"$responseTime", responseTime,
	)
	patternTrns = r.Replace(log.patternTrns)
	log.Debug.Println(fmt.Sprintf(patternTrns))
}

func (log Log) PrintEndpointTrns(sourceSystemID, sessionID, trnsID, subrnumb, requestIP, serviceName, funcName, serviceType, endpointServiceName, endpointStatusType, endpointStatusCode, endpointErrCode, responseTime string) {
	hostname, _ := os.Hostname()
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	dateNow := newTimeFmt(time.Now().In(bangkok), "2006-01-02T15:04:05")
	logTrns := logTrnsEndpointPattern{
		Date:                dateNow.String(),
		HostName:            hostname,
		SourceSystemID:      sourceSystemID,
		SessionID:           sessionID,
		TrnsID:              trnsID,
		Subrnumb:            subrnumb,
		RequestIP:           requestIP,
		ServiceName:         serviceName,
		FuncName:            funcName,
		ServiceType:         serviceType,
		EndpointServiceName: endpointServiceName,
		EndpointStatusType:  endpointStatusType,
		EndpointStatusCode:  endpointStatusCode,
		EndpointErrCode:     endpointErrCode,
		ResponseTime:        responseTime,
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

func InitDebuglog(filename, sessionID, trnsID, srcSysName string) LogLevel {
	var logs LogLevel
	fileDebug := CreateLogFile(filename + "_Debug.log")
	logs.Debuglog = NewDebugLog(fileDebug, "DEBUG", srcSysName)
	logs.Infolog = NewDebugLog(fileDebug, "INFO ", srcSysName)
	logs.Errorlog = NewDebugLog(fileDebug, "ERROR", srcSysName)
	logs.Debuglog.SetInitlogDetail(sessionID, trnsID)
	logs.Infolog.SetInitlogDetail(sessionID, trnsID)
	logs.Errorlog.SetInitlogDetail(sessionID, trnsID)
	return logs
}

func InitTrnslog(filename, sourceSystemID, sessionID, trnsID, subrnumb, requestIP, serviceName, funcName string) Log {
	fileTrns := CreateLogFile(filename + "_Transaction.log")
	trnslog := NewTrnsLog(fileTrns, sourceSystemID, sessionID, trnsID, subrnumb, requestIP, serviceName, funcName)
	return trnslog
}

func InitEndpointTrnslog(filename string) Log {
	fileTrns := CreateLogFile(filename + "_Endpoint.log")
	trnslog := NewEndpointTrnsLog(fileTrns)
	return trnslog
}
