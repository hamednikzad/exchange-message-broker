package main

import (
	"log"

	"github.com/asaskevich/EventBus"
)

func calculator(a int, b int) {
	log.Printf("%d\n", a+b)
}

func main2() {
	bus := EventBus.New()
	bus.Subscribe("main:calculator", calculator)

	bus.Publish("main:calculator", 20, 40)
	bus.Unsubscribe("main:calculator", calculator)
}
