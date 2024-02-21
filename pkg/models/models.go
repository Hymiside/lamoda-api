package models

type ConfigPostgres struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type ConfigServer struct {
	Host string
	Port string
}