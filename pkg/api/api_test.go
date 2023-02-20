package api_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"userhistory/pkg/api"

	"github.com/gin-gonic/gin"
)

const notaExitosa = `{"message":"Nota creada correctamente"}`

func autoRequest(r *gin.Engine, method, url, expect string, body []byte) error {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	if string(responseData) != expect {
		return errors.New(fmt.Sprintf("Esperado %s, Obtenido %s", expect, string(responseData)))
	}
	if http.StatusOK != w.Code {
		return errors.New(fmt.Sprintf("Esperado %d, Obtenido  %d", http.StatusOK, w.Code))
	}
	return nil
}

func TestCreateNote(t *testing.T) {
	var jsonNote = []byte(`{"cuerpo":"Nota creada a traves de la API","titulo":"Nota de Prueba"}`)
	var jsonVoidNote = []byte(`{"cuerpo":"","titulo":"Nota de Prueba2"}`)
	var jsonVoidTitleNote = []byte(`{"cuerpo":"Nota creada a traves de la API","titulo":""}`)
	r := gin.Default()
	r.POST("/CreateNote", api.CreateNote)
	if err := autoRequest(r, "POST", "/CreateNote", notaExitosa, jsonNote); err != nil {
		t.Errorf("Error: %v", err)
	}
	//this should fail
	if err := autoRequest(r, "POST", "/CreateNote", notaExitosa, jsonVoidNote); err == nil {
		t.Errorf("Error expected")
	}
	if err := autoRequest(r, "POST", "/CreateNote", notaExitosa, jsonVoidTitleNote); err == nil {
		t.Errorf("Error expected")
	}
}
