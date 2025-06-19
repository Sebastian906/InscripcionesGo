package repository

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./inscripciones.db")
	if err != nil {
		return nil, err
	}

	// Crear tablas si no existen
	queries := []string{
		`CREATE TABLE IF NOT EXISTS estudiantes (
            cedula TEXT PRIMARY KEY,
            nombre TEXT NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS materias (
            codigo TEXT PRIMARY KEY,
            nombre TEXT NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS inscripciones (
            estudiante_cedula TEXT,
            materia_codigo TEXT,
            FOREIGN KEY(estudiante_cedula) REFERENCES estudiantes(cedula),
            FOREIGN KEY(materia_codigo) REFERENCES materias(codigo),
            PRIMARY KEY(estudiante_cedula, materia_codigo)
        )`,
	}

	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
