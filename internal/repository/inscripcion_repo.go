package repository

import (
	"database/sql"
	"inscripciones/internal/domain"
)

type InscripcionRepository interface {
	Create(estudianteCedula, materiaCodigo string) error
	GetByEstudiante(cedula string) ([]*domain.Materia, error)
	GetByMateria(codigo string) ([]*domain.Estudiante, error)
	CountByEstudiante(cedula string) (int, error)
	Exists(estudianteCedula, materiaCodigo string) (bool, error)
}

type inscripcionRepo struct {
	db *sql.DB
}

func NewInscripcionRepository(db *sql.DB) InscripcionRepository {
	return &inscripcionRepo{db: db}
}

func (r *inscripcionRepo) Create(estudianteCedula, materiaCodigo string) error {
	_, err := r.db.Exec(
		"INSERT INTO inscripciones (estudiante_cedula, materia_codigo) VALUES (?, ?)",
		estudianteCedula,
		materiaCodigo,
	)
	return err
}

func (r *inscripcionRepo) GetByEstudiante(cedula string) ([]*domain.Materia, error) {
	rows, err := r.db.Query(`
		SELECT m.codigo, m.nombre 
		FROM materias m
		JOIN inscripciones i ON m.codigo = i.materia_codigo
		WHERE i.estudiante_cedula = ?
	`, cedula)
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

func (r *inscripcionRepo) GetByMateria(codigo string) ([]*domain.Estudiante, error) {
	rows, err := r.db.Query(`
		SELECT e.cedula, e.nombre 
		FROM estudiantes e
		JOIN inscripciones i ON e.cedula = i.estudiante_cedula
		WHERE i.materia_codigo = ?
	`, codigo)
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

func (r *inscripcionRepo) CountByEstudiante(cedula string) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM inscripciones 
		WHERE estudiante_cedula = ?
	`, cedula).Scan(&count)
	return count, err
}

func (r *inscripcionRepo) Exists(estudianteCedula, materiaCodigo string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM inscripciones WHERE estudiante_cedula = ? AND materia_codigo = ?)",
		estudianteCedula,
		materiaCodigo,
	).Scan(&exists)
	return exists, err
}