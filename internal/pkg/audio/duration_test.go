package audio

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var estimator *Estimator

func initTest(t *testing.T) {
	var err error
	estimator, err = NewEstimator()
	assert.Nil(t, err)
}

func TestFile_Wav(t *testing.T) {
	initTest(t)
	var s string
	estimator.estFunc = func(cmd string) (float64, error) {
		s = cmd
		return 10, nil
	}
	d, err := estimator.Seconds("file.wav")
	assert.Equal(t, "sox --i -D file.wav", s)
	assert.Equal(t, float64(10), d)
	assert.Nil(t, err)
}

func TestFile_Mp3(t *testing.T) {
	initTest(t)
	var s string
	estimator.estFunc = func(cmd string) (float64, error) {
		s = cmd
		return 11, nil
	}
	d, err := estimator.Seconds("file.mp3")
	assert.Equal(t, "ffprobe -i file.mp3 -show_entries format=duration -v quiet -of csv=p=0", s)
	assert.Equal(t, float64(11), d)
	assert.Nil(t, err)
}

func TestFile_Fail(t *testing.T) {
	initTest(t)
	estimator.estFunc = func(cmd string) (float64, error) {
		return 0, errors.New("olia")
	}
	_, err := estimator.Seconds("file.mp3")
	assert.NotNil(t, err)
}

func TestEstimate(t *testing.T) {
	d, err := estimate("echo 2.0")
	assert.Equal(t, float64(2), d)
	assert.Nil(t, err)
}

func TestEstimate_Fails(t *testing.T) {
	_, err := estimate("xxxxecho 2.0")
	assert.NotNil(t, err)
}
