package main

import (
	"github.com/airenas/audio-len-service/internal/pkg/audio"

	"github.com/airenas/audio-len-service/internal/pkg/file"

	"github.com/airenas/audio-len-service/internal/pkg/service"
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
)

func main() {
	goapp.StartWithDefault()

	data := service.Data{}
	data.Port = goapp.Config.GetInt("port")

	var err error
	data.Saver, err = file.NewSaver(goapp.Config.GetString("tempDir"))
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "Can't init file saver"))
	}
	data.Estimator, err = audio.NewEstimator()
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "Can't init audio duration estimator"))
	}

	err = service.StartWebServer(&data)
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "Can't start the service"))
	}
}
