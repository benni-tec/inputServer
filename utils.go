package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
	"time"
)

var config = func() *viper.Viper {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("./")
	vp.AddConfigPath("./data")

	if err := vp.ReadInConfig(); err != nil {
		fmt.Printf("----------> Error reading config file: %v", err)
	}

	vp.SetDefault("port", ":25139")
	vp.SetDefault("verbosity", 0)
	vp.SetDefault("log-file", "")

	return vp
}()

type Error struct {
	error
	Code        string
	Location    string
	Description string
	Cause       string
}

func (err Error) Set(e error) Error {
	err.error = e
	return err
}

func (err Error) Text(location ...bool) string {
	var out string
	if len(location) > 0 && location[0] {
		out += err.Location + " "
	}
	out += "(" + err.Code + ") "
	out += err.Description + " "
	if err.Cause != "" {
		out += "-> " + err.Cause
	}
	return out
}

func (err Error) Throw(dontLog ...bool) {
	BaseLog("error", err.Location, err.Text(), VerbosityError, dontLog...)
}

func (err Error) Encode() []byte {
	err.Throw()
	if err.Cause == "" {
		err.Cause = err.Error()
	}
	out, _ := json.Marshal(err)
	return out
}

const (
	VerbosityError = 0
	VerbosityWarn  = 1
	VerbosityInfo  = 2
)

// Method to log events
func LogInfo(location string, str string) {
	BaseLog("info", location, str, VerbosityInfo)
}

func LogWarn(location string, str string) {
	BaseLog("warn", location, str, VerbosityWarn)
}

func LogOk(location string, str string) {
	BaseLog("ok", location, str, VerbosityError)
}

var lockTerminal sync.Mutex

func BaseLog(tag string, location string, str string, verbosity int, dontLog ...bool) {
	if verbosity <= config.GetInt("verbosity") { // if verbosity level is equal or below set
		t := time.Now()                                                                                                         // getting current time
		timestamp := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()) // creating timestamp

		lockTerminal.Lock()

		fmt.Print(timestamp + " [")

		switch strings.ToUpper(tag) {
		case "WARN":
			color.Set(color.FgYellow, color.Bold)
		case "ERROR":
			color.Set(color.FgRed, color.Bold)
			tag = "EROR"
		case "OK":
			color.Set(color.FgGreen, color.Bold)
			tag = " OK "
		default:
			color.Set(color.Bold)
		}

		fmt.Print(strings.ToUpper(tag))

		color.Unset()

		fmt.Println("] [" + location + "] " + str)

		lockTerminal.Unlock()

		v := timestamp + " [" + strings.ToUpper(tag) + "] [" + location + "] " + str // creating

		// open file
		logFile := config.GetString("log-file")
		if logFile != "" && !(len(dontLog) > 0 && dontLog[0]) {
			file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
			if err != nil {
				Error{
					Code:        "file-open",
					Location:    "BaseLog",
					Description: "Unable to open log file!",
				}.Set(err).Throw(true)
			} else {
				_, err2 := fmt.Fprintln(file, v)
				if err2 != nil {
					Error{
						Code:        "file-write",
						Location:    "BaseLog",
						Description: "Unable to write to log file!",
					}.Set(err2).Throw(true)
				}
			}

			err3 := file.Close()
			if err3 != nil {
				Error{
					Code:        "file-close",
					Location:    "BaseLog",
					Description: "Error closing log file!",
				}.Set(err3).Throw(true)
				return
			}
		}
	}
}
