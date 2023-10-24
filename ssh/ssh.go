package ssh

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/crypto/ssh"
)

func GenConfig() *ssh.ClientConfig {
	env := GetSshEnv()
	config := &ssh.ClientConfig{
		User: env.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(env.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config
}
func Connection() *ssh.Client {
	env := GetSshEnv()
	config := GenConfig()
	coon, err := ssh.Dial("tcp", env.Host, config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	return coon
}

func printOutput(outputReader io.Reader, w *http.ResponseWriter, isResponseWrite bool) {
	flusher, ok := (*w).(http.Flusher)
	if !ok {
		http.Error(*w, "Streaming not supported", http.StatusInternalServerError)
		return
	}
	scanner := bufio.NewScanner(outputReader)
	for scanner.Scan() {
		line := scanner.Text()
		(*w).Write([]byte(line + "\n"))
		if isResponseWrite == true {
			fmt.Println(line)
			flusher.Flush()
		}
	}
	if scanner.Err() != nil {
		http.Error(*w, scanner.Err().Error(), http.StatusInternalServerError)
	}
}

func CreateSession(conn *ssh.Client) *ssh.Session {
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
		return nil
	}
	return session
}

func Deploy(rootPath *string, shellCommand *sql.NullString, w *http.ResponseWriter, r *http.Request) {
	conn := Connection()
	defer conn.Close()
	session := CreateSession(conn)
	if session == nil {
		http.Error(*w, "Failed to create session", http.StatusInternalServerError)
	}
	defer session.Close()
	configShell := ""

	if shellCommand.Valid == true {
		configShell = shellCommand.String
	}

	var command = "cd " + *rootPath
	if configShell != "" {
		command = command + " && " + configShell
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Start(command)

	ctx := r.Context()
	ch := make(chan struct{})

	// 打印标准输出
	go printOutput(stdout, w, true)

	// 打印标准错误输出
	go printOutput(stderr, w, false)

	go func(ch chan struct{}) {
		err = session.Wait()
		ch <- struct{}{}
	}(ch)

	select {
	case <-ch:
	case <-ctx.Done():
		err := ctx.Err()
		http.Error(*w, err.Error(), http.StatusInternalServerError)
	}
}
