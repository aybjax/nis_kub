package service

import (
	"context"
	"errors"
	"fmt"
	"nis_courses/app"
	"nis_courses/app_db"

	"github.com/aybjax/nis_lib/helper"
	"github.com/aybjax/nis_lib/pbdto"
	"github.com/go-kit/log"
)

type ServiceInstance struct {
	db         app_db.DB
	grpcClient pbdto.GetCourseStudentsClient
	queue      app.Queue
	logger     log.Logger
}

func NewService(db app_db.DB,
	grpcClient pbdto.GetCourseStudentsClient, queue app.Queue, logger log.Logger,
) app.CourseService {
	s := &ServiceInstance{
		db:         db,
		grpcClient: grpcClient,
		queue:      queue,
		logger:     logger,
	}

	return s
}

func (s *ServiceInstance) GetAll(ctx context.Context) ([]*pbdto.Course, error) {
	s.logger.Log(
		"identity", "ServiceInstance",
		"method", "GetAll",
	)
	result, err := s.db.ReadAll(ctx)

	if err != nil {
		s.logger.Log(
			"identity", "service",
			"method", "GetAll",
			"err", fmt.Sprint(err),
		)
		return []*pbdto.Course{}, helper.NewMapError(err)
	}

	return result, nil
}

func (s *ServiceInstance) Get(ctx context.Context, id string) (*pbdto.Course, error) {
	s.logger.Log(
		"identity", "service",
		"method", "Get",
	)
	result, err := s.db.ReadById(ctx, id)

	if err != nil {
		s.logger.Log(
			"identity", "service",
			"method", "Get",
			"err", fmt.Sprint(err),
		)
		return &pbdto.Course{}, helper.NewMapError(err)
	}

	return result, nil
}

func (s *ServiceInstance) GetStudents(ctx context.Context, id string) ([]*pbdto.Student, error) {
	s.logger.Log(
		"identity", "service",
		"method", "GetStudents",
	)
	resp, err := s.grpcClient.GetStudents(ctx, &pbdto.Request{
		Id: id,
	})

	if err != nil {
		s.logger.Log(
			"identity", "service",
			"method", "GetStudents",
			"err", fmt.Sprint(err),
		)
		return nil, err
	}

	return resp.Students, nil
}

func (s *ServiceInstance) GetCourses(ctx context.Context, studentId string) ([]*pbdto.Course, error) {
	s.logger.Log(
		"identity", "service",
		"method", "GetCourses",
	)
	result, err := s.db.ReadByStudentId(ctx, studentId)

	if err != nil {
		s.logger.Log(
			"identity", "service",
			"method", "GetCourses",
			"err", fmt.Sprint(err),
		)
		return []*pbdto.Course{}, helper.NewMapError(err)
	}

	return result, nil
}

func (s *ServiceInstance) Post(ctx context.Context, payload *pbdto.Course) (interface{}, error) {
	s.logger.Log(
		"identity", "service",
		"method", "Post",
	)
	id, studentIds, err := s.db.Create(ctx, payload)

	if err != nil {
		s.logger.Log(
			"identity", "service",
			"method", "Post",
			"err", fmt.Sprint(err),
		)
		return "", helper.NewMapError(err)
	}

	if err := s.queue.EnqueueTransformNotify(&pbdto.DiffIds{
		Id:     id,
		NewIds: studentIds,
	}); err != nil {
		s.logger.Log(
			"identity", "service",
			"method", "Post",
			"errorType", "data transformation",
			"err", fmt.Sprint(err),
		)
	}

	return id, nil
}

func (s *ServiceInstance) Put(ctx context.Context, id string, payload *pbdto.Course) (interface{}, error) {
	s.logger.Log("identity", "service",
		"method", "Put",
	)
	newIds, oldIds, err := s.db.Update(ctx, id, payload)

	if err != nil {
		s.logger.Log("identity", "service",
			"method", "Put",
			"err", fmt.Sprint(err),
		)
		return "", helper.NewMapError(err)
	}

	if err := s.queue.EnqueueTransformNotify(&pbdto.DiffIds{
		Id:     id,
		NewIds: newIds,
		OldIds: oldIds,
	}); err != nil {
		s.logger.Log("identity", "service",
			"method", "Put",
			"errorType", "data transformation",
			"err", fmt.Sprint(err),
		)
	}

	return id, nil
}

func (s *ServiceInstance) Delete(ctx context.Context, id string) (interface{}, error) {
	s.logger.Log("identity", "service",
		"method", "Delete",
	)
	oldIds, err := s.db.Delete(ctx, id)

	if err != nil {
		s.logger.Log("identity", "service",
			"method", "Delete",
			"err", fmt.Sprint(err),
		)
		return "", helper.NewMapError(err)
	}

	if err := s.queue.EnqueueTransformNotify(&pbdto.DiffIds{
		Id:     id,
		OldIds: oldIds,
	}); err != nil {
		s.logger.Log("identity", "service",
			"method", "Delete",
			"errorType", "data transformation",
			"err", fmt.Sprint(err),
		)
	}

	return "OK", nil
}

func (s *ServiceInstance) StudentModifiedListener(ctx context.Context, updateInfo *pbdto.UpdateEmbedded) error {
	s.logger.Log("identity", "service",
		"method", "studentModifiedListener",
	)
	var errs []error
	var exists bool
	var err error

	if updateInfo.Type == pbdto.UpdateType_Add {
		exists, err = s.db.AddStudentIdTo(context.Background(), updateInfo.Id, updateInfo.PayloadId)
		s.logger.Log("identity", "service",
			"method", "studentModifiedListener",
			"updateType", "add",
			"err", fmt.Sprint(err),
		)
	} else if updateInfo.Type == pbdto.UpdateType_Delete {
		exists, err = s.db.DeleteStudentIdFrom(context.Background(), updateInfo.Id, updateInfo.PayloadId)
		s.logger.Log("identity", "service",
			"method", "studentModifiedListener",
			"updateType", "delete",
			"err", fmt.Sprint(err),
		)
	}

	errs = append(errs, err)

	if !exists && updateInfo.Type == pbdto.UpdateType_Add {
		s.logger.Log("identity", "service",
			"method", "studentModifiedListener",
			"updateType", "not exist",
			"err", fmt.Sprint(errors.Join(errs...)),
		)
		if err := s.queue.NotifyCourseChanged(&pbdto.UpdateEmbedded{
			Type:      pbdto.UpdateType_Delete,
			Id:        updateInfo.PayloadId,
			PayloadId: updateInfo.Id,
		}); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (s *ServiceInstance) CourseModifiedListener(ctx context.Context, diffIds *pbdto.DiffIds) error {
	s.logger.Log("identity", "service",
		"method", "courseModifiedListener",
	)
	if len(diffIds.NewIds) == 0 && len(diffIds.OldIds) == 0 {
		s.logger.Log("identity", "service",
			"method", "courseModifiedListener",
			"err", "empty data",
		)
		return errors.New("Empty ids")
	}

	var errs []error

	if len(diffIds.OldIds) == 0 {
		s.logger.Log("identity", "service",
			"method", "courseModifiedListener",
			"updateType", "add",
		)
		for m := range helper.GenerateUpdateMessage(diffIds.Id, diffIds.NewIds, pbdto.UpdateType_Add) {
			errs = append(
				errs,
				s.queue.NotifyCourseChanged(m),
			)
		}

		return errors.Join(errs...)
	}

	if len(diffIds.NewIds) == 0 {
		s.logger.Log("identity", "service",
			"method", "courseModifiedListener",
			"updateType", "delete",
		)
		for m := range helper.GenerateUpdateMessage(diffIds.Id, diffIds.OldIds, pbdto.UpdateType_Delete) {
			errs = append(
				errs,
				s.queue.NotifyCourseChanged(m),
			)
		}

		return errors.Join(errs...)
	}

	s.logger.Log("identity", "service",
		"method", "courseModifiedListener",
		"updateType", "mixed",
	)
	deletedIds := helper.SetDiff(diffIds.OldIds, diffIds.NewIds)
	for m := range helper.GenerateUpdateMessage(diffIds.Id, deletedIds, pbdto.UpdateType_Delete) {
		errs = append(
			errs,
			s.queue.NotifyCourseChanged(m),
		)
	}

	addedIds := helper.SetDiff(diffIds.NewIds, diffIds.OldIds)
	for m := range helper.GenerateUpdateMessage(diffIds.Id, addedIds, pbdto.UpdateType_Add) {
		errs = append(
			errs,
			s.queue.NotifyCourseChanged(m),
		)
	}

	return errors.Join(errs...)
}
