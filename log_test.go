package log

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestDebugLog(t *testing.T) {
	w := bytes.NewBuffer([]byte{})

	level := "DEBUG"
	logdebug := NewDebugLog(w, level, "EVOUCHER")
	subrnumb := "66987654321"
	sessionID := "SessionID"
	trnsID := "trnsID"
	logdebug.SetInitlogDetail(sessionID, trnsID)
	logdebug.SetRecordDetail(subrnumb)
	logdebug.Println("msg")

	expected := "|DEBUG|EVOUCHER|SessionID|trnsID|66987654321|msg\n"

	if string(w.String()[19:]) != expected {
		t.Errorf("Expected %#v \n        but got  %#v", expected, string(w.String()[19:]))
		t.Errorf("Expected %v but got %v", len(expected), len(w.String()))
	}
}
func TestInfoLog(t *testing.T) {

	w := bytes.NewBuffer([]byte{})
	level := "INFO"
	loginfo := NewDebugLog(w, level, "EVOUCHER")
	subrnumb := "66987654321"
	sessionID := "SessionID"
	trnsID := "trnsID"
	loginfo.SetInitlogDetail(sessionID, trnsID)
	loginfo.SetRecordDetail(subrnumb)
	loginfo.Println("msg")

	expected := "|INFO|EVOUCHER|SessionID|trnsID|66987654321|msg\n"

	if string(w.String()[19:]) != expected {
		t.Errorf("Expected %#v \n        but got  %#v", expected, w.String())
		t.Errorf("Expected %v but got %v", len(expected), len(w.String()))
	}
}
func TestErrorLog(t *testing.T) {

	w := bytes.NewBuffer([]byte{})
	level := "ERROR"
	logerror := NewDebugLog(w, level, "EVOUCHER")
	subrnumb := "66987654321"
	sessionID := "SessionID"
	trnsID := "trnsID"
	logerror.SetInitlogDetail(sessionID, trnsID)
	logerror.SetRecordDetail(subrnumb)
	logerror.Println("msg")

	expected := "|ERROR|EVOUCHER|SessionID|trnsID|66987654321|msg\n"

	if string(w.String()[19:]) != expected {
		t.Errorf("Expected %#v \n        but got  %#v", expected, w.String())
		t.Errorf("Expected %v but got %v", len(expected), len(w.String()))
	}
}

func TestPrintf(t *testing.T) {

	w := bytes.NewBuffer([]byte{})
	level := "ERROR"
	logerror := NewDebugLog(w, level, "EVOUCHER")
	subrnumb := "66987654321"
	sessionID := "SessionID"
	trnsID := "trnsID"
	logerror.SetInitlogDetail(sessionID, trnsID)
	logerror.SetRecordDetail(subrnumb)
	test := 1
	logerror.Printf("msg: %#v", test)

	expected := "|ERROR|EVOUCHER|SessionID|trnsID|66987654321|msg: 1\n"

	if string(w.String()[19:]) != expected {
		t.Errorf("Expected %#v \n        but got  %#v", expected, w.String()[19:])
		t.Errorf("Expected %v but got %v", len(expected), len(w.String()))
	}
}
func TestPrintTrns(t *testing.T) {

	w := bytes.NewBuffer([]byte{})
	sourceSystemID := "EVOUCHER"
	sessionID := "12345"
	trnsID := "1233244"
	subrnumb := "66987654321"
	requestIP := "2313121"
	serviceName := "Test"
	funcName := "TestFunc"
	logtrns := NewTrnsLog(w, sourceSystemID, sessionID, trnsID, requestIP, serviceName, funcName)
	statusType := "S"
	errCode := "0"
	errMsg := "message"
	endpointErrCode := "0"
	responseTime := "234"
	logtrns.PrintTrns(subrnumb, statusType, errCode, errMsg, endpointErrCode, responseTime)
	hostname, _ := os.Hostname()
	expected := "|" + hostname + "|EVOUCHER|12345|1233244|66987654321|2313121|Test|TestFunc|S|0|message|0|234\n"
	if string(w.String()[19:]) != expected {
		t.Errorf("Expected %#v \n        but got  %#v", expected, w.String()[19:])
		t.Errorf("Expected %v but got %v", len(expected), len(w.String()[19:]))
	}
}

func TestPrintEndpointTrns(t *testing.T) {

	w := bytes.NewBuffer([]byte{})

	sourceSystemID := "EVOUCHER"
	sessionID := "12345"
	trnsID := "1233244"
	subrnumb := "66987654321"
	requestIP := "2313121"
	serviceName := "Test"
	funcName := "TestFunc"
	logtrns := NewEndpointTrnsLog(w, sourceSystemID, sessionID, trnsID, requestIP, serviceName, funcName)
	serviceType := "REWARD"
	endpointServiceName := "enquiryPrivilege"
	endpointStatusType := ""
	endpointStatusCode := ""
	endpointErrCode := "0"
	responseTime := "234"
	logtrns.PrintEndpointTrns(subrnumb, serviceType, endpointServiceName, endpointStatusType, endpointStatusCode, endpointErrCode, responseTime)
	hostname, _ := os.Hostname()
	expected := "|" + hostname + "|EVOUCHER|12345|1233244|66987654321|2313121|Test|TestFunc|REWARD|enquiryPrivilege|||0|234\n"
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
