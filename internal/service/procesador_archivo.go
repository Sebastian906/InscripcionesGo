package service

import (
	"fmt"
	"strings"
	"inscripciones/internal/domain"
	"inscripciones/internal/repository"
	"inscripciones/pkg/fileutil"
)

type ProcesadorArchivo struct {
	lector          fileutil.LectorArchivo
	estudianteRepo  repository.EstudianteRepository
	materiaRepo     repository.MateriaRepository
	inscripcionRepo repository.InscripcionRepository
}

func NewProcesadorArchivo(
	lector fileutil.LectorArchivo,
	estudianteRepo repository.EstudianteRepository,
	materiaRepo repository.MateriaRepository,
	inscripcionRepo repository.InscripcionRepository,
) *ProcesadorArchivo {
	return &ProcesadorArchivo{
		lector:          lector,
		estudianteRepo:  estudianteRepo,
		materiaRepo:     materiaRepo,
		inscripcionRepo: inscripcionRepo,
	}
}

func (p *ProcesadorArchivo) ProcesarArchivo(ruta string) (*domain.ConsolidadoInscripciones, error) {
	lineas, err := p.lector.ObtenerLineas(ruta)
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo: %w", err)
	}

	consolidado := domain.NewConsolidadoInscripciones()
	lineasValidas := []string{}

	// Validar y procesar cada línea
	for i, linea := range lineas {
		if err := p.validarLinea(linea, i+1); err != nil {
			fmt.Printf("Advertencia línea %d: %v\n", i+1, err)
			continue // Saltar líneas inválidas
		}

		campos := strings.Split(linea, ",")
		cedula := strings.TrimSpace(campos[0])
		nombreEstudiante := strings.TrimSpace(campos[1])
		codigoMateria := strings.TrimSpace(campos[2])
		nombreMateria := strings.TrimSpace(campos[3])

		// Agregar estudiante si no existe en el consolidado
		if _, ok := consolidado.Estudiantes[cedula]; !ok {
			consolidado.Estudiantes[cedula] = domain.NewEstudiante(cedula, nombreEstudiante)
		}

		// Agregar materia si no existe en el consolidado
		if _, ok := consolidado.Materias[codigoMateria]; !ok {
			consolidado.Materias[codigoMateria] = domain.NewMateria(codigoMateria, nombreMateria)
		}

		lineasValidas = append(lineasValidas, linea)
	}

	if len(lineasValidas) == 0 {
		return nil, fmt.Errorf("no se encontraron líneas válidas en el archivo")
	}

	// Guardar en base de datos
	err = p.guardarEnBaseDatos(consolidado, lineasValidas)
	if err != nil {
		return nil, fmt.Errorf("error al guardar en base de datos: %w", err)
	}

	return consolidado, nil
}

func (p *ProcesadorArchivo) validarLinea(linea string, _ int) error {
	if strings.TrimSpace(linea) == "" {
		return fmt.Errorf("línea vacía")
	}

	campos := strings.Split(linea, ",")
	if len(campos) != 4 {
		return fmt.Errorf("formato incorrecto - se esperan 4 campos separados por coma, encontrados %d", len(campos))
	}

	// Validar que ningún campo esté vacío
	for i, campo := range campos {
		if strings.TrimSpace(campo) == "" {
			return fmt.Errorf("campo %d está vacío", i+1)
		}
	}

	cedula := strings.TrimSpace(campos[0])
	nombreEstudiante := strings.TrimSpace(campos[1])
	codigoMateria := strings.TrimSpace(campos[2])
	nombreMateria := strings.TrimSpace(campos[3])

	// Validaciones adicionales
	if len(cedula) < 6 || len(cedula) > 12 {
		return fmt.Errorf("cédula '%s' debe tener entre 6 y 12 caracteres", cedula)
	}

	if len(nombreEstudiante) < 2 {
		return fmt.Errorf("nombre del estudiante '%s' debe tener al menos 2 caracteres", nombreEstudiante)
	}

	if len(codigoMateria) < 2 {
		return fmt.Errorf("código de materia '%s' debe tener al menos 2 caracteres", codigoMateria)
	}

	if len(nombreMateria) < 2 {
		return fmt.Errorf("nombre de materia '%s' debe tener al menos 2 caracteres", nombreMateria)
	}

	return nil
}

func (p *ProcesadorArchivo) guardarEnBaseDatos(consolidado *domain.ConsolidadoInscripciones, lineas []string) error {
	// Guardar estudiantes
	for _, estudiante := range consolidado.Estudiantes {
		exists, err := p.estudianteRepo.Exists(estudiante.Cedula)
		if err != nil {
			return fmt.Errorf("error al verificar existencia del estudiante %s: %w", estudiante.Cedula, err)
		}
		if !exists {
			err = p.estudianteRepo.Create(estudiante)
			if err != nil {
				return fmt.Errorf("error al crear estudiante %s: %w", estudiante.Cedula, err)
			}
		}
	}

	// Guardar materias
	for _, materia := range consolidado.Materias {
		exists, err := p.materiaRepo.Exists(materia.Codigo)
		if err != nil {
			return fmt.Errorf("error al verificar existencia de la materia %s: %w", materia.Codigo, err)
		}
		if !exists {
			err = p.materiaRepo.Create(materia)
			if err != nil {
				return fmt.Errorf("error al crear materia %s: %w", materia.Codigo, err)
			}
		}
	}

	// Procesar y guardar inscripciones
	for _, linea := range lineas {
		campos := strings.Split(linea, ",")
		if len(campos) != 4 {
			continue // Esta validación ya se hizo antes, pero por seguridad
		}

		cedula := strings.TrimSpace(campos[0])
		codigoMateria := strings.TrimSpace(campos[2])

		// Verificar si la inscripción ya existe
		exists, err := p.inscripcionRepo.Exists(cedula, codigoMateria)
		if err != nil {
			return fmt.Errorf("error al verificar existencia de inscripción %s-%s: %w", cedula, codigoMateria, err)
		}

		if !exists {
			err = p.inscripcionRepo.Create(cedula, codigoMateria)
			if err != nil {
				return fmt.Errorf("error al crear inscripción %s-%s: %w", cedula, codigoMateria, err)
			}
		}
	}

	return nil
}