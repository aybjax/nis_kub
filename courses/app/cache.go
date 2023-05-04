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

func (c *Cache) WriteAll(data []*pbdto.Course) error {
	bs, err := proto.Marshal(&pbdto.Courses{
		Data: data,
	})

	if err == nil {
		return c.client.Set("courses.all", bs)
	}

	return err
}

func (c *Cache) RetriveAll() ([]*pbdto.Course, error) {
	val, err := c.client.Get("courses.all")

	if err != nil {
		return nil, err
	} else if len(val) == 0 {
		return nil, fmt.Errorf("No data in cache by key = %s", "courses.all")
	}

	data := &pbdto.Courses{}

	if err := proto.Unmarshal(val, data); err == nil {
		return data.Data, nil
	} else {
		return nil, err
	}
}

func (c *Cache) WriteOneById(id string, data *pbdto.Course) error {
	if bs, err := proto.Marshal(data); err == nil {
		c.client.Set(fmt.Sprintf("courses.%s", id), bs)
	}

	return nil
}

func (c *Cache) RetriveOneById(id string) (*pbdto.Course, error) {
	val, err := c.client.Get(fmt.Sprintf("courses.%s", id))

	if err != nil {
		return nil, err
	}

	if len(val) == 0 {
		return nil, fmt.Errorf("No data in cache by key = %s", fmt.Sprintf("courses.%s", id))
	}

	data := &pbdto.Course{}
	if err := proto.Unmarshal(val, data); err == nil {
		return data, nil
	} else {
		return nil, err
	}
}

func (c *Cache) WriteByStudentId(student_id string, data []*pbdto.Course) error {
	bs, err := proto.Marshal(&pbdto.Courses{
		Data: data,
	})

	if err == nil {
		return c.client.Set(fmt.Sprintf("courses.by_student.%s", student_id), bs)
	}

	return err
}

func (c *Cache) RetrieveByStudentId(student_id string) ([]*pbdto.Course, error) {
	val, err := c.client.Get(fmt.Sprintf("courses.by_student.%s", student_id))

	if err != nil {
		return nil, err
	} else if len(val) == 0 {
		return nil, fmt.Errorf("No data in cache by key = %s", fmt.Sprintf("courses.by_student.%s", student_id))
	}

	data := &pbdto.Courses{}

	if err := proto.Unmarshal(val, data); err == nil {
		return data.Data, nil
	} else {
		return nil, err
	}
}

func (c *Cache) InvalidateCreated() error {
	return c.client.Delete("courses.all")
}

func (c *Cache) InvalidateUpdated(c_id string, newStudentIds []string, oldStudentIds []string) error {
	var errs []error

	for _, s_id := range helper.SetDiff(newStudentIds, oldStudentIds) {
		errs = append(
			errs,
			c.client.Delete(fmt.Sprintf("courses.by_student.%s", s_id)),
		)
	}

	for _, s_id := range helper.SetDiff(oldStudentIds, newStudentIds) {
		errs = append(
			errs,
			c.client.Delete(fmt.Sprintf("courses.by_student.%s", s_id)),
		)
	}

	errs = append(
		errs,
		c.client.Delete("courses.all"),
	)

	errs = append(
		errs,
		c.client.Delete(fmt.Sprintf("courses.%s", c_id)),
	)

	return errors.Join(errs...)
}

func (c *Cache) InvalidateDeleted(c_id string, oldStudentIds []string) error {
	var errs []error

	for _, s_id := range oldStudentIds {
		errs = append(
			errs,
			c.client.Delete(fmt.Sprintf("courses.by_student.%s", s_id)),
		)
	}

	errs = append(
		errs,
		c.client.Delete("courses.all"),
	)

	errs = append(
		errs,
		c.client.Delete(fmt.Sprintf("courses.%s", c_id)),
	)

	return errors.Join(errs...)
}
