package dbconfig

type Config struct {
	DBMysqlConfig       map[string]string `mapstructure:"DBMysqlConfig"`
	DBPostgresConfig    map[string]string `mapstructure:"DBPostgresConfig"`
	EnableAutoMigration bool              `mapstructure:"enableAutoMigration"`
}
