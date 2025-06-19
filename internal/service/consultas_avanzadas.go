package service

import (
	"fmt"
	"inscripciones/internal/domain"
	"inscripciones/internal/repository"
)

type ConsultasAvanzadasService struct {
	estudianteRepo  repository.EstudianteRepository
	materiaRepo     repository.MateriaRepository
	inscripcionRepo repository.InscripcionRepository
}

func NewConsultasAvanzadasService(
	estudianteRepo repository.EstudianteRepository,
	materiaRepo repository.MateriaRepository,
	inscripcionRepo repository.InscripcionRepository,
) *ConsultasAvanzadasService {
	return &ConsultasAvanzadasService{
		estudianteRepo:  estudianteRepo,
		materiaRepo:     materiaRepo,
		inscripcionRepo: inscripcionRepo,
	}
}

// EstadisticasGenerales representa las estadísticas del sistema
type EstadisticasGenerales struct {
	TotalEstudiantes    int
	TotalMaterias       int
	TotalInscripciones  int
	EstudiantesConMasMaterias []EstudianteConMaterias
	MateriasConMasEstudiantes []MateriaConEstudiantes
}

type EstudianteConMaterias struct {
	Estudiante     *domain.Estudiante
	CantidadMaterias int
}

type MateriaConEstudiantes struct {
	Materia          *domain.Materia
	CantidadEstudiantes int
}

// RegistroCompleto representa un registro completo de inscripción
type RegistroCompleto struct {
	Estudiante *domain.Estudiante
	Materia    *domain.Materia
}

// BuscarEstudiantePorCedula busca un estudiante por su cédula y retorna información completa
func (s *ConsultasAvanzadasService) BuscarEstudiantePorCedula(cedula string) (*domain.Estudiante, []*domain.Materia, error) {
	// Buscar estudiante
	estudiante, err := s.estudianteRepo.GetByCedula(cedula)
	if err != nil {
		return nil, nil, fmt.Errorf("error al buscar estudiante: %w", err)
	}
	
	if estudiante == nil {
		return nil, nil, nil // Estudiante no encontrado
	}
	
	// Obtener materias del estudiante
	materias, err := s.inscripcionRepo.GetByEstudiante(cedula)
	if err != nil {
		return nil, nil, fmt.Errorf("error al obtener materias del estudiante: %w", err)
	}
	
	return estudiante, materias, nil
}

// ObtenerEstadisticasGenerales genera estadísticas completas del sistema
func (s *ConsultasAvanzadasService) ObtenerEstadisticasGenerales() (*EstadisticasGenerales, error) {
	estadisticas := &EstadisticasGenerales{}
	
	// Obtener todos los estudiantes
	estudiantes, err := s.estudianteRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error al obtener estudiantes: %w", err)
	}
	estadisticas.TotalEstudiantes = len(estudiantes)
	
	// Obtener todas las materias
	materias, err := s.materiaRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error al obtener materias: %w", err)
	}
	estadisticas.TotalMaterias = len(materias)
	
	// Calcular total de inscripciones y estudiantes con más materias
	totalInscripciones := 0
	var estudiantesConMaterias []EstudianteConMaterias
	
	for _, estudiante := range estudiantes {
		count, err := s.inscripcionRepo.CountByEstudiante(estudiante.Cedula)
		if err != nil {
			continue // Continuar con el siguiente estudiante en caso de error
		}
		
		totalInscripciones += count
		if count > 0 {
			estudiantesConMaterias = append(estudiantesConMaterias, EstudianteConMaterias{
				Estudiante:       estudiante,
				CantidadMaterias: count,
			})
		}
	}
	estadisticas.TotalInscripciones = totalInscripciones
	
	// Ordenar estudiantes por cantidad de materias (top 5)
	estadisticas.EstudiantesConMasMaterias = s.obtenerTop5EstudiantesConMasMaterias(estudiantesConMaterias)
	
	// Calcular materias con más estudiantes
	var materiasConEstudiantes []MateriaConEstudiantes
	
	for _, materia := range materias {
		estudiantes, err := s.inscripcionRepo.GetByMateria(materia.Codigo)
		if err != nil {
			continue // Continuar con la siguiente materia en caso de error
		}
		
		if len(estudiantes) > 0 {
			materiasConEstudiantes = append(materiasConEstudiantes, MateriaConEstudiantes{
				Materia:             materia,
				CantidadEstudiantes: len(estudiantes),
			})
		}
	}
	
	// Ordenar materias por cantidad de estudiantes (top 5)
	estadisticas.MateriasConMasEstudiantes = s.obtenerTop5MateriasConMasEstudiantes(materiasConEstudiantes)
	
	return estadisticas, nil
}

// InsertarNuevoRegistro permite insertar un nuevo registro de inscripción
func (s *ConsultasAvanzadasService) InsertarNuevoRegistro(cedula, nombreEstudiante, codigoMateria, nombreMateria string) error {
	// Verificar si el estudiante existe, si no, crearlo
	existe, err := s.estudianteRepo.Exists(cedula)
	if err != nil {
		return fmt.Errorf("error al verificar existencia del estudiante: %w", err)
	}
	
	if !existe {
		estudiante := domain.NewEstudiante(cedula, nombreEstudiante)
		err = s.estudianteRepo.Create(estudiante)
		if err != nil {
			return fmt.Errorf("error al crear estudiante: %w", err)
		}
	}
	
	// Verificar si la materia existe, si no, crearla
	existe, err = s.materiaRepo.Exists(codigoMateria)
	if err != nil {
		return fmt.Errorf("error al verificar existencia de la materia: %w", err)
	}
	
	if !existe {
		materia := domain.NewMateria(codigoMateria, nombreMateria)
		err = s.materiaRepo.Create(materia)
		if err != nil {
			return fmt.Errorf("error al crear materia: %w", err)
		}
	}
	
	// Verificar si la inscripción ya existe
	existe, err = s.inscripcionRepo.Exists(cedula, codigoMateria)
	if err != nil {
		return fmt.Errorf("error al verificar existencia de la inscripción: %w", err)
	}
	
	if existe {
		return fmt.Errorf("el estudiante ya está inscrito en esta materia")
	}
	
	// Crear la inscripción
	err = s.inscripcionRepo.Create(cedula, codigoMateria)
	if err != nil {
		return fmt.Errorf("error al crear inscripción: %w", err)
	}
	
	return nil
}

// ObtenerTodosLosRegistros obtiene todos los registros de inscripciones
func (s *ConsultasAvanzadasService) ObtenerTodosLosRegistros() ([]RegistroCompleto, error) {
	var registros []RegistroCompleto
	
	// Obtener todos los estudiantes
	estudiantes, err := s.estudianteRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error al obtener estudiantes: %w", err)
	}
	
	// Para cada estudiante, obtener sus materias
	for _, estudiante := range estudiantes {
		materias, err := s.inscripcionRepo.GetByEstudiante(estudiante.Cedula)
		if err != nil {
			continue // Continuar con el siguiente estudiante en caso de error
		}
		
		// Crear un registro por cada materia del estudiante
		for _, materia := range materias {
			registros = append(registros, RegistroCompleto{
				Estudiante: estudiante,
				Materia:    materia,
			})
		}
	}
	
	return registros, nil
}

// Función auxiliar para obtener top 5 estudiantes con más materias
func (s *ConsultasAvanzadasService) obtenerTop5EstudiantesConMasMaterias(estudiantes []EstudianteConMaterias) []EstudianteConMaterias {
	// Ordenamiento burbuja simple (para mantener simplicidad)
	for i := 0; i < len(estudiantes)-1; i++ {
		for j := 0; j < len(estudiantes)-i-1; j++ {
			if estudiantes[j].CantidadMaterias < estudiantes[j+1].CantidadMaterias {
				estudiantes[j], estudiantes[j+1] = estudiantes[j+1], estudiantes[j]
			}
		}
	}
	
	// Retornar máximo 5
	if len(estudiantes) > 5 {
		return estudiantes[:5]
	}
	return estudiantes
}

// Función auxiliar para obtener top 5 materias con más estudiantes
func (s *ConsultasAvanzadasService) obtenerTop5MateriasConMasEstudiantes(materias []MateriaConEstudiantes) []MateriaConEstudiantes {
	// Ordenamiento burbuja simple (para mantener simplicidad)
	for i := 0; i < len(materias)-1; i++ {
		for j := 0; j < len(materias)-i-1; j++ {
			if materias[j].CantidadEstudiantes < materias[j+1].CantidadEstudiantes {
				materias[j], materias[j+1] = materias[j+1], materias[j]
			}
		}
	}
	
	// Retornar máximo 5
	if len(materias) > 5 {
		return materias[:5]
	}
	return materias
}