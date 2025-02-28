package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

// The service implementation code
type TimerService struct {
	Duration       string `json:"duration"`
	RespondSubject string `json:"respond"`
}

func Set(nc *nats.Conn, duration time.Duration, respondSubject string) {
	// Create a really simple timer service go routine that sleeps, then sends a message
	go func(nc *nats.Conn) {
		time.Sleep(duration)
		nc.Publish(respondSubject, []byte("⏰"))
	}(nc)
}

// The service boilerplate
func main() {
	nc, err := nats.Connect("nats://demo.nats.io", nats.Name("TimerService"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	timerHandler := func(req micro.Request) {
		var timer TimerService
		err := json.Unmarshal(req.Data(), &timer)
		if err != nil {
			req.Respond([]byte(err.Error()))
			return
		}

		// conver the duration 1s/1m/1h into duration
		duration, err := time.ParseDuration(timer.Duration)
		if err != nil {
			req.Respond([]byte(err.Error()))
			return
		}

		Set(nc, duration, timer.RespondSubject)

		fmt.Println("Timer set for", timer.Duration, "responding to", timer.RespondSubject)

		req.Respond([]byte("ok"))
	}

	srv, err := micro.AddService(nc, micro.Config{
		Name:        "TimerService",
		Version:     "0.1.0",
		Description: "Cron timer service",
	})

	if err != nil {
		log.Fatal(err)
	}

	root := srv.AddGroup("timer")

	root.AddEndpoint("set", micro.HandlerFunc(timerHandler))

	defer srv.Stop()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	fmt.Println("bye")
}
