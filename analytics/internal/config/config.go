package config

import "os"

type Config struct {
    DBURL       string
    ServiceName string
}

func Load() Config {
    return Config{
        DBURL:       getenv("ATONIX_DB_URL", "postgres://analytics:analytics@localhost:5432/analytics?sslmode=disable"),
        ServiceName: getenv("ATONIX_SERVICE_NAME", "atonixcorp-analytics"),
    }
}

func getenv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}
