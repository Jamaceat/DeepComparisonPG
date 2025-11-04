# 🔍 Deep Database Comparator

Una aplicación avanzada en Go para comparar profundamente datos entre dos bases de datos PostgreSQL, incluyendo análisis completo de foreign keys, detección de diferencias a nivel de registro, y análisis de referencias cruzadas.

> **📁 Organización de Archivos**: Todos los archivos generados (JSON, SQL) se almacenan automáticamente en la carpeta `generated/` para mantener el proyecto organizado.

## 📚 Índice

### 🚀 **Inicio Rápido**
- [✨ Características Principales](#-características-principales)
- [🏗️ Arquitectura del Proyecto](#️-arquitectura-del-proyecto)
- [🚀 Instalación y Configuración](#-instalación-y-configuración)
- [⚙️ Configuración](#configuración)

### 🔧 **Uso de la Aplicación**
- [🆔 Decodificación Automática de UUIDs](#-decodificación-automática-de-uuids)
- [� Uso del Sistema](#-uso-del-sistema)
- [�📋 Opciones Disponibles](#-opciones-disponibles)
- [📚 Ejemplos de Uso](#-ejemplos-de-uso)
- [�️ Exclusión de Columnas por Archivo](#️-exclusión-de-columnas-por-archivo)
- [�🚀 Optimización de Rendimiento](#-optimización-de-rendimiento)

### 📊 **Análisis y Resultados**
- [� Formato de Salida JSON](#-formato-de-salida-json)
- [� Análisis de Referencias](#-análisis-de-referencias-nuevo)
- [📄 Estructura del Archivo comparison_result_references.json](#-estructura-del-archivo-comparison_result_referencesjson)
- [🎯 Casos de Uso del Análisis de Referencias](#-casos-de-uso-del-análisis-de-referencias)

### 🛠️ **Configuración Avanzada**
- [Ejemplo de Escenario](#ejemplo-de-escenario)
- [Algoritmo de Matching](#algoritmo-de-matching)
- [Desarrollo](#desarrollo)
- [⚠️ Consideraciones y Limitaciones](#️-consideraciones-y-limitaciones)

### 🆘 **Soporte y Resolución de Problemas**
- [🛠️ Troubleshooting y Mejores Prácticas](#️-troubleshooting-y-mejores-prácticas)
- [🤝 Contribuir al Proyecto](#-contribuir-al-proyecto)
- [📋 Roadmap](#-roadmap)
- [📜 Versionado](#-versionado)
- [📄 Licencia](#-licencia)
- [🆘 Soporte y Comunidad](#-soporte-y-comunidad)

### ⚡ **Navegación Rápida**
- **🏃‍♂️ [Empezar YA](#-instalación)** - Instalación y primer uso
- **💡 [Ejemplos Prácticos](#-comparación-de-datos)** - Comandos listos para usar
- **🔍 [Análisis de Referencias](#-análisis-de-referencias)** - Nueva funcionalidad
- **🚀 [Rendimiento](#-procesamiento-concurrente)** - Optimización con workers
- **❓ [Problemas](#️-troubleshooting-y-mejores-prácticas)** - Solución de errores

### 🎯 **Casos de Uso Frecuentes**
| Necesidad | Ir a Sección | Comando Rápido |
|-----------|--------------|----------------|
| **Comparar tabla básica** | [Comparación de Datos](#-comparación-de-datos) | `./deepComparator -table=mi_tabla -verbose` |
| **Excluir columnas audit** | [Exclusión de Columnas](#️-exclusión-de-columnas-por-archivo) | `./deepComparator -table=mi_tabla -exclude-from-file` |
| **Ver qué referencia una tabla** | [Análisis de Referencias](#-análisis-de-referencias) | `./deepComparator -find-references -table=mi_tabla` |
| **🆕 Encontrar dónde se usa un ID** | [FK References](#-análisis-de-fk-references-nuevo) | `./deepComparator -table=concepts -id="89" -analyze-fk-references` |
| **🆕 Generar script UPDATE FK** | [Update Script](#-generación-de-scripts-update-nuevo) | `./deepComparator -table=concepts -source-db=db1 -id-target=89 -id-destination=90 -generate-update-script` |
| **UUIDs legibles** | [Decodificación UUID](#-decodificación-automática-de-uuids) | `./deepComparator -table=mi_tabla -decode-uuids=true` |
| **Mejorar rendimiento** | [Optimización](#-optimización-de-rendimiento) | `./deepComparator -table=mi_tabla -max-workers=8` |
| **Solucionar errores** | [Troubleshooting](#️-troubleshooting-y-mejores-prácticas) | Ver sección de errores comunes |

**🆙 [Volver arriba](#-deep-database-comparator)** ↑

---

## ✨ Características Principales

### **🔄 Comparación Profunda de Datos**
- **Matching Inteligente**: Empareja registros basándose en el contenido, no en IDs que pueden diferir
- **Análisis de Foreign Keys**: Incluye datos completos de tablas referenciadas en los resultados
- **Detección Granular**: Identifica diferencias específicas por columna con contexto completo
- **Exclusión Inteligente**: Sistema configurable para omitir columnas de auditoría o metadatos
- **🆔 Decodificación UUID**: Convierte automáticamente UUIDs codificados en Base64 a formato legible para facilitar búsquedas en BD

### **🎯 Análisis de Referencias (Nuevo)**
- **Mapeo Completo**: Encuentra todas las tablas que referencian una tabla/columna específica  
- **Análisis Cruzado**: Compara valores referenciados entre ambas bases de datos
- **Categorización**: Clasifica referencias como comunes, solo en DB1, o solo en DB2
- **Auditoría de Integridad**: Detecta referencias huérfanas o inconsistencias

### **🔍 Análisis de FK References (Nuevo)**
- **Búsqueda por ID**: Encuentra todas las tablas que referencian un ID específico como foreign key
- **Soporte Universal**: Funciona con IDs numéricos y UUIDs  
- **Conteo Preciso**: Cuenta matches exactos en ambas bases de datos
- **Muestras de Datos**: Incluye samples de las referencias encontradas
- **Salida Específica**: Archivo `id_matches_tables.json` dedicado

### **🔧 Generación de Scripts UPDATE (Nuevo)**
- **Scripts SQL Seguros**: Genera transacciones completas con UPDATE y DELETE
- **Actualización de FK**: Reemplaza automáticamente todas las referencias de foreign key
- **Selección de BD**: Elige qué base de datos usar como fuente para el análisis
- **Validación Completa**: Verifica constraints antes de generar el script
- **Documentación Integrada**: Scripts autocomentados con información de origen y destino

### **⚙️ Configuración Avanzada**
- **Exclusión por Archivos**: Sistema basado en archivos para omitir columnas específicas
- **Criterios Personalizados**: Control granular sobre qué columnas incluir/excluir
- **Múltiples Esquemas**: Soporte para esquemas específicos de PostgreSQL
- **Configuración por Entorno**: Archivos .env separados para diferentes ambientes
- **Organización Automática**: Todos los archivos generados se almacenan en `generated/`

### **📊 Salida Estructurada**
- **JSON Detallado**: Reportes completos con información de foreign keys
- **Organización Automática**: Todos los archivos se generan en carpeta `generated/`
- **Dos Modos**: Comparación normal y análisis de referencias
- **Métricas Completas**: Estadísticas detalladas de matches, diferencias y referencias
- **Información Contextual**: Datos completos de filas referenciadas, no solo IDs

## 🏗️ Arquitectura del Proyecto

```
deepComparator/
├── cmd/
│   └── main.go                 # CLI y punto de entrada principal
├── pkg/
│   ├── config/                 # 🔧 Manejo de configuración
│   │   └── config.go          #     Variables de entorno y validación
│   ├── database/              # 🗄️  Conexión y operaciones de PostgreSQL
│   │   └── database.go        #     Conexiones, queries, y manejo de FK
│   ├── comparator/            # 🔍 Lógica de comparación y análisis
│   │   └── comparator.go      #     Algoritmos de matching y análisis FK
│   └── models/                # 📋 Estructuras de datos y tipos
│       └── models.go          #     Modelos para comparación y referencias
├── exclude_columns.txt        # 📝 Columnas a excluir por defecto
├── .env.example              # ⚙️  Ejemplo de configuración
├── go.mod                    # 📦 Dependencias de Go
├── go.sum                    # 🔒 Checksums de dependencias
└── README.md                 # 📖 Esta documentación
```

### **🧩 Componentes Clave**

- **`cmd/main.go`**: CLI con dos modos (comparación y análisis de referencias)
- **`pkg/comparator/`**: Algoritmos de matching inteligente y análisis profundo de FKs
- **`pkg/database/`**: Conexiones PostgreSQL y queries optimizadas para metadatos
- **`pkg/models/`**: Estructuras para comparación, referencias y configuración
- **`exclude_columns.txt`**: Lista configurable de columnas a omitir (auditoría, etc.)

## 🚀 Instalación y Configuración

### **Prerrequisitos**

- **Go 1.21+**: Lenguaje de programación
- **PostgreSQL**: Acceso a dos bases de datos PostgreSQL
- **Permisos**: Lectura en las tablas a comparar y esquemas `information_schema`

### **Instalación**

```bash
# Clonar el repositorio
git clone <repository-url>
cd deepComparator

# Instalar dependencias
go mod tidy

# Compilar la aplicación
go build -o deepComparator ./cmd

# Verificar instalación
./deepComparator --help
```

### **Configuración Inicial**

```bash
# Copiar archivo de configuración
cp .env.example .env

# Editar configuración (ver sección siguiente)
nano .env

# Verificar conexiones
./deepComparator -table=pg_tables -schema=information_schema
```

## Configuración

### Variables de Entorno

Copia el archivo `.env.example` a `.env` y configura las conexiones a tus bases de datos:

```bash
cp .env.example .env
```

Edita el archivo `.env`:

```env
# Database 1 Configuration
DB1_HOST=localhost
DB1_PORT=5432
DB1_DATABASE=database1
DB1_USERNAME=postgres
DB1_PASSWORD=password
DB1_SSL_MODE=disable

# Database 2 Configuration  
DB2_HOST=localhost
DB2_PORT=5432
DB2_DATABASE=database2
DB2_USERNAME=postgres
DB2_PASSWORD=password
DB2_SSL_MODE=disable

# Output Configuration
OUTPUT_FORMAT=json
OUTPUT_FILE=comparison_result.json
LOG_LEVEL=info
```

## 💻 Uso del Sistema

### **🔄 Modo Comparación (Por Defecto)**
Compara datos entre dos bases de datos con análisis profundo de foreign keys.

```bash
./deepComparator -table=<nombre_tabla> [opciones]
```

### **🔍 Modo Análisis de Referencias**
Encuentra todas las tablas que referencian una tabla/columna específica.

```bash
./deepComparator -table=<nombre_tabla> -find-references [opciones]
```

### **🆔 Modo Análisis de FK References (Nuevo)**
Encuentra todas las tablas que referencian un ID específico como foreign key.

```bash
./deepComparator -table=<nombre_tabla> -id=<valor_id> -analyze-fk-references [opciones]
```

### **🔧 Modo Generación de Scripts UPDATE (Nuevo)**
Genera un script SQL para actualizar todas las foreign keys de un ID objetivo a un ID destino y eliminar el registro original.

```bash
./deepComparator -table=<nombre_tabla> [-source-db=<db1|db2>] -id-target=<id_origen> -id-destination=<id_destino> -generate-update-script [opciones]
```

### **📋 Opciones Disponibles**

| Opción | Descripción | Valor por defecto |
|--------|-------------|-------------------|
| `-table` | **Requerido**: Nombre de la tabla a comparar | - |
| `-schema` | Esquema de la tabla | `public` |
| `-env` | Archivo de configuración de entorno | `.env` |
| `-output` | Archivo de salida (sobrescribe la configuración del .env) | - |
| `-exclude` | Columnas a excluir de la comparación (separadas por comas) | - |
| `-include` | Columnas específicas a incluir (separadas por comas) | - |
| `-include-pk` | Incluir columnas de clave primaria en la comparación | `false` |
| `-exclude-from-file` | Excluir columnas desde archivo | `true` |
| `-exclude-file` | Archivo con columnas a excluir (una por línea) | `exclude_columns.txt` |
| `-show-exclude-columns` | Mostrar lista de columnas desde archivo de exclusión y salir | `false` |
| `-verbose` | Habilitar logging detallado | `false` |
| `-find-references` | **Nuevo**: Encontrar todas las referencias a una tabla/columna | `false` |
| `-target-column` | **Nuevo**: Columna objetivo para análisis de referencias | `id` |
| `-analyze-fk-references` | **🆕 Nuevo**: Encontrar tablas que referencian un ID específico | `false` |
| `-id` | **🆕 Nuevo**: ID específico a buscar en referencias FK (numérico o UUID) | - |
| `-generate-update-script` | **🆕 Nuevo**: Generar script SQL para actualizar FK y eliminar registro | `false` |
| `-source-db` | **🆕 Nuevo**: Base de datos fuente ('db1' o 'db2') para análisis de script | `db1` |
| `-id-target` | **🆕 Nuevo**: ID objetivo que será reemplazado | - |
| `-id-destination` | **🆕 Nuevo**: ID destino que reemplazará al objetivo | - |
| `-max-workers` | **Nuevo**: Número máximo de workers concurrentes | `4` |
| `-decode-uuids` | **Nuevo**: Decodificar UUIDs Base64 para facilitar búsquedas en BD | `true` |

### **📚 Ejemplos de Uso**

#### **🔧 Configuración y Verificación**

```bash
# Ver columnas que se excluyen automáticamente
./deepComparator -show-exclude-columns

# Verificar conexión a las bases de datos
./deepComparator -table=pg_tables -schema=information_schema -verbose
```

#### **🔄 Comparación de Datos**

```bash
# Comparación básica con exclusiones automáticas (→ generated/comparison_result.json)
./deepComparator -table=billing_model -verbose

# Incluir todas las columnas (sin exclusiones)
./deepComparator -table=billing_model -exclude-from-file=false -verbose

# Usar archivo de exclusiones personalizado
./deepComparator -table=billing_model -exclude-file="custom_exclude.txt" -verbose

# Excluir columnas específicas adicionales
./deepComparator -table=billing_model -exclude="notes,temp_field" -verbose

# Comparar solo columnas específicas
./deepComparator -table=billing_model -include="name,status,amount" -verbose

# Incluir claves primarias en la comparación
./deepComparator -table=billing_model -include-pk=true -verbose

# Optimización de rendimiento con workers concurrentes
./deepComparator -table=billing_model -max-workers=8 -verbose

# Comparación rápida para bases de datos grandes
./deepComparator -table=large_table -max-workers=16 -exclude-from-file=true -verbose

# Especificar esquema y archivo de salida
./deepComparator -table=users -schema=auth -output=user_comparison.json -verbose
```

#### **🔍 Análisis de Referencias**

```bash
# Encontrar todas las tablas que referencian concepts.id
./deepComparator -table=concepts -find-references -verbose

# Analizar referencias a una columna específica
./deepComparator -table=users -target-column=user_id -find-references -verbose

# Guardar análisis en archivo específico (→ generated/formula_refs.json)
./deepComparator -table=formula -find-references -output=formula_refs.json

# Analizar referencias en esquema específico
./deepComparator -table=categories -schema=catalog -find-references
```

#### **🎯 Casos de Uso Avanzados**

```bash
# Migración: verificar integridad antes de deploy
./deepComparator -table=products -exclude-from-file=false -include-pk=true

# Auditoría: encontrar diferencias en configuraciones
./deepComparator -table=system_config -include="key,value,enabled" -verbose

# Limpieza: encontrar referencias antes de eliminar datos
./deepComparator -table=old_categories -find-references -verbose

# Debug: comparar con todas las columnas para troubleshooting  
./deepComparator -table=transactions -exclude-file="/dev/null" -verbose
```

## 🛡️ Exclusión de Columnas por Archivo

### ¿Qué Columnas Excluir?

Puedes excluir cualquier columna que no sea relevante para tu comparación. Comúnmente se excluyen:

- **Columnas de auditoría**: `created_at`, `updated_at`, `created_by`, `updated_by`
- **Columnas de versioning**: `version`, `revision`, `row_version`
- **Columnas del sistema**: `last_login`, `session_id`, `ip_address`
- **Columnas temporales**: `temp_field`, `migration_flag`, `batch_id`
- **Cualquier columna que definas**: Tienes control total

### Archivo de Configuración

**Por defecto**, la aplicación usa el archivo `exclude_columns.txt` que contiene más de 50 columnas comunes que normalmente no son relevantes para comparaciones de datos. **TÚ PUEDES MODIFICAR ESTE ARCHIVO** según tus necesidades.

### Ver Qué Columnas se Excluyen

```bash
# Ver todas las columnas que se excluyen desde el archivo
./deepComparator -show-exclude-columns
```

### Personalización Total

**Por defecto**, la aplicación lee las columnas desde el archivo `exclude_columns.txt`. Este archivo es **completamente editable** y puedes:

- ✅ **Agregar** cualquier columna específica de tu proyecto
- ✅ **Quitar** columnas que sí quieres comparar  
- ✅ **Crear** múltiples archivos para diferentes tipos de tablas
- ✅ **Usar archivos vacíos** para no excluir nada

### Ver y Personalizar Columnas

```bash
# Ver qué columnas se excluyen actualmente
./deepComparator -show-exclude-columns

# Usar tu propio archivo personalizado
./deepComparator -table=billing_model -exclude-file="mi_archivo.txt"

# No excluir nada del archivo
./deepComparator -table=billing_model -exclude-from-file=false

# Editar el archivo por defecto
nano audit_columns.txt

# Usar un archivo personalizado
./deepComparator -table=billing_model -audit-file="mi_config.txt"

# No excluir ninguna columna (archivo vacío)
./deepComparator -table=billing_model -audit-file="/dev/null"
```

### Formato del Archivo

```
# audit_columns.txt
# Líneas que empiecen con # son comentarios
# Una columna por línea

created_at
updated_at
created_by
# Agregar las columnas específicas de tu proyecto
mi_campo_auditoria
batch_processed_at
```

### Ventajas

- **Comparaciones más relevantes**: Se enfoca en datos de negocio, no en metadatos técnicos
- **Menos ruido**: Evita falsos positivos por diferencias en timestamps o versioning
- **Configuración flexible**: Puedes agregar tus propias columnas de auditoría
- **Control total**: Puedes desactivar la exclusión cuando sea necesario

## **🆔 Decodificación Automática de UUIDs**
*📚 [Volver al índice](#-índice)*

### **¿Qué es la Decodificación UUID?**

PostgreSQL a menudo almacena UUIDs en formato Base64, lo que los hace difíciles de leer y buscar. Por ejemplo:
- **Base64**: `MDA5NTZjNGYtNDgzNS1iNjk4LTJkM2QtMDVlNWRjYzNlNzBl`
- **UUID Legible**: `00956c4f-4835-b698-2d3d-05e5dcc3e70e`

### **Funcionalidad Automática**

**Por defecto**, la aplicación detecta y decodifica automáticamente UUIDs codificados en Base64 tanto en:

✅ **Diferencias de Columnas**: Los valores en `column_differences` aparecen como UUIDs legibles
✅ **Referencias de Foreign Keys**: Los valores en `db1_references`, `db2_references`, etc. se decodifican automáticamente
✅ **Análisis de Referencias**: Todos los valores referenciados se muestran en formato UUID estándar

### **Configuración**

```bash
# Por defecto está ACTIVADO (recomendado)
./deepComparator -table=billing_model -verbose

# Deshabilitar decodificación (para debugging)
./deepComparator -table=billing_model -decode-uuids=false -verbose

# Análisis de referencias con decodificación
./deepComparator -find-references -table=concepts -decode-uuids=true -verbose
```

### **Beneficios**

1. **🔍 Búsquedas Fáciles**: Puedes copiar los UUIDs directamente del JSON y usarlos en consultas SQL
2. **📋 Legibilidad**: Los reportes son más fáciles de entender y revisar
3. **🛠️ Debugging**: Facilita la identificación y resolución de problemas de datos
4. **📊 Auditoría**: Los UUIDs legibles simplifican los procesos de auditoría

### **Ejemplo Práctico**

**Antes (con `-decode-uuids=false`):**
```json
{
  "column_differences": [
    {
      "column_name": "id",
      "db1_value": "MDA5NTZjNGYtNDgzNS1iNjk4LTJkM2QtMDVlNWRjYzNlNzBl",
      "db2_value": "MTVkNjZhZDctOTllNC1iY2Q5LWFiMzUtZjY4YTEwMDc5YmJh"
    }
  ]
}
```

**Después (con `-decode-uuids=true`, por defecto):**
```json
{
  "column_differences": [
    {
      "column_name": "id", 
      "db1_value": "00956c4f-4835-b698-2d3d-05e5dcc3e70e",
      "db2_value": "15d66ad7-99e4-bcd9-ab35-f68a10079bba"
    }
  ]
}
```

**Consulta SQL Directa:**
```sql
-- Ahora puedes buscar directamente con el UUID legible
SELECT * FROM billing_model WHERE id = '00956c4f-4835-b698-2d3d-05e5dcc3e70e';
```

## **🚀 Optimización de Rendimiento**
*📚 [Volver al índice](#-índice)*

### **Procesamiento Concurrente**

La aplicación incluye un sistema de **workers concurrentes** que mejora significativamente el rendimiento, especialmente para:

- 📊 **Bases de datos grandes** con miles/millones de registros
- 🔗 **Múltiples foreign keys** que requieren análisis paralelo
- 📈 **Análisis de referencias** en múltiples tablas simultáneamente

### **Configuración de Workers**

```bash
# Configuración por defecto (4 workers)
./deepComparator -table=billing_model -verbose

# Optimización para bases de datos pequeñas (1-2 workers)
./deepComparator -table=billing_model -max-workers=2 -verbose

# Optimización para bases de datos medianas (4-8 workers)
./deepComparator -table=billing_model -max-workers=8 -verbose

# Optimización para bases de datos grandes (8-16 workers)
./deepComparator -table=large_table -max-workers=16 -verbose

# Análisis de referencias con alta concurrencia
./deepComparator -find-references -table=billing_model -max-workers=12 -verbose

# Comparación con decodificación UUID deshabilitada (para debugging)
./deepComparator -table=billing_model -decode-uuids=false -verbose
```

### **Operaciones Paralelas**

El sistema concurrente paraleliza las siguientes operaciones:

1. **📥 Fetch de datos**: Obtención simultánea de datos de ambas bases de datos
2. **🔗 Análisis de Foreign Keys**: Procesamiento paralelo de múltiples relaciones
3. **📋 Análisis de referencias**: Búsqueda concurrente en múltiples tablas referenciadoras
4. **⚡ Categorización de valores**: Procesamiento paralelo de comparaciones complejas

### **Recomendaciones de Rendimiento**

| Escenario | Tamaño de DB | Workers Recomendados | Comando |
|-----------|--------------|---------------------|----------|
| **Pequeña** | < 1K registros | 1-2 | `-max-workers=2` |
| **Mediana** | 1K-100K registros | 4-8 | `-max-workers=8` |
| **Grande** | 100K-1M registros | 8-16 | `-max-workers=16` |
| **Muy Grande** | > 1M registros | 12-24 | `-max-workers=24` |

### **Prueba de Rendimiento**

Incluye una herramienta de testing de rendimiento:

```bash
# Compilar y ejecutar test de rendimiento
cd cmd/performance_test
go run main.go
```

Este test compara el rendimiento con diferentes números de workers para ayudarte a encontrar la configuración óptima para tu entorno.

## Ejemplo de Escenario

Considerando el ejemplo que mencionaste:

**Tabla 1 - billing_model (ID: 26)**
```json
{
  "id": 26,
  "description": "Get interest discount about number 693 31 of JUL of 2025",
  "order": 3,
  "status": "ac",
  "end_date": "2025-12-29",
  "start_date": "2025-09-08",
  "concept_id": 6
}
```

**Tabla 2 - billing_model (ID: 27)**
```json
{
  "id": 27,
  "order": 3,
  "status": "ac", 
  "end_date": "2025-12-29",
  "concept_id": 6,
  "start_date": "2025-09-08"
}
```

### Comando de Comparación
```bash
./deepComparator -table=billing_model -exclude="id" -verbose
```

### Resultado Esperado

La aplicación:

1. **Emparejará** estos registros por contenido similar (excluyendo el ID)
2. **Detectará** que la `description` falta en el segundo registro
3. **Analizará** la foreign key `concept_id=6` en ambas bases de datos
4. **Comparará** los datos de la tabla referenciada para `concept_id=6`
5. **Generará** un reporte detallado con las diferencias encontradas

## 📊 Formato de Salida JSON

El resultado se genera en formato JSON estructurado. A continuación se explica cada sección:

### **Estructura Principal**

```json
{
  "table_name": "billing_model",           // Nombre de la tabla comparada
  "schema": "public",                      // Esquema de la tabla
  "timestamp": "2025-10-28T15:30:00Z",    // Momento de la comparación
  "total_rows_db1": 29,                   // Total de filas en DB1
  "total_rows_db2": 33,                   // Total de filas en DB2
  "matched_rows": 0,                      // Filas que hacen match entre DB1 y DB2
  "unmatched_rows": 62,                   // Filas que NO hacen match (29+33)
  "only_in_db1": [...],                   // Filas que solo están en DB1
  "only_in_db2": [...],                   // Filas que solo están en DB2
  "differences": [...],                   // Filas que hacen match pero tienen diferencias
  "foreign_key_results": [...]            // Resultados del análisis de foreign keys
}
```

### **🆔 Formato de Salida FK References - id_matches_tables.json**

Para el análisis de FK References (`-analyze-fk-references`), se genera un archivo específico:

```json
{
  "analysis_info": {
    "table_name": "concepts",
    "schema": "public", 
    "target_id": "89",
    "target_id_type": "integer",
    "timestamp": "2024-01-15T10:30:00Z",
    "execution_time_ms": 1250,
    "fk_constraints_found": 3
  },
  "fk_constraints": [
    {
      "table": "transactions",
      "schema": "public", 
      "column": "concept_id",
      "references": "concepts.id",
      "constraint_name": "fk_transactions_concept_id"
    }
  ],
  "reference_results": [
    {
      "referencing_table": {
        "name": "transactions",
        "schema": "public",
        "column": "concept_id"
      },
      "matches_found": {
        "total_db1": 7,
        "total_db2": 7,
        "matching_references": 7,
        "different_references": 0
      },
      "sample_matches": [
        {
          "db1_row": {
            "id": 145,
            "concept_id": 89,
            "amount": 1500.00,
            "status": "completed"
          },
          "db2_row": {
            "id": 145, 
            "concept_id": 89,
            "amount": 1500.00,
            "status": "completed"
          },
          "is_identical": true
        }
      ]
    }
  ],
  "summary": {
    "total_referencing_tables": 1,
    "total_references_found": 7,
    "all_references_match": true,
    "has_differences": false
  }
}
```

### **Sección `only_in_db1` / `only_in_db2`**

Contienen las filas completas que existen solo en una base de datos:

```json
"only_in_db1": [
  {
    "id": 26,
    "description": "Get interest discount about number 693",
    "order": 3,
    "status": "ac",
    "concept_id": 6,
    "created_at": "2023-02-27T21:00:31.059667Z"
  }
]
```

### **Sección `differences`**

Aparece cuando hay filas que hacen match pero tienen diferencias en algunas columnas:

```json
"differences": [
  {
    "row_identifier": "order:3|status:ac",     // Clave única generada para la fila
    "db1_row": { /* fila completa de DB1 */ },
    "db2_row": { /* fila completa de DB2 */ },
    "column_differences": [
      {
        "column_name": "description",
        "db1_value": "Descripción vieja",
        "db2_value": "Descripción nueva",
        "is_foreign_key": false,               // ¿Es esta columna una FK?
        "foreign_key_reference": {...}         // Datos referenciados (si es FK)
      }
    ]
  }
]
```

### **Sección `foreign_key_results`** - Análisis Profundo de FKs

Esta es la sección más importante para entender las relaciones:

```json
"foreign_key_results": [
  {
    "foreign_key": {
      "column_name": "formula_id",              // Columna FK en tabla principal
      "referenced_table": "formula",           // Tabla referenciada
      "referenced_schema": "public",           // Esquema de la tabla referenciada
      "referenced_column_name": "id",          // Columna referenciada (PK)
      "constraint_name": "fk_billing_formula"  // Nombre de la restricción FK
    },
    "comparison_result": {
      // Estadísticas de comparación de la tabla referenciada
      "table_name": "formula",
      "matched_rows": 5,        // Cuántas filas de 'formula' hacen match
      "unmatched_rows": 10,     // Cuántas filas no hacen match
      "only_in_db1": [...],     // Filas de 'formula' solo en DB1
      "only_in_db2": [...],     // Filas de 'formula' solo en DB2
      "differences": [...]      // Diferencias en filas de 'formula' que sí hacen match
    },
    "fk_references": [...]      // ¡DATOS REALES de las tablas referenciadas!
  }
]
```

### **Sección `fk_references`** - ⭐ La Más Importante

Contiene los **datos completos** de las filas referenciadas por las foreign keys:

```json
"fk_references": [
  {
    "foreign_key": {
      "column_name": "formula_id",
      "referenced_table": "formula",
      "constraint_name": "fk_billing_formula"
    },
    "db1_referenced": {
      // DATOS COMPLETOS de la fila referenciada en DB1
      "id": 19,
      "name": "Formula Saldo Concepto", 
      "description": "Obtiene el saldo actual del concepto",
      "formula": "SELECT balance FROM accounts WHERE id = ?",
      "version": "1.0",
      "created_at": "2021-09-29T14:33:41.06Z"
    },
    "db2_referenced": {
      // DATOS COMPLETOS de la fila referenciada en DB2
      "id": 19,
      "name": "Formula Saldo Concepto",
      "description": "Obtiene el saldo actual del concepto", 
      "formula": "SELECT balance * 1.1 FROM accounts WHERE id = ?",  // ¡Cambió!
      "version": "1.1",                                            // ¡Nueva versión!
      "created_at": "2021-09-29T14:33:41.06Z"
    },
    "referenced_diff": true    // ¿Son diferentes los datos referenciados?
  }
]
```

### **Significado de `referenced_diff`**

- **`"referenced_diff": false`** = Los datos de la fila referenciada son **IDÉNTICOS** en ambas bases de datos
- **`"referenced_diff": true`** = Los datos de la fila referenciada son **DIFERENTES** entre las bases de datos

### **¿Por Qué es Útil `fk_references`?**

Imagina que tienes:
- **Tabla principal**: `billing_model` con `formula_id = 19` en ambas DBs
- **A primera vista**: Parece que es la misma fórmula
- **En realidad**: La lógica de la fórmula cambió en DB2

Sin este análisis profundo, no te darías cuenta que aunque el ID es el mismo, **la fórmula actual es diferente** y puede producir resultados distintos.

### **Ejemplo Práctico Completo**

```json
{
  "table_name": "billing_model",
  "matched_rows": 25,
  "differences": [
    {
      "row_identifier": "order:3|status:ac",
      "column_differences": [
        {
          "column_name": "formula_id",
          "db1_value": 19,
          "db2_value": 19,
          "is_foreign_key": true,
          "foreign_key_reference": {
            "foreign_key": {
              "column_name": "formula_id",
              "referenced_table": "formula"
            },
            "db1_referenced": {
              "id": 19,
              "formula": "SELECT balance FROM accounts"
            },
            "db2_referenced": {
              "id": 19, 
              "formula": "SELECT balance * 1.1 FROM accounts"  // ¡Diferente!
            },
            "referenced_diff": true
          }
        }
      ]
    }
  ],
  "foreign_key_results": [
    {
      "foreign_key": {
        "column_name": "formula_id",
        "referenced_table": "formula"
      },
      "comparison_result": {
        "matched_rows": 1,
        "differences": [
          // Aquí verías las diferencias en la tabla 'formula'
        ]
      },
      "fk_references": [
        // Datos completos de todas las fórmulas referenciadas
      ]
    }
  ]
}
```

### **Casos de Uso**

1. **Migración de Datos**: Verificar que las relaciones se migraron correctamente
2. **Sincronización**: Detectar diferencias entre entornos de desarrollo y producción  
3. **Auditoría**: Encontrar cambios en configuraciones o datos maestros
4. **Debugging**: Entender por qué dos registros aparentemente iguales producen resultados diferentes

La clave está en que no solo comparamos los IDs de las foreign keys, sino **los datos reales** a los que apuntan esas FKs.

## 🔍 Análisis de Referencias (Nuevo)
*📚 [Volver al índice](#-índice)*

Además de la comparación de tablas, el sistema incluye una funcionalidad para encontrar todas las referencias a una tabla/columna específica.

### **¿Qué hace?**

Encuentra todas las tablas que tienen foreign keys apuntando a una tabla/columna específica y analiza los valores referenciados en ambas bases de datos.

### **Uso**

```bash
./deepComparator -table=<tabla> -target-column=<columna> -find-references [opciones]
```

### **Ejemplos**

```bash
# Encontrar todas las tablas que referencian concepts.id
./deepComparator -table=concepts -target-column=id -find-references -verbose

# Encontrar referencias a formula.id
./deepComparator -table=formula -find-references -verbose

# Especificar archivo de salida personalizado
./deepComparator -table=concepts -find-references -output=my_references.json

# Usar esquema específico
./deepComparator -table=users -schema=auth -find-references -verbose

# Análisis de referencias con optimización de rendimiento
./deepComparator -table=concepts -find-references -max-workers=12 -verbose

# Análisis masivo para tablas con muchas referencias
./deepComparator -table=main_catalog -find-references -max-workers=16 -verbose
```

#### **🆔 Análisis de FK References (Nuevo)**

```bash
# Encontrar todas las tablas que referencian el ID 89 de concepts
./deepComparator -table=concepts -id="89" -analyze-fk-references -verbose

# Análisis con UUID específico
./deepComparator -table=users -id="550e8400-e29b-41d4-a716-446655440000" -analyze-fk-references

# Guardar resultado en archivo específico  
./deepComparator -table=products -id="123" -analyze-fk-references -output=product_references.json

# Análisis con esquema específico
./deepComparator -table=categories -schema=catalog -id="45" -analyze-fk-references

# Optimización para tablas con muchas referencias FK
./deepComparator -table=concepts -id="89" -analyze-fk-references -max-workers=8 -verbose

# Análisis sin decodificación UUID (para debugging)
./deepComparator -table=accounts -id="encoded_uuid" -analyze-fk-references -decode-uuids=false
```

#### **🔧 Generación de Scripts UPDATE (Nuevo)**

```bash
# Generar script para mover FK del concepto ID 89 al ID 90 (→ generated/update_fk_references.sql)
./deepComparator -table=concepts -id-target=89 -id-destination=90 -generate-update-script -verbose

# Con base de datos específica
./deepComparator -table=concepts -source-db=db2 -id-target=89 -id-destination=90 -generate-update-script -verbose

# Script para reemplazar usuario con UUID específico 
./deepComparator -table=users -source-db=db2 -id-target="550e8400-e29b-41d4-a716-446655440000" -id-destination="660f9500-f39c-52e5-c827-116f6ee4f81f" -generate-update-script

# Generar script con nombre personalizado
./deepComparator -table=categories -source-db=db1 -id-target=25 -id-destination=30 -generate-update-script -output=migrate_category_25.sql

# Script para esquema específico
./deepComparator -table=products -schema=catalog -source-db=db1 -id-target=100 -id-destination=105 -generate-update-script -verbose

# Análisis complejo con optimización de workers
./deepComparator -table=main_entities -source-db=db1 -id-target=999 -id-destination=1000 -generate-update-script -max-workers=8 -verbose
```

### **Opciones Específicas**

| Opción | Descripción | Valor por defecto |
|--------|-------------|-------------------|
| `-find-references` | Activar modo de análisis de referencias | `false` |
| `-target-column` | Columna objetivo para encontrar referencias | `id` |

### **Archivo de Salida**

Por defecto genera `comparison_result_references.json` (o `<archivo_especificado>_references.json` si se usa `-output`).

### **📄 Estructura del Archivo `comparison_result_references.json`**

#### **Metadatos Principales**
```json
{
  "target_table": "banks",              // Tabla objetivo analizada
  "target_schema": "public",            // Esquema de la tabla objetivo
  "target_column": "id",               // Columna objetivo (por defecto "id")
  "timestamp": "2025-10-29T10:56:46Z", // Momento del análisis
  "total_references": 0,               // Total de valores referenciados activos
  "referencing_tables": 14,            // Número de tablas que tienen FK a la objetivo
  "references": [...]                  // Array detallado de cada referencia
}
```

#### **Detalle de Referencias**
Cada elemento en el array `references` contiene:

```json
{
  "table_name": "bank_accounts",                    // Tabla que referencia
  "schema": "public",                              // Esquema de la tabla referenciadora
  "column_name": "bank_id",                        // Columna FK en tabla referenciadora
  "constraint_name": "fk_bank_id",                 // Nombre de la restricción FK
  "db1_references": [1, 2, 3, 15, 22],            // Valores FK únicos en DB1
  "db2_references": [1, 2, 3, 15, 22, 25],        // Valores FK únicos en DB2
  "common_references": [1, 2, 3, 15, 22],         // Valores presentes en AMBAS DBs
  "only_in_db1": [],                               // Valores FK solo en DB1
  "only_in_db2": [25]                              // Valores FK solo en DB2
}
```

### **🔍 Significado de Cada Campo**

#### **Campos de Metadatos:**
| Campo | Tipo | Descripción | Ejemplo |
|-------|------|-------------|---------|
| `target_table` | string | Tabla principal que se está analizando | `"banks"` |
| `target_schema` | string | Esquema donde está la tabla objetivo | `"public"` |
| `target_column` | string | Columna objetivo (normalmente PK) | `"id"` |
| `timestamp` | string | Marca de tiempo ISO 8601 del análisis | `"2025-10-29T10:56:46Z"` |
| `total_references` | number | Suma de todos los valores referenciados activos | `847` |
| `referencing_tables` | number | Cantidad de tablas que tienen FK hacia la objetivo | `14` |

#### **Campos por Referencia:**
| Campo | Tipo | Descripción | Cuándo aparece |
|-------|------|-------------|----------------|
| `table_name` | string | Nombre de la tabla que contiene la FK | Siempre |
| `schema` | string | Esquema de la tabla referenciadora | Siempre |
| `column_name` | string | Nombre de la columna FK | Siempre |
| `constraint_name` | string | Nombre de la restricción de foreign key | Siempre |
| `db1_references` | array/null | Valores únicos de la FK encontrados en DB1 | `null` si tabla no existe en DB1 |
| `db2_references` | array/null | Valores únicos de la FK encontrados en DB2 | `null` si tabla no existe en DB2 |
| `common_references` | array/null | Valores FK que existen en AMBAS bases de datos | `null` si no hay coincidencias |
| `only_in_db1` | array/null | Valores FK que SOLO están en DB1 | `null` si no hay exclusivos |
| `only_in_db2` | array/null | Valores FK que SOLO están en DB2 | `null` si no hay exclusivos |

### **📊 Estados Posibles de los Datos**

#### **✅ Escenario Normal (Consistente)**
```json
{
  "db1_references": [1, 2, 3, 5],
  "db2_references": [1, 2, 3, 5],
  "common_references": [1, 2, 3, 5],
  "only_in_db1": [],
  "only_in_db2": []
}
```
**Interpretación**: Ambas DBs tienen exactamente las mismas referencias. ✅ Consistente.

#### **⚠️ Escenario con Diferencias**
```json
{
  "db1_references": [1, 2, 3, 5, 8],
  "db2_references": [1, 2, 3, 5, 9, 10],
  "common_references": [1, 2, 3, 5],
  "only_in_db1": [8],
  "only_in_db2": [9, 10]
}
```
**Interpretación**: 
- DB1 tiene referencia al ID `8` que no existe en DB2
- DB2 tiene referencias a IDs `9, 10` que no existen en DB1
- Ambas comparten referencias a IDs `1, 2, 3, 5`

#### **❌ Escenario de Tabla Inexistente**
```json
{
  "db1_references": null,
  "db2_references": [1, 2, 3],
  "common_references": null,
  "only_in_db1": null,
  "only_in_db2": null
}
```
**Interpretación**: La tabla referenciadora no existe en DB1, solo en DB2.

### **🎯 Casos de Uso del Análisis de Referencias**

#### **1. 🔍 Auditoría de Integridad Referencial**
```bash
# Verificar consistencia de referencias a tabla principal
./deepComparator -table=banks -find-references -verbose
```
**Objetivo**: Detectar referencias huérfanas o inconsistentes entre ambientes.

#### **2. 🚚 Migración Segura de Datos**
```bash
# Antes de eliminar registros, verificar impacto
./deepComparator -table=concepts -target-column=id -find-references
```
**Objetivo**: Identificar todas las tablas que se verían afectadas por cambios.

#### **3. 🧹 Limpieza de Datos**
```bash
# Encontrar referencias no utilizadas
./deepComparator -table=categories -find-references -max-workers=8
```
**Objetivo**: Detectar IDs que existen en tabla principal pero no tienen referencias.

#### **4. 📊 Análisis de Consistencia entre Ambientes**
```bash
# Comparar referencias entre producción y desarrollo
./deepComparator -table=users -schema=auth -find-references
```
**Objetivo**: Verificar que ambos ambientes tienen la misma estructura referencial.

#### **5. 🔄 Sincronización de Referencias**
```bash
# Análisis masivo de múltiples tablas
./deepComparator -table=main_catalog -find-references -max-workers=16
```
**Objetivo**: Identificar discrepancias para proceso de sincronización.

### **📋 Ejemplos Prácticos de Interpretación**

#### **Ejemplo 1: Análisis de Bancos**
```bash
./deepComparator -table=banks -find-references -verbose
```

**Resultado esperado:**
```json
{
  "target_table": "banks",
  "referencing_tables": 14,
  "references": [
    {
      "table_name": "bank_accounts",
      "column_name": "bank_id",
      "db1_references": [1, 2, 5, 8],
      "db2_references": [1, 2, 5, 8, 12],
      "only_in_db2": [12]
    }
  ]
}
```

**📖 Interpretación:**
- ✅ **14 tablas** referencian la tabla `banks`
- ⚠️ **Inconsistencia detectada**: DB2 tiene cuentas bancarias (`bank_accounts`) que referencian al banco ID `12`, pero ese banco puede no existir en DB1
- 🔧 **Acción requerida**: Verificar si el banco ID `12` debe ser migrado a DB1

#### **Ejemplo 2: Migración Segura**
```bash
# Antes de eliminar el concepto ID 25
./deepComparator -table=concepts -target-column=id -find-references
```

**Si el resultado muestra:**
```json
{
  "references": [
    {
      "table_name": "billing_model", 
      "only_in_db1": [25],
      "only_in_db2": []
    },
    {
      "table_name": "settlement_formulas",
      "common_references": [25]
    }
  ]
}
```

**📖 Interpretación:**
- ❌ **NO ELIMINAR** el concepto 25 todavía
- 📋 **Acción requerida**: 
  1. Limpiar referencia en `billing_model` de DB1
  2. Coordinar eliminación en `settlement_formulas` de ambas DBs
  3. Solo entonces eliminar el concepto 25

#### **Ejemplo 3: Detección de Datos Huérfanos**
```json
{
  "table_name": "user_permissions",
  "db1_references": [1, 5, 99],
  "db2_references": [1, 5],
  "only_in_db1": [99]
}
```

**📖 Interpretación:**
- ⚠️ **Posible dato huérfano**: El usuario ID `99` tiene permisos en DB1 pero no en DB2
- 🔍 **Investigar**: ¿El usuario 99 fue eliminado de DB2? ¿Debe eliminarse de DB1?

### **🚨 Señales de Alerta a Buscar**

| Patrón en JSON | Significado | Acción Recomendada |
|----------------|-------------|-------------------|
| `"only_in_db1": [...]` con valores | Referencias huérfanas en DB1 | Investigar y limpiar |
| `"only_in_db2": [...]` con valores | Referencias huérfanas en DB2 | Investigar y sincronizar |
| `"db1_references": null` | Tabla no existe en DB1 | Verificar migración de esquema |
| `"db2_references": null` | Tabla no existe en DB2 | Verificar despliegue de esquema |
| `"total_references": 0` | No hay datos activos | Normal para tablas vacías |
| `"referencing_tables": 0` | No hay FKs apuntando a tabla | Verificar estructura de FKs |

## 🔧 Generación de Scripts UPDATE (Nuevo)
*📚 [Volver al índice](#-índice)*

Esta funcionalidad genera scripts SQL seguros para actualizar todas las referencias de foreign key de un ID específico a otro ID y eliminar el registro original.

### **¿Para qué sirve?**

- **Consolidación de Datos**: Fusionar registros duplicados actualizando todas sus referencias
- **Migración de IDs**: Cambiar IDs manteniendo integridad referencial
- **Limpieza de Datos**: Reemplazar registros obsoletos por versiones actualizadas
- **Reasignación**: Cambiar la propiedad de registros a otros elementos

### **Funcionamiento**

1. **Descubre** todas las foreign keys que apuntan a la tabla objetivo
2. **Genera** comandos UPDATE para cada tabla que referencia el ID objetivo  
3. **Incluye** el comando DELETE del registro original
4. **Envuelve** todo en una transacción segura (BEGIN/COMMIT)

### **Uso Básico**

```bash
./deepComparator -table=<tabla> -source-db=<db1|db2> -id-target=<id_origen> -id-destination=<id_destino> -generate-update-script [opciones]
```

### **Parámetros Requeridos**

| Parámetro | Descripción | Ejemplo | Por defecto |
|-----------|-------------|---------|-------------|
| `-table` | Tabla que contiene el registro a migrar | `concepts` | - |
| `-source-db` | Base de datos para análisis ('db1' o 'db2') | `db1` | `db1` |
| `-id-target` | ID que será reemplazado | `89` | - |
| `-id-destination` | ID que reemplazará al objetivo | `90` | - |

### **Ejemplo Completo**

```bash
# Generar script para migrar concepto ID 89 -> 90 (usa db1 por defecto)
./deepComparator -table=concepts -id-target=89 -id-destination=90 -generate-update-script -verbose

# Con base de datos específica
./deepComparator -table=concepts -source-db=db2 -id-target=89 -id-destination=90 -generate-update-script -verbose
```

### **Resultado del Script**

El script generado (`update_fk_references.sql`) contendrá:

```sql
-- Generated FK Update Script
-- Target table: public.concepts
-- Update FK references from ID 89 to ID 90
-- Generated at: 2025-10-30 10:00:02
-- WARNING: Review this script before execution!

BEGIN;

-- Update foreign key references
-- Table: public.transactions, Column: concept_id
UPDATE public.transactions SET concept_id = 90 WHERE concept_id = 89;

-- Table: public.bill_items, Column: concept_id  
UPDATE public.bill_items SET concept_id = 90 WHERE concept_id = 89;

-- Delete original record
DELETE FROM public.concepts WHERE id = 89;

COMMIT;

-- Script execution completed
-- Verify results and check referential integrity
```

### **Opciones Adicionales**

| Opción | Descripción | Ejemplo |
|--------|-------------|---------|
| `-output` | Nombre personalizado para el archivo SQL | `-output=migrate_concept_89.sql` |
| `-schema` | Esquema específico (por defecto: public) | `-schema=catalog` |
| `-verbose` | Mostrar información detallada | `-verbose` |
| `-max-workers` | Workers para análisis (por defecto: 4) | `-max-workers=8` |

### **Casos de Uso Reales**

#### **🔄 1. Fusión de Registros Duplicados**

```bash
# Problema: Tienes dos conceptos que representan lo mismo
# Solución: Fusionar concepto 45 en concepto 42
./deepComparator -table=concepts -source-db=db1 -id-target=45 -id-destination=42 -generate-update-script
```

#### **🏗️ 2. Migración de Estructura**

```bash
# Problema: Necesitas cambiar la numeración de categorías
# Solución: Migrar categoría 100 -> 200
./deepComparator -table=categories -source-db=db2 -id-target=100 -id-destination=200 -generate-update-script -output=migrate_categories.sql
```

#### **👤 3. Reasignación de Propiedad**

```bash
# Problema: Un usuario se va y necesitas reasignar sus datos
# Solución: Transferir usuario A -> usuario B
./deepComparator -table=users -source-db=db1 -id-target="user-uuid-a" -id-destination="user-uuid-b" -generate-update-script
```

### **⚠️ Consideraciones de Seguridad**

#### **Antes de Ejecutar el Script:**

1. **✅ Respalda la base de datos** completamente
2. **✅ Revisa el script generado** línea por línea
3. **✅ Verifica que el ID destino existe** en la tabla objetivo
4. **✅ Confirma que no hay constraints** que puedan fallar
5. **✅ Ejecuta en un entorno de prueba** primero

#### **Durante la Ejecución:**

```bash
# Ejecutar el script en la base de datos
psql -d mi_base_datos -f update_fk_references.sql

# Verificar resultados
SELECT COUNT(*) FROM concepts WHERE id = 89; -- Debe ser 0
SELECT COUNT(*) FROM transactions WHERE concept_id = 90; -- Debe incluir las migraciones
```

#### **Después de la Ejecución:**

```bash
# Verificar integridad referencial
SELECT 
  t.table_name,
  COUNT(*) as orphaned_references
FROM information_schema.tables t
WHERE t.table_name IN ('transactions', 'bill_items')
  AND NOT EXISTS (
    SELECT 1 FROM concepts c 
    WHERE c.id = (SELECT concept_id FROM t.table_name LIMIT 1)
  );
```

### **🚨 Limitaciones**

- **Transacciones**: El script usa transacciones, pero en tablas muy grandes puede ser lento
- **Constraints**: No maneja constraints complejos automáticamente  
- **Cascadas**: No detecta DELETE/UPDATE CASCADE automáticos
- **Validación**: No valida que el ID destino exista antes de generar el script

### **📊 Información del Archivo Generado**

Por defecto, el archivo se llama `update_fk_references.sql`, pero puedes cambiarlo:

```bash
# Archivo por defecto
-generate-update-script                          → generated/update_fk_references.sql

# Archivo personalizado
-generate-update-script -output=migrate_data     → generated/migrate_data.sql

# Con nombre específico
-generate-update-script -output=migration_123    → generated/migration_123.sql
```

**📁 Organización de Archivos**: Todos los archivos se generan automáticamente en la carpeta `generated/` que se excluye del control de versiones via `.gitignore`.

## Algoritmo de Matching

La aplicación utiliza un algoritmo inteligente de matching que:

1. **Excluye automáticamente** las columnas de clave primaria (a menos que se especifique `-include-pk`)
2. **Genera una clave** basada en el contenido de las columnas relevantes
3. **Empareja registros** con claves idénticas
4. **Identifica diferencias** en registros emparejados
5. **Analiza foreign keys** recursivamente

## Desarrollo

### Agregar Nuevas Funcionalidades

El proyecto está estructurado de manera modular:

- **`pkg/models/`**: Agregar nuevas estructuras de datos
- **`pkg/database/`**: Extender funcionalidades de base de datos
- **`pkg/comparator/`**: Implementar nuevos algoritmos de comparación
- **`pkg/config/`**: Agregar nuevas opciones de configuración

### Ejecutar Tests (cuando estén disponibles)

```bash
go test ./...
```

### Formato y Linting

```bash
go fmt ./...
go vet ./...
```

## ⚠️ Consideraciones y Limitaciones

### **🎯 Compatibilidad**
- **Base de Datos**: Solo PostgreSQL (versión 9.1+)
- **Esquemas**: Requiere acceso al esquema `information_schema`
- **Tipos de Datos**: Soporte completo excepto tipos binarios complejos (bytea grandes)

### **🚀 Rendimiento**
- **Tablas Pequeñas** (< 1K filas): Instantáneo
- **Tablas Medianas** (1K - 100K filas): 1-10 segundos  
- **Tablas Grandes** (> 100K filas): Usar exclusiones para optimizar
- **Foreign Keys**: El análisis profundo puede incrementar el tiempo en tablas con muchas FKs

### **🔍 Profundidad de Análisis**
- **Foreign Keys**: Un nivel de profundidad (no recursivo infinito)
- **Matching**: Basado en contenido, requiere datos similares para emparejamiento
- **Exclusiones**: Columnas excluidas pueden afectar la precisión del matching

### **💾 Memoria**
- **Uso**: Carga tablas completas en memoria para comparación
- **Optimización**: Usa exclusión de columnas para reducir uso de memoria
- **Límite**: Recomendado para tablas que caben en RAM disponible

### **🔒 Seguridad**
- **Permisos**: Solo requiere permisos de lectura
- **Conexión**: Soporta SSL/TLS para conexiones seguras
- **Datos**: No modifica datos, solo lectura y análisis

## 🛠️ Troubleshooting y Mejores Prácticas
*📚 [Volver al índice](#-índice)*

### **🚨 Problemas Comunes**

#### **"No matched rows found"**
```bash
# Verificar que las columnas de matching existen en ambas DBs
./deepComparator -table=mytable -include-pk=true -verbose

# Revisar exclusiones automáticas
./deepComparator -show-exclude-columns

# Usar exclusión mínima para debugging
./deepComparator -table=mytable -exclude-from-file=false -verbose
```

#### **"Error connecting to database"**
```bash
# Verificar configuración
cat .env

# Probar conexión manual
psql -h $DB1_HOST -p $DB1_PORT -U $DB1_USERNAME -d $DB1_DATABASE

# Verificar SSL/TLS
export DB1_SSL_MODE=require
```

#### **"Too many differences found"**
```bash
# Usar exclusiones más agresivas
./deepComparator -table=mytable -exclude="created_at,updated_at,version"

# Comparar solo columnas críticas
./deepComparator -table=mytable -include="name,status,key_field"
```

#### **"UUIDs aparecen codificados en Base64"**
```bash
# Verificar que decodificación esté habilitada (por defecto lo está)
./deepComparator -table=mytable -decode-uuids=true -verbose

# Para debugging, deshabilitar decodificación temporalmente
./deepComparator -table=mytable -decode-uuids=false -verbose

# En análisis de referencias
./deepComparator -find-references -table=mytable -decode-uuids=true
```

### **✅ Mejores Prácticas**

#### **📊 Para Migraciones**
1. **Ejecutar con `-include-pk=true`** para verificar IDs
2. **Usar `-exclude-from-file=false`** para comparación completa
3. **Analizar referencias** antes de migrar datos maestros
4. **Documentar diferencias** encontradas para seguimiento

#### **🔍 Para Auditorías**
1. **Configurar exclusiones** específicas por tipo de tabla

#### **🆔 Para Análisis FK References**
1. **Identificar registros críticos**: Usar con registros maestros importantes
2. **Verificar integridad**: Asegurar que todas las referencias existen
3. **Análisis de impacto**: Ver qué se afecta antes de eliminar datos
4. **Debugging**: Encontrar dónde se usa un ID específico

## 🎯 Casos de Uso Prácticos

### **🔍 1. Análisis de Impacto antes de Eliminar**

```bash
# Antes de eliminar el concepto ID=89, ver qué tablas lo referencian
./deepComparator -table=concepts -id="89" -analyze-fk-references -verbose

# Revisar el resultado
cat id_matches_tables.json | jq '.reference_results[].referencing_table.name'
```

### **🚚 2. Migración de Datos Maestros**

```bash
# 1. Verificar consistencia del maestro
./deepComparator -table=concepts -verbose

# 2. Analizar impacto de registros específicos
./deepComparator -table=concepts -id="key_concept_id" -analyze-fk-references

# 3. Comparar tablas dependientes
./deepComparator -table=transactions -include="concept_id,amount,status" -verbose
```

### **🔎 3. Debugging de Integridad Referencial**

```bash
# Encontrar todas las referencias a un UUID específico
./deepComparator -table=users -id="550e8400-e29b-41d4-a716-446655440000" -analyze-fk-references -decode-uuids=true

# Ver resultado estructurado
cat id_matches_tables.json | jq '.reference_results[] | {table: .referencing_table.name, matches: .matches_found.total_db1}'
```

### **📊 4. Auditoría Completa de Sistema**

```bash
#!/bin/bash
# Script para auditar múltiples tablas críticas

# 1. Tablas maestras principales
for table in "concepts" "users" "categories"; do
    echo "Comparing $table..."
    ./deepComparator -table=$table -verbose -output="${table}_comparison.json"
done

# 2. Análisis de referencias para registros clave
./deepComparator -table=concepts -id="main_concept" -analyze-fk-references -output="main_concept_references.json"

# 3. Verificar foreign keys
./deepComparator -table=transactions -find-references -verbose
```

### **🔧 5. Fusión de Registros Duplicados**

```bash
# 1. Analizar impacto del registro duplicado
./deepComparator -table=concepts -id="89" -analyze-fk-references -verbose

# 2. Generar script de fusión
./deepComparator -table=concepts -source-db=db1 -id-target=89 -id-destination=90 -generate-update-script -output=merge_concept_89_to_90.sql

# 3. Revisar y ejecutar (después de backup!)
cat merge_concept_89_to_90.sql
psql -d database -f merge_concept_89_to_90.sql
```

### **🏗️ 6. Migración Segura de Datos**

```bash
#!/bin/bash
# Workflow completo para migrar usuario UUID

OLD_USER="550e8400-e29b-41d4-a716-446655440000"  
NEW_USER="660f9500-f39c-52e5-c827-116f6ee4f81f"

# 1. Verificar que el usuario destino existe
echo "Checking destination user exists..."
./deepComparator -table=users -id="$NEW_USER" -analyze-fk-references

# 2. Analizar impacto del usuario origen
echo "Analyzing source user impact..."
./deepComparator -table=users -id="$OLD_USER" -analyze-fk-references -output="user_${OLD_USER}_impact.json"

# 3. Generar script de migración
echo "Generating migration script..."
./deepComparator -table=users -source-db=db1 -id-target="$OLD_USER" -id-destination="$NEW_USER" -generate-update-script -output="migrate_user_${OLD_USER}.sql"

echo "Migration script ready: migrate_user_${OLD_USER}.sql"
echo "Review before execution!"
```
2. **Usar análisis de referencias** para mapear dependencias
3. **Ejecutar comparaciones regulares** en datos críticos
4. **Archivar resultados** para análisis histórico

#### **🚀 Para Rendimiento**
1. **Excluir columnas innecesarias** (logs, timestamps, etc.)
2. **Usar esquemas específicos** en lugar de `public`
3. **Ejecutar en horarios de baja carga** para tablas grandes
4. **Monitorear uso de memoria** en tablas muy grandes

#### **🎯 Para Debugging**
1. **Empezar con `-verbose`** para entender el proceso
2. **Usar exclusión mínima** para encontrar problemas de matching
3. **Comparar tablas pequeñas primero** para validar configuración
4. **Revisar logs de PostgreSQL** si hay errores de conexión

## 🤝 Contribuir al Proyecto

### **🛠️ Desarrollo Local**

```bash
# Fork y clonar el repositorio
git clone https://github.com/tu-usuario/deepComparator.git
cd deepComparator

# Instalar dependencias de desarrollo
go mod tidy

# Ejecutar tests (cuando estén disponibles)
go test ./...

# Verificar formato y linting
go fmt ./...
go vet ./...
```

### **📝 Proceso de Contribución**

1. **Fork** el repositorio
2. **Crear rama** temática: `git checkout -b feature/nueva-funcionalidad`
3. **Desarrollar** con tests y documentación
4. **Commit** siguiendo convenciones: `git commit -m "feat: agregar análisis de índices"`
5. **Push** a tu fork: `git push origin feature/nueva-funcionalidad`  
6. **Crear Pull Request** con descripción detallada

### **🎯 Áreas de Contribución**

- **🔍 Nuevos Tipos de Análisis**: Índices, triggers, funciones almacenadas
- **🚀 Optimizaciones**: Algoritmos de matching más eficientes  
- **🧪 Testing**: Suite de tests unitarios y de integración
- **📊 Formatos**: Exportación a Excel, CSV, HTML
- **🗄️ Bases de Datos**: Soporte para MySQL, SQL Server
- **🎨 UI**: Interfaz web o desktop para el comparador

## 📋 Roadmap

### **v2.0 - Próximas Funcionalidades**
- [ ] **Validación de Scripts**: Verificar que ID destino existe antes de generar script
- [ ] **Scripts con Rollback**: Generar scripts de reversión automáticos
- [ ] **Batch Operations**: Procesar múltiples migraciones en lote  
- [ ] **Análisis de Índices**: Comparar índices, constraints y triggers
- [ ] **Comparación Incremental**: Solo analizar cambios desde última ejecución
- [ ] **Cache Inteligente**: Almacenar resultados para re-ejecuciones rápidas
- [ ] **Filtros Avanzados**: Condiciones WHERE para limitar datos a comparar
- [ ] **Reportes HTML**: Salida visual para presentaciones

### **v2.1 - Integraciones**
- [ ] **CI/CD**: Plugins para GitLab, GitHub Actions, Jenkins
- [ ] **APIs REST**: Endpoint HTTP para integraciones
- [ ] **Webhooks**: Notificaciones automáticas de diferencias
- [ ] **Slack/Teams**: Integración con herramientas de comunicación

## 📜 Versionado

Este proyecto usa [Semantic Versioning](https://semver.org/):

- **MAJOR**: Cambios incompatibles en la API
- **MINOR**: Nueva funcionalidad compatible con versiones anteriores  
- **PATCH**: Corrección de bugs compatibles

**Versión Actual**: `v1.4.1`
- ✅ Comparación profunda de datos con foreign keys
- ✅ Análisis de referencias cruzadas  
- ✅ Exclusión configurable de columnas
- ✅ Salida JSON estructurada
- ✅ Decodificación automática de UUIDs Base64 para facilitar búsquedas en BD
- ✅ Procesamiento concurrente optimizado con workers configurables
- ✅ **🆕 Nuevo**: Análisis de FK References por ID específico con conteo preciso
- ✅ **🆕 Nuevo**: Generación de scripts SQL UPDATE para migrar foreign keys de manera segura
- ✅ **🗂️ Nuevo**: Organización automática de archivos en carpeta `generated/`

## 📄 Licencia

Este proyecto está bajo la **Licencia MIT**. Ver el archivo `LICENSE` para más detalles.

## 🆘 Soporte y Comunidad

### **📞 Obtener Ayuda**
- **Issues**: [GitHub Issues](https://github.com/owner/deepComparator/issues) para bugs y features
- **Discusiones**: [GitHub Discussions](https://github.com/owner/deepComparator/discussions) para preguntas
- **Email**: [soporte@deepcomparator.com](mailto:soporte@deepcomparator.com)

### **🐛 Reportar Bugs**
Incluir en el issue:
1. **Versión** de deepComparator (`./deepComparator --version`)
2. **Sistema operativo** y versión de Go
3. **Configuración** (sin passwords): `.env` y exclusiones usadas
4. **Comando exacto** ejecutado
5. **Salida completa** del error
6. **Comportamiento esperado** vs actual

### **💡 Solicitar Funcionalidades**
1. **Describir el caso de uso** detalladamente
2. **Explicar el beneficio** para otros usuarios
3. **Proponer implementación** si tienes ideas técnicas
4. **Agregar ejemplos** de cómo se usaría

---

### 🧭 **Navegación Final**
- **🏠 [Volver al Inicio](#-deep-database-comparator)** ↑
- **📚 [Ver Índice Completo](#-índice)** 📋
- **🚀 [Instalación Rápida](#-instalación-y-configuración)** ⚡
- **📖 [Ejemplos de Uso](#-ejemplos-de-uso)** 💡
- **🛠️ [Troubleshooting](#️-troubleshooting-y-mejores-prácticas)** 🔧

**⭐ Si este proyecto te es útil, considera darle una estrella en GitHub para ayudar a otros desarrolladores a encontrarlo.**