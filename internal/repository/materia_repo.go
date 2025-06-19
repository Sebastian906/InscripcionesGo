package repository

import (
	"database/sql"
	"inscripciones/internal/domain"
)

type MateriaRepository interface {
	Create(materia *domain.Materia) error
	GetByCodigo(codigo string) (*domain.Materia, error)
	GetAll() ([]*domain.Materia, error)  // Agregar este método a la interfaz
	Exists(codigo string) (bool, error)
}

type materiaRepo struct {
	db *sql.DB
}

func NewMateriaRepository(db *sql.DB) MateriaRepository {
	return &materiaRepo{db: db}
}

func (r *materiaRepo) Create(materia *domain.Materia) error {
	_, err := r.db.Exec(
		"INSERT INTO materias (codigo, nombre) VALUES (?, ?)",
		materia.Codigo,
		materia.Nombre,
	)
	return err
}

func (r *materiaRepo) GetByCodigo(codigo string) (*domain.Materia, error) {
	row := r.db.QueryRow("SELECT codigo, nombre FROM materias WHERE codigo = ?", codigo)

	var m domain.Materia
	err := row.Scan(&m.Codigo, &m.Nombre)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r *materiaRepo) Exists(codigo string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM materias WHERE codigo = ?)",
		codigo,
	).Scan(&exists)
	return exists, err
}

func (r *materiaRepo) GetAll() ([]*domain.Materia, error) {
	rows, err := r.db.Query("SELECT codigo, nombre FROM materias")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materias []*domain.Materia
	for rows.Next() {
		var m domain.Materia
		if err := rows.Scan(&m.Codigo, &m.Nombre); err != nil {
			return nil, err
		}
		materias = append(materias, &m)
	}
	return materias, nil
}