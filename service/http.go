package service

import (
	"encoding/json"
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
		w.Header().Set("Content-Type", "application/json")

		// http 获取 get后面的Query参数
		if r.Method != http.MethodGet {
			http.Error(w, "request method failed", http.StatusInternalServerError)
			return
		}
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
		logger := ssh.Deploy(&result.RootPath, &result.ShellCommand)
		if logger == nil {
			fmt.Println("logger get failed :", logger)
		}

		response := SuccessMsg{
			Code: 200,
			Msg:  fmt.Sprintf("%s", logger),
		}
		fmt.Println("response:", response)
		jsonData, err := json.Marshal(response)
		fmt.Println("json,", jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(jsonData)

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
