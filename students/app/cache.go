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
	_CACHE_KEY_ALL_STUDENTS = "students.all"
)

func _CACHE_KEY_ID_STUDENTS(id string) string {
	return fmt.Sprintf("students.%s", id)
}
func _CACHE_KEY_STUDENTS_BY_COURSE_ID(course_id string) string {
	return fmt.Sprintf("students.by_course.%s", course_id)
}

//go:generate mockgen -source=./cache.go -destination=./mock_app/mock_cache.go
type Cache interface {
	WriteAll(data []*pbdto.Student) error
	RetriveAll() ([]*pbdto.Student, error)
	WriteOneById(id string, data *pbdto.Student) error
	RetriveOneById(id string) (*pbdto.Student, error)
	WriteByCourseId(course_id string, data []*pbdto.Student) error
	RetrieveByCourseId(course_id string) ([]*pbdto.Student, error)
	InvalidateCreated() error
	InvalidateUpdated(c_id string, newCourseIds []string, oldCourseIds []string) error
	InvalidateDeleted(c_id string, oldCourseIds []string) error
}

type CacheImpl struct {
	client cmntypes.AppCache
}

func NewCache(engine cmntypes.AppCache) Cache {
	return &CacheImpl{
		engine,
	}
}

func (c *CacheImpl) WriteAll(data []*pbdto.Student) error {
	bs, err := proto.Marshal(&pbdto.Students{
		Data: data,
	})

	if err == nil {
		return c.client.Set(_CACHE_KEY_ALL_STUDENTS, bs)
	}

	return err
}

func (c *CacheImpl) RetriveAll() ([]*pbdto.Student, error) {
	val, err := c.client.Get(_CACHE_KEY_ALL_STUDENTS)

	if err != nil {
		return nil, err
	} else if len(val) == 0 {
		return nil, fmt.Errorf("No data in cache by key = %s", _CACHE_KEY_ALL_STUDENTS)
	}

	data := &pbdto.Students{}

	if err := proto.Unmarshal(val, data); err == nil {
		return data.Data, nil
	} else {
		return nil, err
	}
}

func (c *CacheImpl) WriteOneById(id string, data *pbdto.Student) error {
	if bs, err := proto.Marshal(data); err == nil {
		c.client.Set(_CACHE_KEY_ID_STUDENTS(id), bs)
	}

	return nil
}

func (c *CacheImpl) RetriveOneById(id string) (*pbdto.Student, error) {
	val, err := c.client.Get(_CACHE_KEY_ID_STUDENTS(id))

	if err != nil {
		return nil, err
	}

	if len(val) == 0 {
		return nil, fmt.Errorf("No data in cache by key = %s", _CACHE_KEY_ID_STUDENTS(id))
	}

	data := &pbdto.Student{}
	if err := proto.Unmarshal(val, data); err == nil {
		return data, nil
	} else {
		return nil, err
	}
}

func (c *CacheImpl) WriteByCourseId(course_id string, data []*pbdto.Student) error {
	bs, err := proto.Marshal(&pbdto.Students{
		Data: data,
	})

	if err == nil {
		return c.client.Set(_CACHE_KEY_STUDENTS_BY_COURSE_ID(course_id), bs)
	}

	return err
}

func (c *CacheImpl) RetrieveByCourseId(course_id string) ([]*pbdto.Student, error) {
	val, err := c.client.Get(_CACHE_KEY_STUDENTS_BY_COURSE_ID(course_id))

	if err != nil {
		return nil, err
	} else if len(val) == 0 {
		return nil, fmt.Errorf("No data in cache by key = %s", _CACHE_KEY_STUDENTS_BY_COURSE_ID(course_id))
	}

	data := &pbdto.Students{}

	if err := proto.Unmarshal(val, data); err == nil {
		return data.Data, nil
	} else {
		return nil, err
	}
}

func (c *CacheImpl) InvalidateCreated() error {
	return c.client.Delete(_CACHE_KEY_ALL_STUDENTS)
}

func (c *CacheImpl) InvalidateUpdated(s_id string, newCourseIds []string, oldCourseIds []string) error {
	var errs []error

	for _, c_id := range helper.SetDiff(newCourseIds, oldCourseIds) {
		errs = append(
			errs,
			c.client.Delete(_CACHE_KEY_STUDENTS_BY_COURSE_ID(c_id)),
		)
	}

	for _, c_id := range helper.SetDiff(oldCourseIds, newCourseIds) {
		errs = append(
			errs,
			c.client.Delete(_CACHE_KEY_STUDENTS_BY_COURSE_ID(c_id)),
		)
	}

	errs = append(
		errs,
		c.client.Delete(_CACHE_KEY_ALL_STUDENTS),
	)

	errs = append(
		errs,
		c.client.Delete(_CACHE_KEY_ID_STUDENTS(s_id)),
	)

	return errors.Join(errs...)
}

func (c *CacheImpl) InvalidateDeleted(s_id string, oldCourseIds []string) error {
	var errs []error

	for _, c_id := range oldCourseIds {
		errs = append(
			errs,
			c.client.Delete(_CACHE_KEY_STUDENTS_BY_COURSE_ID(c_id)),
		)
	}

	errs = append(
		errs,
		c.client.Delete(_CACHE_KEY_ALL_STUDENTS),
	)

	errs = append(
		errs,
		c.client.Delete(_CACHE_KEY_ID_STUDENTS(s_id)),
	)

	return errors.Join(errs...)
}
