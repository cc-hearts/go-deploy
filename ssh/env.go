package ssh

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"sync"
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
var once sync.Once

func GetSshEnv() *InnerEnv {
	once.Do(func() {
		sshEnv = &InnerEnv{
			Host:      os.Getenv("SSH_HOST"),
			User:      os.Getenv("SSH_USER"),
			Password:  os.Getenv("SSH_PASSWORD"),
			MysqlHost: os.Getenv("SSH_MYSQL_HOST"),
			Port:      os.Getenv("PORT"),
		}
	})
	return sshEnv
}
