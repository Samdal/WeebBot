package main

import (
	"fmt"

	"./bot" //pulls bot folder
)

func main() {
	fmt.Println("Starting bot")

	bot.Start()

	<-make(chan struct{}) //keeps the program running
}
