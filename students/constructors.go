package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"nis_students/app"
	"nis_students/app/adapter"
	"nis_students/app_db"
	"nis_students/service"
	"time"

	"github.com/aybjax/nis_lib/pbdto"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type GRPCClientClose func() error

func NewDB(logger log.Logger) (*app_db.DB, error) {
	var db *sql.DB
	{
		var err error
		logger.Log("postgres", adapter.ENV.GetPostgresConnectionString())
		db, err = sql.Open("postgres", adapter.ENV.GetPostgresConnectionString())
		if err != nil {
			return nil, err
		}

		if err = db.Ping(); err != nil {
			return nil, err
		}
	}

	app_db.Migrate(db)

	return app_db.NewAppDB(db, logger), nil
}

func NewGRPCClient() (pbdto.GetStudentCoursesClient, GRPCClientClose, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s",
		adapter.ENV.GRPC_COURSE_URL,
		adapter.ENV.GRPC_COURSE_PORT),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pbdto.NewGetStudentCoursesClient(conn)

	return c, conn.Close, nil
}

func NewHTTPServer(svc app.StudentService) (*http.Server, error) {
	r := mux.NewRouter()
	r.Use(service.HttpMiddleware)
	hs := service.GetHTTPHandlers(svc)
	r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())
	r.Methods("GET").Path("/api/students/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `{"message": "students alive"}`)
	})
	r.Methods("GET").Path("/api/students/{id}/courses").Handler(hs.GetCourses)
	r.Methods("GET").Path("/api/students/{id}").Handler(hs.Get)
	r.Methods("GET").Path("/api/students").Handler(hs.GetAll)
	r.Methods("POST").Path("/api/students").Handler(hs.Post)
	r.Methods("PUT").Path("/api/students/{id}").Handler(hs.Put)
	r.Methods("DELETE").Path("/api/students/{id}").Handler(hs.Delete)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%s", adapter.EXPOSE.HTTP_PORT),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv, nil
}

func NewGRPCServer(svc app.StudentService) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", adapter.EXPOSE.GRPC_PORT))
	if err != nil {
		return nil, nil, err
	}

	g := grpc.NewServer()
	handlers := service.GetGRPCHandlers(svc)
	pbdto.RegisterGetCourseStudentsServer(g, handlers)
	reflection.Register(g)

	return g, lis, nil
}
