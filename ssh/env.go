package ssh

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type InnerEnv struct {
	Host      string
	User      string
	Password  string
	MysqlHost string
	Port      string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("env file is not found")
	}
}

var sshEnv *InnerEnv

func GetSshEnv() *InnerEnv {
	if sshEnv == nil {
		sshEnv = &InnerEnv{
			Host:      os.Getenv("SSH_HOST"),
			User:      os.Getenv("SSH_USER"),
			Password:  os.Getenv("SSH_PASSWORD"),
			MysqlHost: os.Getenv("SSH_MYSQL_HOST"),
			Port:      os.Getenv("PORT"),
		}
	}
	return sshEnv
}
