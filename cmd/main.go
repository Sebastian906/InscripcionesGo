package main

import (
	"fmt"
	"log"
	"inscripciones/internal/repository"
	"inscripciones/internal/service"
	"inscripciones/internal/ui"
	"inscripciones/pkg/fileutil"
)

func main() {
	// Mostrar mensaje de bienvenida
	fmt.Println("Sistema de Inscripciones Universitarias")
	fmt.Println("======================================")
	fmt.Println()

	// Inicializar base de datos
	fmt.Println("Inicializando base de datos...")
	db, err := repository.InitDB()
	if err != nil {
		log.Fatal("Error al inicializar base de datos:", err)
	}
	defer db.Close()
	fmt.Println("✓ Base de datos inicializada correctamente")

	// Crear repositorios
	estudianteRepo := repository.NewEstudianteRepository(db)
	materiaRepo := repository.NewMateriaRepository(db)
	inscripcionRepo := repository.NewInscripcionRepository(db)

	// Crear servicios
	lectorArchivo := &fileutil.LectorArchivoTexto{}
	procesadorArchivo := service.NewProcesadorArchivo(
		lectorArchivo,
		estudianteRepo,
		materiaRepo,
		inscripcionRepo,
	)

	inscripcionService := service.NewInscripcionService(
		estudianteRepo,
		materiaRepo,
		inscripcionRepo,
	)

	consultasAvanzadasService := service.NewConsultasAvanzadasService(
		estudianteRepo,
		materiaRepo,
		inscripcionRepo,
	)

	// Crear interfaz de usuario
	consoleUI := ui.NewConsoleUI(
		procesadorArchivo,
		inscripcionService,
		consultasAvanzadasService,
	)

	fmt.Println("✓ Servicios inicializados correctamente")
	fmt.Println()
	fmt.Println("INSTRUCCIONES:")
	fmt.Println("- Para cargar archivos desde testdata/, solo escriba el nombre del archivo")
	fmt.Println("- Ejemplo: 'inscripciones_validas.txt' en lugar de la ruta completa")
	fmt.Println("- Los archivos de exportación se guardarán en el directorio actual")
	fmt.Println()

	// Iniciar la interfaz de usuario
	consoleUI.MostrarMenu()
}