package adapter

import (
	"fmt"
	"os"
)

var ENV *env
var EXPOSE *expose

func init() {
	ENV = &env{
		NATS_URL:          os.Getenv("NATS_URL"),
		NATS_PORT:         os.Getenv("NATS_PORT"),
		REDIS_URL:         os.Getenv("REDIS_URL"),
		REDIS_PORT:        os.Getenv("REDIS_PORT"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_URL:      os.Getenv("POSTGRES_URL"),
		POSTGRES_PORT:     os.Getenv("POSTGRES_PORT"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
		GRPC_COURSE_URL:   os.Getenv("GRPC_COURSE_URL"),
		GRPC_COURSE_PORT:  os.Getenv("GRPC_COURSE_PORT"),
	}

	EXPOSE = &expose{
		HTTP_PORT: "8000",
		GRPC_PORT: "50051",
	}
}

type env struct {
	NATS_URL          string
	NATS_PORT         string
	REDIS_URL         string
	REDIS_PORT        string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_URL      string
	POSTGRES_PORT     string
	POSTGRES_DB       string
	GRPC_COURSE_URL   string
	GRPC_COURSE_PORT  string
}

func (e *env) GetPostgresConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		e.POSTGRES_USER,
		e.POSTGRES_PASSWORD,
		e.POSTGRES_URL,
		e.POSTGRES_PORT,
		e.POSTGRES_DB,
	)
}

type expose struct {
	HTTP_PORT string
	GRPC_PORT string
}
