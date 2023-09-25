package service

import (
	"encoding/json"
	"fmt"
	"go-deploy/db"
	"go-deploy/ssh"
	"net/http"
)

type SuccessMsg struct {
	Code int    `json:"Code"`
	Msg  string `json:"Msg"`
}

func initDeployHttp() {
	http.HandleFunc("/deploy", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// http 获取 get后面的Query参数
		if r.Method != http.MethodGet {
			http.Error(w, "request method failed", http.StatusInternalServerError)
			return
		}
		result, err := db.QueryConfig(1)
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

		fmt.Println(string(jsonData))
		w.Write(*logger)

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
