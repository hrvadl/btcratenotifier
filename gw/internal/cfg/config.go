package cfg

type Config struct {
	MailerAddr      string
	RateWatcherAddr string
	Port            string
	LogLevel        string
}

func NewFromEnv() (Config, error) {
	return Config{}, nil
}
