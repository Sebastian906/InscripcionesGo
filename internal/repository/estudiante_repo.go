package repository

import (
	"database/sql"
	"inscripciones/internal/domain"
)

type EstudianteRepository interface {
	Create(estudiante *domain.Estudiante) error
	GetByCedula(cedula string) (*domain.Estudiante, error)
	GetAll() ([]*domain.Estudiante, error)  // Agregar este método a la interfaz
	Exists(cedula string) (bool, error)
}

type estudianteRepo struct {
	db *sql.DB
}

func NewEstudianteRepository(db *sql.DB) EstudianteRepository {
	return &estudianteRepo{db: db}
}

func (r *estudianteRepo) Create(estudiante *domain.Estudiante) error {
	_, err := r.db.Exec(
		"INSERT INTO estudiantes (cedula, nombre) VALUES (?, ?)",
		estudiante.Cedula,
		estudiante.Nombre,
	)
	return err
}

func (r *estudianteRepo) GetByCedula(cedula string) (*domain.Estudiante, error) {
	row := r.db.QueryRow("SELECT cedula, nombre FROM estudiantes WHERE cedula = ?", cedula)

	var e domain.Estudiante
	err := row.Scan(&e.Cedula, &e.Nombre)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &e, nil
}

func (r *estudianteRepo) Exists(cedula string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM estudiantes WHERE cedula = ?)",
		cedula,
	).Scan(&exists)
	return exists, err
}

func (r *estudianteRepo) GetAll() ([]*domain.Estudiante, error) {
	rows, err := r.db.Query("SELECT cedula, nombre FROM estudiantes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var estudiantes []*domain.Estudiante
	for rows.Next() {
		var e domain.Estudiante
		if err := rows.Scan(&e.Cedula, &e.Nombre); err != nil {
			return nil, err
		}
		estudiantes = append(estudiantes, &e)
	}
	return estudiantes, nil
}