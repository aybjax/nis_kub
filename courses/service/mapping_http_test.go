package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aybjax/nis_lib/pbdto"
	"github.com/gorilla/mux"
)

func TestDecodeIdRequest(t *testing.T) {
	req, _ := http.NewRequest("", "", nil)
	result, _ := decodeIdRequest(context.TODO(), mux.SetURLVars(req, map[string]string{"id": "id"}))
	data := result.(string)

	if data != "id" {
		t.Errorf("TestDecodeIdRequest")
	}
}

// func decodeIdPayloadRequest(_ context.Context, r *http.Request) (interface{}, error) {
// 	payload := &pbdto.Student{}

// 	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
// 		return nil, err
// 	}

// 	vars := mux.Vars(r)

// 	return idPayloadRequest{
// 		Id:   vars["id"],
// 		Data: payload,
// 	}, nil
// }

func TestDecodeIdPayloadRequest(t *testing.T) {
	payload := pbdto.Student{
		Id: "payloadId",
	}
	payloadData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("", "", bytes.NewReader(payloadData))
	result, _ := decodeIdPayloadRequest(context.TODO(), mux.SetURLVars(req, map[string]string{"id": "id"}))
	data := result.(idPayloadRequest)

	if data.Id != "id" && data.Data.Id != "payloadId" {
		t.Errorf("TestDecodeIdRequest")
	}
}
