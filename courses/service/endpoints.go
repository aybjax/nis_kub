package service

import (
	"context"
	"nis_courses/app"

	"github.com/aybjax/nis_lib/pbdto"
	"github.com/go-kit/kit/endpoint"
)

func makeGetAllEndpoint(svc app.CourseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		courses, err := svc.GetAll(ctx)

		return map[string]interface{}{
			"data": courses,
		}, err
	}
}

func makeGetEndpoint(svc app.CourseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)

		course, err := svc.Get(ctx, req)

		return map[string]interface{}{
			"data": course,
		}, err
	}
}

func makeGetStudentsEndpoint(svc app.CourseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)

		students, err := svc.GetStudents(ctx, req)

		return map[string]interface{}{
			"data": students,
		}, err
	}
}

func makeGetCoursesEndpoint(svc app.CourseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)

		courses, err := svc.GetCourses(ctx, req)

		if err != nil {
			return nil, err
		}

		return &pbdto.CoursesResponse{
			Courses: courses,
		}, err
	}
}

func makePostEndpoint(svc app.CourseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(idPayloadRequest)

		course, err := svc.Post(ctx, req.Data)

		return map[string]interface{}{
			"data": course,
		}, err
	}
}

func makePutEndpoint(svc app.CourseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(idPayloadRequest)

		course, err := svc.Put(ctx, req.Id, req.Data)

		return map[string]interface{}{
			"data": course,
		}, err
	}
}

func makeDeleteEndpoint(svc app.CourseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)

		id, err := svc.Delete(ctx, req)

		return map[string]interface{}{
			"data": id,
		}, err
	}
}
