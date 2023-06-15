package service

import (
	"admin/pkg/kafka"
	"log"
)

type KafkaService struct {
	producer   *kafka.Producer
	consumer   *kafka.Consumer
	ResponseCh chan []byte
}

func NewKafkaSerivce(producer *kafka.Producer, consumer *kafka.Consumer) *KafkaService {
	return &KafkaService{
		producer:   producer,
		consumer:   consumer,
		ResponseCh: make(chan []byte),
	}
}

func (s *KafkaService) SendMessages(topic string, message string) error {
	if err := s.producer.SendMessage(topic, message); err != nil {
		return err
	}

	log.Println("Message sent to Kafka:", message)
	return nil

}

func (s *KafkaService) ConsumeMessages(topic string, handler func(message string)) error {
	return s.consumer.ConsumeMessages(topic, handler)
}

func (s *KafkaService) Close() {
	_ = s.consumer.Close()
	_ = s.producer.Close()
}
