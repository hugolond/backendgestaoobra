package queue

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func Connect() (conn *amqp.Connection) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}
	var localhost = os.Getenv("LOCALHOST_RABBIT")
	var user = os.Getenv("USER_RABBIT")
	var password = os.Getenv("PASSWORD_RABBIT")
	//fmt.Println(user)
	dsn := "amqp://" + user + ":" + password + "@" + localhost + ":5672/"
	conn, err = amqp.Dial(dsn)
	if err != nil {
		panic(err.Error())
	}
	return
}

func Notify(payload []byte, exchange string, routingKey string) {

	conn := Connect()
	ch, err := conn.Channel()
	if err != nil {
		panic(err.Error())
	}
	defer ch.Close()
	defer conn.Close()

	err = ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(payload),
		})

	if err != nil {
		panic(err.Error())
	}
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | Message sent RabbitMq : " + string(payload))
}

func StartCosuming(ch *amqp.Channel, in chan []byte) {

	//defer ch.Close()
	/*q, err := ch.QueueDeclare(
		"delivered_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err.Error())
	}*/
	msg, err := ch.Consume(
		"delivered_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err.Error())
	}

	go func() {
		for m := range msg {
			in <- []byte(m.Body)
		}
		close(in)
	}()
	fmt.Println("Waiting for messages...")
}
