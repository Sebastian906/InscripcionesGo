package domain

type Estudiante struct {
    Cedula string
    Nombre string
}

func NewEstudiante(cedula, nombre string) *Estudiante {
    return &Estudiante{
        Cedula: cedula,
        Nombre: nombre,
    }
}