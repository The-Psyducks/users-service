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

func sendMessage(queue *amqp.Channel, msg []byte) error {
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        msg,
		DeliveryMode: amqp.Persistent,
	}

	fmt.Println("Sending login attempt to queue: ", message)
	err := queue.Publish("", os.Getenv("CLOUDAMQP_QUEUE"), false, false, message)

	if err != nil {
		return fmt.Errorf("error publishing login attempt message to queue: %w", err)
	}
	return nil
}

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

	return sendMessage(u.amqpQueue, loginAttempt)
}

func (u *User) sendNewRegistryMessage(id string, provider *string) error {
	if provider == nil {
		fiero := constants.InternalProvider
		provider = &fiero
	}

	queueMsg := model.QueueMessage{
		MessageType: constants.NewRegistry,
		Message: model.NewRegistry{
			RegistrationId:     id,
			Provider:   *provider,
			Timestamp:  time.Now().GoString(),
		},
	}

	newRegistry, err := json.Marshal(queueMsg)
	if err != nil {
		return fmt.Errorf("error marshalling login attempt message for rabbit: %w", err)
	}

	return sendMessage(u.amqpQueue, newRegistry)
}

func (u *User) sendUserBlockedMessage(id string, reason string) error {
	queueMsg := model.QueueMessage{
		MessageType: constants.UserBlocked,
		Message: model.UserBlocked{
			UserId:     id,
			Reason:     reason,
			Action:     constants.BlockActionBlock,
			Timestamp:  time.Now().GoString(),
		},
	}

	userBlocked, err := json.Marshal(queueMsg)
	if err != nil {
		return fmt.Errorf("error marshalling login attempt message for rabbit: %w", err)
	}

	return sendMessage(u.amqpQueue, userBlocked)
}

func (u *User) sendUserUnblockedMessage(id string) error {
	queueMsg := model.QueueMessage{
		MessageType: constants.UserUnblocked,
		Message: model.UserUnblocked{
			UserId:     id,
			Timestamp:  time.Now().GoString(),
		},
	}

	userBlocked, err := json.Marshal(queueMsg)
	if err != nil {
		return fmt.Errorf("error marshalling login attempt message for rabbit: %w", err)
	}

	return sendMessage(u.amqpQueue, userBlocked)
}

func (u *User) sendNewUserMessage(userId string, location string, oldRegistrationId string) error {
	queueMsg := model.QueueMessage{
		MessageType: constants.NewUser,
		Message: model.NewUser{
			UserId:             userId,
			Location:           location,
			OldRegistrationId:  oldRegistrationId,
		},
	}

	newUser, err := json.Marshal(queueMsg)
	if err != nil {
		return fmt.Errorf("error marshalling login attempt message for rabbit: %w", err)
	}

	return sendMessage(u.amqpQueue, newUser)
}