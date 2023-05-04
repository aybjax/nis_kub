package app_db

import (
	"context"
	"fmt"

	"github.com/aybjax/nis_lib/pbdto"
	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	client           *mongo.Client
	courseCollection *mongo.Collection
	logger           log.Logger
}

func NewAppDB(client *mongo.Client, logger log.Logger) *DB {
	courseCollection := client.Database("testdb").Collection("courses")
	logger.Log(
		"DB.method", "NewAppDB",
		"msg", "courseCollection retrieved",
	)
	return &DB{
		client:           client,
		courseCollection: courseCollection,
		logger:           logger,
	}
}

func (db *DB) Close() error {
	db.logger.Log(
		"DB.method", "Close",
		"msg", "closing db",
	)
	return db.client.Disconnect(context.Background())
}

func (db *DB) ReadAll(ctx context.Context) ([]*pbdto.Course, error) {
	var docs []*CourseDB

	c, err := db.courseCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		db.logger.Log(
			"DB.method", "ReadAll",
			"msg", "retrive all error",
			"err", fmt.Sprint(err),
		)
		return []*pbdto.Course{}, err
	}

	err = c.All(ctx, &docs)

	if err != nil {
		db.logger.Log(
			"DB.method", "ReadAll",
			"msg", "mapping error",
			"err", fmt.Sprint(err),
		)
		return []*pbdto.Course{}, err
	}

	result := make([]*pbdto.Course, len(docs))

	for i, v := range docs {
		result[i] = v.ToProto()
	}

	return result, err
}

func (db *DB) ReadById(ctx context.Context, id string) (*pbdto.Course, error) {
	var data CourseDB

	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		db.logger.Log(
			"DB.method", "ReadById",
			"msg", "id to hex error",
			"err", fmt.Sprint(err),
		)
		return &pbdto.Course{}, err
	}

	err = db.courseCollection.FindOne(ctx, bson.M{"_id": _id}).Decode(&data)

	if err != nil {
		db.logger.Log(
			"DB.method", "ReadById",
			"msg", "find 1 by id error",
			"err", fmt.Sprint(err),
		)
		return &pbdto.Course{}, err
	}

	result := data.ToProto()

	return result, nil
}

func (db *DB) ReadByStudentId(ctx context.Context, studentId string) ([]*pbdto.Course, error) {
	docs := make([]CourseDB, 0)

	c, err := db.courseCollection.Find(context.TODO(), bson.M{
		"student_ids": studentId,
	})

	if err != nil {
		db.logger.Log(
			"DB.method", "ReadByStudentId",
			"msg", "find 1 by id error",
			"err", fmt.Sprint(err),
		)
		return []*pbdto.Course{}, err
	}

	c.All(ctx, &docs)

	result := make([]*pbdto.Course, len(docs))

	for i, v := range docs {
		result[i] = v.ToProto()
	}

	return result, nil
}

func (db *DB) Create(ctx context.Context, payload *pbdto.Course) (string, []string, error) {
	data, err := CourseDB{}.FromProto(payload)

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

	course, err := db.courseCollection.InsertOne(ctx, &data.CourseDBPayload)

	if err != nil {
		db.logger.Log(
			"DB.method", "Create",
			"msg", "db insertion error",
			"err", fmt.Sprint(err),
		)
		return "", nil, err
	}

	id := course.InsertedID.(primitive.ObjectID)

	return id.Hex(), data.StudentIds, nil
}

func (db *DB) Update(ctx context.Context, id string, payload *pbdto.Course) ([]string, []string, error) {
	var oldDoc CourseDB
	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		db.logger.Log(
			"DB.method", "Update",
			"msg", "id to hex error",
			"err", fmt.Sprint(err),
		)
		return nil, nil, err
	}

	data, err := CourseDB{}.FromProto(payload)

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

	err = db.courseCollection.FindOneAndUpdate(ctx, bson.M{"_id": _id},
		bson.M{
			"$set": &data.CourseDBPayload,
		},
	).Decode(&oldDoc)

	return data.StudentIds, oldDoc.StudentIds, err
}

func (db *DB) Delete(ctx context.Context, id string) ([]string, error) {
	var oldDoc CourseDB
	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		db.logger.Log(
			"DB.method", "Delete",
			"msg", "id to hex error",
			"err", fmt.Sprint(err),
		)
		return nil, err
	}

	err = db.courseCollection.FindOneAndDelete(ctx, bson.M{"_id": _id}).Decode(&oldDoc)

	if err != nil {
		db.logger.Log(
			"DB.method", "Delete",
			"msg", "find and delete error",
			"err", fmt.Sprint(err),
		)
		return nil, err
	}

	return oldDoc.StudentIds, nil
}

func (db *DB) AddStudentIdTo(ctx context.Context, id string, courseId string) (bool, error) {
	var errs []error
	var exists bool
	oldData := CourseDB{}
	_id, _ := primitive.ObjectIDFromHex(id)

	err := db.courseCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": _id},
		bson.M{
			"$push": bson.M{
				"student_ids": courseId,
			},
		},
	).Decode(&oldData)

	if err != nil {
		db.logger.Log(
			"DB.method", "AddStudentIdTo",
			"msg", "find and update error",
			"err", fmt.Sprint(err),
		)
		errs = append(errs, err)
	}

	exists = !oldData.Id.IsZero()

	return exists, err
}

func (db *DB) DeleteStudentIdFrom(ctx context.Context, id string, courseId string) (bool, error) {
	var errs []error
	var exists bool
	oldData := CourseDB{}
	_id, _ := primitive.ObjectIDFromHex(id)

	err := db.courseCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": _id},
		bson.M{
			"$pull": bson.M{
				"student_ids": courseId,
			},
		},
	).Decode(&oldData)

	if err != nil {
		db.logger.Log(
			"DB.method", "DeleteStudentIdFrom",
			"msg", "find and update error",
			"err", fmt.Sprint(err),
		)
		errs = append(errs, err)
	}

	exists = !oldData.Id.IsZero()

	return exists, err
}

func (db *DB) GetStudentIds(id string) ([]string, error) {
	oldData := CourseDB{}
	_id, _ := primitive.ObjectIDFromHex(id)

	err := db.courseCollection.FindOne(
		context.Background(),
		bson.M{"_id": _id},
	).Decode(&oldData)

	if err != nil {
		db.logger.Log(
			"DB.method", "GetStudentIds",
			"msg", "find one error",
			"err", fmt.Sprint(err),
		)
		return nil, err
	}

	return oldData.StudentIds, nil
}
