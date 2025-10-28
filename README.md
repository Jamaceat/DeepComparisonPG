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

## Formato de Salida

El resultado se genera en formato JSON con la siguiente estructura:

```json
{
  "table_name": "billing_model",
  "schema": "public", 
  "timestamp": "2025-10-28T15:30:00Z",
  "total_rows_db1": 100,
  "total_rows_db2": 98,
  "matched_rows": 95,
  "unmatched_rows": 5,
  "only_in_db1": [...],
  "only_in_db2": [...],
  "differences": [
    {
      "row_identifier": "order:3,status:ac,end_date:2025-12-29...",
      "db1_row": {...},
      "db2_row": {...},
      "column_differences": [
        {
          "column_name": "description",
          "db1_value": "Get interest discount about number 693 31 of JUL of 2025",
          "db2_value": null
        }
      ]
    }
  ],
  "foreign_key_results": [
    {
      "foreign_key": {
        "column_name": "concept_id",
        "referenced_table": "concepts",
        "referenced_schema": "public",
        "referenced_column_name": "id"
      },
      "comparison_result": {
        "matched_rows": 1,
        "differences": [...]
      }
    }
  ]
}
```

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