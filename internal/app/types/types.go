package types

// Config конфиг приложения
type Config struct {
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerPort    string `env:"SERVER_PORT" envDefault:"8080"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	DBPath        string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/short-url-db.json"`
	DatabaseDsn   string `env:"DATABASE_DSN" envDefault:"host=localhost port=5432 user=postgres password=123 dbname=db sslmode=disable"`
}

// URL - структура для url
type URL struct {
	UUID     string `db:"uuid"`
	Hash     string `db:"hash"`
	URL      string `db:"url"`
	ShortURL string `db:"short_url"`
}
