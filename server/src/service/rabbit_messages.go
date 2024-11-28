package service

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"users-service/src/constants"
	"users-service/src/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (u *User) sendLogInAttemptMessage(id string, succesful bool, provider *string) error {
	if provider == nil {
		fiero := constants.InternalProvider
		provider = &fiero
	}

	queueMsg := model.QueueMessage{
		MessageType: constants.LoginAttempt,
		Message: model.LoginAttempt{
			Succesfull: succesful,
			UserId:     id,
			Provider:   *provider,
			Timestamp:  time.Now().GoString(),
		},
	}

	loginAttempt, err := json.Marshal(queueMsg)
	if err != nil {
		return fmt.Errorf("error marshalling login attempt message for rabbit: %w", err)
	}

	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        loginAttempt,
		DeliveryMode: amqp.Persistent,
	}

	fmt.Println("Sending login attempt to queue: ", message)
	err = u.amqpQueue.Publish("", os.Getenv("CLOUDAMQP_QUEUE"), false, false, message)

	if err != nil {
		return fmt.Errorf("error publishing login attempt message to queue: %w", err)
	}
	return nil
}