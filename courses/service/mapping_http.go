package service

import (
	"context"
	"encoding/json"
	"github.com/aybjax/nis_lib/pbdto"
	"github.com/gorilla/mux"
	"net/http"
)

func decodeEmptyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	return vars["id"], nil
}

func decodeIdPayloadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var payload pbdto.Course

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)

	return idPayloadRequest{
		Id:   vars["id"],
		Data: &payload,
	}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
