package util

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type DbInfo struct {
	UserName string
	Database string
	Password string
	Host     string
	Adapter  string
}

type VerifoneAPIConfig struct {
	UserName string
	Password string
	Host     string
}

type AmqpConfig struct {
	UserName string
	Password string
	Host     string
	Port     int
}

type Config struct {
	DbConfig DbInfo
	ApiConfig VerifoneAPIConfig
	AmqpConfig AmqpConfig
}

func get_environment() string {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}
	return environment
}

func (config *Config) init_db_config() {
	data, err := ioutil.ReadFile("config/database.yml")
	if err != nil {
		panic(err)
	}

	m := make(map[string]DbInfo)
	err = yaml.Unmarshal(data, &m)
	config.DbConfig = m[get_environment()]
}

func (config *Config) init_api_config() {
	data, err := ioutil.ReadFile("config/verifone_api_creds.yml")
	if err != nil {
		panic(err)
	}

	m := make(map[string]VerifoneAPIConfig)
	err = yaml.Unmarshal(data, &m)
	config.ApiConfig = m[get_environment()]
}

func (config *Config) init_amqp_config() {
	data, err := ioutil.ReadFile("config/amqp.yml")
	if err != nil {
		panic(err)
	}

	m := make(map[string]AmqpConfig)
	err = yaml.Unmarshal(data, &m)
	config.AmqpConfig = m[get_environment()]
}

func (config *Config) Init() {
	config.init_db_config()
	config.init_api_config()
	config.init_amqp_config()
}