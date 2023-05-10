package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/aybjax/nis_lib/cmntypes"
	"github.com/aybjax/nis_lib/pbdto"
	"google.golang.org/protobuf/proto"
)

const (
	_QUEUE_CHANNEL_TO_MODIFY             = "student.to_modify"
	_QUEUE_TOPIC_TO_MODIFY               = ""
	_QUEUE_TOPIC_STUDENT_TO_MODIFY_QUEUE = "STUDENT_TO_MODIFY_QUEUE"
	_QUEUE_CHANNEL_STUDENT_MODIFIED      = "student.modified"
	_QUEUE_TOPIC_STUDENT_MODIFIED        = ""
)

//go:generate mockgen -source=./queue.go -destination=./mock_app/mock_queue.go
type Queue interface {
	EnqueueTransformNotify(queueMsg *pbdto.DiffIds) error
	TransformNotifyListener(cb func(context.Context, *pbdto.DiffIds) error)
	NotifyStudentChanged(_ context.Context, update *pbdto.UpdateEmbedded) error
	CourseChangedListener(cb func(ctx context.Context, update *pbdto.UpdateEmbedded) error)
}

type QueueImpl struct {
	client cmntypes.AppQueue
}

func NewQueue(engine cmntypes.AppQueue) Queue {
	return &QueueImpl{engine}
}

func (q *QueueImpl) EnqueueTransformNotify(queueMsg *pbdto.DiffIds) error {
	msg, err := proto.Marshal(queueMsg)

	if err != nil {
		return err
	}

	return q.client.Publish(msg, _QUEUE_CHANNEL_TO_MODIFY, _QUEUE_TOPIC_TO_MODIFY)
}

func (q *QueueImpl) TransformNotifyListener(cb func(context.Context, *pbdto.DiffIds) error) {
	q.client.Subscribe(_QUEUE_CHANNEL_TO_MODIFY, _QUEUE_TOPIC_STUDENT_TO_MODIFY_QUEUE, func(data []byte) error {
		diffIds := &pbdto.DiffIds{}

		proto.Unmarshal(data, diffIds)

		if diffIds.Id == "" {
			return errors.New("Empty data")
		}

		if err := cb(context.Background(), diffIds); err != nil {
			return fmt.Errorf("Could not derialize from %s:%s, data is %s",
				_QUEUE_CHANNEL_TO_MODIFY,
				_QUEUE_TOPIC_STUDENT_TO_MODIFY_QUEUE,
				data)
		}

		return nil
	})
}

func (q *QueueImpl) NotifyStudentChanged(_ context.Context, update *pbdto.UpdateEmbedded) error {
	bs, err := proto.Marshal(update)

	if err != nil {
		return err
	}

	return q.client.Publish(bs, _QUEUE_CHANNEL_STUDENT_MODIFIED, _QUEUE_TOPIC_STUDENT_MODIFIED)
}

func (q *QueueImpl) CourseChangedListener(cb func(ctx context.Context, update *pbdto.UpdateEmbedded) error) {
	q.client.Subscribe("course.modified", "COURSE_MODIFIED_QUEUE", func(data []byte) error {
		updateEmbedded := &pbdto.UpdateEmbedded{}

		if err := proto.Unmarshal(data, updateEmbedded); err != nil {
			return err
		}

		if err := cb(context.Background(), updateEmbedded); err != nil {
			return err
		}

		return nil
	})
}
