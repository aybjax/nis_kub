package service

import (
	"context"
	"errors"
	"nis_courses/app/mock_app"
	"nis_courses/app_db/mock_app_db"
	"reflect"
	"testing"

	"github.com/aybjax/nis_lib/helper"
	"github.com/aybjax/nis_lib/pbdto"
	"github.com/aybjax/nis_lib/pbdto/mock_pbdto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestServiceTest(t *testing.T) {
	suite.Run(t, new(serviceTest))
}

type L struct{}

func (L) Log(keyvals ...interface{}) error {
	return nil
}

var (
	db         *mock_app_db.MockDB
	grpcClient *mock_pbdto.MockGetCourseStudentsClient
	queue      *mock_app.MockQueue
)

type serviceTest struct {
	suite.Suite
	ctrl    *gomock.Controller
	service *ServiceInstance
}

func (s *serviceTest) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	db = mock_app_db.NewMockDB(s.ctrl)
	grpcClient = mock_pbdto.NewMockGetCourseStudentsClient(s.ctrl)
	queue = mock_app.NewMockQueue(s.ctrl)

	s.service = &ServiceInstance{
		db:         db,
		grpcClient: grpcClient,
		queue:      queue,
		logger:     L{},
	}
}

func (s *serviceTest) TearDownTest() {
	s.ctrl.Finish()
}

func (s *serviceTest) TestGetAll() {
	s.T().Run("No db error", func(t *testing.T) {
		ctx := context.TODO()

		db.EXPECT().ReadAll(ctx).Return(nil, nil)
		_, err := s.service.GetAll(context.TODO())

		s.Equal(nil, err)
	})
	s.T().Run("Returns error", func(t *testing.T) {
		ctx := context.TODO()

		dbErr := errors.New("db error")
		db.EXPECT().ReadAll(ctx).Return(nil, dbErr)
		_, err := s.service.GetAll(context.TODO())

		if err == nil || !reflect.DeepEqual(err, helper.NewMapError(dbErr)) {
			s.T().Errorf("Db error is not propagated")
		}
	})
}

func (s *serviceTest) TestGet() {
	s.T().Run("No db error", func(t *testing.T) {
		ctx := context.TODO()

		db.EXPECT().ReadById(ctx, "1").Return(nil, nil)
		_, err := s.service.Get(context.TODO(), "1")

		s.Equal(nil, err)
	})
	s.T().Run("Returns error", func(t *testing.T) {
		ctx := context.TODO()

		dbErr := errors.New("db error")
		db.EXPECT().ReadById(ctx, "1").Return(nil, dbErr)
		_, err := s.service.Get(context.TODO(), "1")

		if err == nil || !reflect.DeepEqual(err, helper.NewMapError(dbErr)) {
			s.T().Errorf("Db error is not propagated")
		}
	})
}

func (s *serviceTest) TestGetCourses() {
	s.T().Run("No db error", func(t *testing.T) {
		ctx := context.TODO()

		db.EXPECT().ReadByStudentId(ctx, "1").Return(nil, nil)
		_, err := s.service.GetCourses(context.TODO(), "1")

		s.Equal(nil, err)
	})
	s.T().Run("Returns error", func(t *testing.T) {
		ctx := context.TODO()

		dbErr := errors.New("db error")
		db.EXPECT().ReadByStudentId(ctx, "1").Return(nil, dbErr)
		_, err := s.service.GetCourses(context.TODO(), "1")

		if err == nil || !reflect.DeepEqual(err, helper.NewMapError(dbErr)) {
			s.T().Errorf("Db error is not propagated")
		}
	})
}

func (s *serviceTest) TestGetStudents() {
	s.T().Run("No grpc error", func(t *testing.T) {
		ctx := context.TODO()

		grpcClient.EXPECT().GetStudents(ctx, &pbdto.Request{
			Id: "1",
		}).Return(&pbdto.StudentsResponse{
			Students: nil,
		}, nil)
		_, err := s.service.GetStudents(ctx, "1")

		s.Equal(nil, err)
	})
	s.T().Run("Returns error", func(t *testing.T) {
		grpcErr := errors.New("grpc error")
		ctx := context.TODO()
		grpcClient.EXPECT().GetStudents(ctx, &pbdto.Request{
			Id: "1",
		}).Return(nil, grpcErr)
		_, err := s.service.GetStudents(ctx, "1")

		s.Equal(err, grpcErr)
	})
}

func (s *serviceTest) TestPost() {
	s.T().Run("No db error", func(t *testing.T) {
		ctx := context.TODO()

		payload := &pbdto.Course{
			Id: "PayloadId",
		}
		db.EXPECT().Create(ctx, payload).Return("1", []string{"2"}, nil)
		queue.EXPECT().EnqueueTransformNotify(&pbdto.DiffIds{
			Id:     "1",
			NewIds: []string{"2"},
		})
		id, err := s.service.Post(ctx, payload)
		s.Equal(map[string]interface{}{
			"id": "1",
		}, id)
		s.Equal(nil, err)
	})
	s.T().Run("Returns error", func(t *testing.T) {
		ctx := context.TODO()

		payload := &pbdto.Course{
			Id: "PayloadId",
		}
		dbErr := errors.New("db error")
		db.EXPECT().Create(ctx, payload).Return("", nil, dbErr)
		_, err := s.service.Post(ctx, payload)

		if err == nil || !reflect.DeepEqual(err, helper.NewMapError(dbErr)) {
			s.T().Errorf("Db error is not propagated")
		}
	})
}

func (s *serviceTest) TestPut() {
	s.T().Run("No db error", func(t *testing.T) {
		ctx := context.TODO()

		id := "PutId"
		payload := &pbdto.Course{
			Id: "PayloadId",
		}
		db.EXPECT().Update(ctx, id, payload).Return([]string{"1"}, []string{"2"}, nil)
		queue.EXPECT().EnqueueTransformNotify(&pbdto.DiffIds{
			Id:     id,
			NewIds: []string{"1"},
			OldIds: []string{"2"},
		})
		idInf, err := s.service.Put(ctx, id, payload)
		s.Equal("OK", idInf)
		s.Equal(nil, err)
	})
	s.T().Run("Returns error", func(t *testing.T) {
		ctx := context.TODO()

		id := "PutId"
		payload := &pbdto.Course{
			Id: "PayloadId",
		}
		dbErr := errors.New("db error")
		db.EXPECT().Update(ctx, id, payload).Return(nil, nil, dbErr)
		_, err := s.service.Put(ctx, id, payload)

		if err == nil || !reflect.DeepEqual(err, helper.NewMapError(dbErr)) {
			s.T().Errorf("Db error is not propagated")
		}
	})
}

func (s *serviceTest) TestDelete() {
	s.T().Run("No db error", func(t *testing.T) {
		ctx := context.TODO()

		id := "PutId"
		db.EXPECT().Delete(ctx, id).Return([]string{"1"}, nil)
		queue.EXPECT().EnqueueTransformNotify(&pbdto.DiffIds{
			Id:     id,
			OldIds: []string{"1"},
		})
		idInf, err := s.service.Delete(ctx, id)
		s.Equal("OK", idInf)
		s.Equal(nil, err)
	})
	s.T().Run("Returns error", func(t *testing.T) {
		ctx := context.TODO()

		id := "PutId"
		dbErr := errors.New("db error")
		db.EXPECT().Delete(ctx, id).Return(nil, dbErr)
		_, err := s.service.Delete(ctx, id)

		if err == nil || !reflect.DeepEqual(err, helper.NewMapError(dbErr)) {
			s.T().Errorf("Db error is not propagated")
		}
	})
}

func (s *serviceTest) TestStudentModifiedListener() {
	s.T().Run("Add and exists", func(t *testing.T) {
		updateInfo := &pbdto.UpdateEmbedded{
			Id:        "UpdatedId",
			PayloadId: "PayloadId",
			Type:      pbdto.UpdateType_Add,
		}
		ctx := context.TODO()

		db.EXPECT().AddStudentIdTo(ctx, updateInfo.Id,
			updateInfo.PayloadId).
			Return(true, nil)

		err := s.service.StudentModifiedListener(ctx, updateInfo)

		if err != nil {
			t.Errorf("%s", err)
		}
	})
	s.T().Run("Add and does not exist", func(t *testing.T) {
		updateInfo := &pbdto.UpdateEmbedded{
			Id:        "UpdatedId",
			PayloadId: "PayloadId",
			Type:      pbdto.UpdateType_Add,
		}
		ctx := context.TODO()

		db.EXPECT().AddStudentIdTo(ctx, updateInfo.Id,
			updateInfo.PayloadId).
			Return(false, nil)

		queue.EXPECT().NotifyCourseChanged(&pbdto.UpdateEmbedded{
			Type:      pbdto.UpdateType_Delete,
			Id:        updateInfo.PayloadId,
			PayloadId: updateInfo.Id,
		})

		err := s.service.StudentModifiedListener(ctx, updateInfo)

		if err != nil {
			t.Errorf("%s", err)
		}
	})
	s.T().Run("Delete and exists", func(t *testing.T) {
		updateInfo := &pbdto.UpdateEmbedded{
			Id:        "UpdatedId",
			PayloadId: "PayloadId",
			Type:      pbdto.UpdateType_Delete,
		}
		ctx := context.TODO()

		db.EXPECT().DeleteStudentIdFrom(ctx, updateInfo.Id,
			updateInfo.PayloadId).
			Return(false, nil)

		err := s.service.StudentModifiedListener(ctx, updateInfo)

		if err != nil {
			t.Errorf("%s", err)
		}
	})
	s.T().Run("Delete and does not exist", func(t *testing.T) {
		updateInfo := &pbdto.UpdateEmbedded{
			Id:        "UpdatedId",
			PayloadId: "PayloadId",
			Type:      pbdto.UpdateType_Delete,
		}
		ctx := context.TODO()

		db.EXPECT().DeleteStudentIdFrom(ctx, updateInfo.Id,
			updateInfo.PayloadId).
			Return(false, nil)

		err := s.service.StudentModifiedListener(ctx, updateInfo)

		if err != nil {
			t.Errorf("%s", err)
		}
	})
	s.T().Run("Add and every error", func(t *testing.T) {
		updateInfo := &pbdto.UpdateEmbedded{
			Id:        "UpdatedId",
			PayloadId: "PayloadId",
			Type:      pbdto.UpdateType_Add,
		}
		ctx := context.TODO()
		dbErr := errors.New("db error")
		quErr := errors.New("db error")

		db.EXPECT().AddStudentIdTo(ctx, updateInfo.Id,
			updateInfo.PayloadId).
			Return(false, dbErr)
		queue.EXPECT().NotifyCourseChanged(&pbdto.UpdateEmbedded{
			Type:      pbdto.UpdateType_Delete,
			Id:        updateInfo.PayloadId,
			PayloadId: updateInfo.Id,
		}).Return(quErr)

		err := s.service.StudentModifiedListener(ctx, updateInfo)

		s.Equal(err, errors.Join(dbErr, quErr))
	})
	s.T().Run("Delete and every error", func(t *testing.T) {
		updateInfo := &pbdto.UpdateEmbedded{
			Id:        "UpdatedId",
			PayloadId: "PayloadId",
			Type:      pbdto.UpdateType_Delete,
		}
		ctx := context.TODO()
		dbErr := errors.New("db error")

		db.EXPECT().DeleteStudentIdFrom(ctx, updateInfo.Id,
			updateInfo.PayloadId).
			Return(false, dbErr)

		err := s.service.StudentModifiedListener(ctx, updateInfo)

		s.Equal(err, errors.Join(dbErr))
	})
}

func (s *serviceTest) TestCourseModifiedListener() {
	s.T().Run("All new and old ids", func(t *testing.T) {
		payload := &pbdto.DiffIds{
			Id:     "diffId",
			OldIds: []string{"oldId1", "oldId2", "commonId"},
			NewIds: []string{"newId1", "newId2", "commonId"},
		}
		ctx := context.TODO()

		queue.EXPECT().NotifyCourseChanged(&pbdto.UpdateEmbedded{
			Id:        "oldId1",
			PayloadId: "diffId",
			Type:      pbdto.UpdateType_Delete,
		})
		queue.EXPECT().NotifyCourseChanged(&pbdto.UpdateEmbedded{
			Id:        "oldId2",
			PayloadId: "diffId",
			Type:      pbdto.UpdateType_Delete,
		})
		queue.EXPECT().NotifyCourseChanged(&pbdto.UpdateEmbedded{
			Id:        "newId1",
			PayloadId: "diffId",
			Type:      pbdto.UpdateType_Add,
		})
		queue.EXPECT().NotifyCourseChanged(&pbdto.UpdateEmbedded{
			Id:        "newId2",
			PayloadId: "diffId",
			Type:      pbdto.UpdateType_Add,
		})

		s.service.CourseModifiedListener(ctx, payload)
	})

	s.T().Run("All new ids", func(t *testing.T) {
		payload := &pbdto.DiffIds{
			Id:     "diffId",
			NewIds: []string{"newId1", "newId2"},
		}
		ctx := context.TODO()

		queue.EXPECT().NotifyCourseChanged(&pbdto.UpdateEmbedded{
			Id:        "newId1",
			PayloadId: "diffId",
			Type:      pbdto.UpdateType_Add,
		})
		queue.EXPECT().NotifyCourseChanged(&pbdto.UpdateEmbedded{
			Id:        "newId2",
			PayloadId: "diffId",
			Type:      pbdto.UpdateType_Add,
		})

		s.service.CourseModifiedListener(ctx, payload)
	})
	s.T().Run("All old ids", func(t *testing.T) {
		payload := &pbdto.DiffIds{
			Id:     "diffId",
			OldIds: []string{"oldId1", "oldId2"},
		}
		ctx := context.TODO()

		queue.EXPECT().NotifyCourseChanged(&pbdto.UpdateEmbedded{
			Id:        "oldId1",
			PayloadId: "diffId",
			Type:      pbdto.UpdateType_Delete,
		})
		queue.EXPECT().NotifyCourseChanged(&pbdto.UpdateEmbedded{
			Id:        "oldId2",
			PayloadId: "diffId",
			Type:      pbdto.UpdateType_Delete,
		})

		s.service.CourseModifiedListener(ctx, payload)
	})

	s.T().Run("All new and old errors", func(t *testing.T) {
		payload := &pbdto.DiffIds{
			Id:     "diffId",
			OldIds: []string{"oldId1"},
			NewIds: []string{"newId1"},
		}
		ctx := context.TODO()

		err1 := errors.New("err1")
		err2 := errors.New("err2")

		// TODO add context if needed
		queue.EXPECT().NotifyCourseChanged(gomock.Any()).Return(err1)
		queue.EXPECT().NotifyCourseChanged(gomock.Any()).Return(err2)

		err := s.service.CourseModifiedListener(ctx, payload)

		s.Equal(err, errors.Join(err1, err2))
	})

	s.T().Run("All new ids", func(t *testing.T) {
		payload := &pbdto.DiffIds{
			Id:     "diffId",
			NewIds: []string{"newId1", "newId2"},
		}
		ctx := context.TODO()

		err1 := errors.New("err1")
		err2 := errors.New("err2")

		queue.EXPECT().NotifyCourseChanged(gomock.Any()).Return(err1)
		queue.EXPECT().NotifyCourseChanged(gomock.Any()).Return(err2)

		err := s.service.CourseModifiedListener(ctx, payload)

		s.Equal(err, errors.Join(err1, err2))
	})
	s.T().Run("All old ids", func(t *testing.T) {
		payload := &pbdto.DiffIds{
			Id:     "diffId",
			OldIds: []string{"oldId1", "oldId2"},
		}
		ctx := context.TODO()

		err1 := errors.New("err1")
		err2 := errors.New("err2")

		queue.EXPECT().NotifyCourseChanged(gomock.Any()).Return(err1)
		queue.EXPECT().NotifyCourseChanged(gomock.Any()).Return(err2)

		err := s.service.CourseModifiedListener(ctx, payload)

		s.Equal(err, errors.Join(err1, err2))
	})
}
