package main

import (
	"flag"
	"fmt"
	"github.com/kardianos/service"
	log "github.com/kdar/factorlog"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
	"os"
)

const (
	SERVICE_INSTALL   = "install"
	SERVICE_UNINSTALL = "uninstall"
	SERVICE_START     = "start"
	SERVICE_STOP      = "stop"
	SERVICE_RUN       = "run"
)

var (
	logFlag  = flag.String("log", "./cubes_service.log", "set log path")
	portFlag = flag.Int("port", 9100, "set port")

	logger *log.FactorLog
)

type ServiceRunner struct {
}

func NewServiceRunner() *ServiceRunner {
	return &ServiceRunner{}
}

func (r *ServiceRunner) Start(s service.Service) error {
	go r.run()
	logger.Infof("ServiceRunner started.")
	return nil
}
func (r *ServiceRunner) run() {
	logger.Infof("ServiceRunner running.")
	run_gin_server()
}
func (r *ServiceRunner) Stop(s service.Service) error {
	logger.Infof("CubesRunner stopped.")
	return nil
}

func main() {
	flag.Parse()
	logger = utils.SetGlobalLogger(*logFlag)

	logger.Infof("Entering cubes windows service: args = %s, len(args) = %d", utils.Json(os.Args), len(os.Args))

	svcConfig := &service.Config{
		Name:        "CubesService",
		DisplayName: "Cubes Web Server",
		Description: "Windows Service for Cubes Web Server, used for data reports.",
	}

	runner := NewServiceRunner()
	s, err := service.New(runner, svcConfig)
	if err != nil {
		logger.Fatal(err)
	}

	if len(os.Args) > 1 {
		var err error
		verb := os.Args[1]
		switch verb {
		case SERVICE_INSTALL:
			err = s.Install()
			if err != nil {
				fmt.Printf("Failed to install: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" installed.\n", svcConfig.Name)

		case SERVICE_UNINSTALL:
			err = s.Uninstall()
			if err != nil {
				fmt.Printf("Failed to remove: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" removed.\n", svcConfig.Name)

		case SERVICE_RUN:
			s.Run()

		case SERVICE_START:
			err = s.Start()
			if err != nil {
				fmt.Printf("Failed to start: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" started.\n", svcConfig.Name)

		case SERVICE_STOP:
			err = s.Stop()
			if err != nil {
				fmt.Printf("Failed to stop: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" stopped.\n", svcConfig.Name)
		}
		return
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
