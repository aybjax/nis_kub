package service

import (
	"context"
	"github.com/aybjax/nis_lib/pbdto"
	gt "github.com/go-kit/kit/transport/grpc"
	"nis_students/app"
)

type GRPCHandlers struct {
	getStudents gt.Handler
}

func (g *GRPCHandlers) GetStudents(ctx context.Context, req *pbdto.Request) (*pbdto.StudentsResponse, error) {
	_, resp, err := g.getStudents.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pbdto.StudentsResponse), nil
}

func GetGRPCHandlers(svc app.StudentService) pbdto.GetCourseStudentsServer {
	var result GRPCHandlers

	result.getStudents = gt.NewServer(
		makeGetStudentsEndpoint(svc),
		decodeGetStudentsRequest,
		encodeGetStudentsResponse,
	)

	return &result
}
