package ackbar

import (
	"fmt"
	"time"
	"github.com/pborman/uuid"
	"os"
	"net/http"
	"encoding/json"
	"bytes"
	"runtime"
	"strconv"
	"net/http/httputil"
)

type Log struct {
	LogID string
	AppName string
	Date time.Time
	LogLevel LogLevel
	Message string
	Object string
	HttpBodyDump string
	FuncCall string
}

type LogLevel string

const (
	Info LogLevel = "INFO"
	Warn LogLevel = "WARN"
	Error LogLevel = "ERROR"
	Fatal LogLevel = "FATAL"
	Debug LogLevel = "DEBUG"
	Trace LogLevel = "TRACE"
)

func acceptLog(log Log) error {

	body, err := json.Marshal(log)
	object, err := json.Marshal(log.Object)
	log.Object = string(object)

	_, err = http.Post(fmt.Sprintf("%s/%s/%s", os.Getenv("ES_HOST"), os.Getenv("APP_NAME"), log.LogID),
		"application/json", bytes.NewBuffer(body))

	return err
}

func Trap(magnitude LogLevel, message string, object interface{}) error {

	log := Log{}

	pc, file, line, _ := runtime.Caller(1)
	funcCall := runtime.FuncForPC(pc).Name()

	log.LogID = uuid.New()
	log.AppName = os.Getenv("APP_NAME")
	log.Date = time.Now().UTC()
	log.LogLevel = magnitude
	log.Message = message
	log.Object = fmt.Sprintf("%v", object)
	log.FuncCall = fmt.Sprintf("FILE:%s, FUNC:%s, LINE:%s", file, funcCall, strconv.Itoa(line))

	return acceptLog(log)
}

func TrapWithHttpResponse(magnitude LogLevel, message string, object interface{}, resp *http.Response) error {

	log := Log{}

	pc, file, line, _ := runtime.Caller(1)
	funcCall := runtime.FuncForPC(pc).Name()

	body, _ := httputil.DumpResponse(resp, true)
	log.HttpBodyDump = string(body)

	log.LogID = uuid.New()
	log.AppName = os.Getenv("APP_NAME")
	log.Date = time.Now().UTC()
	log.LogLevel = magnitude
	log.Message = message
	log.Object = fmt.Sprintf("%v", object)
	log.FuncCall = fmt.Sprintf("FILE:%s, FUNC:%s, LINE:%s", file, funcCall, strconv.Itoa(line))

	return acceptLog(log)
}

func TrapWithHttpRequest(magnitude LogLevel, message string, object interface{}, resp *http.Request) error {

	log := Log{}

	pc, file, line, _ := runtime.Caller(1)
	funcCall := runtime.FuncForPC(pc).Name()

	body, _ := httputil.DumpRequest(resp, true)
	log.HttpBodyDump = string(body)

	log.LogID = uuid.New()
	log.AppName = os.Getenv("APP_NAME")
	log.Date = time.Now().UTC()
	log.LogLevel = magnitude
	log.Message = message
	log.Object = fmt.Sprintf("%v", object)
	log.FuncCall = fmt.Sprintf("FILE:%s, FUNC:%s, LINE:%s", file, funcCall, strconv.Itoa(line))

	return acceptLog(log)
}