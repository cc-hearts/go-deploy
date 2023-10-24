package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-deploy/ssh"
	"log"
)

type DeployConfig struct {
	Id           int            `db:"id"`
	RootPath     string         `db:"root_path"`
	Config       sql.NullString `db:"config"`
	ShellCommand sql.NullString `db:"shell_command"`
	Name         string         `db:"name"`
}

func connection() (*sql.DB, error) {
	env := ssh.GetSshEnv()
	db, err := sql.Open("mysql", env.MysqlHost)
	if err != nil {
		log.Fatal("mysql connection failed :", err)
		return nil, err
	}
	return db, nil
}

func init() {
	CreateDataBase()
}

func CreateDataBase() {
	db, err := connection()
	if err != nil {
		log.Fatal("mysql connection failed :", err)
		return
	}
	_, err = db.Exec(`create table if not exists deploy_config
											(
											    id           int auto_increment
											        primary key,
											    config       text                                     null comment '配置文件',
											    rootPath     varchar(255)                             not null comment '项目根路径',
											    shellCommand varchar(255)                             null comment 'shell命令',
											    name         varchar(255)                             not null comment '部署名称',
											    createTime   datetime(6) default CURRENT_TIMESTAMP(6) not null,
											    updateTime   datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6)
											)`)
	if err != nil {
		log.Fatal("mysql create table failed :", err)
		return
	}
}
func QueryConfig(id int) (*DeployConfig, error) {
	db, err := connection()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer db.Close()
	result, err := db.Query(`select id,root_path,config,shell_command,name from c_deploy where id = ?`, id)
	defer result.Close()
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
	return &config, nil
}
