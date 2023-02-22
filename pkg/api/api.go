package api

import (
	"net/http"
	"sort"
	"strconv"
	"userhistory/pkg/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//RestServer funcion que se encarga de configurar el servidor para REST
//configura los endpoints y corre el servidor
func RestServer() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}}))
	router.POST("/CreateNote", CreateNote)
	router.POST("/UpdateNota", UpdateNota)
	router.POST("/EliminarNota", EliminarNota)
	router.GET("/TodasLasNotas", TodasLasNotas)
	router.GET("/BuscarNotas", BuscarNotas)
	router.GET("/Nota", Nota)
	router.Run("localhost:8080")
}

//UpdateNota funcion endpoint que recibe id y campos de Nota para modificar dicha Nota
//si hay algun problema con los campos responde o el id es invalido con BadRequest
//si hay algun error al hacer la operacion responde con InternalServerError
//si la modificacion fue realizada correctamente responde con StatusOK
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

//Nota funcion que recibe un ID y retorna la nota que contenga ese ID
//si el ID es invalido retorna BadRequest
//si el ID no existe retorna una nota vacia
//si hay algun error interno retorna InternalServerError
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
		return
	}

	c.JSON(http.StatusOK, nota)

}

//TodasLasNotas funcion que recibe un tipo de ordenamiento y retorna todas la notas
//ordenadas segun el tipo recibido
//si no se recibe ningun tipo o es un tipo invalido entonces retorna BadRequest
//si hay algun error interno retorna InternalServerError
//si no hay notas en la base de datos retorna una lista de notas vacia
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

//CreateNote funcion endpoint que recibe campos de Nota para crear una Nota nueva
//si hay algun problema con los campos responde con BadRequest
//si hay algun error al hacer la operacion responde con InternalServerError
//si la creacion fue realizada correctamente responde con StatusOK
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

//buscarNotas funcion busca todas las notas que en el campo buscado por "tipo" contengan el string de busqueda "q"
//si el campo es invalido o la busqueda esta vacia retorna BadRequest
//si hay algun error interno retorna BadRequest
//retorna todas las notas encontradas
//si no encuentra ninguna nota retorna una lista vacia
func BuscarNotas(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de busqueda invalido"})
		return
	}
	if busqueda == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Busqueda vacia"})
		return
	}
	var notas []models.Nota
	err := models.BuscarNotas(&notas, tipo, busqueda)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, notas)
}

//EliminarNota funcion que recibe un ID y elimina dicho registro
//si el ID es invalido retorna BadRequest
//si el ID no existe retorna NotFound
//si hay algun error interno retorna InternalServerError
//si elimina la nota correctamente retorna OK``
func EliminarNota(c *gin.Context) {
	var nota models.Nota

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id invalido"})
		return
	}
	err = models.EntregarNota(&nota, &id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if nota.Cuerpo == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nota no encontrada"})
		return
	}
	err = models.EliminarNota(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Nota eliminada correctamente"})
}
