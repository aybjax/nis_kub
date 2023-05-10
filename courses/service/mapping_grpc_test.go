package service

import (
	"context"
	"testing"

	"github.com/aybjax/nis_lib/pbdto"
)

func TestDecodeGetCoursesRequest(t *testing.T) {
	const id = "decodeGetCoursesRequestId"
	result, _ := decodeGetCoursesRequest(context.TODO(), &pbdto.Request{
		Id: id,
	})

	data := result.(string)

	if data != id {
		t.Errorf("TestDecodeGetCoursesRequest error")
	}
}
