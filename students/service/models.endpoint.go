package service

import "github.com/aybjax/nis_lib/pbdto"

type idPayloadRequest struct {
	Id   string         `json:"id"`
	Data *pbdto.Student `json:"data"`
}
