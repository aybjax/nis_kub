package app_db

import (
	"errors"
	"github.com/aybjax/nis_lib/pbdto"
	"github.com/go-playground/validator/v10"
)

type StudentDB struct {
	Id               string `bson:"_id,omitempty" json:"id,omitempty"`
	StudentDBPayload `bson:",inline"`
}

type StudentDBPayload struct {
	Name      string   `json:"name,omitempty" validate:"required,min=3,max=20"`
	FirstName string   `json:"first_name,omitempty" validate:"required,min=3,max=20"`
	LastName  string   `json:"last_name,omitempty" validate:"required,min=3,max=20"`
	CourseIds []string `json:"course_ids,omitempty" validate:"max=15"`
}

func (StudentDB) FromProto(data *pbdto.Student) (*StudentDB, error) {
	if data == nil {
		return nil, errors.New("Argument nil")
	}

	return &StudentDB{
		Id: data.Id,
		StudentDBPayload: StudentDBPayload{
			Name:      data.Name,
			FirstName: data.FirstName,
			LastName:  data.LastName,
			CourseIds: data.CourseIds,
		},
	}, nil
}

func (s *StudentDB) GetCourseIds() []string {
	if s.CourseIds == nil {
		return make([]string, 0)
	}

	return s.CourseIds
}

func (s *StudentDB) ToProto() *pbdto.Student {
	return &pbdto.Student{
		Id:        s.Id,
		Name:      s.Name,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		CourseIds: s.CourseIds,
	}
}

func (c *StudentDB) HasValidPayload() error {
	validate := validator.New()

	err := validate.Struct(c.StudentDBPayload)

	if err != nil {
		return err
	}

	return nil
}

func (c *StudentDB) GetPayload() StudentDBPayload {
	return c.StudentDBPayload
}
