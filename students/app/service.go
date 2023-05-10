package app

import (
	"context"

	"github.com/aybjax/nis_lib/pbdto"
)

type StudentService interface {
	GetAll(ctx context.Context) ([]*pbdto.Student, error)
	Get(ctx context.Context, id string) (*pbdto.Student, error)
	GetCourses(ctx context.Context, id string) ([]*pbdto.Course, error)
	GetStudents(ctx context.Context, id string) ([]*pbdto.Student, error)
	Post(ctx context.Context, payload *pbdto.Student) (id interface{}, err error)
	Put(ctx context.Context, id string, payload *pbdto.Student) (newId interface{}, err error)
	Delete(ctx context.Context, id string) (oldId interface{}, err error)
	CourseModifiedListener(context.Context, *pbdto.UpdateEmbedded) error
	StudentModifiedListener(context.Context, *pbdto.DiffIds) error
}
