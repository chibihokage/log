package log

import (
	"bytes"
	"testing"
	"time"
)

func TestDebugLog(t *testing.T) {
	w := bytes.NewBuffer([]byte{})

	level := "DEBUG"
	logdebug := NewDebugLog(w, level)
	refnumb := "esvFEweb2017"
	subrnumb := "66987654321"
	service := "esvFEweb"
	proc := "PROC"
	logdebug.SetInitlogDetail(service, proc)
	logdebug.SetRecordDetail(refnumb, subrnumb)
	logdebug.Println("msg")

	expected := "|DEBUG|esvFEweb2017|66987654321|esvFEweb|PROC|msg\n"

	if string(w.String()[19:]) != expected {
		t.Errorf("Expected %#v \n        but got  %#v", expected, string(w.String()[19:]))
		t.Errorf("Expected %v but got %v", len(expected), len(w.String()))
	}
}
func TestInfoLog(t *testing.T) {

	w := bytes.NewBuffer([]byte{})
	level := "INFO"
	loginfo := NewDebugLog(w, level)
	refnumb := "esvFEweb2017"
	subrnumb := "66987654321"
	service := "esvFEweb"
	proc := "PROC"
	loginfo.SetInitlogDetail(service, proc)
	loginfo.SetRecordDetail(refnumb, subrnumb)
	loginfo.Println("msg")

	expected := "|INFO|esvFEweb2017|66987654321|esvFEweb|PROC|msg\n"

	if string(w.String()[19:]) != expected {
		t.Errorf("Expected %#v \n        but got  %#v", expected, w.String())
		t.Errorf("Expected %v but got %v", len(expected), len(w.String()))
	}
}
func TestErrorLog(t *testing.T) {

	w := bytes.NewBuffer([]byte{})
	level := "ERROR"
	logerror := NewDebugLog(w, level)
	refnumb := "esvFEweb2017"
	subrnumb := "66987654321"
	service := "esvFEweb"
	proc := "PROC"
	logerror.SetInitlogDetail(service, proc)
	logerror.SetRecordDetail(refnumb, subrnumb)
	logerror.Println("msg")

	expected := "|ERROR|esvFEweb2017|66987654321|esvFEweb|PROC|msg\n"

	if string(w.String()[19:]) != expected {
		t.Errorf("Expected %#v \n        but got  %#v", expected, w.String())
		t.Errorf("Expected %v but got %v", len(expected), len(w.String()))
	}
}

func TestPrintf(t *testing.T) {

	w := bytes.NewBuffer([]byte{})
	level := "ERROR"
	logerror := NewDebugLog(w, level)
	refnumb := "esvFEweb2017"
	subrnumb := "66987654321"
	service := "esvFEweb"
	proc := "PROC"
	logerror.SetInitlogDetail(service, proc)
	logerror.SetRecordDetail(refnumb, subrnumb)
	test := 1
	logerror.Printf("msg: %#v", test)

	expected := "|ERROR|esvFEweb2017|66987654321|esvFEweb|PROC|msg: 1\n"

	if string(w.String()[19:]) != expected {
		t.Errorf("Expected %#v \n        but got  %#v", expected, w.String()[19:])
		t.Errorf("Expected %v but got %v", len(expected), len(w.String()))
	}
}
func TestPrintTrns(t *testing.T) {

	w := bytes.NewBuffer([]byte{})

	logtrns := NewTrnsLog(w)
	subrnumb := "66987654321"
	service := "esvFEweb"
	status := "S"
	errcode := "T1000"
	logtrns.PrintTrns(service, subrnumb, status, errcode)

	expected := "|esvFEweb|66987654321|S|T1000\n"

	if string(w.String()[19:]) != expected {
		t.Errorf("Expected %#v \n        but got  %#v", expected, w.String()[19:])
		t.Errorf("Expected %v but got %v", len(expected), len(w.String()[19:]))
	}
}
func TestTimeFormat(t *testing.T) {
	myNow := newTimeFmt(time.Date(2018, 01, 18, 0, 0, 0, 0, time.UTC), "2006-01-02")
	if myNow.String() != "2018-01-18" {
		t.Error("wrong format", myNow)
	}
}
