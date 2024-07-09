package main

import (
	"log"

	"github.com/BrunoTulio/goscanner/bootstrap"
	"github.com/BrunoTulio/goscanner/handler"
	"github.com/BrunoTulio/goscanner/scanner"
	"github.com/BrunoTulio/goscanner/server"
	"github.com/BrunoTulio/goscanner/ui"
)

func main() {

	server := server.NewServerHTTP(bootstrap.PortDefault)
	scanner, err := scanner.NewScannerRuntime()

	if err != nil {
		log.Fatal(err)
	}
	myApp := ui.New(server, scanner)
	handler.Setup(myApp)

	myApp.Run()
	scanner.Close()
}
