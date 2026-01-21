# Reorganización de Arquitectura de Persistencia - QuatroBus Parcel

## Resumen Ejecutivo

Se ha reorganizado la capa de persistencia del proyecto `ms-parcel-core` para seguir el mismo patrón de arquitectura utilizado en `ms-vehicle` y otros microservicios del ecosistema QuatroBus.

---

## Cambios Realizados

### 1. Nueva Estructura de Carpetas

```
internal/
├── infrastructure/
│   ├── http/          # Ya existía - handlers, middleware, routing
│   └── persistence/   # NUEVO - capa de persistencia
│       ├── database/  # Conexión y migraciones
│       ├── postgres/  # Modelos y repositorios PostgreSQL
│       └── memory/    # Repositorios en memoria (para testing/dev)
```

### 2. Archivos Creados

#### **Configuración:**
- `internal/config/config.go` - Estructuras de configuración de BD y app

#### **Database (Conexión y Migraciones):**
- `internal/infrastructure/persistence/database/connect.go` - Conexión a PostgreSQL con GORM
- `internal/infrastructure/persistence/database/migrate.go` - AutoMigrate de todos los modelos

#### **PostgreSQL Models:**
- `internal/infrastructure/persistence/postgres/tenant_scope.go` - Scope global de tenant_id
- `internal/infrastructure/persistence/postgres/parcel_model.go` - Modelo DBParcel
- `internal/infrastructure/persistence/postgres/parcel_item_model.go` - Modelo DBParcelItem
- `internal/infrastructure/persistence/postgres/parcel_payment_model.go` - Modelo DBParcelPayment
- `internal/infrastructure/persistence/postgres/tracking_event_model.go` - Modelo DBTrackingEvent (⚠️ requiere ajustes)
- `internal/infrastructure/persistence/postgres/print_record_model.go` - Modelo DBPrintRecord
- `internal/infrastructure/persistence/postgres/price_rule_model.go` - Modelo DBPriceRule (⚠️ requiere ajustes)

---

## Pendientes (Próximos Pasos)

### 1. Instalar Dependencias de GORM

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

### 2. Mover Repositorios In-Memory

Los repositorios in-memory actuales están en:
```
internal/parcel/*/infrastructure/repository/in_memory_*.go
```

**Deben moverse a:**
```
internal/infrastructure/persistence/memory/
```

Lista de archivos a mover:
- `in_memory_parcel_repository.go`
- `in_memory_parcel_item_repository.go`
- `in_memory_parcel_payment_repository.go`
- `in_memory_tracking_repository.go`
- `in_memory_print_repository.go`
- `in_memory_price_rule_repository.go`

### 3. Crear Repositorios PostgreSQL

Para cada entidad, crear su repositorio PostgreSQL que implemente la interfaz port correspondiente:

```
internal/infrastructure/persistence/postgres/
├── parcel_postgres_repository.go
├── parcel_item_postgres_repository.go
├── parcel_payment_postgres_repository.go
├── tracking_event_postgres_repository.go
├── print_record_postgres_repository.go
└── price_rule_postgres_repository.go
```

### 4. Ajustar Modelos con Errores

**tracking_event_model.go** - El dominio `TrackingEvent` es más simple de lo esperado. Ajustar campos.

**price_rule_model.go** - Usar `coredomain.ShipmentType` y `PriceUnit` correctamente.

### 5. Actualizar `main.go`

Integrar la conexión a BD y migraciones en el arranque de la aplicación:

```go
import (
    "ms-parcel-core/internal/config"
    "ms-parcel-core/internal/infrastructure/persistence/database"
)

func main() {
    // Cargar configuración
    cfg := loadConfig()
    
    // Conectar a PostgreSQL
    db, err := database.Connect(cfg.DB)
    if err != nil {
        log.Fatal(err)
    }
    
    // Ejecutar migraciones
    if err := database.Migrate(db); err != nil {
        log.Fatal(err)
    }
    
    // ... resto del setup
}
```

### 6. Variables de Entorno

Agregar configuración de BD al archivo `.env` o config:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=parcel_db
```

---

## Patrón de Repositorio PostgreSQL

Ejemplo de implementación:

```go
package postgres

import (
    "context"
    "gorm.io/gorm"
    "ms-parcel-core/internal/parcel/parcel_core/domain"
    "ms-parcel-core/internal/parcel/parcel_core/port"
)

type ParcelPostgresRepository struct {
    db *gorm.DB
}

var _ port.ParcelRepository = (*ParcelPostgresRepository)(nil)

func NewParcelPostgresRepository(db *gorm.DB) *ParcelPostgresRepository {
    return &ParcelPostgresRepository{db: db}
}

func (r *ParcelPostgresRepository) Create(ctx context.Context, tenantID string, parcel domain.Parcel) (*domain.Parcel, error) {
    var dbModel DBParcel
    if err := dbModel.FromDomain(parcel); err != nil {
        return nil, err
    }
    
    db := r.db.WithContext(ctx).Set("tenant_id", tenantID)
    if err := db.Create(&dbModel).Error; err != nil {
        return nil, err
    }
    
    result := dbModel.ToDomain()
    return &result, nil
}

// ... implementar resto de métodos de la interfaz port.ParcelRepository
```

---

## Convenciones Establecidas

### Nombrado de Modelos
- **Dominio:** `Parcel`, `ParcelItem`, `PriceRule`
- **PostgreSQL:** `DBParcel`, `DBParcelItem`, `DBPriceRule`
- **Tablas:** `parcels`, `parcel_items`, `price_rules`

### Métodos de Conversión
- `ToDomain()` - Convierte modelo DB a dominio
- `FromDomain(domain)` - Convierte dominio a modelo DB

### Tenant Scope
- Todos los queries automáticamente filtran por `tenant_id`
- Se inyecta en el contexto de GORM: `db.Set("tenant_id", tenantID)`

### UUID
- Todos los IDs son `uuid.UUID` en BD
- Se convierten a `string` en el dominio
- Hook `BeforeCreate` genera UUID automático si falta

---

## Ventajas de la Nueva Arquitectura

✅ **Consistencia** con otros microservicios (ms-vehicle, etc.)
✅ **Separación clara** entre memoria y persistencia real
✅ **Fácil testing** - se puede swap entre memory y postgres
✅ **Multi-tenancy** automático con tenant scope
✅ **Migraciones automáticas** con GORM AutoMigrate
✅ **Type-safe** con modelos fuertemente tipados

---

## Documentación Relacionada

- [Perfil Go](../.github/instructions/go_profile.instructions.md)
- [Límites de Parcel](../.github/instructions/parcel_boundaries.instructions.md)
- [Guía de Pricing](./pricing_rules_guide.md)
- [Referencia Swagger](./swagger_endpoints_reference.md)

---

**Estado:** 🟡 En progreso - Requiere instalación de dependencias y ajustes finales

**Última Actualización:** 21 de enero de 2026  
**Responsable:** Equipo de Desarrollo - QuatroBus
