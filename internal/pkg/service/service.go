package service

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	// FileSaver saves the file with the provided name
	FileSaver interface {
		Save(name string, reader io.Reader) (string, error)
	}

	// DurationEstimator estimates file duration
	DurationEstimator interface {
		Seconds(name string) (float64, error)
	}

	//Data is service operation data
	Data struct {
		Port int

		Saver     FileSaver
		Estimator DurationEstimator
	}
)

//StartWebServer starts the HTTP service and listens for the admin requests
func StartWebServer(data *Data) error {
	goapp.Log.Infof("Starting HTTP audio len service at %d", data.Port)
	r := NewRouter(data)
	http.Handle("/", r)
	portStr := strconv.Itoa(data.Port)
	err := http.ListenAndServe(":"+portStr, nil)

	if err != nil {
		return errors.Wrap(err, "Can't start HTTP listener at port "+portStr)
	}
	return nil
}

//NewRouter creates the router for HTTP service
func NewRouter(data *Data) *mux.Router {
	router := mux.NewRouter()
	router.Methods("POST").Path("/duration").Handler(&durationHandler{data: data})
	return router
}

type durationHandler struct {
	data *Data
}

type durationResult struct {
	Duration float64 `json:"duration"`
}

func (h *durationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	goapp.Log.Debugf("Request from %s", r.RemoteAddr)

	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file", http.StatusBadRequest)
		goapp.Log.Error(err)
		return
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	ext = strings.ToLower(ext)
	if !checkFileExtension(ext) {
		http.Error(w, "Wrong file extension: "+ext, http.StatusBadRequest)
		goapp.Log.Errorf("Wrong file extension: " + ext)
		return
	}

	id := uuid.New().String()
	fileName := id + ext

	fileName, err = h.data.Saver.Save(fileName, file)
	if err != nil {
		http.Error(w, "Can not save file", http.StatusInternalServerError)
		goapp.Log.Error(err)
		return
	}
	defer deleteFile(fileName)

	res := durationResult{}
	res.Duration, err = h.data.Estimator.Seconds(fileName)
	if err != nil {
		http.Error(w, "Can not extract duration", http.StatusInternalServerError)
		goapp.Log.Error(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&res)
	if err != nil {
		http.Error(w, "Can not prepare result", http.StatusInternalServerError)
		logrus.Error(err)
	}
}

func checkFileExtension(ext string) bool {
	return ext == ".wav" || ext == ".mp3" || ext == ".mp4" || ext == ".m4a"
}

func deleteFile(file string) {
	os.RemoveAll(file)
}
