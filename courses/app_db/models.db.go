package app_db

import (
	"errors"

	"github.com/aybjax/nis_lib/pbdto"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Omit empty in bson field => empty values are not updated
type CourseDB struct {
	Id              primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CourseDBPayload `bson:",inline"`
}

type CourseDBPayload struct {
	Name        string   `bson:"name" json:"name,omitempty" validate:"required,min=3,max=20"`
	Description string   `bson:"description" json:"description,omitempty" validate:"required,min=3,max=200"`
	Discipline  string   `bson:"discipline" json:"discipline,omitempty" validate:"required,min=3,max=20"`
	Teacher     string   `bson:"teacher" json:"teacher,omitempty" validate:"required,min=3,max=20"`
	StudentIds  []string `bson:"student_ids" json:"student_ids,omitempty" validate:"max=60"`
}

func (CourseDB) FromProto(data *pbdto.Course) (*CourseDB, error) {
	if data == nil {
		return nil, errors.New("Argument nil")
	}

	objId, _ := primitive.ObjectIDFromHex(data.Id)
	studentIds := make([]string, 0)

	if len(data.StudentIds) > 0 {
		studentIds = data.StudentIds
	}

	return &CourseDB{
		Id: objId,
		CourseDBPayload: CourseDBPayload{
			Name:        data.Name,
			Description: data.Description,
			Discipline:  data.Discipline,
			Teacher:     data.Teacher,
			StudentIds:  studentIds,
		},
	}, nil
}

func (c *CourseDB) ToProto() *pbdto.Course {
	return &pbdto.Course{
		Id:          c.Id.Hex(),
		Name:        c.Name,
		Description: c.Description,
		Discipline:  c.Discipline,
		Teacher:     c.Teacher,
		StudentIds:  c.StudentIds,
	}
}

func (c *CourseDB) HasValidPayload() error {
	validate := validator.New()

	err := validate.Struct(c.CourseDBPayload)

	if err != nil {
		return err
	}

	return nil
}

func (c *CourseDB) GetPayload() CourseDBPayload {
	return c.CourseDBPayload
}
