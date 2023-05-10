package app_db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/aybjax/nis_lib/helper"
	"github.com/aybjax/nis_lib/pbdto"
	"github.com/go-kit/log"
	"github.com/lib/pq"
)

var GET_ALL_ERROR = errors.New("could not retrieve all data")

//go:generate mockgen -source=./db.go -destination=./mock_app_db/mock_db.go
type DB interface {
	Close() error
	ReadAll() ([]*pbdto.Student, error)
	ReadById(id string) (*pbdto.Student, error)
	ReadByCourseId(course_id string) ([]*pbdto.Student, error)
	Create(payload *pbdto.Student) (string, []string, error)
	Update(id string, payload *pbdto.Student) ([]string, []string, error)
	Delete(id string) ([]string, error)
	AddCourseIdTo(id string, courseId string) (bool, error)
	DeleteCourseIdFrom(id string, courseId string) (bool, error)
	GetCourseIds(id string) ([]string, error)
}

type DBImpl struct {
	client *sql.DB
	logger log.Logger
}

func NewAppDB(client *sql.DB, logger log.Logger) DB {
	return &DBImpl{
		client: client,
		logger: logger,
	}
}

func (db *DBImpl) Close() error {
	db.logger.Log(
		"DB.method", "Close",
		"msg", "closing db",
	)
	return db.client.Close()
}

func (db *DBImpl) ReadAll() ([]*pbdto.Student, error) {
	result := make([]*pbdto.Student, 0)

	rows, err := db.client.Query(`
		SELECT id, name, first_name, last_name, course_ids FROM students
	`)

	if err != nil {
		db.logger.Log(
			"DB.method", "ReadAll",
			"msg", "retrive all error",
			"err", fmt.Sprint(err),
		)
		return []*pbdto.Student{}, err
	}

	defer rows.Close()

	for rows.Next() {
		row := &pbdto.Student{}

		err = rows.Scan(&row.Id, &row.Name, &row.FirstName, &row.LastName, pq.Array(&row.CourseIds))

		if err != nil {
			db.logger.Log(
				"DB.method", "ReadAll",
				"msg", "retrive all - scanning error",
				"err", fmt.Sprint(err),
			)
			return []*pbdto.Student{}, err
		}

		result = append(result, row)
	}

	return result, err
}

func (db *DBImpl) ReadById(id string) (*pbdto.Student, error) {
	var data StudentDB

	err := db.client.QueryRow(`
	SELECT id, name, first_name, last_name, course_ids FROM students WHERE id=$1
	`, id).Scan(&data.Id, &data.Name, &data.FirstName, &data.LastName, pq.Array(&data.CourseIds))

	if err != nil {
		db.logger.Log(
			"DB.method", "ReadById",
			"msg", "retrieve error",
			"err", fmt.Sprint(err),
		)
		return &pbdto.Student{}, err
	}

	result := data.ToProto()

	return result, nil
}

func (db *DBImpl) ReadByCourseId(course_id string) ([]*pbdto.Student, error) {
	result := make([]*pbdto.Student, 0)

	rows, err := db.client.Query(`
	SELECT * FROM students WHERE course_ids && $1
	`, pq.Array([]string{course_id}))

	if err != nil {
		db.logger.Log(
			"DB.method", "ReadByCourseId",
			"msg", "find 1 by id error",
			"err", fmt.Sprint(err),
		)
		return []*pbdto.Student{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var row pbdto.Student

		err = rows.Scan(&row.Id, &row.Name, &row.FirstName, &row.LastName, pq.Array(&row.CourseIds))

		if err != nil {
			db.logger.Log(
				"DB.method", "ReadByCourseId",
				"msg", "find 1 by id -- scan error",
				"err", fmt.Sprint(err),
			)
			return []*pbdto.Student{}, err
		}

		result = append(result, &row)
	}

	return result, nil
}

func (db *DBImpl) Create(payload *pbdto.Student) (string, []string, error) {
	data, err := StudentDB{}.FromProto(payload)
	var id string

	if err != nil {
		db.logger.Log(
			"DB.method", "Create",
			"msg", "data mapping error",
			"err", fmt.Sprint(err),
		)
		return "", nil, err
	}

	if err := data.HasValidPayload(); err != nil {
		db.logger.Log(
			"DB.method", "Create",
			"msg", "data validation error",
			"err", fmt.Sprint(err),
		)
		return "", nil, err
	}

	err = db.client.QueryRow(`
	INSERT INTO students (name, first_name, last_name, course_ids)
	VALUES ($1, $2, $3, $4) RETURNING id;
	`, data.Name, data.FirstName, data.LastName, pq.Array(data.GetCourseIds())).Scan(&id)

	if err != nil {
		db.logger.Log(
			"DB.method", "Create",
			"msg", "db insertion error",
			"err", fmt.Sprint(err),
		)
		return "", nil, err
	}

	return id, data.CourseIds, nil
}

func (db *DBImpl) Update(id string, payload *pbdto.Student) ([]string, []string, error) {
	var oldIds []string
	var idExist string
	data, err := StudentDB{}.FromProto(payload)

	if err != nil {
		db.logger.Log(
			"DB.method", "Update",
			"msg", "data mapping error",
			"err", fmt.Sprint(err),
		)
		return nil, nil, err
	}

	if err := data.HasValidPayload(); err != nil {
		db.logger.Log(
			"DB.method", "Update",
			"msg", "data validation error",
			"err", fmt.Sprint(err),
		)
		return nil, nil, err
	}

	err = db.client.QueryRow(`
	UPDATE students
		SET name = $1, first_name = $2, last_name = $3, course_ids = $4
		WHERE id = $5
		RETURNING id, (SELECT course_ids FROM students WHERE id = $5);
	`, data.Name, data.FirstName, data.LastName, pq.Array(data.GetCourseIds()), id).
		Scan(&idExist, pq.Array(&oldIds))

	if err != nil {
		db.logger.Log(
			"DB.method", "Update",
			"msg", "db update error",
			"err", fmt.Sprint(err),
		)
		return nil, nil, err
	}

	if idExist == "" {
		db.logger.Log(
			"DB.method", "Update",
			"msg", "db update error -- no row",
			"err", fmt.Sprint(err),
		)
		return nil, nil, errors.New("not found")
	}

	return data.CourseIds, oldIds, err
}

func (db *DBImpl) Delete(id string) ([]string, error) {
	var oldIds []string
	var idExist string

	err := db.client.QueryRow(`
		DELETE FROM students WHERE id = $1
		RETURNING id, (SELECT course_ids FROM students WHERE id = $1);
	`, id).
		Scan(&idExist, pq.Array(&oldIds))

	if err != nil {
		db.logger.Log(
			"DB.method", "Delete",
			"msg", "db delete error",
			"err", fmt.Sprint(err),
		)
		return nil, helper.NewMapError(err)
	}

	if idExist == "" {
		db.logger.Log(
			"DB.method", "Delete",
			"msg", "db delete error -- no row",
			"err", fmt.Sprint(err),
		)
		return nil, helper.NewMapError(errors.New("Not found"))
	}

	return oldIds, err
}

func (db *DBImpl) AddCourseIdTo(id string, courseId string) (bool, error) {
	var existingId string

	err := db.client.QueryRow(`
					UPDATE students
						SET course_ids = ARRAY_APPEND(course_ids, $1)
						WHERE id = $2
							AND NOT $1 = ANY(course_ids)
					RETURNING id;
					`,
		courseId, id).Scan(&existingId)

	if err != nil {
		db.logger.Log(
			"DB.method", "AddCourseIdTo",
			"msg", "find and update error",
			"err", fmt.Sprint(err),
		)
		db.client.QueryRow(`
							SELECT id FROM students where id = $1 LIMIT 1
					`, id).Scan(&existingId)
	}

	return existingId != "", err
}

func (db *DBImpl) DeleteCourseIdFrom(id string, courseId string) (bool, error) {
	var existingId string

	err := db.client.QueryRow(`
					UPDATE students
						SET course_ids = ARRAY_REMOVE(course_ids, $1)
						WHERE id = $2
							AND $1 = ANY(course_ids)
					RETURNING id;
				`, courseId, id).Scan(&existingId)

	if err != nil {
		db.logger.Log(
			"DB.method", "DeleteCourseIdFrom",
			"msg", "find and update error",
			"err", fmt.Sprint(err),
		)
		db.client.QueryRow(`
							SELECT id FROM students where id = $1 LIMIT 1
					`, id).Scan(&existingId)
	}

	return existingId != "", err
}

func (db *DBImpl) GetCourseIds(id string) ([]string, error) {
	var result []string

	err := db.client.QueryRow(`
		SELECT course_ids FROM students WHERE id = $1 LIMIT 1;
	`, id).Scan(pq.Array(&result))

	if err != nil {
		db.logger.Log(
			"DB.method", "GetCourseIds",
			"msg", "find one error",
			"err", fmt.Sprint(err),
		)
		return result, err
	}

	return result, nil
}
