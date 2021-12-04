package main

import (
	"github.com/sjurtf/tibber_exporter/cmd"
	"github.com/sjurtf/tibber_exporter/tibber"
	"log"
	"os"
)

func main() {

	tibberAccessToken := os.Getenv("TIBBER_ACCESS_TOKEN")
	tibberHomeId := os.Getenv("TIBBER_HOME_ID")

	t, err := tibber.NewTibber(tibberAccessToken, tibberHomeId)
	if err != nil {
		log.Fatalln("unable to start tibber")
	}

	go t.Measurements()

	e := cmd.NewExporter(t)
	e.StartGather()
	e.Listen()

}
