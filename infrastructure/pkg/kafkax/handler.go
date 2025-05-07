package kafkax

import (
	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"go.uber.org/zap"
)

type Handler[T any] struct {
	zap *zap.SugaredLogger
	fn  func(msg *sarama.ConsumerMessage, event T) error
}

func NewHandler[T any](zap *zap.SugaredLogger, fn func(msg *sarama.ConsumerMessage, event T) error) *Handler[T] {
	return &Handler[T]{zap: zap, fn: fn}
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		// 在这里调用业务处理逻辑
		var t T
		err := sonic.Unmarshal(msg.Value, &t)
		if err != nil {
			// 你也可以在这里引入重试的逻辑
			h.zap.Errorf("反序列消息体失败:%v,topic:%s,partition:%d,offset:%d",
				err, msg.Topic, msg.Partition, msg.Offset)
		}
		err = h.fn(msg, t)
		if err != nil {
			h.zap.Errorf("消息处理失败:%v,topic:%s,partition:%d,offset:%d",
				err, msg.Topic, msg.Partition, msg.Offset)
		}
		session.MarkMessage(msg, "消费成功")
	}
	return nil
}

//type Event struct {
//	Aid int64 `json:"aid"`
//	Uid int64 `json:"uid"`
//}
//func (t *Test) Start() error {
//	cg, err := sarama.NewConsumerGroupFromClient("consume_group", t.client)
//	if err != nil {
//		return err
//	}
//	go func() {
//		e := cg.Consume(context.Background(),
//			[]string{"topic"},
//			saramax.NewHandler[Event](t.zap,t.Consume))
//		if e != nil {
//			//日志
//		}
//	}()
//	return err
//}
//
//func (t *Test) Consume(msg *sarama.ConsumerMessage,
//	event Event) error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	return i.repo.AddRecord(ctx, domain.HistoryRecord{
//		BizId: event.Aid,
//		Biz:   "article",
//		Uid:   event.Uid,
//	})
//}
