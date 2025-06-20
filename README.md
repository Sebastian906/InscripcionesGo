# Sistema de Inscripciones Universitarias

Un sistema desarrollado en Go que permite gestionar y procesar las inscripciones de estudiantes universitarios a través de archivos de texto, con almacenamiento en base de datos SQLite y una interfaz de consola interactiva.

## 📋 Tabla de Contenidos

- [Descripción](#descripción)
- [Objetivos](#objetivos)
- [Características](#características)
- [Arquitectura del Proyecto](#arquitectura-del-proyecto)
- [Estructura de Directorios](#estructura-de-directorios)
- [Instalación y Configuración](#instalación-y-configuración)
- [Uso del Sistema](#uso-del-sistema)
- [Formato de Archivos](#formato-de-archivos)
- [Funcionalidades](#funcionalidades)
- [Diagrama de Flujo](#diagrama-de-flujo)
- [Ejemplos de Uso](#ejemplos-de-uso)
- [Testing](#testing)
- [Tecnologías Utilizadas](#tecnologías-utilizadas)
- [Contribución](#contribución)

## 📖 Descripción

El Sistema de Inscripciones Universitarias es una aplicación de consola desarrollada en Go que permite a los coordinadores de programas universitarios procesar archivos de texto con información de inscripciones de estudiantes. El sistema valida, almacena y permite consultar la información de manera eficiente.

## 🎯 Objetivos

- **Comprender la relación** entre diagramas de clase y su implementación en código
- **Evaluar alternativas** para implementar asociaciones en el diseño orientado a objetos
- **Diseñar y desarrollar** un sistema con interfaz de consola y acceso a base de datos
- **Implementar y validar** código con un enfoque práctico

## ✨ Características

- ✅ **Procesamiento de archivos**: Lee y valida archivos CSV con información de inscripciones
- ✅ **Base de datos SQLite**: Almacenamiento persistente de estudiantes, materias e inscripciones
- ✅ **Interfaz de consola**: Menú interactivo para todas las operaciones
- ✅ **Validación de datos**: Detección de errores en formato y duplicados
- ✅ **Exportación**: Generación de reportes en formato JSON y CSV
- ✅ **Consultas avanzadas**: Estadísticas y búsquedas personalizadas
- ✅ **Manejo de errores**: Tratamiento robusto de excepciones

## 🏗️ Arquitectura del Proyecto

El proyecto sigue una **arquitectura hexagonal (Clean Architecture)** con las siguientes capas:

### Capas del Sistema

1. **Domain (Dominio)**: Entidades de negocio (Estudiante, Materia, Inscripción)
2. **Repository (Repositorio)**: Acceso a datos y persistencia
3. **Service (Servicio)**: Lógica de negocio y casos de uso
4. **UI (Interfaz)**: Presentación e interacción con el usuario
5. **Utils (Utilidades)**: Herramientas auxiliares (lectura de archivos)

### Principios de Diseño

- **Separación de responsabilidades**: Cada capa tiene una función específica
- **Inyección de dependencias**: Facilita testing y mantenimiento
- **Interfaces**: Abstracción para facilitar cambios y testing
- **Single Responsibility**: Cada clase/función tiene una única responsabilidad

## 📁 Estructura de Directorios

```
/inscripciones
├── /cmd                     # Punto de entrada de la aplicación
│   ├── inscripciones.db     # Base de datos SQLite
│   └── main.go              # Función principal
├── /internal                # Código interno de la aplicación
│   ├── /domain              # Entidades del dominio
│   │   ├── estudiante.go    # Entidad Estudiante
│   │   ├── materia.go       # Entidad Materia
│   │   └── inscripcion.go   # Entidad Inscripción y Consolidado
│   ├── /repository          # Capa de acceso a datos
│   │   ├── database.go      # Configuración de BD
│   │   ├── estudiante_repo.go
│   │   ├── materia_repo.go
│   │   └── inscripcion_repo.go
│   ├── /service             # Lógica de negocio
│   │   ├── consultas_avanzadas.go
│   │   ├── procesador_archivo.go
│   │   └── inscripcion_service.go
│   └── /ui                  # Interfaz de usuario
│       └── console.go       # Interfaz de consola
├── /pkg                     # Paquetes reutilizables
│   └── /fileutil            # Utilidades para archivos
│       └── lector_archivo.go
├── /testdata                # Archivos de prueba
│   ├── inscripciones_validas.txt
│   └── inscripciones_invalidas.txt
├── go.mod                   # Dependencias del proyecto
├── go.sum                   # Checksums de dependencias
├── inscripciones.csv        # Archivo de salida CSV
├── inscripciones.json       # Archivo de salida JSON
└── inscripciones.db         # Base de datos SQLite
```

## 🚀 Instalación y Configuración

### Prerrequisitos

- Go 1.19 o superior
- Git (opcional, para clonar el repositorio)

### Pasos de Instalación

1. **Clonar el repositorio** (si aplica):
```bash
git clone <url-del-repositorio>
cd inscripciones
```

2. **Instalar dependencias**:
```bash
go mod tidy
```

3. **Compilar el proyecto**:
```bash
go build -o inscripciones ./cmd
```

4. **Ejecutar la aplicación**:
```bash
./inscripciones
```

O directamente con Go:
```bash
go run ./cmd/main.go
```

## 🎮 Uso del Sistema

### Menú Principal

Al ejecutar la aplicación, se presenta el siguiente menú:

```
=== SISTEMA DE INSCRIPCIONES UNIVERSITARIAS ===
1. Cargar archivo de inscripciones
2. Mostrar total de materias por estudiante
3. Filtrar estudiantes por materia
4. Exportar datos a JSON
5. Exportar datos a CSV
6. Consultas avanzadas
7. Salir
```

### Menú de Consultas Avanzadas

```
=== CONSULTAS AVANZADAS ===
1. Buscar estudiante por cédula
2. Ver estadísticas generales
3. Insertar nuevo registro
4. Ver todos los registros
5. Volver al menú principal
```

## 📄 Formato de Archivos

### Archivo de Entrada

Los archivos de inscripciones deben seguir el formato CSV:

```
cedula,nombre_estudiante,codigo_materia,nombre_materia
1234567,Lulú López,1040,Cálculo
9876534,Pepito Pérez,1040,Cálculo
4567766,Calvin Clein,1050,Física I
1234567,Lulú López,1060,Administración
4567766,Calvin Clein,1070,Espíritu Empresarial
```

### Validaciones

- **Formato**: Exactamente 4 campos separados por comas
- **Cédula**: Entre 6 y 12 caracteres
- **Nombres**: Mínimo 2 caracteres
- **Códigos**: Mínimo 2 caracteres
- **Campos vacíos**: No se permiten campos vacíos

## 🔧 Funcionalidades

### 1. Procesamiento de Archivos
- Lectura y validación de archivos CSV
- Detección de errores y líneas inválidas
- Procesamiento robusto con manejo de excepciones

### 2. Gestión de Base de Datos
- Creación automática de tablas SQLite
- Prevención de duplicados
- Consultas optimizadas

### 3. Interfaz de Usuario
- Menú interactivo en consola
- Navegación intuitiva
- Mensajes informativos y de error

### 4. Exportación de Datos
- **JSON**: Formato estructurado para APIs
- **CSV**: Compatible con Excel y otras herramientas

### 5. Consultas y Reportes
- Estadísticas generales del sistema
- Búsquedas por estudiante o materia
- Rankings de estudiantes y materias más populares

## 📊 Diagrama de Flujo

![image](https://github.com/user-attachments/assets/aa279b1e-f9b0-4495-97ee-37550b9a1d85)

## 🔧 Ejemplos de Uso

### 1. Cargar un Archivo de Inscripciones

```bash
# Ejecutar la aplicación
go run ./cmd/main.go

# Seleccionar opción 1
# Ingresar: inscripciones_validas.txt
```

### 2. Consultar Estudiante por Cédula

```bash
# Desde el menú principal, seleccionar opción 6
# Luego opción 1
# Ingresar cédula: 1234567
```

### 3. Exportar Datos

```bash
# Opción 4 para JSON
# Opción 5 para CSV
# Los archivos se generan en el directorio actual
```

## 🧪 Testing

### Archivos de Prueba Incluidos

1. **inscripciones_validas.txt**: Archivo con datos correctos
2. **inscripciones_invalidas.txt**: Archivo con errores para testing

### Casos de Prueba

#### Ejemplo de Datos Válidos
```
1234567,Lulú López,1040,Cálculo
9876534,Pepito Pérez,1040,Cálculo
4567766,Calvin Clein,1050,Física I
```

#### Ejemplo de Datos Inválidos
```
,Estudiante Sin Cédula,1040,Cálculo
1234567,,1040,Cálula
1234567,Estudiante,1040
1234567,Estudiante,1040,Materia,Extra
```

### Ejecutar Pruebas

```bash
# Ejecutar la aplicación
go run ./cmd/main.go

# Probar con archivo válido
# Opción 1 → inscripciones_validas.txt

# Probar con archivo inválido
# Opción 1 → inscripciones_invalidas.txt
```

## 💻 Tecnologías Utilizadas

- **Lenguaje**: Go 1.24.3
- **Base de Datos**: SQLite con driver `github.com/glebarez/go-sqlite`
- **Arquitectura**: Clean Architecture / Hexagonal
- **Patrones**: Repository, Service Layer, Dependency Injection
- **Testing**: Manual con archivos de prueba

## 🤝 Contribución

### Estructura de Commits

- `feat:` Nueva funcionalidad
- `fix:` Corrección de errores
- `docs:` Documentación
- `refactor:` Refactorización de código
- `test:` Pruebas

### Cómo Contribuir

1. Fork el proyecto
2. Crear una rama feature (`git checkout -b feature/AmazingFeature`)
3. Commit los cambios (`git commit -m 'feat: Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request

## 📝 Notas Adicionales

### Manejo de Errores

El sistema implementa un manejo robusto de errores:
- Validación de formato de archivos
- Detección de duplicados
- Manejo de errores de base de datos
- Mensajes informativos al usuario

### Rendimiento

- Uso eficiente de memoria con interfaces
- Consultas SQL optimizadas
- Procesamiento por lotes para archivos grandes

### Seguridad

- Validación de entrada de datos
- Prevención de inyección SQL con prepared statements
- Manejo seguro de archivos

## 📞 Soporte

Para reportar bugs o solicitar nuevas funcionalidades, por favor crear un issue en el repositorio del proyecto.

---

*Desarrollado como proyecto académico para el aprendizaje de arquitecturas de software y programación en Go.*
