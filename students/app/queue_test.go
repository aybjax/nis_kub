package app

import (
	"context"
	"testing"

	"github.com/aybjax/nis_lib/cmntypes/mock_cmntypes"
	"github.com/aybjax/nis_lib/pbdto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestEnqueueTransformNotify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver := mock_cmntypes.NewMockAppQueue(ctrl)
	queue := &QueueImpl{driver}

	payload := &pbdto.DiffIds{
		Id:     "testId",
		OldIds: []string{"testOld"},
		NewIds: []string{"tesNew"},
	}

	bs, _ := proto.Marshal(payload)

	driver.EXPECT().Publish(bs, _QUEUE_CHANNEL_TO_MODIFY, _QUEUE_TOPIC_TO_MODIFY)
	queue.EnqueueTransformNotify(payload)
}

func TestTransformNotifyListener(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver := mock_cmntypes.NewMockAppQueue(ctrl)
	queue := &QueueImpl{driver}

	val := &pbdto.DiffIds{
		Id: "DiffIds",
	}
	bs, _ := proto.Marshal(val)

	cb := func(_ context.Context, data *pbdto.DiffIds) error {
		assert.Equal(t, data.Id, "DiffIds")

		return nil
	}

	driver.EXPECT().Subscribe(_QUEUE_CHANNEL_TO_MODIFY,
		_QUEUE_TOPIC_STUDENT_TO_MODIFY_QUEUE, gomock.Any(),
	).Do(func(_, _ string, listenerCb func(data []byte) error) {
		listenerCb(bs)
	})

	queue.TransformNotifyListener(cb)
}

func TestNotifyStudentChanged(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver := mock_cmntypes.NewMockAppQueue(ctrl)
	queue := &QueueImpl{driver}

	payload := &pbdto.UpdateEmbedded{
		Id:        "testId",
		PayloadId: "PayloadId",
		Type:      pbdto.UpdateType_Delete,
	}

	bs, _ := proto.Marshal(payload)

	driver.EXPECT().Publish(
		bs,
		_QUEUE_CHANNEL_STUDENT_MODIFIED,
		_QUEUE_TOPIC_STUDENT_MODIFIED,
	)
	queue.NotifyStudentChanged(context.TODO(), payload)
}

func TestCourseChangedListener(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver := mock_cmntypes.NewMockAppQueue(ctrl)
	queue := &QueueImpl{driver}

	payload := &pbdto.UpdateEmbedded{
		Id: "UpdateEmbeddedId",
	}
	bs, _ := proto.Marshal(payload)

	cb := func(ctx context.Context, update *pbdto.UpdateEmbedded) error {
		assert.Equal(t, update.Id, "UpdateEmbeddedId")

		return nil
	}

	// TODO make queue name constant
	driver.EXPECT().
		Subscribe("course.modified", "COURSE_MODIFIED_QUEUE", gomock.Any()).
		Do(func(_, _ string, listenerCb func(data []byte) error) {
			listenerCb(bs)
		})
	queue.CourseChangedListener(cb)
}
