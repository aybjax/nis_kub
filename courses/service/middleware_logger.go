package service

import (
	"nis_courses/app"
	"time"

	"github.com/aybjax/nis_lib/pbdto"
	"github.com/go-kit/log"
	"golang.org/x/net/context"
)

type LoggerMiddleware struct {
	logger     log.Logger
	serverType string
	next       app.CourseService
}

func NewLoggerMiddleware(next app.CourseService, logger log.Logger, _type string) app.CourseService {
	return &LoggerMiddleware{
		logger:     logger,
		serverType: _type,
		next:       next,
	}
}

func (lw *LoggerMiddleware) GetAll(ctx context.Context) (output []*pbdto.Course, err error) {
	defer func(begin time.Time) {
		_ = lw.logger.Log(
			"server type", lw.serverType,
			"method", "GetAll",
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = lw.next.GetAll(ctx)

	return
}
func (lw *LoggerMiddleware) Get(ctx context.Context, id string) (output *pbdto.Course, err error) {
	defer func(begin time.Time) {
		_ = lw.logger.Log(
			"server type", lw.serverType,
			"method", "Get",
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = lw.next.Get(ctx, id)

	return
}
func (lw *LoggerMiddleware) GetStudents(ctx context.Context, id string) (output []*pbdto.Student, err error) {
	defer func(begin time.Time) {
		_ = lw.logger.Log(
			"server type", lw.serverType,
			"method", "GetStudents",
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = lw.next.GetStudents(ctx, id)

	return
}
func (lw *LoggerMiddleware) GetCourses(ctx context.Context, id string) (output []*pbdto.Course, err error) {
	defer func(begin time.Time) {
		_ = lw.logger.Log(
			"server type", lw.serverType,
			"method", "GetCourses",
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = lw.next.GetCourses(ctx, id)

	return
}
func (lw *LoggerMiddleware) Post(ctx context.Context, payload *pbdto.Course) (id interface{}, err error) {
	defer func(begin time.Time) {
		_ = lw.logger.Log(
			"server type", lw.serverType,
			"method", "Post",
			"took", time.Since(begin),
		)
	}(time.Now())

	id, err = lw.next.Post(ctx, payload)

	return
}
func (lw *LoggerMiddleware) Put(ctx context.Context, id string, payload *pbdto.Course) (newId interface{}, err error) {
	defer func(begin time.Time) {
		_ = lw.logger.Log(
			"server type", lw.serverType,
			"method", "Put",
			"took", time.Since(begin),
		)
	}(time.Now())

	newId, err = lw.next.Put(ctx, id, payload)

	return
}
func (lw *LoggerMiddleware) Delete(ctx context.Context, id string) (oldId interface{}, err error) {
	defer func(begin time.Time) {
		_ = lw.logger.Log(
			"server type", lw.serverType,
			"method", "Delete",
			"took", time.Since(begin),
		)
	}(time.Now())

	oldId, err = lw.next.Delete(ctx, id)

	return
}

func (lw *LoggerMiddleware) StudentModifiedListener(ctx context.Context, updateInfo *pbdto.UpdateEmbedded) (err error) {
	defer func(begin time.Time) {
		_ = lw.logger.Log(
			"server type", lw.serverType,
			"method", "StudentModifiedListener",
			"took", time.Since(begin),
		)
	}(time.Now())

	err = lw.next.StudentModifiedListener(ctx, updateInfo)

	return
}

func (lw *LoggerMiddleware) CourseModifiedListener(ctx context.Context, diffIds *pbdto.DiffIds) (err error) {
	defer func(begin time.Time) {
		_ = lw.logger.Log(
			"server type", lw.serverType,
			"method", "StudentModifiedListener",
			"took", time.Since(begin),
		)
	}(time.Now())

	err = lw.next.CourseModifiedListener(ctx, diffIds)

	return
}
