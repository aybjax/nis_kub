package service

import (
	"context"
	"fmt"
	"nis_students/app"
	"nis_students/app_db"

	"github.com/aybjax/nis_lib/pbdto"
	"github.com/go-kit/log"
)

type CacheMiddleware struct {
	cache  *app.Cache
	db     *app_db.DB
	next   app.StudentService
	logger log.Logger
}

func NewCacheMiddleware(next app.StudentService, cache *app.Cache, db *app_db.DB, logger log.Logger) app.StudentService {
	return &CacheMiddleware{
		cache:  cache,
		db:     db,
		next:   next,
		logger: logger,
	}
}

func (cw *CacheMiddleware) GetAll(ctx context.Context) (output []*pbdto.Student, err error) {
	if res, err := cw.cache.RetriveAll(); err == nil && len(res) > 0 {
		_ = cw.logger.Log(
			"method", "GetAll",
			"serving", "cache",
		)
		return res, nil
	} else {
		_ = cw.logger.Log(
			"method", "GetAll",
			"msg", "retrieve error",
			"err", fmt.Sprint(err),
		)
	}

	if output, err = cw.next.GetAll(ctx); err == nil {
		if err := cw.cache.WriteAll(output); err != nil {
			_ = cw.logger.Log(
				"method", "GetAll",
				"msg", "write error",
				"err", fmt.Sprint(err),
			)
		}
	}

	return
}
func (cw *CacheMiddleware) Get(ctx context.Context, id string) (output *pbdto.Student, err error) {
	if res, err := cw.cache.RetriveOneById(id); err == nil && res != nil {
		_ = cw.logger.Log(
			"method", "Get",
			"serving", "cache",
		)
		return res, nil
	} else {
		_ = cw.logger.Log(
			"method", "Get",
			"msg", "retrieve error",
			"err", fmt.Sprint(err),
		)
	}

	if output, err = cw.next.Get(ctx, id); err == nil {
		if err := cw.cache.WriteOneById(id, output); err != nil {
			_ = cw.logger.Log(
				"method", "Get",
				"msg", "write error",
				"err", fmt.Sprint(err),
			)
		}
	}

	return
}
func (cw *CacheMiddleware) GetCourses(ctx context.Context, id string) (output []*pbdto.Course, err error) {
	output, err = cw.next.GetCourses(ctx, id)

	return
}

func (cw *CacheMiddleware) GetStudents(ctx context.Context, courseId string) (output []*pbdto.Student, err error) {
	if res, err := cw.cache.RetrieveByCourseId(courseId); err == nil && len(res) > 0 {
		_ = cw.logger.Log(
			"method", "GetStudents",
			"serving", "cache",
		)
		return res, nil
	} else {
		_ = cw.logger.Log(
			"method", "GetStudents",
			"msg", "retrieve error",
			"err", fmt.Sprint(err),
		)
	}

	if output, err = cw.next.GetStudents(ctx, courseId); err == nil {
		if err := cw.cache.WriteByCourseId(courseId, output); err != nil {
			_ = cw.logger.Log(
				"method", "GetStudents",
				"msg", "write error",
				"err", fmt.Sprint(err),
			)
		}
	}

	return
}
func (cw *CacheMiddleware) Post(ctx context.Context, payload *pbdto.Student) (id interface{}, err error) {
	if id, err = cw.next.Post(ctx, payload); err == nil && id != nil {
		if err := cw.cache.InvalidateCreated(); err != nil {
			_ = cw.logger.Log(
				"method", "Post",
				"serving", "invalidation",
			)
		} else {
			_ = cw.logger.Log(
				"method", "Post",
				"msg", "invalidation error",
				"err", fmt.Sprint(err),
			)
		}
	}

	return
}
func (cw *CacheMiddleware) Put(ctx context.Context, id string, payload *pbdto.Student) (newId interface{}, err error) {
	oldIds, oldIdErr := cw.db.GetCourseIds(id)

	if newId, err = cw.next.Put(ctx, id, payload); err == nil && oldIdErr == nil {
		if err := cw.cache.InvalidateUpdated(id, payload.CourseIds, oldIds); err != nil {
			_ = cw.logger.Log(
				"method", "Put",
				"serving", "invalidation",
			)
		} else {
			_ = cw.logger.Log(
				"method", "Put",
				"msg", "invalidation error",
				"err", fmt.Sprint(err),
			)
		}
	}

	return
}
func (cw *CacheMiddleware) Delete(ctx context.Context, id string) (oldId interface{}, err error) {
	oldIds, oldIdErr := cw.db.GetCourseIds(id)

	if oldId, err = cw.next.Delete(ctx, id); err == nil && oldIdErr == nil && oldId != nil {
		if err := cw.cache.InvalidateDeleted(id, oldIds); err != nil {
			_ = cw.logger.Log(
				"method", "Delete",
				"serving", "invalidation",
			)
		} else {
			_ = cw.logger.Log(
				"method", "Delete",
				"msg", "invalidation error",
				"err", fmt.Sprint(err),
			)
		}
	}

	return
}

func (cw *CacheMiddleware) CourseModifiedListener(ctx context.Context, updateInfo *pbdto.UpdateEmbedded) (err error) {
	if err = cw.next.CourseModifiedListener(ctx, updateInfo); err == nil {
		if updateInfo.Type == pbdto.UpdateType_Add {
			cw.cache.InvalidateUpdated(updateInfo.Id, []string{updateInfo.PayloadId}, nil)
		} else if updateInfo.Type == pbdto.UpdateType_Delete {
			cw.cache.InvalidateUpdated(updateInfo.Id, nil, []string{updateInfo.PayloadId})
		}
	}

	return
}

func (cw *CacheMiddleware) StudentModifiedListener(ctx context.Context, diffIds *pbdto.DiffIds) (err error) {
	return cw.next.StudentModifiedListener(ctx, diffIds)
}
