package config

import (
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Load reads configuration
func Load() error {
	_ = godotenv.Load()

	// Define CLI flags
	pflag.String("db-host", "", "PostgreSQL host (enviper: DB_HOST)")
	pflag.String("db-port", "", "PostgreSQL port (env: DB_PORT)")
	pflag.String("db-user", "", "PostgreSQL user (env: DB_USER)")
	pflag.String("db-password", "", "PostgreSQL password (env: DB_PASSWORD)")
	pflag.String("db-name", "", "PostgreSQL database name (env: DB_NAME)")
	pflag.String("db-sslmode", "", "PostgreSQL SSL mode (env: DB_SSLMODE)")
	pflag.String("grpc-port", "", "gRPC server port (env: GRPC_PORT)")
	pflag.String("grinex-url", "", "Grinex rates endpoint (env: GRINEX_URL)")
	pflag.Duration("grinex-timeout", 0, "Grinex HTTP request timeout (env: GRINEX_TIMEOUT)")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "5432")
	viper.SetDefault("db.user", "postgres")
	viper.SetDefault("db.name", "usdt_rate")
	viper.SetDefault("db.sslmode", "disable")
	viper.SetDefault("grpc.port", "50051")
	viper.SetDefault("grinex.url", "https://grinex.io/api/v1/spot/depth?symbol=usdta7a5")
	viper.SetDefault("grinex.timeout", 10*time.Second)
	viper.SetDefault("calculator.avgnm_precision", 8)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	return nil
}

// AvgNMPrecision returns the number of decimal places used when rounding AvgNM results
func AvgNMPrecision() int32 {
	return int32(viper.GetInt("calculator.avgnm_precision"))
}

// DSN returns a PostgreSQL connection string
func DSN() string {
	return "host=" + viper.GetString("db.host") +
		" port=" + viper.GetString("db.port") +
		" user=" + viper.GetString("db.user") +
		" password=" + viper.GetString("db.password") +
		" dbname=" + viper.GetString("db.name") +
		" sslmode=" + viper.GetString("db.sslmode")
}
