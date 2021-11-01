package apiserver

//Confing Api server
type Config struct {
	BinAddr     string `toml:"bin_addr"`
	DatabaseURL string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{
		BinAddr: "localhost:8080",
	}
}
