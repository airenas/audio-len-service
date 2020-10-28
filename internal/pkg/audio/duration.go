package audio

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

//Estimator estimates audio duration
type Estimator struct {
	estFunc func(string) (float64, error)
}

//NewEstimator return new estimator instance
func NewEstimator() (*Estimator, error) {
	res := Estimator{}
	res.estFunc = estimate
	return &res, nil
}

//Seconds returns file length in secods
func (e *Estimator) Seconds(name string) (float64, error) {
	ext := filepath.Ext(name)
	if ext == ".wav" {
		return e.estFunc(fmt.Sprintf("sox --i -D %s", name))
	}
	return e.estFunc(fmt.Sprintf("ffprobe -i %s -show_entries format=duration -v quiet -of csv=p=0", name))
}

func estimate(cmdParams string) (float64, error) {
	cmdArr := strings.Split(cmdParams, " ")
	if len(cmdArr) < 2 {
		return 0, errors.New("Wrong command. No parameter " + cmdParams)
	}

	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer
	err := cmd.Run()
	if err != nil {
		return 0, errors.Wrap(err, "Output: "+string(outputBuffer.Bytes()))
	}
	res := strings.TrimSpace(string(outputBuffer.Bytes()))
	return strconv.ParseFloat(res, 64)
}
