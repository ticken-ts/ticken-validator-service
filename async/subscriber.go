package async

import (
	"encoding/json"
	"fmt"
	"ticken-validator-service/infra"
	"ticken-validator-service/infra/bus"
	"ticken-validator-service/services"
)

type Subscriber struct {
	busSubscriber infra.BusSubscriber

	/********+*** processors ************/
	eventProcessor     *EventSubscriber
	attendantProcessor *AttendantSubscriber
	/************************************/
}

func NewSubscriber(busSubscriber infra.BusSubscriber, serviceProvider services.IProvider) (*Subscriber, error) {
	if !busSubscriber.IsConnected() {
		return nil, fmt.Errorf("bus subscriber is not connected")
	}

	subscriber := &Subscriber{
		busSubscriber:      busSubscriber,
		eventProcessor:     NewEventSubscriber(serviceProvider.GetEventManager()),
		attendantProcessor: NewAttendantSubscriber(serviceProvider.GetAttendantManager()),
	}

	return subscriber, nil
}

func (processor *Subscriber) Start() error {
	err := processor.busSubscriber.Listen(processor.handler)
	if err != nil {
		return err
	}
	return nil
}

func (processor *Subscriber) handler(rawmsg []byte) {
	msg := new(bus.Message)
	err := json.Unmarshal(rawmsg, msg)
	if err != nil {
		println("error processing message")
	}

	var processingError error = nil
	switch msg.Type {
	case NewEventMessageType:
		processingError = processor.eventProcessor.NewEventHandler(msg.Data)
	case NewAttendantMessageType:
		processingError = processor.attendantProcessor.NewAttendantHandler(msg.Data)

	default:
		processingError = fmt.Errorf("message type %s not supportaed\n", msg.Type)
	}

	if processingError != nil {
		fmt.Println(processingError)
	}
}
