package controller

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func WriteRawJSON(data []byte, status int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		return errors.Wrap(err, "Error happened during response writing")
	}
	return nil
}

func WriteJSON(data interface{}, status int, w http.ResponseWriter) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "Error happened during data marshalling")
	}
	err = WriteRawJSON(dataBytes, status, w)
	if err != nil {
		return errors.Wrap(err, "Error happened during raw json writing")
	}
	return nil
}

func WriteError(err error, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	_, _ = w.Write(([]byte)(err.Error()))
}
