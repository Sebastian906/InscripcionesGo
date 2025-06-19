package ui

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"inscripciones/internal/domain"
	"inscripciones/internal/service"
	"os"
	"path/filepath"
	"strings"
)

type ConsoleUI struct {
	procesador         *service.ProcesadorArchivo
	inscripcionSvc     *service.InscripcionService
	consultasAvanzadas *service.ConsultasAvanzadasService
	consolidado        *domain.ConsolidadoInscripciones
	archivoCargado     bool
}

func NewConsoleUI(
	procesador *service.ProcesadorArchivo,
	inscripcionSvc *service.InscripcionService,
	consultasAvanzadas *service.ConsultasAvanzadasService,
) *ConsoleUI {
	return &ConsoleUI{
		procesador:         procesador,
		inscripcionSvc:     inscripcionSvc,
		consultasAvanzadas: consultasAvanzadas,
		consolidado:        domain.NewConsolidadoInscripciones(),
		archivoCargado:     false,
	}
}

func (c *ConsoleUI) MostrarMenu() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n=== SISTEMA DE INSCRIPCIONES UNIVERSITARIAS ===")
		fmt.Println("1. Cargar archivo de inscripciones")
		fmt.Println("2. Mostrar total de materias por estudiante")
		fmt.Println("3. Filtrar estudiantes por materia")
		fmt.Println("4. Exportar datos a JSON")
		fmt.Println("5. Exportar datos a CSV")
		fmt.Println("6. Consultas avanzadas")
		fmt.Println("7. Salir")
		fmt.Print("Seleccione una opción: ")

		scanner.Scan()
		opcion := scanner.Text()

		switch opcion {
		case "1":
			c.cargarArchivo(scanner)
		case "2":
			c.mostrarMateriasPorEstudiante()
		case "3":
			c.filtrarPorMateria(scanner)
		case "4":
			c.exportarJSON()
		case "5":
			c.exportarCSV()
		case "6":
			c.mostrarMenuConsultasAvanzadas(scanner)
		case "7":
			fmt.Println("Saliendo del programa...")
			return
		default:
			fmt.Println("Opción no válida. Intente nuevamente.")
		}
	}
}

func (c *ConsoleUI) mostrarMenuConsultasAvanzadas(scanner *bufio.Scanner) {
	for {
		fmt.Println("\n=== CONSULTAS AVANZADAS ===")
		fmt.Println("1. Buscar estudiante por cédula")
		fmt.Println("2. Ver estadísticas generales")
		fmt.Println("3. Insertar nuevo registro")
		fmt.Println("4. Ver todos los registros")
		fmt.Println("5. Volver al menú principal")
		fmt.Print("Seleccione una opción: ")

		scanner.Scan()
		opcion := scanner.Text()

		switch opcion {
		case "1":
			c.buscarEstudiantePorCedula(scanner)
		case "2":
			c.mostrarEstadisticasGenerales()
		case "3":
			c.insertarNuevoRegistro(scanner)
		case "4":
			c.mostrarTodosLosRegistros()
		case "5":
			return // Volver al menú principal
		default:
			fmt.Println("Opción no válida. Intente nuevamente.")
		}
	}
}

func (c *ConsoleUI) cargarArchivo(scanner *bufio.Scanner) {
	fmt.Print("\nIngrese la ruta del archivo de inscripciones: ")
	scanner.Scan()
	ruta := scanner.Text()

	// Si no se especifica una ruta completa, buscar en testdata
	if !filepath.IsAbs(ruta) && !strings.Contains(ruta, string(filepath.Separator)) {
		ruta = filepath.Join("testdata", ruta)
	}

	consolidado, err := c.procesador.ProcesarArchivo(ruta)
	if err != nil {
		fmt.Printf("\nError al procesar archivo: %v\n", err)
		return
	}

	c.consolidado = consolidado
	c.archivoCargado = true

	fmt.Printf("\nArchivo cargado exitosamente!\n")
	fmt.Printf("Estudiantes registrados: %d\n", len(consolidado.Estudiantes))
	fmt.Printf("Materias registradas: %d\n", len(consolidado.Materias))
}

func (c *ConsoleUI) mostrarMateriasPorEstudiante() {
	if !c.archivoCargado {
		fmt.Println("\nPrimero debe cargar un archivo de inscripciones (Opción 1)")
		return
	}

	fmt.Println("\n=== MATERIAS POR ESTUDIANTE ===")
	for cedula, estudiante := range c.consolidado.Estudiantes {
		count, err := c.inscripcionSvc.ContarMateriasPorEstudiante(cedula)
		if err != nil {
			fmt.Printf("Error al contar materias para %s: %v\n", estudiante.Nombre, err)
			continue
		}
		fmt.Printf("- %s (Cédula: %s): %d materias\n", estudiante.Nombre, cedula, count)
	}
}

func (c *ConsoleUI) filtrarPorMateria(scanner *bufio.Scanner) {
	if !c.archivoCargado {
		fmt.Println("\nPrimero debe cargar un archivo de inscripciones (Opción 1)")
		return
	}

	fmt.Print("\nIngrese el código de la materia: ")
	scanner.Scan()
	codigo := scanner.Text()

	// Verificar si la materia existe
	materia, ok := c.consolidado.Materias[codigo]
	if !ok {
		fmt.Println("Materia no encontrada. Intente con un código válido.")
		return
	}

	estudiantes, err := c.inscripcionSvc.ObtenerEstudiantesPorMateria(codigo)
	if err != nil {
		fmt.Printf("Error al obtener estudiantes: %v\n", err)
		return
	}

	fmt.Printf("\n=== ESTUDIANTES INSCRITOS EN %s (%s) ===\n", materia.Nombre, materia.Codigo)
	if len(estudiantes) == 0 {
		fmt.Println("No hay estudiantes inscritos en esta materia")
		return
	}

	for i, estudiante := range estudiantes {
		fmt.Printf("%d. %s (Cédula: %s)\n", i+1, estudiante.Nombre, estudiante.Cedula)
	}
	fmt.Printf("\nTotal: %d estudiantes\n", len(estudiantes))
}

func (c *ConsoleUI) buscarEstudiantePorCedula(scanner *bufio.Scanner) {
	fmt.Print("\nIngrese la cédula del estudiante: ")
	scanner.Scan()
	cedula := strings.TrimSpace(scanner.Text())

	if cedula == "" {
		fmt.Println("La cédula no puede estar vacía.")
		return
	}

	estudiante, materias, err := c.consultasAvanzadas.BuscarEstudiantePorCedula(cedula)
	if err != nil {
		fmt.Printf("Error al buscar estudiante: %v\n", err)
		return
	}

	if estudiante == nil {
		fmt.Printf("No se encontró un estudiante con cédula: %s\n", cedula)
		return
	}

	fmt.Printf("\n=== INFORMACIÓN DEL ESTUDIANTE ===\n")
	fmt.Printf("Cédula: %s\n", estudiante.Cedula)
	fmt.Printf("Nombre: %s\n", estudiante.Nombre)
	fmt.Printf("Total de materias inscritas: %d\n\n", len(materias))

	if len(materias) > 0 {
		fmt.Println("Materias inscritas:")
		for i, materia := range materias {
			fmt.Printf("%d. %s - %s\n", i+1, materia.Codigo, materia.Nombre)
		}
	} else {
		fmt.Println("El estudiante no tiene materias inscritas.")
	}
}

func (c *ConsoleUI) mostrarEstadisticasGenerales() {
	estadisticas, err := c.consultasAvanzadas.ObtenerEstadisticasGenerales()
	if err != nil {
		fmt.Printf("Error al obtener estadísticas: %v\n", err)
		return
	}

	fmt.Println("\n=== ESTADÍSTICAS GENERALES ===")
	fmt.Printf("Total de estudiantes: %d\n", estadisticas.TotalEstudiantes)
	fmt.Printf("Total de materias: %d\n", estadisticas.TotalMaterias)
	fmt.Printf("Total de inscripciones: %d\n\n", estadisticas.TotalInscripciones)

	if len(estadisticas.EstudiantesConMasMaterias) > 0 {
		fmt.Println("TOP 5 - Estudiantes con más materias:")
		for i, item := range estadisticas.EstudiantesConMasMaterias {
			fmt.Printf("%d. %s (%s): %d materias\n", 
				i+1, item.Estudiante.Nombre, item.Estudiante.Cedula, item.CantidadMaterias)
		}
		fmt.Println()
	}

	if len(estadisticas.MateriasConMasEstudiantes) > 0 {
		fmt.Println("TOP 5 - Materias con más estudiantes:")
		for i, item := range estadisticas.MateriasConMasEstudiantes {
			fmt.Printf("%d. %s (%s): %d estudiantes\n", 
				i+1, item.Materia.Nombre, item.Materia.Codigo, item.CantidadEstudiantes)
		}
	}
}

func (c *ConsoleUI) insertarNuevoRegistro(scanner *bufio.Scanner) {
	fmt.Println("\n=== INSERTAR NUEVO REGISTRO ===")
	
	fmt.Print("Ingrese la cédula del estudiante: ")
	scanner.Scan()
	cedula := strings.TrimSpace(scanner.Text())
	
	fmt.Print("Ingrese el nombre del estudiante: ")
	scanner.Scan()
	nombreEstudiante := strings.TrimSpace(scanner.Text())
	
	fmt.Print("Ingrese el código de la materia: ")
	scanner.Scan()
	codigoMateria := strings.TrimSpace(scanner.Text())
	
	fmt.Print("Ingrese el nombre de la materia: ")
	scanner.Scan()
	nombreMateria := strings.TrimSpace(scanner.Text())

	// Validaciones básicas
	if cedula == "" || nombreEstudiante == "" || codigoMateria == "" || nombreMateria == "" {
		fmt.Println("Todos los campos son obligatorios.")
		return
	}

	err := c.consultasAvanzadas.InsertarNuevoRegistro(cedula, nombreEstudiante, codigoMateria, nombreMateria)
	if err != nil {
		fmt.Printf("Error al insertar registro: %v\n", err)
		return
	}

	fmt.Println("Registro insertado exitosamente!")
}

func (c *ConsoleUI) mostrarTodosLosRegistros() {
	registros, err := c.consultasAvanzadas.ObtenerTodosLosRegistros()
	if err != nil {
		fmt.Printf("Error al obtener registros: %v\n", err)
		return
	}

	if len(registros) == 0 {
		fmt.Println("\nNo hay registros en la base de datos.")
		return
	}

	fmt.Printf("\n=== TODOS LOS REGISTROS (%d) ===\n", len(registros))
	fmt.Printf("%-12s %-25s %-10s %-25s\n", "CÉDULA", "NOMBRE ESTUDIANTE", "CÓD MAT", "NOMBRE MATERIA")
	fmt.Println(strings.Repeat("-", 75))

	for _, registro := range registros {
		fmt.Printf("%-12s %-25s %-10s %-25s\n",
			registro.Estudiante.Cedula,
			c.truncateString(registro.Estudiante.Nombre, 25),
			registro.Materia.Codigo,
			c.truncateString(registro.Materia.Nombre, 25))
	}
}

func (c *ConsoleUI) exportarJSON() {
	if !c.archivoCargado {
		fmt.Println("\nPrimero debe cargar un archivo de inscripciones (Opción 1)")
		return
	}

	type EstudianteExport struct {
		Cedula string `json:"cedula"`
		Nombre string `json:"nombre"`
	}

	type MateriaExport struct {
		Codigo string `json:"codigo"`
		Nombre string `json:"nombre"`
	}

	type InscripcionExport struct {
		Estudiante EstudianteExport `json:"estudiante"`
		Materia    MateriaExport    `json:"materia"`
	}

	var inscripciones []InscripcionExport

	for cedula, estudiante := range c.consolidado.Estudiantes {
		materias, err := c.inscripcionSvc.ObtenerMateriasPorEstudiante(cedula)
		if err != nil {
			fmt.Printf("Error al obtener materias para %s: %v\n", estudiante.Nombre, err)
			continue
		}

		for _, materia := range materias {
			inscripciones = append(inscripciones, InscripcionExport{
				Estudiante: EstudianteExport{
					Cedula: estudiante.Cedula,
					Nombre: estudiante.Nombre,
				},
				Materia: MateriaExport{
					Codigo: materia.Codigo,
					Nombre: materia.Nombre,
				},
			})
		}
	}

	jsonData, err := json.MarshalIndent(inscripciones, "", "  ")
	if err != nil {
		fmt.Printf("Error al generar JSON: %v\n", err)
		return
	}

	filename := "inscripciones.json"
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error al escribir archivo JSON: %v\n", err)
		return
	}

	fmt.Printf("\nDatos exportados exitosamente a %s\n", filename)
}

func (c *ConsoleUI) exportarCSV() {
	if !c.archivoCargado {
		fmt.Println("\nPrimero debe cargar un archivo de inscripciones (Opción 1)")
		return
	}

	filename := "inscripciones.csv"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error al crear archivo CSV: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir encabezados
	headers := []string{"CEDULA", "NOMBRE_ESTUDIANTE", "CODIGO_MATERIA", "NOMBRE_MATERIA"}
	if err := writer.Write(headers); err != nil {
		fmt.Printf("Error al escribir encabezados CSV: %v\n", err)
		return
	}

	// Escribir datos
	for cedula, estudiante := range c.consolidado.Estudiantes {
		materias, err := c.inscripcionSvc.ObtenerMateriasPorEstudiante(cedula)
		if err != nil {
			fmt.Printf("Error al obtener materias para %s: %v\n", estudiante.Nombre, err)
			continue
		}

		for _, materia := range materias {
			record := []string{
				estudiante.Cedula,
				estudiante.Nombre,
				materia.Codigo,
				materia.Nombre,
			}
			if err := writer.Write(record); err != nil {
				fmt.Printf("Error al escribir registro CSV: %v\n", err)
				continue
			}
		}
	}

	fmt.Printf("\nDatos exportados exitosamente a %s\n", filename)
}

// Función auxiliar para truncar strings
func (c *ConsoleUI) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}