package service

import (
	"context"
	"github.com/aybjax/nis_lib/pbdto"
	gt "github.com/go-kit/kit/transport/grpc"
	"nis_courses/app"
)

type GRPCHandlers struct {
	getCourses gt.Handler
}

func (g *GRPCHandlers) GetCourses(ctx context.Context, req *pbdto.Request) (*pbdto.CoursesResponse, error) {
	_, resp, err := g.getCourses.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pbdto.CoursesResponse), nil
}

func GetGRPCHandlers(svc app.CourseService) pbdto.GetStudentCoursesServer {
	var result GRPCHandlers

	result.getCourses = gt.NewServer(
		makeGetCoursesEndpoint(svc),
		decodeGetCoursesRequest,
		encodeGetCoursesResponse,
	)

	return &result
}
