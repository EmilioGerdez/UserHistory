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

func autoRequest(r *gin.Engine, method, url, expect string, body []byte, checkData bool) error {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	if checkData {
		if string(responseData) != expect {
			return errors.New(fmt.Sprintf("Esperado %s, Obtenido %s", expect, string(responseData)))
		}
	}
	if http.StatusOK != w.Code {
		return errors.New(fmt.Sprintf("Esperado %d, Obtenido  %d", http.StatusOK, w.Code))
	}
	return nil
}

func TestCreateNote(t *testing.T) {
	var jsonNote = []byte(`{"cuerpo":"Nota creada a traves de la API","titulo":"Nota de Prueba", "tema":"Principal"}`)
	var jsonVoidNote = []byte(`{"cuerpo":"","titulo":"Nota de Prueba2", "tema":"Principal"}`)
	var jsonVoidTitleNote = []byte(`{"cuerpo":"Nota creada a traves de la API","titulo":"", "tema":"Principal"}`)
	var jsonVoidThemeNote = []byte(`{"cuerpo":"Nota creada a traves de la API","titulo":"Nota de Prueba3", "tema":""}`)
	r := gin.Default()
	r.POST("/CreateNote", api.CreateNote)
	if err := autoRequest(r, "POST", "/CreateNote", notaExitosa, jsonNote, true); err != nil {
		t.Errorf("Error: %v", err)
	}
	//this should fail
	if err := autoRequest(r, "POST", "/CreateNote", notaExitosa, jsonVoidNote, true); err == nil {
		t.Errorf("Error expected")
	}
	if err := autoRequest(r, "POST", "/CreateNote", notaExitosa, jsonVoidTitleNote, true); err == nil {
		t.Errorf("Error expected")
	}
	if err := autoRequest(r, "POST", "/CreateNote", notaExitosa, jsonVoidThemeNote, true); err == nil {
		t.Errorf("Error expected")
	}
}

func TestTodasLasNotas(t *testing.T) {
	var json = []byte(``)
	r := gin.Default()
	r.GET("/TodasLasNotas", api.TodasLasNotas)
	if err := autoRequest(r, "GET", "/TodasLasNotas?sort=T", notaExitosa, json, false); err != nil {
		t.Errorf("Error: %v", err)
	}
	if err := autoRequest(r, "GET", "/TodasLasNotas?sort=t", notaExitosa, json, false); err != nil {
		t.Errorf("Error: %v", err)
	}
	if err := autoRequest(r, "GET", "/TodasLasNotas?sort=d", notaExitosa, json, false); err != nil {
		t.Errorf("Error: %v", err)
	}
	//this should fail
	if err := autoRequest(r, "GET", "/TodasLasNotas?", notaExitosa, json, false); err == nil {
		t.Errorf("Error expected")
	}
	if err := autoRequest(r, "GET", "/TodasLasNotas?sort=1", notaExitosa, json, false); err == nil {
		t.Errorf("Error expected")
	}
}

func TestNota(t *testing.T) {
	var json = []byte(``)
	r := gin.Default()
	r.GET("/Nota", api.Nota)

	if err := autoRequest(r, "GET", "/Nota?id=0", notaExitosa, json, false); err != nil {
		t.Errorf("Error: %v", err)
	}
	//this should fail
	if err := autoRequest(r, "GET", "/Nota?id=string", notaExitosa, json, false); err == nil {
		t.Errorf("Error expected")
	}
	if err := autoRequest(r, "GET", "/Nota", notaExitosa, json, false); err == nil {
		t.Errorf("Error expected")
	}
}
