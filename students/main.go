package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"nis_students/app"
	"nis_students/app/adapter"
	"nis_students/service"
	"os"
	"os/signal"
	"syscall"

	kitlog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	_ "github.com/lib/pq"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

func main() {
	logger := kitlog.NewLogfmtLogger(os.Stderr)
	cacheEngine, err := adapter.NewCacheEngine()
	if err != nil {
		panic(err)
	}
	logger.Log(
		"msg", "Cache engine created",
	)
	cache := app.NewCache(cacheEngine)
	logger.Log(
		"msg", "Cache created",
	)
	queueEngine, err := adapter.NewQueueEngine()
	if err != nil {
		logger.Log("queueEngine", "panic")
		panic(err)
	}
	logger.Log(
		"msg", "Queue engine created",
	)
	queue := app.NewQueue(queueEngine)
	logger.Log(
		"msg", "Queue created",
	)
	db, err := NewDB(logger)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	logger.Log(
		"msg", "App db created",
	)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	grpcClient, closeConn, err := NewGRPCClient()
	if err != nil {
		panic(err)
	}
	logger.Log(
		"msg", "Grpc client created",
	)
	defer closeConn()
	field_keys := []string{"method", "error"}

	requestCountHttp := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "students",
		Subsystem: "http",
		Name:      "request_count",
		Help:      "Number of requests received",
	}, field_keys)
	requestCountGrpc := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "students",
		Subsystem: "grpc",
		Name:      "request_count",
		Help:      "Number of requests received",
	}, field_keys)

	requestLatencyHttp := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "students",
		Subsystem: "http",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds",
	}, field_keys)
	requestLatencyGrpc := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "students",
		Subsystem: "grpc",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds",
	}, field_keys)

	svc := service.NewService(db, grpcClient, queue, logger)

	svc = service.NewCacheMiddleware(svc, cache, db, logger)

	httpSvc := service.NewLoggerMiddleware(svc, logger, "HTTP")
	grpcSvc := service.NewLoggerMiddleware(svc, logger, "GRPC")
	queueSvc := service.NewLoggerMiddleware(svc, logger, "Queue")

	httpSvc = service.NewMetricsMiddleware(httpSvc, requestCountHttp, requestLatencyHttp)
	grpcSvc = service.NewMetricsMiddleware(grpcSvc, requestCountGrpc, requestLatencyGrpc)
	queueSvc = service.NewMetricsMiddleware(queueSvc, requestCountGrpc, requestLatencyGrpc)

	queue.CourseChangedListener(queueSvc.CourseModifiedListener)
	queue.TransformNotifyListener(queueSvc.StudentModifiedListener)
	logger.Log("queue listener", "bootstrapped")

	httpServer, err := NewHTTPServer(httpSvc)
	if err != nil {
		logger.Log("httpServer", "panic")
		panic(err)
	}
	logger.Log("httpServer", "bootstrapped")
	grpcServer, listener, err := NewGRPCServer(grpcSvc)
	if err != nil {
		logger.Log("grpcServer", "panic")
		panic(err)
	}
	logger.Log("grpcServer", "bootstrapped")

	go func(httpServer *http.Server) {
		logger.Log("httpserver", fmt.Sprintf("Serving server at: %s\n", adapter.EXPOSE.HTTP_PORT))
		fmt.Printf("Serving server at: %s\n", adapter.EXPOSE.HTTP_PORT)
		errs <- httpServer.ListenAndServe()
	}(httpServer)

	go func(grpcServer *grpc.Server, listener net.Listener) {
		logger.Log("grpcserver", fmt.Sprintf("Serving server at: %s\n", fmt.Sprintf(":%s", adapter.EXPOSE.GRPC_PORT)))
		fmt.Printf("Serving server at: %s\n", fmt.Sprintf(":%s", adapter.EXPOSE.GRPC_PORT))
		if err := grpcServer.Serve(listener); err != nil {
			errs <- fmt.Errorf("gRPC failed to serve: %v\n", err)
		}
	}(grpcServer, listener)

	logger.Log("main", "waiting")
	fmt.Println("Waiting to exit")
	log.Fatalf("exited: %s\n", <-errs)
	fmt.Println("Exiting")
}
