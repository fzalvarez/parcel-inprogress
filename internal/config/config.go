package config

// DBConfig contiene la configuración de conexión a PostgreSQL
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// Config contiene toda la configuración de la aplicación
type Config struct {
	DB          DBConfig
	ServerPort  string
	Environment string
}
