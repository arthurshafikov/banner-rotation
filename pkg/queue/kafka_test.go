package queue

import (
	"context"
	"testing"
	"time"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	mock_queue "github.com/arthurshafikov/banner-rotation/pkg/queue/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

var (
	value = core.IncrementEvent{
		BannerID:      1,
		SlotID:        1,
		SocialGroupID: 1,
		Datetime:      time.Now(),
	}
	topic = "clicks"
)

func TestAddToQueue(t *testing.T) {
	g, ctx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)
	ctrl := gomock.NewController(t)
	kafkaWriterMock := mock_queue.NewMockKafkaWriter(ctrl)
	queue := NewQueue(ctx, kafkaWriterMock)
	g.Go(func() error {
		return queue.Dispatch()
	})
	gomock.InOrder(
		kafkaWriterMock.EXPECT().WriteMessages(ctx, gomock.Any()).Times(15).Return(nil),
	)

	for i := 0; i < 15; i++ {
		value.BannerID = int64(i)
		value.Datetime = time.Now()
		require.NoError(t, queue.AddToQueue(topic, value))
	}

	for len(queue.elements) > 0 {
		time.Sleep(time.Microsecond)
	}
	cancel()
	require.NoError(t, g.Wait())
}
