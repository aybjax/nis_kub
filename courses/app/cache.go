package app

import (
	"errors"
	"fmt"

	"github.com/aybjax/nis_lib/cmntypes"
	"github.com/aybjax/nis_lib/helper"
	"github.com/aybjax/nis_lib/pbdto"
	"google.golang.org/protobuf/proto"
)

const (
	_CACHE_KEY_ALL_COURSES = "courses.all"
)

func _CACHE_KEY_ID_COURSES(id string) string {
	return fmt.Sprintf("courses.%s", id)
}
func _CACHE_KEY_COURSES_BY_STUDENT_ID(student_id string) string {
	return fmt.Sprintf("courses.by_student.%s", student_id)
}

//go:generate mockgen -source=./cache.go -destination=./mock_app/mock_cache.go

type Cache interface {
	WriteAll(data []*pbdto.Course) error
	RetriveAll() ([]*pbdto.Course, error)
	WriteOneById(id string, data *pbdto.Course) error
	RetriveOneById(id string) (*pbdto.Course, error)
	WriteByStudentId(student_id string, data []*pbdto.Course) error
	RetrieveByStudentId(student_id string) ([]*pbdto.Course, error)
	InvalidateCreated() error
	InvalidateUpdated(c_id string, newStudentIds []string, oldStudentIds []string) error
	InvalidateDeleted(c_id string, oldStudentIds []string) error
}

type CacheImpl struct {
	client cmntypes.AppCache
}

func NewCache(engine cmntypes.AppCache) Cache {
	return &CacheImpl{
		engine,
	}
}

func (c *CacheImpl) WriteAll(data []*pbdto.Course) error {
	bs, err := proto.Marshal(&pbdto.Courses{
		Data: data,
	})

	if err == nil {
		return c.client.Set(_CACHE_KEY_ALL_COURSES, bs)
	}

	return err
}

func (c *CacheImpl) RetriveAll() ([]*pbdto.Course, error) {
	val, err := c.client.Get(_CACHE_KEY_ALL_COURSES)

	if err != nil {
		return nil, err
	} else if len(val) == 0 {
		return nil, fmt.Errorf("No data in cache by key = %s", _CACHE_KEY_ALL_COURSES)
	}

	data := &pbdto.Courses{}

	if err := proto.Unmarshal(val, data); err == nil {
		return data.Data, nil
	} else {
		return nil, err
	}
}

func (c *CacheImpl) WriteOneById(id string, data *pbdto.Course) error {
	if bs, err := proto.Marshal(data); err == nil {
		c.client.Set(_CACHE_KEY_ID_COURSES(id), bs)
	}

	return nil
}

func (c *CacheImpl) RetriveOneById(id string) (*pbdto.Course, error) {
	val, err := c.client.Get(_CACHE_KEY_ID_COURSES(id))

	if err != nil {
		return nil, err
	}

	if len(val) == 0 {
		return nil, fmt.Errorf("No data in cache by key = %s", _CACHE_KEY_ID_COURSES(id))
	}

	data := &pbdto.Course{}
	if err := proto.Unmarshal(val, data); err == nil {
		return data, nil
	} else {
		return nil, err
	}
}

func (c *CacheImpl) WriteByStudentId(student_id string, data []*pbdto.Course) error {
	bs, err := proto.Marshal(&pbdto.Courses{
		Data: data,
	})

	if err == nil {
		return c.client.Set(_CACHE_KEY_COURSES_BY_STUDENT_ID(student_id), bs)
	}

	return err
}

func (c *CacheImpl) RetrieveByStudentId(student_id string) ([]*pbdto.Course, error) {
	val, err := c.client.Get(_CACHE_KEY_COURSES_BY_STUDENT_ID(student_id))

	if err != nil {
		return nil, err
	} else if len(val) == 0 {
		return nil, fmt.Errorf("No data in cache by key = %s", _CACHE_KEY_COURSES_BY_STUDENT_ID(student_id))
	}

	data := &pbdto.Courses{}

	if err := proto.Unmarshal(val, data); err == nil {
		return data.Data, nil
	} else {
		return nil, err
	}
}

func (c *CacheImpl) InvalidateCreated() error {
	return c.client.Delete(_CACHE_KEY_ALL_COURSES)
}

func (c *CacheImpl) InvalidateUpdated(c_id string, newStudentIds []string, oldStudentIds []string) error {
	var errs []error

	for _, s_id := range helper.SetDiff(newStudentIds, oldStudentIds) {
		errs = append(
			errs,
			c.client.Delete(_CACHE_KEY_COURSES_BY_STUDENT_ID(s_id)),
		)
	}

	for _, s_id := range helper.SetDiff(oldStudentIds, newStudentIds) {
		errs = append(
			errs,
			c.client.Delete(_CACHE_KEY_COURSES_BY_STUDENT_ID(s_id)),
		)
	}

	errs = append(
		errs,
		c.client.Delete(_CACHE_KEY_ALL_COURSES),
	)

	errs = append(
		errs,
		c.client.Delete(_CACHE_KEY_ID_COURSES(c_id)),
	)

	return errors.Join(errs...)
}

func (c *CacheImpl) InvalidateDeleted(c_id string, oldStudentIds []string) error {
	var errs []error

	for _, s_id := range oldStudentIds {
		errs = append(
			errs,
			c.client.Delete(_CACHE_KEY_COURSES_BY_STUDENT_ID(s_id)),
		)
	}

	errs = append(
		errs,
		c.client.Delete(_CACHE_KEY_ALL_COURSES),
	)

	errs = append(
		errs,
		c.client.Delete(_CACHE_KEY_ID_COURSES(c_id)),
	)

	return errors.Join(errs...)
}
