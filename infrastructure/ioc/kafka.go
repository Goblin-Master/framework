package ioc

import (
	"framework/infrastructure/config"
	"github.com/IBM/sarama"
)

func InitSaramaClient() sarama.Client {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(config.Conf.Kafka.Addrs, cfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSaramaSyncProducer(client sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return p
}

func InitSaramaAsyncProducer(client sarama.Client) sarama.AsyncProducer {
	p, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return p
}
