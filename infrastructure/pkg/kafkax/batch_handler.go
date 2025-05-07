package kafkax

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"go.uber.org/zap"
	"time"
)

type BatchHandler[T any] struct {
	fn        func(msgs []*sarama.ConsumerMessage, events []T) error
	zap       *zap.SugaredLogger
	batchSize int
}

func NewBatchHandler[T any](zap *zap.SugaredLogger, fn func(msgs []*sarama.ConsumerMessage, events []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{zap: zap, fn: fn, batchSize: 10}
}
func (b *BatchHandler[T]) WithBatchSize(batchSize int) {
	b.batchSize = batchSize
}

func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for {
		// 解决有序性，开 channel，同一个业务发到同一个 channel
		batch := make([]*sarama.ConsumerMessage, 0, b.batchSize)
		events := make([]T, 0, b.batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		var done = false
		for i := 0; i < b.batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				// 超时了
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				var t T
				err := sonic.Unmarshal(msg.Value, &t)
				if err != nil {
					b.zap.Errorf("反序列消息体失败:%v,topic:%s,partition:%d,offset:%d",
						err, msg.Topic, msg.Partition, msg.Offset)
					continue
				}
				batch = append(batch, msg)
				events = append(events, t)
			}
		}
		cancel()
		// 凑够了一批，然后你就处理
		err := b.fn(batch, events)
		if err != nil {
			b.zap.Errorf("消息处理失败:%v", err)
			//记录这一批次消息信息
			for msg := range msgs {
				b.zap.Errorf("topic:%s,partition:%d,offset:%d",
					msg.Topic, msg.Partition, msg.Offset)
			}
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "消费成功")
		}
	}
}
