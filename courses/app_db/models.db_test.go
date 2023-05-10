package app_db

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/aybjax/nis_lib/pbdto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFromProto(t *testing.T) {
	id := primitive.NewObjectID()
	expected := &CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "expectedName",
			Description: "expectedDescription",
			Discipline:  "expectedDiscipline",
			Teacher:     "expectedTeacher",
			StudentIds:  []string{"expectedStudentIds"},
		},
	}

	result, _ := (&CourseDB{}).FromProto(&pbdto.Course{
		Id:          id.Hex(),
		Name:        "expectedName",
		Description: "expectedDescription",
		Discipline:  "expectedDiscipline",
		Teacher:     "expectedTeacher",
		StudentIds:  []string{"expectedStudentIds"},
	})

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("asdasddad")
	}
}

func TestGetCourseIds(t *testing.T) {
	id := primitive.NewObjectID()
	expected := &CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "expectedName",
			Description: "expectedDescription",
			Discipline:  "expectedDiscipline",
			Teacher:     "expectedTeacher",
			StudentIds:  []string{"expectedStudentIds"},
		},
	}

	if !reflect.DeepEqual(expected.StudentIds, []string{"expectedStudentIds"}) {
		t.Errorf("asdasddad")
	}
}

func TestToProto(t *testing.T) {
	id := primitive.NewObjectID()
	expected := (&CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "expectedName",
			Description: "expectedDescription",
			Discipline:  "expectedDiscipline",
			Teacher:     "expectedTeacher",
			StudentIds:  []string{"expectedStudentIds"},
		},
	}).ToProto()

	result := &pbdto.Course{
		Id:          id.Hex(),
		Name:        "expectedName",
		Description: "expectedDescription",
		Discipline:  "expectedDiscipline",
		Teacher:     "expectedTeacher",
		StudentIds:  []string{"expectedStudentIds"},
	}

	expectedJson, _ := json.Marshal(expected)
	resultJson, _ := json.Marshal(result)

	if !reflect.DeepEqual(expectedJson, resultJson) {
		t.Errorf("asdasddad")
	}
}

func TestGetPayload(t *testing.T) {
	id := primitive.NewObjectID()
	expected := &CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "expectedName",
			Description: "expectedDescription",
			Discipline:  "expectedDiscipline",
			Teacher:     "expectedTeacher",
			StudentIds:  []string{"expectedStudentIds"},
		},
	}
	if !reflect.DeepEqual(expected.GetPayload(), CourseDBPayload{
		Name:        "expectedName",
		Description: "expectedDescription",
		Discipline:  "expectedDiscipline",
		Teacher:     "expectedTeacher",
		StudentIds:  []string{"expectedStudentIds"},
	}) {
		t.Errorf("asdasddad")
	}
}

func TestHasValidPayload(t *testing.T) {
	id := primitive.NewObjectID()
	expected := &CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "12",
			Description: "expectedDescription",
			Discipline:  "expectedDiscipline",
			Teacher:     "expectedTeacher",
			StudentIds:  []string{"expectedStudentIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "123456789012345678901",
			Description: "expectedDescription",
			Discipline:  "expectedDiscipline",
			Teacher:     "expectedTeacher",
			StudentIds:  []string{"expectedStudentIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "12345678901234567890",
			Description: "12",
			Discipline:  "expectedDiscipline",
			Teacher:     "expectedTeacher",
			StudentIds:  []string{"expectedStudentIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "12345678901234567890",
			Description: strings.Repeat("0", 201),
			Discipline:  "expectedDiscipline",
			Teacher:     "expectedTeacher",
			StudentIds:  []string{"expectedStudentIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "12345678901234567890",
			Description: strings.Repeat("0", 20),
			Discipline:  "12",
			Teacher:     "expectedTeacher",
			StudentIds:  []string{"expectedStudentIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "12345678901234567890",
			Description: strings.Repeat("0", 200),
			Discipline:  strings.Repeat("0", 21),
			Teacher:     "expectedTeacher",
			StudentIds:  []string{"expectedStudentIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "12345678901234567890",
			Description: strings.Repeat("0", 200),
			Discipline:  strings.Repeat("0", 20),
			Teacher:     "12",
			StudentIds:  []string{"expectedStudentIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "12345678901234567890",
			Description: strings.Repeat("0", 200),
			Discipline:  strings.Repeat("0", 20),
			Teacher:     strings.Repeat("0", 21),
			StudentIds:  []string{"expectedStudentIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	studentIds := make([]string, 61)
	for i := 0; i < 61; i++ {
		studentIds[i] = "1"
	}

	expected = &CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "12345678901234567890",
			Description: strings.Repeat("0", 200),
			Discipline:  strings.Repeat("0", 20),
			Teacher:     strings.Repeat("0", 20),
			StudentIds:  studentIds,
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	studentIds = make([]string, 60)
	for i := 0; i < 60; i++ {
		studentIds[i] = "1"
	}

	expected = &CourseDB{
		Id: id,
		CourseDBPayload: CourseDBPayload{
			Name:        "12345678901234567890",
			Description: strings.Repeat("0", 200),
			Discipline:  strings.Repeat("0", 20),
			Teacher:     strings.Repeat("0", 20),
			StudentIds:  studentIds,
		},
	}

	if err := expected.HasValidPayload(); err != nil {
		t.Errorf(err.Error())
	}
}
