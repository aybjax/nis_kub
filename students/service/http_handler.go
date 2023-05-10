package service

import (
	"net/http"
	"nis_students/app"

	httptransport "github.com/go-kit/kit/transport/http"
)

type HTTPHandlers struct {
	GetAll     http.Handler
	Get        http.Handler
	GetCourses http.Handler
	Post       http.Handler
	Put        http.Handler
	Delete     http.Handler
}

func HttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func GetHTTPHandlers(svc app.StudentService) HTTPHandlers {
	var result HTTPHandlers
	result.GetAll = httptransport.NewServer(
		makeGetAllEndpoint(svc),
		decodeEmptyRequest,
		encodeResponse,
	)
	result.Get = httptransport.NewServer(
		makeGetEndpoint(svc),
		decodeIdRequest,
		encodeResponse,
	)
	result.GetCourses = httptransport.NewServer(
		makeGetCoursesEndpoint(svc),
		decodeIdRequest,
		encodeResponse,
	)
	result.Post = httptransport.NewServer(
		makePostEndpoint(svc),
		decodeIdPayloadRequest,
		encodeResponse,
	)
	result.Put = httptransport.NewServer(
		makePutEndpoint(svc),
		decodeIdPayloadRequest,
		encodeResponse,
	)
	result.Delete = httptransport.NewServer(
		makeDeleteEndpoint(svc),
		decodeIdRequest,
		encodeResponse,
	)

	return result
}
