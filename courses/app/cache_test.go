package app

import (
	"encoding/json"
	"testing"

	"github.com/aybjax/nis_lib/cmntypes/mock_cmntypes"
	"github.com/aybjax/nis_lib/pbdto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestWriteAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver := mock_cmntypes.NewMockAppCache(ctrl)
	cache := &CacheImpl{driver}

	val := []*pbdto.Course{
		{
			Id:          "1111",
			Name:        "1111",
			Description: "1111",
			Discipline:  "1111",
			Teacher:     "1111",
			StudentIds: []string{
				"1111", "2222", "2222",
			},
		},
		{
			Id:          "2222",
			Name:        "2222",
			Description: "2222",
			Discipline:  "2222",
			Teacher:     "2222",
			StudentIds: []string{
				"2222", "3333", "4444",
			},
		},
	}

	bs, _ := proto.Marshal(&pbdto.Courses{
		Data: val,
	})

	driver.EXPECT().Set(_CACHE_KEY_ALL_COURSES, bs)
	cache.WriteAll(val)
}

func TestReadAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver := mock_cmntypes.NewMockAppCache(ctrl)
	cache := &CacheImpl{driver}

	val := []*pbdto.Course{
		{
			Id:          "1111",
			Name:        "1111",
			Description: "1111",
			Discipline:  "1111",
			Teacher:     "1111",
			StudentIds: []string{
				"1111", "2222", "2222",
			},
		},
		{
			Id:          "2222",
			Name:        "2222",
			Description: "2222",
			Discipline:  "2222",
			Teacher:     "2222",
			StudentIds: []string{
				"2222", "3333", "4444",
			},
		},
	}

	bs, _ := proto.Marshal(&pbdto.Courses{
		Data: val,
	})

	driver.EXPECT().Get(_CACHE_KEY_ALL_COURSES).Return(bs, nil)
	result, _ := cache.RetriveAll()

	resultByte, _ := json.Marshal(result)
	valByte, _ := json.Marshal(val)

	assert.Equal(t, resultByte, valByte, "Cache Retrieved is not equal")
}

func TestWrite1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver := mock_cmntypes.NewMockAppCache(ctrl)
	cache := &CacheImpl{driver}

	first := &pbdto.Course{
		Id:          "1111",
		Name:        "1111",
		Description: "1111",
		Discipline:  "1111",
		Teacher:     "1111",
		StudentIds: []string{
			"1111", "2222", "2222",
		},
	}
	bs, _ := proto.Marshal(first)

	driver.EXPECT().Set(_CACHE_KEY_ID_COURSES("0"), bs)
	cache.WriteOneById("0", first)
}

func TestRead1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver := mock_cmntypes.NewMockAppCache(ctrl)
	cache := &CacheImpl{driver}

	first := &pbdto.Course{
		Id:          "1111",
		Name:        "1111",
		Description: "1111",
		Discipline:  "1111",
		Teacher:     "1111",
		StudentIds: []string{
			"1111", "2222", "2222",
		},
	}
	bs, _ := proto.Marshal(first)

	driver.EXPECT().Get(_CACHE_KEY_ID_COURSES("0")).Return(bs, nil)
	result, _ := cache.RetriveOneById("0")

	resultByte, _ := json.Marshal(result)
	valByte, _ := json.Marshal(first)

	assert.Equal(t, resultByte, valByte, "Read 1 data is not equal")
}

func TestWriteByCourseId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver := mock_cmntypes.NewMockAppCache(ctrl)
	cache := &CacheImpl{driver}

	data := []*pbdto.Course{{
		Id:          "1111",
		Name:        "1111",
		Description: "1111",
		Discipline:  "1111",
		Teacher:     "1111",
		StudentIds: []string{
			"1111", "2222", "2222",
		},
	}}
	bs, _ := proto.Marshal(&pbdto.Courses{
		Data: data,
	})

	driver.EXPECT().Set(_CACHE_KEY_COURSES_BY_STUDENT_ID("0"), bs)
	cache.WriteByStudentId("0", data)
}
func TestReadByCourseId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver := mock_cmntypes.NewMockAppCache(ctrl)
	cache := &CacheImpl{driver}

	data := []*pbdto.Course{{
		Id:          "1111",
		Name:        "1111",
		Description: "1111",
		Discipline:  "1111",
		Teacher:     "1111",
		StudentIds: []string{
			"1111", "2222", "2222",
		},
	}}
	bs, _ := proto.Marshal(&pbdto.Courses{
		Data: data,
	})

	driver.EXPECT().Get(_CACHE_KEY_COURSES_BY_STUDENT_ID("0")).Return(bs, nil)
	result, _ := cache.RetrieveByStudentId("0")

	resultByte, _ := json.Marshal(result)
	valByte, _ := json.Marshal(data)

	assert.Equal(t, resultByte, valByte, "Read data by course is not equal")
}

func TestInvalidateCreating(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver := mock_cmntypes.NewMockAppCache(ctrl)
	cache := &CacheImpl{driver}

	driver.EXPECT().Delete(_CACHE_KEY_ALL_COURSES)

	cache.InvalidateCreated()
}

func TestInvalidateUpdated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver := mock_cmntypes.NewMockAppCache(ctrl)
	cache := &CacheImpl{driver}

	driver.EXPECT().Delete(_CACHE_KEY_ALL_COURSES)
	driver.EXPECT().Delete(_CACHE_KEY_ID_COURSES("0"))
	driver.EXPECT().Delete(_CACHE_KEY_COURSES_BY_STUDENT_ID("1"))
	driver.EXPECT().Delete(_CACHE_KEY_COURSES_BY_STUDENT_ID("2"))

	cache.InvalidateUpdated("0", []string{"1"}, []string{"2"})
}

func TestInvalidatDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver := mock_cmntypes.NewMockAppCache(ctrl)
	cache := &CacheImpl{driver}

	driver.EXPECT().Delete(_CACHE_KEY_ALL_COURSES)
	driver.EXPECT().Delete(_CACHE_KEY_ID_COURSES("0"))
	driver.EXPECT().Delete(_CACHE_KEY_COURSES_BY_STUDENT_ID("1"))

	cache.InvalidateDeleted("0", []string{"1"})
}
