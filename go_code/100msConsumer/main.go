package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {

	/*
		Published by David Gee, 2023 me@dave.dev

		This code is an attempt of how to create a consumer loop that can be controlled by a timer.
		Done for an experiment. Maybe some interesting things in here!

		Various ways of consuming from NATS JetStream (upper 2.10) with Go
		cons.Consume()  | blocking call requiring GR handling with an iterator stop. Has Consume Context.
		cons.Messages() | blocking call requiring GR handling. Has Messages Context.
		cons.Next() | cons.Next(jetstream.FetchMaxWait(100 * time.Millisecond)) works
	*/

	// Boilerplate
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Print(err)
	}

	// create jetstream context from nats connection
	js, err := jetstream.New(nc)
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// create new assets
	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "foo",
		Subjects: []string{"foo.>"},
		Storage:  jetstream.MemoryStorage,
		Replicas: 1,
	})
	if err != nil {
		fmt.Print(err)
	}
	// End of boilerplate

	// ====================
	// Part 1
	// No way to go lower than a second on a timeout. This consumer loop will consume from the messages using a consumer callback.
	// You can continue to publish to foo.1 and the messages will be consumed until you cc1.Stop().
	fmt.Println("Part 1")

	cons1, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       "cons1",
		FilterSubject: "foo.1",
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
	})
	if err != nil {
		fmt.Print(err)
	}

	ctxBg := context.Background()

	js.Publish(ctxBg, "foo.1", []byte("test-1.1"))
	js.Publish(ctxBg, "foo.1", []byte("test-1.2"))
	js.Publish(ctxBg, "foo.1", []byte("test-1.3"))
	js.Publish(ctxBg, "foo.1", []byte("test-1.4"))
	js.Publish(ctxBg, "foo.1", []byte("test-1.5"))

	// PullExpiry of one second.
	cc1, err := cons1.Consume(func(msg jetstream.Msg) {
		fmt.Println("1: Received jetstream message: ", string(msg.Data()))
		msg.Ack()
	}, jetstream.PullExpiry(1*time.Second)) // lowest timeout is 1 second sadly, so this is a no go

	if err != nil {
		fmt.Println(err)
	}

	defer cc1.Stop()

	// ====================
	// Part 2
	// cons2 will be used with the Messages() method
	// Messages() is a blocking call, so we run it in a GR and cancel it from a control-loop.

	time.Sleep(500 * time.Millisecond) // to ensure we can see what's going on and console messages appear in the right order

	cons2, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       "cons2",
		FilterSubject: "foo.2",
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
	})
	if err != nil {
		fmt.Print(err)
	}

	var msg jetstream.Msg

	fmt.Println("\nPart 2 - 100mS timer in Next for loop")

	// iter, err := cons2.Messages(jetstream.PullExpiry(10 * time.Millisecond)) // this didn't work?
	iter, err := cons2.Messages()
	if err != nil {
		fmt.Println("2: Messages() error: ", err)
	}

	var message = make(chan string)

	go func(iter jetstream.MessagesContext, message chan string) {
		for ; err == nil; msg, err = iter.Next() {
			if msg != nil {
				message <- "2: Received jetstream message: " + string(msg.Data())
				// fmt.Println("2: Received jetstream message: ", string(msg.Data()))

				msg.Ack()
			}
		}

		if err != nil && err.Error() != "nats: messages iterator closed" {
			fmt.Println("2: Messages() error: ", err)
		}
	}(iter, message)

	// Message sending loop
	go func() {
		for x := 1; x < 6; x++ {
			time.Sleep(time.Duration(x*30) * time.Millisecond)
			nc.Publish("foo.2", []byte("test-2."+fmt.Sprint(x)))
			//used to emulate breaking the 100mS delay
		}
	}()

	// Simple 100mS control loop to enforce the 100mS
	timer := time.NewTimer(100 * time.Millisecond) // by the time the for loop enters, there will be less than 100mS remaining.
	loop := true
	for loop {
		select {
		case <-timer.C:
			if iter != nil {
				fmt.Println("2: Stopping iterator, 100mS timer expired")
				iter.Stop()
				loop = false
			}

		case <-message:
			fmt.Println("2: Received jetstream message: ", string(msg.Data()))
			timer.Reset(100 * time.Millisecond)
		}
	}

	// ====================
	// Part 3
	// Can go lower than 1 second on a timeout without much fuss
	// Cleanest, nicest example
	time.Sleep(500 * time.Millisecond)

	fmt.Println("\nPart 3 - 100mS loop")

	cons3, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       "cons3",
		FilterSubject: "foo.3",
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
	})

	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for x := 1; x < 6; x++ {
			time.Sleep(time.Duration(x*30) * time.Millisecond)
			nc.Publish("foo.3", []byte("test-3."+fmt.Sprint(x)))
			//used to emulate breaking the 100mS delay
		}
	}()

	for {
		msg, err = cons3.Next(jetstream.FetchMaxWait(100 * time.Millisecond))

		if err == nats.ErrTimeout {
			fmt.Println("3: End of 100mS wait")
			break
		}

		fmt.Println("3: Received jetstream message: ", string(msg.Data()))
	}

	// Clean up

	// cleanup
	err = js.DeleteConsumer(ctx, "foo", "cons1")
	if err != nil && err != nats.ErrConsumerNotFound {
		// NATS ocassionally throws this error: 'nats: API error: code=404 err_code=10014 description=consumer not found'
		if !strings.Contains(err.Error(), "consumer not found") {
			fmt.Println("err1: ", err)
		}
	}

	err = js.DeleteConsumer(ctx, "foo", "cons2")
	if err != nil && err != nats.ErrConsumerNotFound {
		if !strings.Contains(err.Error(), "consumer not found") {
			fmt.Println("err2: ", err)
		}
	}

	err = js.DeleteConsumer(ctx, "foo", "cons3")
	if err != nil && err != nats.ErrConsumerNotFound {
		if !strings.Contains(err.Error(), "consumer not found") {
			fmt.Println("err3: ", err)
		}
	}

	err = js.DeleteStream(ctx, "foo")
	if err != nil && err != nats.ErrStreamNotFound {
		fmt.Println("err4: ", err)
	}
}
