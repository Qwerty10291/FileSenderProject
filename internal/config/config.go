package config

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"os"

	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"
)

type Config struct {
	SecretKey string `yaml:"secret_key" json:"secret_key"`
	JwtTTL    int64  `yaml:"jwt_ttl" json:"jwt_ttl"`
	Db        struct {
		Database string `yaml:"database" json:"database"`
		Host     string `yaml:"host" json:"host"`
		Port     string `yaml:"port" json:"port"`
		Login    string `yaml:"login" json:"login"`
		Password string `yaml:"password" json:"password"`
	} `yaml:"db" json:"db"`
}

func NewConfigFromJson(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil{
		return nil, err
	}
	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil{
		return nil, err
	}
	return config, nil
}

func NewConfigFromYaml(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil{
		return nil, err
	}
	decoder := yaml.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil{
		return nil, err
	}
	return config, nil
}

func (c *Config) GetDbClient() (*sqlx.DB, error) {
	db, err := sql.Open(c.Db.Database, fmt.Sprintf("%s://%s:%s@%s:%s", c.Db.Database, c.Db.Host, c.Db.Port, c.Db.Login, c.Db.Password))
	if err != nil{
		return nil, err
	}
	return sqlx.NewDb(db, c.Db.Database), nil
}
