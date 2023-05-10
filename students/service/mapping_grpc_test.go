package service

import (
	"context"
	"testing"

	"github.com/aybjax/nis_lib/pbdto"
)

func TestDecodeGetStudentsRequest(t *testing.T) {
	const id = "decodeGetStudentsRequestId"
	result, _ := decodeGetStudentsRequest(context.TODO(), &pbdto.Request{
		Id: id,
	})

	data := result.(string)

	if data != id {
		t.Errorf("TestDecodeGetStudentsRequest error")
	}
}
