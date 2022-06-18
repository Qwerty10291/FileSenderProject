package config

type Config struct{
	SecretKey string `yaml:"secret_key"`
	JwtTTL int64 `yaml:"jwt_ttl"`
}