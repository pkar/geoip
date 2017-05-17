package main

import (
	"flag"
	"log"
	"os"

	"github.com/pkar/geoip"
)

func main() {
	log.New(os.Stdout, "", log.Lshortfile|log.Ldate|log.Ltime|log.Lmicroseconds)
	listen := flag.String("listen", "localhost:8899", "The host:port interface to listen on")
	location := flag.String("location", "GeoLiteCity-Location.csv", "The path to the GeoLiteCity-Location.csv file")
	blocks := flag.String("blocks", "GeoLiteCity-Blocks.csv", "The path to the GeoLiteCity-Blocks.csv file")
	flag.Parse()
	err := geoip.Run(*listen, *location, *blocks)
	if err != nil {
		log.Fatal(err)
	}
}
