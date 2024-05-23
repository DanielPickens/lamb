package config

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/danielpickens/yeti/common/consts"
)

type yetiServerConfigYaml struct {
	EnableHTTPS          bool   `yaml:"enable_https"`
	Port                 uint   `yaml:"port"`
	SessionSecretKey     string `yaml:"session_secret_key"`
	MigrationDir         string `yaml:"migration_dir"`
	ReadHeaderTimeout    int    `yaml:"read_header_timeout"`
	TransmissionStrategy string `yaml:"transmission_strategy"`
}

type yetiPostgresqlConfigYaml struct {
	Host            string        `yaml:"host"`
	Port            uint          `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Database        string        `yaml:"database"`
	SSLMode         string        `yaml:"sslmode"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

type yetiS3ConfigYaml struct {
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	Region     string `yaml:"region"`
	Secure     bool   `yaml:"secure"`
	BucketName string `yaml:"bucket_name"`
}

type yetiDockerRegistryConfigYaml struct {
	mlRepositoryName string `yaml:"ml_repository_name"`
	ModelRepositoryName string `yaml:"model_repository_name"`
	Server              string `yaml:"server"`
	Username            string `yaml:"username"`
	Password            string `yaml:"password"`
	Secure              bool   `yaml:"secure"`
}

type yetiDockerImageBuilderConfigYaml struct {
	Privileged bool `yaml:"privileged"`
}

type yetiConfigYaml struct {
	IsSaaS              bool                      `yaml:"is_saas"`
	SaasDomainSuffix    string                    `yaml:"saas_domain_suffix"`
	InCluster           bool                      `yaml:"in_cluster"`
	Server              yetiServerConfigYaml     `yaml:"server"`
	Postgresql          yetiPostgresqlConfigYaml `yaml:"postgresql"`
	S3                  *yetiS3ConfigYaml        `yaml:"s3,omitempty"`
	NewsURL             string                    `yaml:"news_url"`
	InitializationToken string                    `yaml:"initialization_token"`
}

var yetiConfig = &yetiConfigYaml{}

func PopulateyetiConfig() error {
	isSaaS, ok := os.LookupEnv(consts.EnvIsSaaS)
	if ok {
		yetiConfig.IsSaaS = isSaaS == "true"
	}

	saasDomainSuffix, ok := os.LookupEnv(consts.EnvSaasDomainSuffix)
	if ok {
		yetiConfig.SaasDomainSuffix = saasDomainSuffix
	}

	pgHost, ok := os.LookupEnv(consts.EnvPgHost)
	if ok {
		yetiConfig.Postgresql.Host = pgHost
	}
	pgPort, ok := os.LookupEnv(consts.EnvPgPort)
	if ok {
		pgPort_, err := strconv.Atoi(pgPort)
		if err != nil {
			return errors.Wrap(err, "convert port from env to int")
		}
		yetiConfig.Postgresql.Port = uint(pgPort_)
	}
	pgUser, ok := os.LookupEnv(consts.EnvPgUser)
	if ok {
		yetiConfig.Postgresql.User = pgUser
	}
	pgPassword, ok := os.LookupEnv(consts.EnvPgPassword)
	if ok {
		yetiConfig.Postgresql.Password = pgPassword
	}
	pgDatabase, ok := os.LookupEnv(consts.EnvPgDatabase)
	if ok {
		yetiConfig.Postgresql.Database = pgDatabase
	}
	pgSSLMode, ok := os.LookupEnv(consts.EnvPgSSLMode)
	if ok {
		yetiConfig.Postgresql.SSLMode = pgSSLMode
	}
	yetiConfig.Postgresql.MaxOpenConns = 10
	pgMaxOpenConns, ok := os.LookupEnv("PG_MAX_OPEN_CONNS")
	if ok {
		pgMaxOpenConns_, err := strconv.Atoi(pgMaxOpenConns)
		if err != nil {
			return errors.Wrap(err, "convert PG_MAX_OPEN_CONNS from env to int")
		}
		yetiConfig.Postgresql.MaxOpenConns = pgMaxOpenConns_
	}
	yetiConfig.Postgresql.MaxIdleConns = 10
	pgMaxIdleConns, ok := os.LookupEnv("PG_MAX_IDLE_CONNS")
	if ok {
		pgMaxIdleConns_, err := strconv.Atoi(pgMaxIdleConns)
		if err != nil {
			return errors.Wrap(err, "convert PG_MAX_IDLE_CONNS from env to int")
		}
		yetiConfig.Postgresql.MaxIdleConns = pgMaxIdleConns_
	}
	yetiConfig.Postgresql.ConnMaxLifetime = 15 * time.Minute
	pgConnMaxLifetime, ok := os.LookupEnv("PG_CONN_MAX_LIFETIME")
	if ok {
		pgConnMaxLifetime_, err := time.ParseDuration(pgConnMaxLifetime)
		if err != nil {
			return errors.Wrap(err, "convert PG_CONN_MAX_LIFETIME from env to time.Duration")
		}
		yetiConfig.Postgresql.ConnMaxLifetime = pgConnMaxLifetime_
	}
	migrationDir, ok := os.LookupEnv(consts.EnvMigrationDir)
	if ok {
		yetiConfig.Server.MigrationDir = migrationDir
	}
	sessionSecretKey, ok := os.LookupEnv(consts.EnvSessionSecretKey)
	if ok {
		yetiConfig.Server.SessionSecretKey = sessionSecretKey
	}
	if yetiConfig.Server.Port == 0 {
		yetiConfig.Server.Port = 7777
	}

	readHeaderTimeout, ok := os.LookupEnv(consts.EnvReadHeaderTimeout)
	if ok {
		readHeaderTimeout_, err := strconv.Atoi(readHeaderTimeout)
		if err != nil {
			return errors.Wrapf(err, "convert %s from env to int", consts.EnvReadHeaderTimeout)
		}
		yetiConfig.Server.ReadHeaderTimeout = readHeaderTimeout_
	}

	transmissionStrategy, ok := os.LookupEnv(consts.EnvTransmissionStrategy)
	if ok {
		yetiConfig.Server.TransmissionStrategy = transmissionStrategy
	}

	initializationToken, ok := os.LookupEnv(consts.EnvInitializationToken)
	if ok {
		yetiConfig.InitializationToken = initializationToken
	}
	makesureS3IsNotNil := func() {
		if yetiConfig.S3 == nil {
			yetiConfig.S3 = &yetiS3ConfigYaml{}
		}
	}
	s3Endpoint, ok := os.LookupEnv(consts.EnvS3Endpoint)
	if ok {
		makesureS3IsNotNil()
		yetiConfig.S3.Endpoint = s3Endpoint
	}
	s3AccessKey, ok := os.LookupEnv(consts.EnvS3AccessKey)
	if ok {
		makesureS3IsNotNil()
		yetiConfig.S3.AccessKey = s3AccessKey
	}
	s3SecretKey, ok := os.LookupEnv(consts.EnvS3SecretKey)
	if ok {
		makesureS3IsNotNil()
		yetiConfig.S3.SecretKey = s3SecretKey
	}
	s3Region, ok := os.LookupEnv(consts.EnvS3Region)
	if ok {
		makesureS3IsNotNil()
		yetiConfig.S3.Region = s3Region
	}
	s3Secure, ok := os.LookupEnv(consts.EnvS3Secure)
	if ok {
		makesureS3IsNotNil()
		s3Secure_, err := strconv.ParseBool(s3Secure)
		if err != nil {
			return errors.Wrap(err, "convert s3_secure from env to bool")
		}
		yetiConfig.S3.Secure = s3Secure_
	}
	s3BucketName, ok := os.LookupEnv(consts.EnvS3BucketName)
	if ok {
		makesureS3IsNotNil()
		yetiConfig.S3.BucketName = s3BucketName
	}
	return nil
}
