package ssh

import (
	"database/sql"
	"fmt"
	"log"

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

func CreateSession(conn *ssh.Client) *ssh.Session {
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
		return nil
	}
	return session
}

func Deploy(rootPath *string, shellCommand *sql.NullString) *[]byte {
	conn := Connection()
	defer conn.Close()
	session := CreateSession(conn)
	if session == nil {
		log.Fatal("session is nil")
		return nil
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
	fmt.Println(command)
	output, err := session.CombinedOutput(command)
	if err != nil {
		log.Fatalf("Failed to execute command: %s", err)
		return nil
	}
	return &output
}
