package service

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

// StartWebServer starts the HTTP service and listens for the admin requests
func StartWebServer(data *Data) error {
	goapp.Log.Info().Int("port", data.Port).Msg("Starting HTTP audio len service")
	r := NewRouter(data)
	http.Handle("/", r)
	portStr := strconv.Itoa(data.Port)
	err := http.ListenAndServe(":"+portStr, nil)

	if err != nil {
		return fmt.Errorf("can't start HTTP listener at port %s: %w", portStr, err)
	}
	return nil
}

// NewRouter creates the router for HTTP service
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
	goapp.Log.Debug().Str("remote", r.RemoteAddr).Msg("Request")

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Can't parse form data", http.StatusBadRequest)
		goapp.Log.Error().Err(err).Send()
		return
	}
	defer cleanFiles(r.MultipartForm)
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file", http.StatusBadRequest)
		goapp.Log.Error().Err(err).Send()
		return
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	ext = strings.ToLower(ext)

	id := uuid.New().String()
	fileName := id + ext

	fileName, err = h.data.Saver.Save(fileName, file)
	if err != nil {
		http.Error(w, "Can not save file", http.StatusInternalServerError)
		goapp.Log.Error().Err(err).Send()
		return
	}
	defer deleteFile(fileName)

	res := durationResult{}
	res.Duration, err = h.data.Estimator.Seconds(fileName)
	if err != nil {
		http.Error(w, "Can not extract duration", http.StatusInternalServerError)
		goapp.Log.Error().Err(err).Send()
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&res)
	if err != nil {
		http.Error(w, "Can not prepare result", http.StatusInternalServerError)
		goapp.Log.Error().Err(err).Send()
	}
}

func deleteFile(file string) {
	os.RemoveAll(file)
}

func cleanFiles(f *multipart.Form) {
	if f != nil {
		f.RemoveAll()
	}
}
