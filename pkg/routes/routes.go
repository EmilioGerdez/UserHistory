package routes

import (
	"net/http"
	"sort"
	"strconv"
	"userhistory/pkg/models"

	"github.com/gin-gonic/gin"
)

const ContentTypeHTML = "text/html; charset=utf-8"

func Server() {
	router := gin.Default()
	router.GET("/", index)
	router.GET("/CrearNota", crearNota)
	router.GET("/NuevaNota", NuevaNota)
	router.GET("/Notas", lasNotas)
	router.GET("/Buscar", buscarNotas)
	router.GET("/Modificar", Modificar)
	router.GET("/ModificarNota", ModificarNota)
	router.GET("/Eliminar", Eliminar)
	router.LoadHTMLGlob("../../src/templates/*")
	router.Static("/assets", "../../src/assets")
	router.Run("localhost:9090")
}
func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func crearNota(c *gin.Context) {
	c.HTML(http.StatusOK, "crearNota.html", gin.H{})
}

func ModificarNota(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "ID invalido"})
		return
	}
	var nota models.Nota
	err = models.EntregarNota(&nota, &id)
	if err != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{"error": err.Error()})
		return

	}
	c.HTML(http.StatusOK, "detalles.html", gin.H{"Nota": nota})
	return
}
func Modificar(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "ID invalido"})
		return
	}
	var nota models.Nota
	err = models.EntregarNota(&nota, &id)
	if err != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{"error": err.Error()})
		return
	}
	titulo := c.Query("title")
	tema := c.Query("theme")
	cuerpo := c.Query("cuerpo")
	if titulo == "" || tema == "" || cuerpo == "" {
		c.HTML(http.StatusBadRequest, "detalles.html", gin.H{"error": "por favor no deje celdas vacias", "Nota": nota})
		return
	}
	nota.Titulo = titulo
	nota.Tema = tema
	nota.Cuerpo = cuerpo
	err = models.ModificarNota(&nota, &id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "detalles.html", gin.H{"message": "nota modificada correctamente", "Nota": nota})
	return
}

func NuevaNota(c *gin.Context) {
	titulo := c.Query("title")
	tema := c.Query("theme")
	cuerpo := c.Query("cuerpo")
	if titulo == "" || tema == "" || cuerpo == "" {
		c.HTML(http.StatusBadRequest, "crearNota.html", gin.H{"error": "por favor no deje celdas vacias"})
		return
	}
	nota := models.Nota{
		Titulo: titulo,
		Tema:   tema,
		Cuerpo: cuerpo,
	}
	if err := models.CrearNota(&nota); err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "detalles.html", gin.H{"Nota": nota})
	return

}

func buscarNotas(c *gin.Context) {
	busqueda := c.Query("q")
	var tipo string
	switch c.Query("sort") {
	case "T":
		tipo = "Titulo"
	case "t":
		tipo = "Tema"
	case "c":
		tipo = "Cuerpo"
	default:
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "Tipo invalido"})
		return
	}
	if busqueda == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "Busqueda vacia"})
		return
	}
	var notas []models.Nota
	err := models.BuscarNotas(&notas, tipo, busqueda)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error()})
		return
	}

	if len(notas) == 0 {
		c.HTML(http.StatusNotFound, "index.html", gin.H{"error": " notas no encontradas"})
		return
	}
	c.HTML(http.StatusOK, "Notas.html", gin.H{"Notas": notas})
	return
}

func lasNotas(c *gin.Context) {
	orden := c.Query("sort")
	var notas []models.Nota
	models.EntregarNotas(&notas)
	switch orden {
	case "T":
		sort.Slice(
			notas, func(i, j int) bool {
				return notas[i].Titulo < notas[j].Titulo
			},
		)

	case "t":
		sort.Slice(
			notas, func(i, j int) bool {
				return notas[i].Tema < notas[j].Tema
			},
		)

	case "d":
		sort.Slice(
			notas, func(i, j int) bool {
				return notas[i].CreatedAt.Unix() < notas[j].CreatedAt.Unix()
			},
		)
	default:
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "Orden invalido"})
		return
	}
	if len(notas) == 0 {
		c.HTML(http.StatusNotFound, "index.html", gin.H{"error": " notas no encontradas"})
		return
	}
	c.HTML(http.StatusOK, "Notas.html", gin.H{"Notas": notas})
	return
}

func Eliminar(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "ID invalido"})
		return
	}
	var nota models.Nota
	err = models.EntregarNota(&nota, &id)
	if err != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{"error": err.Error()})
		return
	}
	err = models.EliminarNota(&id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{"message": "nota eliminada correctamente"})
	return
}
