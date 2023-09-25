package db

import (
	"go-deploy/ssh"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

type DeployConfig struct {
	Id           int
	RootPath     string
	Config       sql.NullString
	ShellCommand sql.NullString
	Name         string
}

func connection() (*sql.DB, error) {
	env := ssh.GetSshEnv()
	fmt.Println(env.Host)
	db, err := sql.Open("mysql", env.MysqlHost)
	if err != nil {
		log.Fatal("mysql connection failed :", err)
		return nil, err
	}
	return db, nil
}

func QueryConfig(id int) (*DeployConfig, error) {
	db, err := connection()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer db.Close()
	result, err := db.Query("select id,rootPath,config,shellCommand,name from deploy_config where id = ?", strconv.Itoa(id))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var config = DeployConfig{}
	for result.Next() {
		err = result.Scan(&config.Id, &config.RootPath, &config.Config, &config.ShellCommand, &config.Name)
		if err != nil {
			log.Fatal("result parse failed:", err)
			return nil, err
		}
	}
	log.Fatalf("query successful: %v\n", config)
	return &config, nil
}
