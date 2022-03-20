package mocks

type QueueMock struct{}

func NewQueueMock() *QueueMock {
	return &QueueMock{}
}

func (q *QueueMock) AddToQueue(topic string, value interface{}) error {
	return nil
}
