package env

type Config struct {
	EnvFile   string    `mapstructure:"envFile"`
	EnvMode   string    `mapstructure:"envMode"`
	AppConfig EnvConfig `mapstructure:"envLib"` // renamed from envconfig to envLib
}
