package env

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultPostBodyLimit   = 280
	defaultShutdownTimeout = time.Minute
)

type Env struct {
	Server     Server
	Database   Database
	JWT        JWT
	Pagination Pagination
	Posts      Posts
}

type Server struct {
	Addr            string
	CorsOrigin      string
	ShutdownTimeout time.Duration
}

type Database struct {
	DSN string
}

type JWT struct {
	Secret   string
	Lifetime time.Duration
}

type Pagination struct {
	DefaultPostsPageSize  int
	DefaultFollowPageSize int
}

type Posts struct {
	BodyLimit int
}

func Load() (*Env, error) {
	addr, err := lookupEnv("ADDR")
	if err != nil {
		return nil, err
	}

	corsOrigin, err := lookupEnv("CORS_ORIGIN")
	if err != nil {
		return nil, err
	}

	shutdownTimeoutString, set := os.LookupEnv("SHUTDOWN_TIMEOUT")
	shutdownTimeout := defaultShutdownTimeout
	if set {
		var err error
		shutdownTimeout, err = time.ParseDuration(shutdownTimeoutString)
		if err != nil {
			shutdownTimeout = defaultShutdownTimeout
		}
	}

	dsn, err := lookupEnv("DSN")
	if err != nil {
		return nil, err
	}

	jwtSecret, err := lookupEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	jwtLifetimeString, err := lookupEnv("JWT_LIFETIME")
	if err != nil {
		return nil, err
	}
	jwtLifetime, err := time.ParseDuration(jwtLifetimeString)
	if err != nil {
		return nil, err
	}

	defaultPostsPageSizeString, err := lookupEnv("DEFAULT_POSTS_PAGE_SIZE")
	if err != nil {
		return nil, err
	}
	defaultPostsPageSize, err := strconv.Atoi(defaultPostsPageSizeString)
	if err != nil {
		return nil, err
	}

	defaultFollowsPageSizeString, err := lookupEnv("DEFAULT_FOLLOWS_PAGE_SIZE")
	if err != nil {
		return nil, err
	}
	defaultFollowsPageSize, err := strconv.Atoi(defaultFollowsPageSizeString)
	if err != nil {
		return nil, err
	}

	postsBodyLimitString, set := os.LookupEnv("POSTS_BODY_LIMIT")
	postsBodyLimit := defaultPostBodyLimit
	if set {
		var err error
		postsBodyLimit, err = strconv.Atoi(postsBodyLimitString)
		if err != nil {
			return nil, err
		}
		if postsBodyLimit <= 0 {
			postsBodyLimit = defaultFollowsPageSize
		}
	}

	return &Env{
		Server: Server{
			Addr:            addr,
			CorsOrigin:      corsOrigin,
			ShutdownTimeout: shutdownTimeout,
		},
		Database: Database{
			DSN: dsn,
		},
		JWT: JWT{
			Secret:   jwtSecret,
			Lifetime: jwtLifetime,
		},
		Pagination: Pagination{
			DefaultPostsPageSize:  defaultPostsPageSize,
			DefaultFollowPageSize: defaultFollowsPageSize,
		},
		Posts: Posts{
			BodyLimit: postsBodyLimit,
		},
	}, nil
}

func lookupEnv(key string) (value string, err error) {
	value, set := os.LookupEnv(key)
	if !set {
		return "", errEnvNotSet(key)
	}
	return value, nil
}

func errEnvNotSet(name string) error {
	return fmt.Errorf("%s environment variable is not set", name)
}
