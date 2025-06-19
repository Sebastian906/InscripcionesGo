package domain

type Inscripcion struct {
    Estudiante *Estudiante
    Materia    *Materia
}

type ConsolidadoInscripciones struct {
    Estudiantes map[string]*Estudiante
    Materias    map[string]*Materia
}

func NewConsolidadoInscripciones() *ConsolidadoInscripciones {
    return &ConsolidadoInscripciones{
        Estudiantes: make(map[string]*Estudiante),
        Materias:    make(map[string]*Materia),
    }
}