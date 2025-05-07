package kafkax

import (
	"context"
	"framework/infrastructure/config"
	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(config.Conf.Kafka.Addrs, cfg)
	defer producer.Close()
	assert.NoError(t, err)
	bizMsg := MyBizMsg{
		Value: "这是kafka的同步发送",
	}
	data, _ := sonic.Marshal(bizMsg)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test",
		Value: sarama.ByteEncoder(data),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace_id"),
				Value: []byte("123456"),
			},
		},
		Metadata: "只作用于发送过程",
	})
	assert.NoError(t, err)
}

func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(config.Conf.Kafka.Addrs, cfg)
	defer producer.Close()
	assert.NoError(t, err)
	bizMsg := MyBizMsg{
		Value: "这是kafka的异步发送",
	}
	data, _ := sonic.Marshal(bizMsg)
	msgChan := producer.Input()
	msgChan <- &sarama.ProducerMessage{
		Topic: "test",
		Value: sarama.ByteEncoder(data),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace_id"),
				Value: []byte("123456"),
			},
		},
		Metadata: "只作用于发送过程",
	}
	errChan := producer.Errors()
	sussChan := producer.Successes()
	select {
	case err := <-errChan:
		t.Error("发送出了问题", err.Err, err.Msg.Value)
	case suc := <-sussChan:
		t.Log("发送成功", suc.Value)
	case <-time.After(5 * time.Second): // 等待5秒超时
		t.Fatal("等待消息处理超时")
	}
}

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	// 启用自动偏移量提交
	//cfg.Consumer.Offsets.AutoCommit.Enable = true
	// 设置自动提交偏移量的间隔
	//cfg.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	// 设置自动偏移量重置为 OffsetNewest，这确保消费者组在第一次启动时从最新消息开始消费
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumer, err := sarama.NewConsumerGroup(config.Conf.Kafka.Addrs, "test_group", cfg)
	assert.NoError(t, err)
	defer consumer.Close()
	start := time.Now()
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//defer cancel()
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*10, func() {
		cancel()
	})
	err = consumer.Consume(ctx,
		[]string{"test"}, &testConsumerGroupHandler{})
	t.Log(err, time.Since(start).String())
}

type testConsumerGroupHandler struct {
}

func (t *testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("setup")
	//针对消费者组中的消费者
	//partitions := session.Claims()["test"]
	//for _, partition := range partitions {
	//	session.ResetOffset("test", partition, sarama.OffsetNewest, "重置偏移量")
	//}
	return nil
}
func (t *testConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("cleanup")
	return nil
}
func (t *testConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	//同步消费
	for msg := range claim.Messages() {
		var bizMsg MyBizMsg
		err := sonic.Unmarshal(msg.Value, &bizMsg)
		if err != nil {
			log.Println("解析消息出错", err)
			log.Println("消息内容", string(msg.Value))
			continue
		}
		log.Println("收到消息", bizMsg.Value)
		session.MarkMessage(msg, "已消费")
	}
	//异步消费
	//asyncConsume(session, claim)
	return nil
}
func asyncConsume(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) {
	const batchSize = 10
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		var done bool
		var eg errgroup.Group
		for i := 0; i < batchSize && !done; i++ {
			select {
			case msg, ok := <-claim.Messages():
				if !ok {
					//消费者被关闭
					cancel()
					break
				}

				eg.Go(func() error {
					var bizMsg MyBizMsg
					err := sonic.Unmarshal(msg.Value, &bizMsg)
					if err != nil {
						log.Println("解析消息出错", err)
						return err
					}
					log.Println("收到消息", bizMsg.Value)
					return nil
				})
				cancel()
				err := eg.Wait()
				if err != nil {
					log.Println("消费出错", err)
				}
			case <-ctx.Done():
				done = true
			}
		}
		for msg := range claim.Messages() {
			session.MarkMessage(msg, "消费成功")
		}
	}
}

type MyBizMsg struct {
	Value string `json:"value"`
}
