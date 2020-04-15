package smtp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

const (
	TimeFormat                 = "Jan  2 15:04:05"
	TimeRegexpFormat           = `([A-Za-z]{3}\s*[0-9]{1,2} [0-9]{2}:[0-9]{2}:[0-9]{2})`
	HostRegexpFormat           = `(\w+)-(\w+)`
	senderRegexpFormat         = `(from=<)([\w\.-]+@[\w\.-]+)(>)`
	receiveRegexpFormat        = `(to=<)([\w\.-]+@[\w\.-]+)(>)`
	clientHostnameRegexpFormat = `client=(?:[-\w.]|(?:%[\da-fA-F]{2}))+`
	ipaddrRegexpFormat         = `\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`
	messageRegexpFormat        = `(message-id=<)([\w\.-]+@[\w\.-]+)(>)`
	statusRegexpFormat         = `status=(\w+)`
	elhoRegexpFormat           = `ehlo=(\w+)`
	starttlsRegexpFormat       = `starttls=(\w+)`
	rcptRegexpFormat           = `rcpt=(\w+)`
	dataRegexpFormat           = `data=(\w+)`
	quitRegexpFormat           = `quit=(\w+)`
	commandRegexpFormat        = `commands=(\w+)`
)

type (
	LogFormat struct {
		Time     string `json:"time"`
		Hostname string `json:"hostname"`
		From     string `json:"from"`

		ClientHostname string `json:"client_hostname"`
		ClinetIP       string `json:"client_ip"`
		MessageID      string `json:"message_id"`
		To             string `json:"to"`
		Status         string `json:"status"`
		Ehlo           string `json:"ehlo"`
		Starttls       string `json:"starttls"`
		Rcpt           string `json:"rcpt"`
		Data           string `json:"data"`
		Quit           string `json:"quit"`
		Commands       string `json:"commands"`
	}
)

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "    ")
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func initLogParse(filename string) {

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	logFormat, err := Parse(bytes)
	if err != nil {
		fmt.Errorf("parse error")
		os.Exit(1)
	}

	data, _ := JSONMarshal(logFormat)

	fmt.Print(string(data))
}

func Parse(text []byte) (LogFormat, error) {
	logFormat := LogFormat{
		Time:           string(regexp.MustCompile(TimeRegexpFormat).Find(text)),
		Hostname:       string(regexp.MustCompile(HostRegexpFormat).Find(text)),
		ClientHostname: string(regexp.MustCompile(clientHostnameRegexpFormat).Find(text)),
		ClinetIP:       string(regexp.MustCompile(ipaddrRegexpFormat).Find(text)),
		From:           string(regexp.MustCompile(senderRegexpFormat).Find(text)),
		To:             string(regexp.MustCompile(receiveRegexpFormat).Find(text)),
		MessageID:      string(regexp.MustCompile(messageRegexpFormat).Find(text)),
		Status:         string(regexp.MustCompile(statusRegexpFormat).Find(text)),
		Ehlo:           string(regexp.MustCompile(elhoRegexpFormat).Find(text)),
		Starttls:       string(regexp.MustCompile(starttlsRegexpFormat).Find(text)),
		Rcpt:           string(regexp.MustCompile(rcptRegexpFormat).Find(text)),
		Data:           string(regexp.MustCompile(dataRegexpFormat).Find(text)),
		Quit:           string(regexp.MustCompile(quitRegexpFormat).Find(text)),
		Commands:       string(regexp.MustCompile(commandRegexpFormat).Find(text)),
	}

	return logFormat, nil
}