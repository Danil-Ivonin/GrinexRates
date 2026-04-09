package config

import (
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	// PostgreSQL
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// gRPC
	GRPCPort string

	// Grinex HTTP client
	GrinexURL     string
	GrinexTimeout time.Duration
}

// Load reads configuration
func Load() *Config {
	_ = godotenv.Load()

	v := viper.New()

	// Define CLI flags
	pflag.String("db-host", "", "PostgreSQL host (env: DB_HOST)")
	pflag.String("db-port", "", "PostgreSQL port (env: DB_PORT)")
	pflag.String("db-user", "", "PostgreSQL user (env: DB_USER)")
	pflag.String("db-password", "", "PostgreSQL password (env: DB_PASSWORD)")
	pflag.String("db-name", "", "PostgreSQL database name (env: DB_NAME)")
	pflag.String("db-sslmode", "", "PostgreSQL SSL mode (env: DB_SSLMODE)")
	pflag.String("grpc-port", "", "gRPC server port (env: GRPC_PORT)")
	pflag.String("grinex-url", "", "Grinex rates endpoint (env: GRINEX_URL)")
	pflag.Duration("grinex-timeout", 0, "Grinex HTTP request timeout (env: GRINEX_TIMEOUT)")
	pflag.Parse()
	_ = v.BindPFlags(pflag.CommandLine)

	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(".")

	v.SetDefault("db.host", "localhost")
	v.SetDefault("db.port", "5432")
	v.SetDefault("db.user", "postgres")
	v.SetDefault("db.name", "usdt_rate")
	v.SetDefault("db.sslmode", "disable")
	v.SetDefault("grpc.port", "50051")
	v.SetDefault("grinex.url", "https://grinex.io/api/v1/spot/depth?symbol=usdta7a5")
	v.SetDefault("grinex.timeout", 10*time.Second)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	_ = v.ReadInConfig()

	return &Config{
		DBHost:        v.GetString("db.host"),
		DBPort:        v.GetString("db.port"),
		DBUser:        v.GetString("db.user"),
		DBPassword:    v.GetString("db.password"),
		DBName:        v.GetString("db.name"),
		DBSSLMode:     v.GetString("db.sslmode"),
		GRPCPort:      v.GetString("grpc.port"),
		GrinexURL:     v.GetString("grinex.url"),
		GrinexTimeout: v.GetDuration("grinex.timeout"),
	}
}

// DSN returns a PostgreSQL connection string
func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}
