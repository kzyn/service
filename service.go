package main

import (
	"log"
	"os"
	"flag"
	"path/filepath"
	"encoding/json"

	"github.com/kardianos/service"
)

var logger service.Logger

type Config struct {
	Name, DisplayName, Description string
/*
	Dir  string
	Exec string
	Args []string
	Env  []string

	Stderr, Stdout string
*/
}

type program struct {
}
func main() {
	var mode string

	flag.StringVar(&mode, "mode", "", "Control the system service.")
	flag.Parse()

	configPath, err := getConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	config, err := getConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	svcConfig := &service.Config{
		Name:        config.Name,
		DisplayName: config.DisplayName,
		Description: config.Description,
	}

	prg := &program{}
	srv, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = srv.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	switch mode {
	case "install":
		err := srv.Install()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	case "remove":
		err := srv.Uninstall()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	case "stop":
		srv.Stop()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	case "run":
		srv.Run()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	}
}

func getConfigPath() (string, error) {
	fullExecPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, execName := filepath.Split(fullExecPath)
	ext := filepath.Ext(execName )
	name := execName [:len(execName )-len(ext)]

	return filepath.Join(dir, name+".json"), nil
}

func getConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &Config{}

	r := json.NewDecoder(f)
	err = r.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

