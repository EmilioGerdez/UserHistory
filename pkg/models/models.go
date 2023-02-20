package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Nota struct {
	gorm.Model
	Titulo string `json:"titulo,omitempty"`
	Cuerpo string `json:"cuerpo,omitempty"`
}

var DB *gorm.DB

func init() {
	Connection()
	DB.AutoMigrate(&Nota{})
}
func Connection() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Println("Error en la conexion a la base de datos: ", err)
		return err
	}
	log.Println("Conexion a la base de datos establecida correctamente")
	return nil
}

func CrearNota(nota *Nota) error {
	result := DB.Create(nota)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
