# 🔍 Deep Database Comparator

Una aplicación avanzada en Go para comparar profundamente datos entre dos bases de datos PostgreSQL, incluyendo análisis completo de foreign keys, detección de diferencias a nivel de registro, y análisis de referencias cruzadas.

## ✨ Características Principales

### **🔄 Comparación Profunda de Datos**
- **Matching Inteligente**: Empareja registros basándose en el contenido, no en IDs que pueden diferir
- **Análisis de Foreign Keys**: Incluye datos completos de tablas referenciadas en los resultados
- **Detección Granular**: Identifica diferencias específicas por columna con contexto completo
- **Exclusión Inteligente**: Sistema configurable para omitir columnas de auditoría o metadatos

### **🎯 Análisis de Referencias (Nuevo)**
- **Mapeo Completo**: Encuentra todas las tablas que referencian una tabla/columna específica  
- **Análisis Cruzado**: Compara valores referenciados entre ambas bases de datos
- **Categorización**: Clasifica referencias como comunes, solo en DB1, o solo en DB2
- **Auditoría de Integridad**: Detecta referencias huérfanas o inconsistencias

### **⚙️ Configuración Avanzada**
- **Exclusión por Archivos**: Sistema basado en archivos para omitir columnas específicas
- **Criterios Personalizados**: Control granular sobre qué columnas incluir/excluir
- **Múltiples Esquemas**: Soporte para esquemas específicos de PostgreSQL
- **Configuración por Entorno**: Archivos .env separados para diferentes ambientes

### **📊 Salida Estructurada**
- **JSON Detallado**: Reportes completos con información de foreign keys
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
# Comparación básica con exclusiones automáticas
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

# Especificar esquema y archivo de salida
./deepComparator -table=users -schema=auth -output=user_comparison.json -verbose
```

#### **🔍 Análisis de Referencias**

```bash
# Encontrar todas las tablas que referencian concepts.id
./deepComparator -table=concepts -find-references -verbose

# Analizar referencias a una columna específica
./deepComparator -table=users -target-column=user_id -find-references -verbose

# Guardar análisis en archivo específico
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
```

### **Opciones Específicas**

| Opción | Descripción | Valor por defecto |
|--------|-------------|-------------------|
| `-find-references` | Activar modo de análisis de referencias | `false` |
| `-target-column` | Columna objetivo para encontrar referencias | `id` |

### **Archivo de Salida**

Por defecto genera `match_reference_result.json` (o `<nombre>_references.json` si se especifica un archivo base).

### **Formato JSON de Salida**

```json
{
  "target_table": "concepts",
  "target_schema": "public", 
  "target_column": "id",
  "timestamp": "2025-10-29T09:51:11Z",
  "total_references": 510,
  "referencing_tables": 6,
  "references": [
    {
      "table_name": "related_concepts",
      "schema": "public",
      "column_name": "concept_id",
      "constraint_name": "fk_related_concepts_c",
      "db1_references": [1, 2, 6, 7, 13, 15, 16, 18, 20, 21],
      "db2_references": [1, 2, 6, 7, 13, 15, 16, 18, 20, 21],
      "common_references": [1, 2, 6, 7, 13, 15, 16, 18, 20, 21],
      "only_in_db1": [],
      "only_in_db2": []
    }
  ]
}
```

### **Campos del JSON**

- **`target_table/schema/column`**: Tabla/esquema/columna objetivo analizada
- **`total_references`**: Total de valores referenciados encontrados
- **`referencing_tables`**: Número de tablas que referencian la tabla objetivo
- **`references`**: Array con detalles de cada tabla referenciadora

**Por cada referencia:**
- **`table_name/schema`**: Tabla que contiene la foreign key
- **`column_name`**: Columna que es foreign key
- **`constraint_name`**: Nombre de la restricción FK
- **`db1_references`**: Valores únicos encontrados en DB1
- **`db2_references`**: Valores únicos encontrados en DB2  
- **`common_references`**: Valores que existen en ambas DBs
- **`only_in_db1`**: Valores que solo están en DB1
- **`only_in_db2`**: Valores que solo están en DB2

### **Casos de Uso del Análisis de Referencias**

1. **Auditoría de Datos**: Verificar qué IDs se están usando y dónde
2. **Migración Segura**: Antes de eliminar registros, ver qué los referencia
3. **Limpieza de Datos**: Encontrar referencias huérfanas o no utilizadas
4. **Análisis de Impacto**: Entender el alcance de cambios en datos maestros
5. **Sincronización**: Verificar consistencia de referencias entre ambientes

### **Ejemplo Práctico**

Si necesitas eliminar un concepto con `id = 25`, primero ejecutas:

```bash
./deepComparator -table=concepts -target-column=id -find-references
```

El resultado te mostrará:
- **`related_concepts`** tiene 3 referencias al concepto 25
- **`billing_model`** tiene 1 referencia al concepto 25  
- **`settlement_concepts_formula`** tiene 2 referencias al concepto 25

Esto te permite:
1. **Planificar** la limpieza de referencias antes de eliminar
2. **Verificar** que las referencias son consistentes entre DBs
3. **Documentar** el impacto del cambio

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

### **✅ Mejores Prácticas**

#### **📊 Para Migraciones**
1. **Ejecutar con `-include-pk=true`** para verificar IDs
2. **Usar `-exclude-from-file=false`** para comparación completa
3. **Analizar referencias** antes de migrar datos maestros
4. **Documentar diferencias** encontradas para seguimiento

#### **🔍 Para Auditorías**
1. **Configurar exclusiones** específicas por tipo de tabla
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
- [ ] **Análisis de Índices**: Comparar índices, constraints y triggers
- [ ] **Comparación Incremental**: Solo analizar cambios desde última ejecución
- [ ] **Paralelización**: Procesamiento concurrente para tablas grandes
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

**Versión Actual**: `v1.2.0`
- ✅ Comparación profunda de datos con foreign keys
- ✅ Análisis de referencias cruzadas  
- ✅ Exclusión configurable de columnas
- ✅ Salida JSON estructurada

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

**⭐ Si este proyecto te es útil, considera darle una estrella en GitHub para ayudar a otros desarrolladores a encontrarlo.**