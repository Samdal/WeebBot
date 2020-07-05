package main

import (
	"fmt"
	"time"

	"./bot"    //pulls bot folder
	"./config" //pulls config folder
)

func main() {
	fmt.Println("Starting bot")
	start := time.Now()

	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start(start)

	<-make(chan struct{}) //keeps the program running
}
