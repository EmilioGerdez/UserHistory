package api

import (
	"net/http"
	"sort"
	"strconv"
	"userhistory/pkg/models"

	"github.com/gin-gonic/gin"
)

func RestServer() {
	router := gin.Default()
	router.POST("/CreateNote", CreateNote)
	router.POST("/UpdateNota", UpdateNota)
	router.GET("/TodasLasNotas", TodasLasNotas)
	router.GET("/Nota", Nota)
	router.Run("localhost:8080")
}

func UpdateNota(c *gin.Context) {
	var nota models.Nota
	if err := c.ShouldBindJSON(&nota); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if nota.Cuerpo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nota vacia"})
		return
	}

	if nota.Titulo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nota sin Titulo"})
		return
	}

	if nota.Tema == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nota sin Tema"})
		return
	}

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id invalido"})
		return
	}
	if err := models.ModificarNota(&nota, &id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Nota modificada correctamente"})
	return
}
func Nota(c *gin.Context) {
	var nota models.Nota

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id invalido"})
		return
	}

	err = models.EntregarNota(&nota, &id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, nota)

}

func TodasLasNotas(c *gin.Context) {
	var notas []models.Nota
	orden := c.Query("sort")
	if orden != "T" && orden != "t" && orden != "d" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Orden invalido"})
		return
	}
	err := models.EntregarNotas(&notas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
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
	}
	c.JSON(http.StatusOK, notas)
}

func CreateNote(c *gin.Context) {
	var nota models.Nota
	if err := c.ShouldBindJSON(&nota); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if nota.Cuerpo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nota vacia"})
		return
	}

	if nota.Titulo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nota sin Titulo"})
		return
	}

	if nota.Tema == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nota sin Tema"})
		return
	}
	if err := models.CrearNota(&nota); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Nota creada correctamente"})
	return
}
