package controller

import (
	"net/http"
)

type DocsController struct {
	Spec string
}

func (c *DocsController) GetSpec(w http.ResponseWriter, _ *http.Request) {
	WriteRawJSON([]byte(c.Spec), http.StatusOK, w)
}
