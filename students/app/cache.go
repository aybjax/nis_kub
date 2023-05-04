package app

import (
	"errors"
	"fmt"

	"github.com/aybjax/nis_lib/cmntypes"
	"github.com/aybjax/nis_lib/helper"
	"github.com/aybjax/nis_lib/pbdto"
	"google.golang.org/protobuf/proto"
)

type Cache struct {
	client cmntypes.AppCache
}

func NewCache(engine cmntypes.AppCache) *Cache {
	return &Cache{
		engine,
	}
}

func (c *Cache) WriteAll(data []*pbdto.Student) error {
	bs, err := proto.Marshal(&pbdto.Students{
		Data: data,
	})

	if err == nil {
		return c.client.Set("students.all", bs)
	}

	return err
}

func (c *Cache) RetriveAll() ([]*pbdto.Student, error) {
	val, err := c.client.Get("students.all")

	if err != nil {
		return nil, err
	} else if len(val) == 0 {
		return nil, fmt.Errorf("No data in cache by key = %s", "students.all")
	}

	data := &pbdto.Students{}

	if err := proto.Unmarshal(val, data); err == nil {
		return data.Data, nil
	} else {
		return nil, err
	}
}

func (c *Cache) WriteOneById(id string, data *pbdto.Student) error {
	if bs, err := proto.Marshal(data); err == nil {
		c.client.Set(fmt.Sprintf("students.%s", id), bs)
	}

	return nil
}

func (c *Cache) RetriveOneById(id string) (*pbdto.Student, error) {
	val, err := c.client.Get(fmt.Sprintf("students.%s", id))

	if err != nil {
		return nil, err
	}

	if len(val) == 0 {
		return nil, fmt.Errorf("No data in cache by key = %s", fmt.Sprintf("students.%s", id))
	}

	data := &pbdto.Student{}
	if err := proto.Unmarshal(val, data); err == nil {
		return data, nil
	} else {
		return nil, err
	}
}

func (c *Cache) WriteByCourseId(course_id string, data []*pbdto.Student) error {
	bs, err := proto.Marshal(&pbdto.Students{
		Data: data,
	})

	if err == nil {
		return c.client.Set(fmt.Sprintf("students.by_course.%s", course_id), bs)
	}

	return err
}

func (c *Cache) RetrieveByCourseId(course_id string) ([]*pbdto.Student, error) {
	val, err := c.client.Get(fmt.Sprintf("students.by_course.%s", course_id))

	if err != nil {
		return nil, err
	} else if len(val) == 0 {
		return nil, fmt.Errorf("No data in cache by key = %s", fmt.Sprintf("students.by_course.%s", course_id))
	}

	data := &pbdto.Students{}

	if err := proto.Unmarshal(val, data); err == nil {
		return data.Data, nil
	} else {
		return nil, err
	}
}

func (c *Cache) InvalidateCreated() error {
	return c.client.Delete("students.all")
}

func (c *Cache) InvalidateUpdated(c_id string, newCourseIds []string, oldCourseIds []string) error {
	var errs []error

	for _, c_id := range helper.SetDiff(newCourseIds, oldCourseIds) {
		errs = append(
			errs,
			c.client.Delete(fmt.Sprintf("students.by_course.%s", c_id)),
		)
	}

	for _, c_id := range helper.SetDiff(oldCourseIds, newCourseIds) {
		errs = append(
			errs,
			c.client.Delete(fmt.Sprintf("students.by_course.%s", c_id)),
		)
	}

	errs = append(
		errs,
		c.client.Delete("students.all"),
	)

	errs = append(
		errs,
		c.client.Delete(fmt.Sprintf("students.%s", c_id)),
	)

	return errors.Join(errs...)
}

func (c *Cache) InvalidateDeleted(c_id string, oldCourseIds []string) error {
	var errs []error

	for _, c_id := range oldCourseIds {
		errs = append(
			errs,
			c.client.Delete(fmt.Sprintf("students.by_course.%s", c_id)),
		)
	}

	errs = append(
		errs,
		c.client.Delete("students.all"),
	)

	errs = append(
		errs,
		c.client.Delete(fmt.Sprintf("students.%s", c_id)),
	)

	return errors.Join(errs...)
}
