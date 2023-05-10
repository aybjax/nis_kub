package service

import (
	"context"
	"testing"

	"github.com/aybjax/nis_lib/pbdto"
)

type testService struct{}

func (t testService) GetAll(ctx context.Context) ([]*pbdto.Student, error) {
	return []*pbdto.Student{
		&pbdto.Student{
			Id: "GetAllId",
		},
	}, nil
}
func (t testService) Get(ctx context.Context, id string) (*pbdto.Student, error) {
	return &pbdto.Student{
		Id: "GetId",
	}, nil
}
func (t testService) GetCourses(ctx context.Context, id string) ([]*pbdto.Course, error) {
	return []*pbdto.Course{
		{
			Id: "GetCoursesId",
		},
	}, nil
}
func (t testService) GetStudents(ctx context.Context, id string) ([]*pbdto.Student, error) {
	return []*pbdto.Student{
		&pbdto.Student{
			Id: "GetStudentsId",
		},
	}, nil
}
func (t testService) Post(ctx context.Context, payload *pbdto.Student) (id interface{}, err error) {
	return "PostId", nil
}
func (t testService) Put(ctx context.Context, id string, payload *pbdto.Student) (newId interface{}, err error) {
	return "PutId", nil
}
func (t testService) Delete(ctx context.Context, id string) (oldId interface{}, err error) {
	return "DeleteId", nil
}
func (t testService) CourseModifiedListener(ctx context.Context, payload *pbdto.UpdateEmbedded) error {
	payload.Id = "CourseModifiedListenerId"

	return nil
}
func (t testService) StudentModifiedListener(ctx context.Context, payload *pbdto.DiffIds) error {
	payload.Id = "StudentModifiedListenerId"

	return nil
}

func TestMakeGetAllEndpoint(t *testing.T) {
	handler := makeGetAllEndpoint(testService{})

	resp, _ := handler(context.TODO(), nil)

	result := resp.(map[string]interface{})

	data := result["data"].([]*pbdto.Student)

	if len(data) != 1 || data[0].Id != "GetAllId" {
		t.Errorf("TestMakeGetAllEndpoint")
	}
}

func TestMakeGetEndpoint(t *testing.T) {
	handler := makeGetEndpoint(testService{})

	resp, _ := handler(context.TODO(), "request should be string")

	result := resp.(map[string]interface{})

	data := result["data"].(*pbdto.Student)

	if data.Id != "GetId" {
		t.Errorf("TestMakeGetEndpoint")
	}
}

func TestMakeGetCoursesEndpoint(t *testing.T) {
	handler := makeGetCoursesEndpoint(testService{})

	resp, _ := handler(context.TODO(), "Arguemnt should be string")

	result := resp.(map[string]interface{})

	data := result["data"].([]*pbdto.Course)

	if len(data) != 1 || data[0].Id != "GetCoursesId" {
		t.Errorf("TestMakeGetCoursesEndpoint")
	}
}

func TestMakeGetStudentsEndpoint(t *testing.T) {
	handler := makeGetStudentsEndpoint(testService{})

	resp, _ := handler(context.TODO(), "Argument should be string")

	result := resp.(*pbdto.StudentsResponse)

	// data := result["data"].([]*pbdto.Student)
	data := result.Students

	if len(data) != 1 || data[0].Id != "GetStudentsId" {
		t.Errorf("TestMakeGetStudentsEndpoint")
	}
}

func TestMakePostEndpoint(t *testing.T) {
	handler := makePostEndpoint(testService{})

	resp, _ := handler(context.TODO(), idPayloadRequest{})

	result := resp.(map[string]interface{})
	data := result["data"].(string)

	if data != "PostId" {
		t.Errorf("TestMakePostEndpoint")
	}
}

func TestMakePutEndpoint(t *testing.T) {
	handler := makePutEndpoint(testService{})

	resp, _ := handler(context.TODO(), idPayloadRequest{})

	result := resp.(map[string]interface{})
	data := result["data"].(string)

	if data != "PutId" {
		t.Errorf("TestMakePutEndpoint")
	}
}

func TestMakeDeleteEndpoint(t *testing.T) {
	handler := makeDeleteEndpoint(testService{})

	resp, _ := handler(context.TODO(), "Should be string")

	result := resp.(map[string]interface{})
	data := result["data"].(string)

	if data != "DeleteId" {
		t.Errorf("TestMakePutEndpoint")
	}
}
