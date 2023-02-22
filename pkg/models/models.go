package models

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//modelo de Nota, compuesto por gorm.Model para valores basicos de una tabla
type Nota struct {
	gorm.Model
	Titulo string `json:"titulo,omitempty"`
	Cuerpo string `json:"cuerpo,omitempty"`
	Tema   string `json:"tema,omitempty"`
}

var DB *gorm.DB

//Funcion init se ejecuta cada vez que se importa el package models
//establece la conexion con la base de datos
//migra las tablas de ser necesario
func init() {
	Connection()
	DB.AutoMigrate(&Nota{})
}

//Connection funcion que establece conexion con la base de datos
//retorna error si no logra hacer la conexion
func Connection() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("../../src/gorm.db"), &gorm.Config{})
	if err != nil {
		log.Println("Error en la conexion a la base de datos: ", err)
		return err
	}
	log.Println("Conexion a la base de datos establecida correctamente")
	return nil
}

//CrearNota funcion que crea un nuevo registro
//si no logra crearlo retorna error
func CrearNota(nota *Nota) error {
	result := DB.Create(nota)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

//ModificarNota funcion nque modifica un registr existente
//retorna error si no lo consigue
//retorna error si no puede guardar el cambio
func ModificarNota(nota *Nota, id *int) error {
	var note Nota
	result := DB.First(&note, id)
	if result.Error != nil {
		return result.Error
	}
	note.Cuerpo = nota.Cuerpo
	note.Tema = nota.Tema
	note.Titulo = nota.Titulo
	result = DB.Save(&note)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

//EntregarNota funcion que busca una nota con su ID
//puede retornar nota vacia si no la encunetra
//retorna error si no puede hacer la busqueda
func EntregarNota(nota *Nota, id *int) error {
	result := DB.First(nota, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//EntregarNotas funcion que retorna todas las notas de la base de datos
//puede entregar una lista vacia
//retorna error si no puede hacer la busqueda
func EntregarNotas(notas *[]Nota) error {
	result := DB.Find(notas)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//BuscarNotas funcion que busca entre todas las notas
//las busquedas se hacen en un campo especifico
//dicho campo debe tener el string de busqueda en alguna parte del mismo
//puede retornar una lista vacia
//retorna error si no logra hacer la busqueda
func BuscarNotas(notas *[]Nota, tipo, busqueda string) error {
	result := DB.Where(fmt.Sprintf("%v LIKE ?", tipo), fmt.Sprintf("%%%v%%", busqueda)).Find(&notas)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//EliminarNota funcion que elimina una nota usando su ID
//retorna error si no lo logra
func EliminarNota(id *int) error {
	result := DB.Delete(&Nota{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
