package service

import (
	"context"
	"fmt"
	"nis_courses/app"
	"time"

	"github.com/aybjax/nis_lib/pbdto"
	"github.com/go-kit/kit/metrics"
)

type MetricsMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           app.CourseService
}

func NewMetricsMiddleware(next app.CourseService, requestCount metrics.Counter, requestLatency metrics.Histogram) app.CourseService {
	return &MetricsMiddleware{
		next:           next,
		requestCount:   requestCount,
		requestLatency: requestLatency,
	}
}

func (mw *MetricsMiddleware) GetAll(ctx context.Context) (output []*pbdto.Course, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetAll", "error", fmt.Sprint(err != nil)}

		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.next.GetAll(ctx)

	return
}
func (mw *MetricsMiddleware) Get(ctx context.Context, id string) (output *pbdto.Course, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Get", "error", fmt.Sprint(err != nil)}

		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.next.Get(ctx, id)

	return
}
func (mw *MetricsMiddleware) GetStudents(ctx context.Context, id string) (output []*pbdto.Student, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetStudents", "error", fmt.Sprint(err != nil)}

		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.next.GetStudents(ctx, id)

	return
}
func (mw *MetricsMiddleware) GetCourses(ctx context.Context, id string) (output []*pbdto.Course, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetCourses", "error", fmt.Sprint(err != nil)}

		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.next.GetCourses(ctx, id)

	return
}
func (mw *MetricsMiddleware) Post(ctx context.Context, payload *pbdto.Course) (id interface{}, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Post", "error", fmt.Sprint(err != nil)}

		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	id, err = mw.next.Post(ctx, payload)

	return
}
func (mw *MetricsMiddleware) Put(ctx context.Context, id string, payload *pbdto.Course) (newId interface{}, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Put", "error", fmt.Sprint(err != nil)}

		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	newId, err = mw.next.Put(ctx, id, payload)

	return
}
func (mw *MetricsMiddleware) Delete(ctx context.Context, id string) (oldId interface{}, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Delete", "error", fmt.Sprint(err != nil)}

		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	oldId, err = mw.next.Delete(ctx, id)

	return
}

func (mw *MetricsMiddleware) StudentModifiedListener(ctx context.Context, updateInfo *pbdto.UpdateEmbedded) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "StudentModifiedListener", "error", fmt.Sprint(err != nil)}

		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.StudentModifiedListener(ctx, updateInfo)

	return
}

func (mw *MetricsMiddleware) CourseModifiedListener(ctx context.Context, diffIds *pbdto.DiffIds) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CourseModifiedListener", "error", fmt.Sprint(err != nil)}

		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.CourseModifiedListener(ctx, diffIds)

	return
}
