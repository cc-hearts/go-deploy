package service

import (
	"fmt"
	"go-deploy/db"
	"go-deploy/ssh"
	"net/http"
	"strconv"
)

type SuccessMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func initDeployHttp() {
	http.HandleFunc("/deploy", func(w http.ResponseWriter, r *http.Request) {
		// http 获取 get后面的Query参数
		if r.Method != http.MethodGet {
			http.Error(w, "request method failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		queryParams := r.URL.Query()
		queryId := queryParams.Get("id")
		id, err := strconv.Atoi(queryId)
		if err != nil {
			http.Error(w, "id field is invalid", http.StatusInternalServerError)
			return
		}
		result, err := db.QueryConfig(id)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		ssh.Deploy(&result.RootPath, &result.ShellCommand, &w, r)
	})
}

func HttpBootStrap() {
	env := ssh.GetSshEnv()
	if env.Port == "" {
		panic("port is empty")
	}
	initDeployHttp()
	err := http.ListenAndServe(env.Port, nil)
	if err != nil {
		panic(err)
	}
}
