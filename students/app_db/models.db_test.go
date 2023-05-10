package app_db

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/aybjax/nis_lib/pbdto"
)

func TestFromProto(t *testing.T) {
	expected := &StudentDB{
		Id: "expectedId",
		StudentDBPayload: StudentDBPayload{
			Name:      "expectedName",
			FirstName: "expectedFirstName",
			LastName:  "expectedLastName",
			CourseIds: []string{"expectedCourseIds"},
		},
	}

	result, _ := (&StudentDB{}).FromProto(&pbdto.Student{
		Id:        "expectedId",
		Name:      "expectedName",
		FirstName: "expectedFirstName",
		LastName:  "expectedLastName",
		CourseIds: []string{"expectedCourseIds"},
	})

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("asdasddad")
	}
}

func TestGetCourseIds(t *testing.T) {
	expected := &StudentDB{
		Id: "expectedId",
		StudentDBPayload: StudentDBPayload{
			Name:      "expectedName",
			FirstName: "expectedFirstName",
			LastName:  "expectedLastName",
			CourseIds: []string{"expectedCourseIds"},
		},
	}

	if !reflect.DeepEqual(expected.CourseIds, []string{"expectedCourseIds"}) {
		t.Errorf("asdasddad")
	}
}

func TestToProto(t *testing.T) {
	expected := (&StudentDB{
		Id: "expectedId",
		StudentDBPayload: StudentDBPayload{
			Name:      "expectedName",
			FirstName: "expectedFirstName",
			LastName:  "expectedLastName",
			CourseIds: []string{"expectedCourseIds"},
		},
	}).ToProto()

	result := &pbdto.Student{
		Id:        "expectedId",
		Name:      "expectedName",
		FirstName: "expectedFirstName",
		LastName:  "expectedLastName",
		CourseIds: []string{"expectedCourseIds"},
	}

	expectedJson, _ := json.Marshal(expected)
	resultJson, _ := json.Marshal(result)

	if !reflect.DeepEqual(expectedJson, resultJson) {
		t.Errorf("asdasddad")
	}
}

func TestGetPayload(t *testing.T) {
	expected := &StudentDB{
		Id: "expectedId",
		StudentDBPayload: StudentDBPayload{
			Name:      "expectedName",
			FirstName: "expectedFirstName",
			LastName:  "expectedLastName",
			CourseIds: []string{"expectedCourseIds"},
		},
	}
	if !reflect.DeepEqual(expected.GetPayload(), StudentDBPayload{
		Name:      "expectedName",
		FirstName: "expectedFirstName",
		LastName:  "expectedLastName",
		CourseIds: []string{"expectedCourseIds"},
	}) {
		t.Errorf("asdasddad")
	}
}

func TestHasValidPayload(t *testing.T) {
	expected := &StudentDB{
		Id: "expectedId",
		StudentDBPayload: StudentDBPayload{
			Name:      "12",
			FirstName: "expectedFirstName",
			LastName:  "expectedLastName",
			CourseIds: []string{"expectedCourseIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &StudentDB{
		Id: "expectedId",
		StudentDBPayload: StudentDBPayload{
			Name:      "123456789012345678901",
			FirstName: "expectedFirstName",
			LastName:  "expectedLastName",
			CourseIds: []string{"expectedCourseIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &StudentDB{
		Id: "expectedId",
		StudentDBPayload: StudentDBPayload{
			Name:      "12345678901234567890",
			FirstName: "12",
			LastName:  "expectedLastName",
			CourseIds: []string{"expectedCourseIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &StudentDB{
		Id: "expectedId",
		StudentDBPayload: StudentDBPayload{
			Name:      "12345678901234567890",
			FirstName: "123456789012345678901",
			LastName:  "expectedLastName",
			CourseIds: []string{"expectedCourseIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &StudentDB{
		Id: "expectedId",
		StudentDBPayload: StudentDBPayload{
			Name:      "12345678901234567890",
			FirstName: "12345678901234567890",
			LastName:  "12",
			CourseIds: []string{"expectedCourseIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &StudentDB{
		Id: "expectedId",
		StudentDBPayload: StudentDBPayload{
			Name:      "12345678901234567890",
			FirstName: "12345678901234567890",
			LastName:  "123456789012345678901",
			CourseIds: []string{"expectedCourseIds"},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &StudentDB{
		Id: "expectedId",
		StudentDBPayload: StudentDBPayload{
			Name:      "12345678901234567890",
			FirstName: "12345678901234567890",
			LastName:  "12345678901234567890",
			CourseIds: []string{
				"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16",
			},
		},
	}

	if err := expected.HasValidPayload(); err == nil {
		t.Errorf("asdasds")
	}

	expected = &StudentDB{
		Id: "expectedId",
		StudentDBPayload: StudentDBPayload{
			Name:      "12345678901234567890",
			FirstName: "12345678901234567890",
			LastName:  "12345678901234567890",
			CourseIds: []string{
				"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15",
			},
		},
	}

	if err := expected.HasValidPayload(); err != nil {
		t.Errorf("asdasds")
	}
}
