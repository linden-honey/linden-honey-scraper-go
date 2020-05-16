package controller

import (
	"net/http"
)

type docsController struct {
	Spec string
}

func NewDocsController(spec string) *docsController {
	return &docsController{
		Spec: spec,
	}
}

func (c *docsController) GetSpec(w http.ResponseWriter, _ *http.Request) {
	WriteRawJSON([]byte(c.Spec), http.StatusOK, w)
}
