package app

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aybjax/nis_lib/pbdto"
)

type mock struct {
	Data map[string][]byte
}

func (m *mock) Get(key string) ([]byte, error) {
	return m.Data[key], nil
}

func (m *mock) Set(key string, data []byte) error {
	m.Data[key] = data

	return nil
}

func (m *mock) Delete(key string) error {
	delete(m.Data, key)

	return nil
}

func getCache() *Cache {
	return &Cache{
		&mock{
			Data: map[string][]byte{},
		},
	}
}

func TestWriteAllReadAll(t *testing.T) {
	cache := getCache()

	val := []*pbdto.Student{
		{
			Id:        "1111",
			Name:      "1111",
			FirstName: "1111",
			LastName:  "1111",
			CourseIds: []string{
				"1111", "2222", "2222",
			},
		},
		{
			Id:        "2222",
			Name:      "2222",
			FirstName: "2222",
			LastName:  "2222",
			CourseIds: []string{
				"2222", "3333", "4444",
			},
		},
	}

	cache.WriteAll(val)

	result, _ := cache.RetriveAll()

	if !reflect.DeepEqual(val, result) {
		fmt.Println(val)
		fmt.Println(result)
		t.Errorf("Written data is not equal to read")
	}
}

func TestWrite1Read1(t *testing.T) {
	cache := getCache()

	first := &pbdto.Student{
		Id:        "1111",
		Name:      "1111",
		FirstName: "1111",
		LastName:  "1111",
		CourseIds: []string{
			"1111", "2222", "2222",
		},
	}

	cache.WriteOneById("0", &pbdto.Student{
		Id:        "2222",
		Name:      "2222",
		FirstName: "2222",
		LastName:  "2222",
		CourseIds: []string{
			"2222", "3333", "4444",
		},
	})
	cache.WriteOneById("1", first)
	cache.WriteOneById("2", &pbdto.Student{
		Id:        "2222",
		Name:      "2222",
		FirstName: "2222",
		LastName:  "2222",
		CourseIds: []string{
			"2222", "3333", "4444",
		},
	})

	result, _ := cache.client.Get("1")

	if !reflect.DeepEqual(first, result) {
		t.Errorf("Written 1 data is not equal to read")
	}
}

func TestWriteReadByCourseId(t *testing.T) {
	cache := getCache()

	first := []*pbdto.Student{{
		Id:        "1111",
		Name:      "1111",
		FirstName: "1111",
		LastName:  "1111",
		CourseIds: []string{
			"1111", "2222", "2222",
		},
	},
	}

	cache.WriteByCourseId("0", []*pbdto.Student{{
		Id:        "2222",
		Name:      "2222",
		FirstName: "2222",
		LastName:  "2222",
		CourseIds: []string{
			"2222", "3333", "4444",
		},
	}})
	cache.WriteByCourseId("1", first)
	cache.WriteByCourseId("2", []*pbdto.Student{{
		Id:        "2222",
		Name:      "2222",
		FirstName: "2222",
		LastName:  "2222",
		CourseIds: []string{
			"2222", "3333", "4444",
		},
	}})

	result, _ := cache.client.Get("1")

	if !reflect.DeepEqual(first, result) {
		t.Errorf("Written 1 data is not e")
	}
}

func TestInvalidateCreating(t *testing.T) {
	cache := getCache()

	val := []*pbdto.Student{
		{
			Id:        "1111",
			Name:      "1111",
			FirstName: "1111",
			LastName:  "1111",
			CourseIds: []string{
				"1111", "2222", "2222",
			},
		},
		{
			Id:        "2222",
			Name:      "2222",
			FirstName: "2222",
			LastName:  "2222",
			CourseIds: []string{
				"2222", "3333", "4444",
			},
		},
	}

	cache.WriteAll(val)

	cache.InvalidateCreated()

	result, _ := cache.RetriveAll()

	if len(result) != 0 {
		t.Errorf("Invalidate created did not delete all key")
	}
}

func TestInvalidateUpdated(t *testing.T) {
	cache := getCache()

	first := []*pbdto.Student{{
		Id:        "1111",
		Name:      "1111",
		FirstName: "1111",
		LastName:  "1111",
		CourseIds: []string{
			"1111", "2222", "2222",
		},
	},
	}

	cache.WriteByCourseId("0", []*pbdto.Student{{
		Id:        "2222",
		Name:      "2222",
		FirstName: "2222",
		LastName:  "2222",
		CourseIds: []string{
			"2222", "3333", "4444",
		},
	}})
	cache.WriteByCourseId("1", first)
	cache.WriteByCourseId("2", []*pbdto.Student{{
		Id:        "2222",
		Name:      "2222",
		FirstName: "2222",
		LastName:  "2222",
		CourseIds: []string{
			"2222", "3333", "4444",
		},
	}})
}
