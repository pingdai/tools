package gormx

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Gormx struct {
	DB           *gorm.DB
	MaxOpenConns int    `json:"max_open_conns"`
	User         string `json:"user"`
	Password     string `json:"password"`
	DBName       string `json:"db_name"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	IsSetLogger  bool   `json:"is_set_logger"` // 是否打印sql日志

	connInfo string
	init     bool
}

func (gormx *Gormx) MarshalDefaults() {
	if gormx.MaxOpenConns == 0 {
		gormx.MaxOpenConns = 20
	}

	if gormx.User == "" {
		panic("DB user cannot be empty")
	}
	if gormx.Password == "" {
		panic("DB pwd cannot be empty")
	}
	if gormx.DBName == "" {
		panic("DB name cannot be empty")
	}
	if gormx.Host == "" {
		gormx.Host = "127.0.0.1"
	}
	if gormx.Port == 0 {
		gormx.Port = 3306
	}
	gormx.connInfo = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		gormx.User, gormx.Password, gormx.Host, gormx.DBName)
}

func (gormx *Gormx) Init() {
	if !gormx.init {
		gormx.New()
		gormx.init = true
	}
}

func (gormx *Gormx) New() {
	if gormx.init {
		return
	}

	gormx.MarshalDefaults()

	db, err := gorm.Open("mysql", gormx.connInfo)
	if err != nil {
		panic(fmt.Sprintf("Gorm open mysql err:%v", err))
	}

	db.DB().SetMaxOpenConns(gormx.MaxOpenConns)
	db.DB().SetMaxIdleConns(gormx.MaxOpenConns)
	if gormx.IsSetLogger {
		db.LogMode(true)
		db.SetLogger(logrus.StandardLogger()) // 默认级别为info
	}
	db.SingularTable(true)

	if err = db.DB().Ping(); err != nil {
		panic(fmt.Sprintf("DB ping err:%v", err))
	}

	gormx.DB = db
}
