package service

import (
	"context"
	"github.com/aybjax/nis_lib/pbdto"
)

func decodeGetStudentsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pbdto.Request)

	// return *req, nil

	return req.Id, nil
}

func encodeGetStudentsResponse(_ context.Context, response interface{}) (interface{}, error) {
	// resp := response.(app.Course)

	// return &resp, nil
	return response, nil
}
