package log

import (
	"encoding/json"
	"os"
	"time"
)

type errline struct {
	Time string
	Type string

	Message string
}

var e = json.NewEncoder(os.Stdout)

func Err(err error) {
	e.Encode(errline{
		Time:    time.Now().Format(time.RFC3339Nano),
		Type:    "error",
		Message: err.Error(),
	})
}
