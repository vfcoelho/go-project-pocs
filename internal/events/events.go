package events

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/billziss-gh/netchan/netchan"
)

type Producer[T any] struct {
	channel   chan []byte
	waitGroup *sync.WaitGroup
}

type Consumer[T any] struct {
	channel chan []byte
}

func NewProducer[T any]() *Producer[T] {
	wg := &sync.WaitGroup{}

	channel := make(chan []byte)
	wg.Add(1)

	errch := make(chan error, 1)
	err := netchan.Bind("tcp://127.0.0.1/events", channel, errch)
	if nil != err {
		panic(err)
	}

	go func() {
		for {
			err = <-errch
			fmt.Println(err.Error())
		}
	}()

	return &Producer[T]{
		channel: channel,
	}
}

func NewConsumer[T any]() *Consumer[T] {
	channel := make(chan []byte)
	err := netchan.Expose("events", channel) //listener
	if nil != err {
		panic(err)
	}
	return &Consumer[T]{channel: channel}
}

func (p *Producer[T]) Send(event T) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	p.channel <- data
	return nil
}
func (p *Producer[T]) Close() {
	close(p.channel)
	p.waitGroup.Done()
}

func (c *Consumer[T]) Consume(handlers ...Handler) error {
	for {
		data := <-c.channel
		ctx := &ConsumerCtx{message: data, handlers: handlers, values: make(map[string]any)}
		err := ctx.Next()
		if err != nil {
			return err
		}
	}
}
func (c *Consumer[T]) Close() {
	close(c.channel)
}
