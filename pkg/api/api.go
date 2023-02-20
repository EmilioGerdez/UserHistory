package api

import (
	"net/http"
	"userhistory/pkg/models"

	"github.com/gin-gonic/gin"
)

func RestServer() {
	router := gin.Default()
	router.POST("/CreateNote", CreateNote)
	router.Run("localhost:8080")

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
	if err := models.CrearNota(&nota); err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Nota creada correctamente"})
	return
}
