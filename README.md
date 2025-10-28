# Deep Database Comparator

Una aplicación en Go para comparar profundamente datos entre dos bases de datos PostgreSQL, incluyendo análisis de foreign keys y detección de diferencias a nivel de registro.

## Características

- **Comparación Profunda**: Compara no solo los datos principales sino también las relaciones de foreign keys
- **Matching Inteligente**: Empareja registros basándose en el contenido, no en IDs que pueden diferir
- **Análisis de Foreign Keys**: Detecta diferencias en tablas referenciadas automáticamente
- **Configuración Flexible**: Permite incluir/excluir columnas específicas en la comparación
- **Salida Estructurada**: Genera reportes detallados en formato JSON
- **Modular**: Código organizado en packages para fácil mantenimiento y extensibilidad

## Estructura del Proyecto

```
deepComparator/
├── cmd/
│   └── main.go                 # Aplicación principal
├── pkg/
│   ├── config/                 # Manejo de configuración
│   │   └── config.go
│   ├── database/               # Conexión y operaciones de BD
│   │   └── database.go
│   ├── comparator/             # Lógica de comparación
│   │   └── comparator.go
│   └── models/                 # Estructuras de datos
│       └── models.go
├── .env.example               # Ejemplo de configuración
├── go.mod                     # Dependencias de Go
├── go.sum                     # Checksums de dependencias
└── README.md                  # Esta documentación
```

## Instalación

### Prerrequisitos

- Go 1.21 o superior
- Acceso a dos bases de datos PostgreSQL
- Permisos de lectura en las tablas a comparar

### Clonar y Compilar

```bash
# Clonar el repositorio (si aplica) o navegar al directorio del proyecto
cd deepComparator

# Instalar dependencias
go mod tidy

# Compilar la aplicación
go build -o deepComparator ./cmd
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

## Uso

### Sintaxis Básica

```bash
./deepComparator -table=<nombre_tabla> [opciones]
```

### Opciones Disponibles

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

### Ejemplos de Uso

#### Comparación Básica (Excluye columnas de auditoría automáticamente)
```bash
./deepComparator -table=billing_model -verbose
```

#### Ver Columnas que se Excluyen desde Archivo
```bash
./deepComparator -show-exclude-columns
```

#### Incluir Todas las Columnas (no excluir desde archivo)
```bash
./deepComparator -table=billing_model -exclude-from-file=false -verbose
```

#### Usar Archivo de Columnas Personalizado
```bash
./deepComparator -table=billing_model -exclude-file="my_exclude_columns.txt" -verbose
```

#### Usar Archivo Vacío (no excluir ninguna columna)
```bash
./deepComparator -table=billing_model -exclude-file="/dev/null" -verbose
```

#### Excluir Columnas Específicas Adicionales
```bash
./deepComparator -table=billing_model -exclude="notes,comments,description" -verbose
```

#### Comparar Solo Columnas Específicas (Ignora exclusiones de auditoría)
```bash
./deepComparator -table=billing_model -include="description,order,status,concept_id" -verbose
```

#### Incluir Claves Primarias en la Comparación
```bash
./deepComparator -table=billing_model -include-pk=true -exclude-audit=false -verbose
```

#### Especificar Esquema y Archivo de Salida
```bash
./deepComparator -table=billing_model -schema=public -output=results.json -verbose
```

#### Usar Archivo de Configuración Personalizado
```bash
./deepComparator -table=billing_model -env=production.env -verbose
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

## Limitaciones Conocidas

- Solo soporta PostgreSQL
- La comparación de foreign keys es limitada a un nivel de profundidad
- Registros muy grandes pueden afectar el rendimiento
- No incluye soporte para tipos de datos binarios complejos

## Contribuir

1. Hacer fork del proyecto
2. Crear una rama para la funcionalidad: `git checkout -b feature/nueva-funcionalidad`
3. Commit los cambios: `git commit -am 'Agregar nueva funcionalidad'`
4. Push a la rama: `git push origin feature/nueva-funcionalidad`
5. Crear un Pull Request

## Licencia

Este proyecto está bajo la licencia MIT. Ver el archivo `LICENSE` para más detalles.

## Soporte

Para reportar bugs o solicitar funcionalidades, por favor crear un issue en el repositorio del proyecto.