package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/airenas/audio-len-service/internal/pkg/audio"

	"github.com/airenas/audio-len-service/internal/pkg/file"

	"github.com/airenas/audio-len-service/internal/pkg/cmdapp"
	"github.com/airenas/audio-len-service/internal/pkg/service"
	"github.com/pkg/errors"
)

func main() {
	cFile := flag.String("c", "", "Config yml file")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:[params] \n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	err := cmdapp.InitConfig(*cFile)
	if err != nil {
		cmdapp.Log.Fatal(errors.Wrap(err, "Can't init app"))
	}
	data := service.Data{}
	data.Port = cmdapp.Config.GetInt("port")

	data.Saver, err = file.NewSaver(cmdapp.Config.GetString("tempDir"))
	if err != nil {
		cmdapp.Log.Fatal(errors.Wrap(err, "Can't init file saver"))
	}
	data.Estimator, err = audio.NewEstimator()
	if err != nil {
		cmdapp.Log.Fatal(errors.Wrap(err, "Can't init audio duration estimator"))
	}

	err = service.StartWebServer(&data)
	if err != nil {
		cmdapp.Log.Fatal(errors.Wrap(err, "Can't start the service"))
	}
}
