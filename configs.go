package main

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

var dbConfig DbInfo
var apiConfig VerifoneAPIConfig
var amqpConfig AmqpConfig

func get_environment() string {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}
	return environment
}

func init_db_config() {
	data, err := ioutil.ReadFile("config/database.yml")
	if err != nil {
		panic(err)
	}

	m := make(map[string]DbInfo)
	err = yaml.Unmarshal(data, &m)
	dbConfig = m[get_environment()]
}

func init_api_config() {
	data, err := ioutil.ReadFile("config/verifone_api_creds.yml")
	if err != nil {
		panic(err)
	}

	m := make(map[string]VerifoneAPIConfig)
	err = yaml.Unmarshal(data, &m)
	apiConfig = m[get_environment()]
}

func init_amqp_config() {
	data, err := ioutil.ReadFile("config/amqp.yml")
	if err != nil {
		panic(err)
	}

	m := make(map[string]AmqpConfig)
	err = yaml.Unmarshal(data, &m)
	amqpConfig = m[get_environment()]
}

func init_configs() {
	init_db_config()
	init_api_config()
	init_amqp_config()
}

func get_db_config() DbInfo {
	return dbConfig
}

func get_api_config() VerifoneAPIConfig {
	return apiConfig
}

func get_amqp_config() AmqpConfig {
	return amqpConfig
}
