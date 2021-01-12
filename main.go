package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"net/http"
	"strings"
)

func main() {
	mux := MyServerMux{http.NewServeMux()}
	mux.MyHandleFunc("/tap", tap)
	mux.MyHandleFunc("/print", myPrint)
	mux.MyHandleFunc("/", ping)

	LogOk("main", "Starting server...")
	_ = http.ListenAndServe(config.GetString("port"), mux)
}

func ping(req *http.Request) (int, error) {
	return http.StatusAccepted, nil
}

func tap(req *http.Request) (int, error) {
	if err := req.ParseForm(); err != nil {
		return http.StatusInternalServerError, Error{Code: "http-parse-form", Location: "tap", Description: "Unable to parse Query!"}.Set(err)
	}

	var args []string
	raw := strings.ReplaceAll(req.Form.Get("args"), " ", "")
	if raw != "" {
		args = strings.Split(raw, ",")
	}
	key := req.Form.Get("key")
	LogInfo("tap", fmt.Sprintf("Tap \"%v\" with %#v", key, args))

	if len(args) > 0 {
		robotgo.KeyTap(key, args)
	} else {
		robotgo.KeyTap(key)
	}
	return http.StatusOK, nil
}

func myPrint(req *http.Request) (int, error) {
	if err := req.ParseForm(); err != nil {
		return http.StatusInternalServerError, Error{Code: "http-parse-form", Location: "print", Description: "Unable to parse Query!"}.Set(err)
	}

	value := req.Form.Get("value")
	LogInfo("print", fmt.Sprintf("Print \"%v\"", value))

	robotgo.TypeStr(value)
	return http.StatusOK, nil
}

type MyServerMux struct{ *http.ServeMux }

func (mux *MyServerMux) MyHandleFunc(pattern string, f func(req *http.Request) (int, error)) {
	mux.HandleFunc(pattern, func(res http.ResponseWriter, req *http.Request) {
		code, err := func() (code int, e error) {
			defer func() {
				if ee := recover(); ee != nil {
					code = http.StatusInternalServerError
					e = Error{Code: "unexpected", Location: "*MyServerMux/MyHandleFunc", Description: "Recovered error", Cause: fmt.Sprint(ee)}
				}
			}()
			if req.URL.Path != pattern {
				LogWarn("*MyServerMux/MyHandleFunc", "404 Not Found -> "+req.URL.Path)
				return http.StatusNotFound, nil
			}
			return f(req)
		}()

		res.WriteHeader(code)
		if err != nil {
			switch err.(type) {
			case Error:
				res.Write(err.(Error).Encode())
				break
			default:
				res.Write(Error{
					Code:        "unexpected",
					Location:    "*MyServerMux/MyHandleFunc",
					Description: "Unexpected error occured",
				}.Set(err).Encode())
			}
		}
	})
}
