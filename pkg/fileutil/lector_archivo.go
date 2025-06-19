package fileutil

import (
	"bufio"
	"os"
)

type LectorArchivo interface {
	ObtenerLineas(ruta string) ([]string, error)
}

type LectorArchivoTexto struct{}

func (l *LectorArchivoTexto) ObtenerLineas(ruta string) ([]string, error) {
	file, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lineas []string
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		linea := scanner.Text()
		if linea != "" { // Ignorar líneas vacías
			lineas = append(lineas, linea)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lineas, nil
}