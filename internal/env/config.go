package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type ServerConfig struct {
  AppEnv           string `env:"APP_ENV"         envDefault:"dev"`
  Host             string `env:"APP_HOST"        envDefault:"localhost"`
  Port             string `env:"APP_PORT"         envDefault:"8080"`
  Prefix           string `env:"APP_API_PREFIX"  envDefault:"/api/v1"`
  AppVersion       string `env:"APP_VERSION"     envDefault:"0.0.1"`
}

type JWTConfig struct {
  AuthJWTSecretKey string `env:"AUTH_JWT_SECRET_KEY" envDefault:"secret"`
  AuthJWTExpTime   int64  `env:"AUTH_JWT_EXP_TIME" envDefault:"31536000"`
  AuthJWTIssuer    string `env:"AUTH_JWT_ISSUER" envDefault:"ecommerce"`
}

type DBConfig struct {
  ConfigDBUrl      string `env:"CONFIG_DB_URL"`
  ConfigDBHost     string `env:"CONFIG_DB_HOST" envDefault:"localhost"`
  ConfigDBPort     string `env:"CONFIG_DB_PORT" envDefault:"5432"`
  ConfigDBUser     string `env:"CONFIG_DB_USER" envDefault:"root"`
  ConfigDBPassword string `env:"CONFIG_DB_PASSWORD" envDefault:""`
  ConfigDBName     string `env:"CONFIG_DB_NAME" envDefault:"postgres"`
}

type StoreConfig struct {
  StoreDBUrl      string `env:"STORE_DB_URL"`
  StoreDBHost     string `env:"STORE_DB_HOST" envDefault:"localhost"`
  StoreDBPort     string `env:"STORE_DB_PORT" envDefault:"5432"`
  StoreDBUser     string `env:"STORE_DB_USER" envDefault:"postgres"`
  StoreDBPassword string `env:"STORE_DB_PASSWORD" envDefault:"postgres"`
  StoreDBName     string `env:"STORE_DB_NAME" envDefault:"postgres"`
}

type RedisConfig struct {
  RedisUrl      string `env:"REDIS_URL"`
  RedisHost     string `env:"REDIS_HOST" envDefault:"localhost"`
  RedisPort     string `env:"REDIS_PORT" envDefault:"6379"`
  RedisPassword string `env:"REDIS_PASSWORD" envDefault:"redis"`
  RedisDB       string `env:"REDIS_DB" envDefault:"0"`
}

type RabbitmqConfig struct {
  RabbitmqHost     string `env:"RABBITMQ_HOST" envDefault:"localhost"`
  RabbitmqPort     string `env:"RABBITMQ_PORT" envDefault:"5672"`
  RabbitmqUser     string `env:"RABBITMQ_USER" envDefault:"guest"`
  RabbitmqPassword string `env:"RABBITMQ_PASSWORD" envDefault:"guest"`
  RabbitmqVhost    string `env:"RABBITMQ_VHOST" envDefault:"/"`
}

type KafkaConfig struct {
  KafkaUrl  string `env:"KAFKA_URL"`
  KafkaHost string `env:"KAFKA_HOST" envDefault:"localhost"`
  KafkaPort string `env:"KAFKA_PORT" envDefault:"9092"`
}

// Uncomment when needed
// type GRPCConfig struct {
//   GRPCUrl  string `env:"GRPC_URL"`
//   GRPCHost string `env:"GRPC_HOST" envDefault:"localhost"`
//   GRPCPort string `env:"GRPC_PORT" envDefault:"50051"`
// }

type LokiConfig struct {
  LokiUrl  string `env:"LOKI_URL"`
  LokiHost string `env:"LOKI_HOST" envDefault:"localhost"`
  LokiPort string `env:"LOKI_PORT" envDefault:"3100"`
}

type SentryConfig struct {
  SentryUrl  string `env:"SENTRY_URL"`
  SentryHost string `env:"SENTRY_HOST" envDefault:"localhost"`
  SentryPort string `env:"SENTRY_PORT" envDefault:"9000"`
}

type LogConfig struct {
  LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

type ShopeeConfig struct {
  ShopeeApiVersion       string `env:"SHOPEE_API_VERSION"`
  ShopeeApiBaseUrl       string `env:"SHOPEE_API_BASE_URL"`
  ShopeeApiBasePrefix    string `env:"SHOPEE_API_BASE_PREFIX" envDefault:"/api/v2/"`
  ShopeeShopId           string `env:"SHOPEE_SHOP_ID"`
  ShopeeApiUrl           string `env:"SHOPEE_API_URL"`
  ShopeePartnerId        string `env:"SHOPEE_PARTNER_ID"`
  ShopeePartnerSecretKey string `env:"SHOPEE_PARTNER_SECRET_KEY"`
}

type Config struct {
  Server *ServerConfig
  JWT    *JWTConfig
  DB     *DBConfig
  Store  *StoreConfig
  Redis  *RedisConfig
  Rabbitmq *RabbitmqConfig
  Kafka  *KafkaConfig
  Loki   *LokiConfig
  Sentry *SentryConfig
  Log    *LogConfig
  Shopee *ShopeeConfig
}

func LoadEnv(envSet string, logger *zap.Logger) (*Config,error) {

  path := fmt.Sprintf("internal/env/.%s.env",envSet)
  if err := godotenv.Load(path); err != nil {
    logger.Error("failed to load env", zap.String("path", path), zap.Error(err))
    return nil, fmt.Errorf("Env Loader: %w" , err)
  }

  server := &ServerConfig{}
  if err := env.Parse(server); err != nil {
    return nil,fmt.Errorf("failed to parse env: %w", err)
  }

  jwt := &JWTConfig{}
  if err := env.Parse(jwt); err != nil {
    return nil,fmt.Errorf("failed to parse env: %w", err)
  }

  db := &DBConfig{}
  if err := env.Parse(db); err != nil {
    return nil,fmt.Errorf("failed to parse env: %w", err)
  }

  store := &StoreConfig{}
  if err := env.Parse(store); err != nil {
    return nil,fmt.Errorf("failed to parse env: %w", err)
  }

  redis := &RedisConfig{}
  if err := env.Parse(redis); err != nil {
    return nil,fmt.Errorf("failed to parse env: %w", err)
  }

  rabbitmq := &RabbitmqConfig{}
  if err := env.Parse(rabbitmq); err != nil {
    return nil,fmt.Errorf("failed to parse env: %w", err)
  }

  kafka := &KafkaConfig{}
  if err := env.Parse(kafka); err != nil {
    return nil,fmt.Errorf("failed to parse env: %w", err)
  }

  loki := &LokiConfig{}
  if err := env.Parse(loki); err != nil {
    return nil,fmt.Errorf("failed to parse env: %w", err)
  }

  sentry := &SentryConfig{}
  if err := env.Parse(sentry); err != nil {
    return nil,fmt.Errorf("failed to parse env: %w", err)
  }

  log := &LogConfig{}
  if err := env.Parse(log); err != nil {
    return nil,fmt.Errorf("failed to parse env: %w", err)
  }

  shopee := &ShopeeConfig{}
  if err := env.Parse(shopee); err != nil {
    return nil,fmt.Errorf("failed to parse env: %w", err)
  }

  // logger.Sugar().Infow("Env loaded successfully", "env", envSet)
  return &Config{
    Server: server,
    JWT: jwt,
    DB: db,
    Store: store,
    Redis: redis,
    Rabbitmq: rabbitmq,
    Kafka: kafka,
    Loki: loki,
    Sentry: sentry,
    Log: log,
    Shopee:shopee,
  }, nil 
}
