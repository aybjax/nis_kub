package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/aybjax/nis_lib/cmntypes"
	"github.com/aybjax/nis_lib/pbdto"
	"google.golang.org/protobuf/proto"
)

type Queue struct {
	client cmntypes.AppQueue
}

func NewQueue(engine cmntypes.AppQueue) *Queue {
	return &Queue{engine}
}

func (q *Queue) EnqueueTransformNotify(queueMsg *pbdto.DiffIds) error {
	msg, err := proto.Marshal(queueMsg)

	if err != nil {
		return err
	}

	return q.client.Publish(msg, "course.to_modify", "")
}

func (q *Queue) TransformNotifyListener(cb func(context.Context, *pbdto.DiffIds) error) {
	q.client.Subscribe("course.to_modify", "COURSE_TO_MODIFY_QUEUE", func(data []byte) error {
		diffIds := &pbdto.DiffIds{}

		proto.Unmarshal(data, diffIds)

		if diffIds.Id == "" {
			return errors.New("Empty data")
		}

		if err := cb(context.Background(), diffIds); err != nil {
			return fmt.Errorf("Could not derialize from %s:%s, data is %s",
				"course.to_modify",
				"COURSE_TO_MODIFY_QUEUE",
				data)
		}

		return nil
	})
}

func (q *Queue) NotifyCourseChanged(update *pbdto.UpdateEmbedded) error {
	bs, err := proto.Marshal(update)

	if err != nil {
		return err
	}

	return q.client.Publish(bs, "course.modified", "")
}

func (q *Queue) StudentChangedListener(cb func(ctx context.Context, update *pbdto.UpdateEmbedded) error) {
	q.client.Subscribe("student.modified", "STUDENT_MODIFIED_QUEUE", func(data []byte) error {
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
