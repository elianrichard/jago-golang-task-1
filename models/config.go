package models

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DATABASE_URL"`
}
