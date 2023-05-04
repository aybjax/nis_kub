package adapter

import (
	"fmt"
	"os"
)

var ENV *env
var EXPOSE *expose

func init() {
	ENV = &env{
		REDIS_URL:         os.Getenv("REDIS_URL"),
		REDIS_PORT:        os.Getenv("REDIS_PORT"),
		NATS_URL:          os.Getenv("NATS_URL"),
		NATS_PORT:         os.Getenv("NATS_PORT"),
		MONGO_USER:        os.Getenv("MONGO_USER"),
		MONGO_PASSWORD:    os.Getenv("MONGO_PASSWORD"),
		MONGO_URL:         os.Getenv("MONGO_URL"),
		MONGO_PORT:        os.Getenv("MONGO_PORT"),
		GRPC_STUDENT_URL:  os.Getenv("GRPC_STUDENT_URL"),
		GRPC_STUDENT_PORT: os.Getenv("GRPC_STUDENT_PORT"),
	}

	EXPOSE = &expose{
		HTTP_PORT: "8000",
		GRPC_PORT: "50051",
	}
}

type env struct {
	REDIS_URL         string
	REDIS_PORT        string
	NATS_URL          string
	NATS_PORT         string
	GRPC_STUDENT_URL  string
	GRPC_STUDENT_PORT string
	MONGO_USER        string
	MONGO_PASSWORD    string
	MONGO_URL         string
	MONGO_PORT        string
}

func (e *env) GetConnectionString() string {
	if e.MONGO_USER == "" && e.MONGO_PASSWORD == "" {
		return fmt.Sprintf("mongodb://%s:%s", e.MONGO_URL, e.MONGO_PORT)
	}

	return fmt.Sprintf("mongodb://%s:%s@%s:%s",
		e.MONGO_USER,
		e.MONGO_PASSWORD,
		e.MONGO_URL,
		e.MONGO_PORT,
	)
}

type expose struct {
	HTTP_PORT string
	GRPC_PORT string
}
