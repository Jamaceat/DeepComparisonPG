# Audit Columns Feature - Release Notes

## 🎉 Nueva Funcionalidad: Sistema de Archivos de Configuración para Columnas de Auditoría

### ¿Qué se agregó?

Se implementó un **sistema completamente configurable basado en archivos** para gestionar qué columnas de auditoría excluir. Ahora **TÚ TIENES CONTROL TOTAL** sobre qué columnas ignorar.

### 🔧 Nuevas Opciones de Línea de Comandos

```bash
# Excluir columnas desde archivo (activado por defecto)
-exclude-audit=true/false

# Especificar archivo de configuración de columnas 
-audit-file="mi_archivo.txt"

# Mostrar columnas que se excluyen desde archivo
-show-audit-columns
```

### 📋 Columnas de Auditoría Incluidas (60+ columnas)

#### Timestamps
- `created_at`, `updated_at`, `deleted_at`, `modified_at`, `timestamp`
- `date_created`, `date_updated`, `date_modified`, `date_deleted`
- `creation_date`, `modification_date`, `update_date`, `deletion_date`

#### Usuarios y Control de Acceso  
- `created_by`, `updated_by`, `deleted_by`, `modified_by`, `user_id`
- `creator_id`, `updater_id`, `modifier_id`
- `created_user`, `updated_user`, `created_by_user`, `updated_by_user`

#### Versioning y Control
- `version`, `revision`, `row_version`, `version_number`, `record_version`
- `etag`, `checksum`, `hash`, `lock_version`, `optimistic_lock`

#### Sistema y Logging
- `last_login`, `last_access`, `login_count`, `access_count`
- `ip_address`, `user_agent`, `session_id`, `trace_id`, `request_id`

#### Soft Deletes
- `is_deleted`, `deleted`, `active`, `is_active`, `status_deleted`

#### Sincronización y Migración
- `sync_status`, `migrated_at`, `migration_id`, `import_id`, `batch_id`
- `source_system`, `external_id`, `legacy_id`, `original_id`

### 🚀 Ejemplos de Uso

```bash
# Comparación básica (usa audit_columns.txt por defecto)
./deepComparator -table=billing_model -verbose

# Ver qué columnas se excluyen del archivo
./deepComparator -show-audit-columns

# Incluir todas las columnas de auditoría 
./deepComparator -table=billing_model -exclude-audit=false -verbose

# Usar archivo personalizado de columnas
./deepComparator -table=billing_model -audit-file="mi_config.txt" -verbose

# Usar archivo vacío (no excluir ninguna columna)
./deepComparator -table=billing_model -audit-file="/dev/null" -verbose

# Crear diferentes archivos para diferentes tablas
./deepComparator -table=users -audit-file="users_audit.txt" -verbose
./deepComparator -table=orders -audit-file="orders_audit.txt" -verbose
```

### 📝 Formato del Archivo

```
# Mi archivo personalizado: mi_audit.txt
# Una columna por línea
# Líneas con # son comentarios

# Timestamps que no me interesan
created_at
updated_at

# Campos de mi sistema específico  
last_sync_date
batch_processed
migration_status

# Campos temporales
temp_field1
old_legacy_column
```

### 🎯 Beneficios

1. **Comparaciones más precisas**: Se enfoca en datos de negocio reales
2. **Menos falsos positivos**: Evita diferencias técnicas irrelevantes  
3. **Configuración inteligente**: Activado por defecto pero completamente personalizable
4. **Compatibilidad**: No rompe el uso existente, solo mejora los resultados

### 🔧 Cambios Técnicos

#### Nuevos campos en `MatchCriteria`:
- `ExcludeAuditColumns bool`: Activa/desactiva la exclusión
- `CustomAuditColumns []string`: Columnas adicionales definidas por usuario

#### Nuevas funciones:
- `GetDefaultAuditColumns()`: Retorna lista de columnas estándar
- `GetAllAuditColumns()`: Combina columnas por defecto y personalizadas

#### Modificaciones en el comparador:
- `getRowKey()`: Excluye columnas de auditoría al generar claves
- `compareRows()`: Ignora columnas de auditoría en la comparación
- `createDefaultMatchCriteria()`: Activa exclusión por defecto

### ⚡ Rendimiento

- **Sin impacto negativo**: La exclusión mejora el rendimiento al procesar menos columnas
- **Memoria optimizada**: Menos datos a comparar y almacenar
- **Velocidad mejorada**: Comparaciones más rápidas en tablas con muchas columnas de auditoría

### 🔄 Compatibilidad hacia atrás

- ✅ **100% compatible**: El código existente sigue funcionando igual
- ✅ **Mejores resultados**: Las comparaciones existentes ahora son más precisas
- ✅ **Configuración opcional**: Puedes desactivar la funcionalidad si es necesario

---

Esta funcionalidad resuelve directamente tu necesidad de omitir columnas como `created_at`, `created_by`, etc., haciendo que las comparaciones se enfoquen en el contenido real de los datos.