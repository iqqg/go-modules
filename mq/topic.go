package mq

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func topicExample(params []string) {
	if len(params) < 1 {
		log.Fatalln("params error")
		os.Exit(0)
	}

	if params[0] == "recv" {
		recvTopic(params[1:])
	} else {
		sendTopic(params)
	}
}

func recvTopic(keys []string) {
	conn, err := amqp.Dial("amqp://pi:shine@192.168.1.4:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	for i := range keys {
		err = ch.QueueBind(
			q.Name,       // queue name
			keys[i],      // routing key
			"logs_topic", // exchange
			false,
			nil)
		failOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func sendTopic(keys []string) {
	conn, err := amqp.Dial("amqp://pi:shine@192.168.1.4:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	for i := 0; i < 5; i++ {
		body := fmt.Sprintf("Hello world! %d", i)
		for i := range keys {
			err = ch.Publish(
				"logs_topic", // exchange
				keys[i],      // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			failOnError(err, "Failed to publish a message")
			log.Printf(" [x] Sent %s", body)
		}
	}
}
