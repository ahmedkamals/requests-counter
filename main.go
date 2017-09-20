package main

import (
	"flag"
	"io/ioutil"
	"os"
	"requests-counter/app"
	"requests-counter/config"
	"requests-counter/services"
)

func main() {
	serviceLocator := services.NewLocator()
	logger := serviceLocator.Logger()

	configPath := flag.String("config.path", "config/main.json", "Base config.")
	flag.Parse()

	data, err := ioutil.ReadFile(*configPath)

	if nil != err {
		logger.Error(err.Error())
		os.Exit(0)
	}

	config.Configuration, err = serviceLocator.LoadConfig(data)

	if err != nil {
		logger.Critical(err.Error())
		os.Exit(0)
	}

	server := app.NewServer(app.NewDispatcher(app.NewSaver(config.Configuration.Dispatcher.Backup)).Run())
	server.Start()
	serviceLocator.BlockIndefinitely()
	server.Stop()

	logger.Success([]string{`
 _     _                         _                _                     _           _           _
| |   | |           _           | |              (_)     _             | |         | |         | |
| |__ | | ____  ___| |_  ____   | | ____    _   _ _  ___| |_  ____     | | _   ____| | _  _   _| |
|  __)| |/ _  |/___)  _)/ _  |  | |/ _  |  | | | | |/___)  _)/ _  |    | || \ / _  | || \| | | |_|
| |   | ( ( | |___ | |_( ( | |  | ( ( | |   \ V /| |___ | |_( ( | |_   | |_) | ( | | |_) ) |_| |_
|_|   |_|\_||_(___/ \___)_||_|  |_|\_||_|    \_/ |_(___/ \___)_||_( )  |____/ \_||_|____/ \__  |_|
                                                                  |/                     (____/
	`}[0])
}
