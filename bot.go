package main

import (
	"fmt"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/sbaildon/quincy/handlers"

	_ "github.com/sbaildon/quincy/commands/hello"
	_ "github.com/sbaildon/quincy/commands/role"
	_ "github.com/sbaildon/quincy/commands/help"
)

func main() {
	var (
		Token = flag.String("token", "", "Discord auth token")
	)
	flag.Parse()

	dg, err := discordgo.New(*Token)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.SetupHandlers(dg)

	fmt.Println("Opening websocket connection")
	err = dg.Open()
	if err != nil {
		fmt.Println("Unable to open discord connection")
		return
	}

	time.Sleep(1000 * time.Millisecond)

	attempts := 0
	for !dg.DataReady {
		attempts++
		if attempts == 10 {
			fmt.Println("Timeout exceeded")
			return
		}
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("Waiting for connection to be established")
	}

	fmt.Println("Connected")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
