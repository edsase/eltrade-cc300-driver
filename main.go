//go:generate pkger

package main

import (
	"context"
	"fmt"
	"github.com/geeckmc/eltrade-cc300-driver/server"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/kardianos/service"
	"log"
	"net/http"
	"os"
	"time"
)

var logger service.Logger
var httpServer *http.Server

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	httpServer = server.Serve()
}
func (p *program) Stop(s service.Service) error {
	if httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
	}
	return nil
}

func main() {
	loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
	svcConfig := &service.Config{
		Name:        "eltradeCC300Driver",
		DisplayName: "Eltrade CC330 Drvier",
		Description: "Driver for Eltrade Tax control device",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.SystemLogger(nil)
	if err != nil {
		log.Fatal(err)
	}

	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "install":
			s.Install()
			s.Start()
			return
		case "uninstall":
			fmt.Printf("uninstall")
			s.Stop()
			s.Uninstall()
			return
		}
	}
	s.Run()
}
