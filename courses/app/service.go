package app

import (
	"context"

	"github.com/aybjax/nis_lib/pbdto"
)

type CourseService interface {
	GetAll(ctx context.Context) ([]*pbdto.Course, error)
	Get(ctx context.Context, id string) (*pbdto.Course, error)
	GetStudents(ctx context.Context, id string) ([]*pbdto.Student, error)
	GetCourses(ctx context.Context, id string) ([]*pbdto.Course, error)
	Post(ctx context.Context, payload *pbdto.Course) (id interface{}, err error)
	Put(ctx context.Context, id string, payload *pbdto.Course) (newId interface{}, err error)
	Delete(ctx context.Context, id string) (oldId interface{}, err error)
	StudentModifiedListener(ctx context.Context, updateInfo *pbdto.UpdateEmbedded) error
	CourseModifiedListener(ctx context.Context, diffIds *pbdto.DiffIds) error
}
