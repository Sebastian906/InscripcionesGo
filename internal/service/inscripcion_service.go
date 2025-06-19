package service

import (
	"inscripciones/internal/domain"
	"inscripciones/internal/repository"
)

type InscripcionService struct {
	EstudianteRepo  repository.EstudianteRepository  // Hacer público
	MateriaRepo     repository.MateriaRepository     // Hacer público
	InscripcionRepo repository.InscripcionRepository // Hacer público
}

func NewInscripcionService(
	estudianteRepo repository.EstudianteRepository,
	materiaRepo repository.MateriaRepository,
	inscripcionRepo repository.InscripcionRepository,
) *InscripcionService {
	return &InscripcionService{
		EstudianteRepo:  estudianteRepo,
		MateriaRepo:     materiaRepo,
		InscripcionRepo: inscripcionRepo,
	}
}

func (s *InscripcionService) ObtenerEstudiantesPorMateria(codigoMateria string) ([]*domain.Estudiante, error) {
	return s.InscripcionRepo.GetByMateria(codigoMateria)
}

func (s *InscripcionService) ObtenerMateriasPorEstudiante(cedula string) ([]*domain.Materia, error) {
	return s.InscripcionRepo.GetByEstudiante(cedula)
}

func (s *InscripcionService) ContarMateriasPorEstudiante(cedula string) (int, error) {
	return s.InscripcionRepo.CountByEstudiante(cedula)
}

func (s *InscripcionService) ExportarDatos() (*domain.ConsolidadoInscripciones, error) {
	// Obtener todos los estudiantes
	estudiantes, err := s.EstudianteRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Obtener todas las materias
	materias, err := s.MateriaRepo.GetAll()
	if err != nil {
		return nil, err
	}

	consolidado := domain.NewConsolidadoInscripciones()

	// Llenar el consolidado con estudiantes
	for _, estudiante := range estudiantes {
		consolidado.Estudiantes[estudiante.Cedula] = estudiante
	}

	// Llenar el consolidado con materias
	for _, materia := range materias {
		consolidado.Materias[materia.Codigo] = materia
	}

	return consolidado, nil
}