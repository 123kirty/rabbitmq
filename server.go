package eatask

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	AMQP_USERNAME = "admin"
	AMQP_PASSWORD = "password"
	AMQP_HOST     = "localhost"
	AMQP_PORT     = "5672"

	DB_USERNAME = "root"
	DB_PASSWORD = "Xfc.e/Mq4lT5"
	DB_HOST     = "127.0.0.1"
	DB_DBNAME   = "eatask"

	EaTaskQueueName = "eataskqueue"
)

type (
	EaServer struct {
		Db     *gorm.DB
		Config *Config
	}

	Config struct {
		Receiver struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Host     string `json:"host"`
			Port     string `json:"port"`
		} `json:"receiver"`
		Db struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Host     string `json:"host"`
			DbName   string `json:"db_name"`
		}
	}
)

func NewEaServer() (eaServer *EaServer, err error) {

	eaServer = &EaServer{}

	if eaServer.Config, err = NewConfig(); err != nil {
		return
	}

	if eaServer.Db, err = eaServer.prepareDb(); err != nil {
		return
	}

	return
}

func NewConfig() (config *Config, err error) {

	config = &Config{}

	config.Receiver.Username = AMQP_USERNAME
	config.Receiver.Password = AMQP_PASSWORD
	config.Receiver.Host = AMQP_HOST
	config.Receiver.Port = AMQP_PORT

	config.Db.Username = DB_USERNAME
	config.Db.Password = DB_PASSWORD
	config.Db.Host = DB_HOST
	config.Db.DbName = DB_DBNAME

	return
}

func (eaServer *EaServer) GetQueueUrl() (url string) {

	url = fmt.Sprintf("amqp://%s:%s@%s:%s/",
		eaServer.Config.Receiver.Username,
		eaServer.Config.Receiver.Password,
		eaServer.Config.Receiver.Host,
		eaServer.Config.Receiver.Port,
	)
	return
}

func (eaServer *EaServer) prepareDb() (db *gorm.DB, err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		eaServer.Config.Db.Username,
		eaServer.Config.Db.Password,
		eaServer.Config.Db.Host,
		eaServer.Config.Db.DbName)

	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return
	}

	db.AutoMigrate(&Hotel{},
		&Room{},
		&RatePlan{})

	return
}

func (eaServer *EaServer) Run() (err error) {

	eaServer.startServer()

	return
}
