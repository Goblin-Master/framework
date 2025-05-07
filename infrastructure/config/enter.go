package config

import "fmt"

var Conf = new(Config)

type Config struct {
	App   App   `mapstructure:"app"`
	DB    DB    `mapstructure:"db"`
	Redis Redis `mapstructure:"redis"`
	QiNiu QiNiu `mapstructure:"qiNiu"`
	Kafka Kafka `mapstructure:"kafka"`
}
type App struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Env  string `mapstructure:"env"`
	Log  string `mapstructure:"log"`
}

func (app *App) Link() string {
	return fmt.Sprintf("%s:%s", app.Host, app.Port)
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

func (db *DB) DSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", db.User, db.Password, db.Host, db.Port, db.DBName, "5s")
	return dsn
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Enable   bool   `mapstructure:"enable"`
}

func (redis *Redis) DSN() string {
	dsn := fmt.Sprintf("%s:%d", redis.Host, redis.Port)
	return dsn
}

type QiNiu struct {
	Enable    bool   `mapstructure:"enable"`
	AccessKey string `mapstructure:"accessKey"`
	Bucket    string `mapstructure:"bucket"`
	SecretKey string `mapstructure:"secretKey"`
	Url       string `mapstructure:"url"`
	Region    string `mapstructure:"region"`
	Prefix    string `mapstructure:"prefix"`
}

type Kafka struct {
	Enable bool     `mapstructure:"enable"`
	Addrs  []string `mapstructure:"addrs"`
}
