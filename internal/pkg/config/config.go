package config

type SApps struct {
	BasicAuthStatic string `env:"MIDDLEWARE_AUTH"`
	Port            string `env:"PORT"`
}

type SRedis struct {
	Host        string `env:"IN_MEMORY_CACHE_URL"`
	Auth        string `env:"IN_MEMORY_CACHE_AUTH_URL"`
	CacheExpiry int    `env:"IN_MEMORY_CACHE_EXPIRY"`
	Port        int    `env:"IN_MEMORY_CACHE_PORT"`
}

type Config struct {
	Redis SRedis
	Apps  SApps
}
