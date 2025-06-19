package domain

type Materia struct {
    Codigo string
    Nombre string
}

func NewMateria(codigo, nombre string) *Materia {
    return &Materia{
        Codigo: codigo,
        Nombre: nombre,
    }
}