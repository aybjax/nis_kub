package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"nis_courses/app"
	"nis_courses/app/adapter"
	"nis_courses/service"
	"time"

	"github.com/aybjax/nis_lib/pbdto"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type GRPCClientClose func() error

func NewDB() (*mongo.Client, error) {
	var client *mongo.Client
	{
		var err error
		// client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
		client, err = mongo.Connect(context.Background(), options.Client().
			ApplyURI(adapter.ENV.GetConnectionString()))
		if err != nil {
			return nil, err
		}

		err = client.Ping(context.Background(), nil)
		if err != nil {
			panic(err)
		}
	}

	return client, nil
}
func NewCollection(db *mongo.Client) *mongo.Collection {
	return db.Database("nis").Collection("courses")
}

func NewGRPCClient() (pbdto.GetCourseStudentsClient, GRPCClientClose, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s",
		adapter.ENV.GRPC_STUDENT_URL,
		adapter.ENV.GRPC_STUDENT_PORT),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		// panic(err)
	}

	c := pbdto.NewGetCourseStudentsClient(conn)

	return c, conn.Close, nil
}

func NewHTTPServer(svc app.CourseService) (*http.Server, error) {
	r := mux.NewRouter()
	r.Use(service.HttpMiddleware)
	hs := service.GetHTTPHandlers(svc)
	r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())
	r.Methods("GET").Path("/api/courses/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `{"message": "courses alive"}`)
	})
	r.Methods("GET").Path("/api/courses/{id}/students").Handler(hs.GetStudents)
	r.Methods("GET").Path("/api/courses/{id}").Handler(hs.Get)
	r.Methods("GET").Path("/api/courses").Handler(hs.GetAll)
	r.Methods("POST").Path("/api/courses").Handler(hs.Post)
	r.Methods("PUT").Path("/api/courses/{id}").Handler(hs.Put)
	r.Methods("DELETE").Path("/api/courses/{id}").Handler(hs.Delete)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%s", adapter.EXPOSE.HTTP_PORT),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv, nil
}

func NewGRPCServer(svc app.CourseService) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", adapter.EXPOSE.GRPC_PORT))
	if err != nil {
		return nil, nil, err
	}

	g := grpc.NewServer()
	handlers := service.GetGRPCHandlers(svc)
	pbdto.RegisterGetStudentCoursesServer(g, handlers)
	reflection.Register(g)

	return g, lis, nil
}
