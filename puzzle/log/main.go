package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const logFormat = "2006-01-02 15:04:05"

var count int
var locked bool
var f *os.File

type LogEntry struct {
	Time time.Time `json:"time",yaml:"time"`
	Msg  string    `json:"msg",yaml:"msg"`
	Open bool      `json:"open",yaml:"open"`
}

func main() {
	startTime := time.Date(2020, 12, 12, 0, 0, 0, 0, time.UTC)
	actTime := startTime.Add(1 * time.Minute)
	endTime := time.Date(2020, 12, 14, 0, 0, 0, 0, time.UTC)
	locked := true

	var logs []LogEntry
	data, err := ioutil.ReadFile("./puzzle/log/log.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(data, &logs)
	if err != nil {
		fmt.Printf("%v\r\n", err)
	}

	f, err = os.Create("./lock.log")
	if err != nil {
		fmt.Printf("%v\r\n", err)
	}
	defer f.Close()

	count = 0
	for {
		for _, l := range logs {
			if l.Time.After(startTime) && l.Time.Before(actTime) {
				info(l.Time, "event", l.Msg)
				locked = !l.Open
			}
		}
		startTime = actTime
		if locked {
			info(actTime, "info", "health check: central lock: ok, locked")
			actTime = actTime.Add(1 * time.Minute)
		} else {
			info(actTime, "alert", "health check: central lock: ok, unlocked")
			actTime = actTime.Add(10 * time.Second)
		}
		count++
		if actTime.After(endTime) {
			break
		}
	}
}

func info(t time.Time, level, msg string) {
	line := fmt.Sprintf("%04d %s: %s: %s\r\n", count, t.Format(logFormat), level, msg)
	fmt.Print(line)
	f.WriteString(line)
}
