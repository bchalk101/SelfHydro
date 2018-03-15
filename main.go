package main

import (
	"os"
	"time"
	"os/signal"
	"log"
	"syscall"
)

type Time struct {
	Hh int // Hours.
	Mm int // Minutes.
	Ss int // Seconds.
}

func main() {
	f, err := os.OpenFile("/selfhydro/selfhydro-logs", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	log.SetOutput(f)

	log.Println("Starting up SelfHydro")

	controller := NewRaspberryPi()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		s := <-sigs
		if s != syscall.SIGPIPE {
			log.Println("Exiting program...")
			log.Println("RECEIVED SIGNAL: ", s)

			controller.StopSystem()

			os.Exit(0)
		}

	}()
	controller.StartHydroponics()

	for {
		time.Sleep(time.Second)
	}

}

func handleExit(controller *RaspberryPi) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		s := <-sigs
		log.Println("Exiting program...")
		log.Println("RECEIVED SIGNAL: ", s)

		controller.StopSystem()

		os.Exit(0)
	}()
}
