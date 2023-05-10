package service

import (
	"context"
	"nis_students/app"

	"github.com/aybjax/nis_lib/pbdto"
	"github.com/go-kit/kit/endpoint"
)

func makeGetStudentsEndpoint(svc app.StudentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)

		students, err := svc.GetStudents(ctx, req)

		if err != nil {
			return nil, err
		}

		return &pbdto.StudentsResponse{
			Students: students,
		}, err
	}
}

func makeGetAllEndpoint(svc app.StudentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		students, err := svc.GetAll(ctx)

		return map[string]interface{}{
			"data": students,
		}, err
	}
}

func makeGetEndpoint(svc app.StudentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)

		student, err := svc.Get(ctx, req)

		return map[string]interface{}{
			"data": student,
		}, err
	}
}

func makeGetCoursesEndpoint(svc app.StudentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)

		courses, err := svc.GetCourses(ctx, req)

		return map[string]interface{}{
			"data": courses,
		}, err
	}
}

func makePostEndpoint(svc app.StudentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(idPayloadRequest)

		student, err := svc.Post(ctx, req.Data)

		return map[string]interface{}{
			"data": student,
		}, err
	}
}

func makePutEndpoint(svc app.StudentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(idPayloadRequest)

		student, err := svc.Put(ctx, req.Id, req.Data)

		return map[string]interface{}{
			"data": student,
		}, err
	}
}

func makeDeleteEndpoint(svc app.StudentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)

		student, err := svc.Delete(ctx, req)

		return map[string]interface{}{
			"data": student,
		}, err
	}
}
