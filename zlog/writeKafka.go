package zlog

import (
	"fmt"
	"github.com/Shopify/sarama"
	"go.uber.org/zap/zapcore"
)

type LogKafka struct {
	Producer sarama.SyncProducer
	Topic    string
}

func (lk *LogKafka) Write(p []byte) (n int, err error) {

	msg := &sarama.ProducerMessage{}
	msg.Topic = lk.Topic
	msg.Value = sarama.ByteEncoder(p)
	_, _, err = lk.Producer.SendMessage(msg)

	if err != nil {
		fmt.Println(err)
		return
	}
	return

}

func initKafkaWriter(c *Config) (zapcore.WriteSyncer, error) {
	var (
		kl  LogKafka
		err error
	)
	kl.Topic = "logs"
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.NoResponse
	//随机的分区类型
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	kl.Producer, err = sarama.NewSyncProducer([]string{c.KafkaAddr}, config)
	if err != nil {
		return nil, err
	}

	return zapcore.AddSync(&kl), err
}
