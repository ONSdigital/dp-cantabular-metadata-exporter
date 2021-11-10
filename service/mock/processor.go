// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/service"
	kafka "github.com/ONSdigital/dp-kafka/v2"
	"sync"
)

// Ensure, that ProcessorMock does implement service.Processor.
// If this is not the case, regenerate this file with moq.
var _ service.Processor = &ProcessorMock{}

// ProcessorMock is a mock implementation of service.Processor.
//
// 	func TestSomethingThatUsesProcessor(t *testing.T) {
//
// 		// make and configure a mocked service.Processor
// 		mockedProcessor := &ProcessorMock{
// 			ConsumeFunc: func(contextMoqParam context.Context, iConsumerGroup kafka.IConsumerGroup, handler event.Handler)  {
// 				panic("mock out the Consume method")
// 			},
// 		}
//
// 		// use mockedProcessor in code that requires service.Processor
// 		// and then make assertions.
//
// 	}
type ProcessorMock struct {
	// ConsumeFunc mocks the Consume method.
	ConsumeFunc func(contextMoqParam context.Context, iConsumerGroup kafka.IConsumerGroup, handler event.Handler)

	// calls tracks calls to the methods.
	calls struct {
		// Consume holds details about calls to the Consume method.
		Consume []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// IConsumerGroup is the iConsumerGroup argument value.
			IConsumerGroup kafka.IConsumerGroup
			// Handler is the handler argument value.
			Handler event.Handler
		}
	}
	lockConsume sync.RWMutex
}

// Consume calls ConsumeFunc.
func (mock *ProcessorMock) Consume(contextMoqParam context.Context, iConsumerGroup kafka.IConsumerGroup, handler event.Handler) {
	if mock.ConsumeFunc == nil {
		panic("ProcessorMock.ConsumeFunc: method is nil but Processor.Consume was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		IConsumerGroup  kafka.IConsumerGroup
		Handler         event.Handler
	}{
		ContextMoqParam: contextMoqParam,
		IConsumerGroup:  iConsumerGroup,
		Handler:         handler,
	}
	mock.lockConsume.Lock()
	mock.calls.Consume = append(mock.calls.Consume, callInfo)
	mock.lockConsume.Unlock()
	mock.ConsumeFunc(contextMoqParam, iConsumerGroup, handler)
}

// ConsumeCalls gets all the calls that were made to Consume.
// Check the length with:
//     len(mockedProcessor.ConsumeCalls())
func (mock *ProcessorMock) ConsumeCalls() []struct {
	ContextMoqParam context.Context
	IConsumerGroup  kafka.IConsumerGroup
	Handler         event.Handler
} {
	var calls []struct {
		ContextMoqParam context.Context
		IConsumerGroup  kafka.IConsumerGroup
		Handler         event.Handler
	}
	mock.lockConsume.RLock()
	calls = mock.calls.Consume
	mock.lockConsume.RUnlock()
	return calls
}
