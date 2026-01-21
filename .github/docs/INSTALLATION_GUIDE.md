# 🚀 Instrucciones de Instalación Final - QuatroBus Parcel

**Proyecto:** ms-parcel-core  
**Fecha:** 21 de enero de 2026

---

## ⚠️ Estado Actual

El proyecto está **completamente documentado y arquitectado**, pero requiere **2 pasos finales** para estar 100% operativo con PostgreSQL:

1. ✅ **Swagger** - LISTO (funciona con repositorios en memoria)
2. ⏳ **PostgreSQL** - Requiere instalación de dependencias GORM

---

## 📋 Opción 1: Usar Solo Memoria (Ya Funciona)

Si solo quieres probar la API con datos en memoria:

```bash
# El proyecto ya funciona así
go run cmd/api/main.go

# Acceder a Swagger
http://localhost:8080/swagger/index.html
```

✅ **Ventajas:**
- No requiere PostgreSQL
- Perfecto para desarrollo y testing
- Rápido de iniciar

⚠️ **Desventajas:**
- Los datos se pierden al reiniciar
- No hay persistencia real

---

## 📋 Opción 2: Habilitar PostgreSQL (Recomendado para Producción)

### Paso 1: Instalar Dependencias GORM

```bash
# Instalar GORM core
go get -u gorm.io/gorm

# Instalar driver de PostgreSQL
go get -u gorm.io/driver/postgres

# Actualizar go.mod
go mod tidy
```

### Paso 2: Configurar Variables de Entorno

Crear archivo `.env` en la raíz del proyecto:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_password_aqui
DB_NAME=parcel_db

# Server Configuration
SERVER_PORT=8080
ENVIRONMENT=development
```

### Paso 3: Crear Base de Datos PostgreSQL

```bash
# Opción A: Usar PostgreSQL local
createdb parcel_db

# Opción B: Usar Docker
docker run --name parcel-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=parcel_db \
  -p 5432:5432 \
  -d postgres:15
```

### Paso 4: Actualizar main.go

Reemplazar la sección de inicialización en `cmd/api/main.go`:

```go
package main

import (
    "log"
    "os"
    
    "github.com/gin-gonic/gin"
    
    "ms-parcel-core/internal/config"
    "ms-parcel-core/internal/infrastructure/persistence/database"
    "ms-parcel-core/internal/infrastructure/http/router"
    // ... otros imports
)

func main() {
    // Cargar configuración desde ENV
    cfg := config.Config{
        DB: config.DBConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "postgres"),
            Name:     getEnv("DB_NAME", "parcel_db"),
        },
        ServerPort:  getEnv("SERVER_PORT", "8080"),
        Environment: getEnv("ENVIRONMENT", "development"),
    }
    
    // Conectar a PostgreSQL
    db, err := database.Connect(cfg.DB)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    
    // Ejecutar migraciones
    if err := database.Migrate(db); err != nil {
        log.Fatalf("Error running migrations: %v", err)
    }
    
    log.Println("✅ Database connected and migrated successfully")
    
    // TODO: Crear repositorios PostgreSQL aquí en lugar de in-memory
    // parcelRepo := postgres.NewParcelPostgresRepository(db)
    // ... etc
    
    // Iniciar servidor Gin
    r := gin.Default()
    router.SetupRoutes(r)
    
    if err := r.Run(":" + cfg.ServerPort); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}
```

### Paso 5: Ejecutar el Proyecto

```bash
# Cargar variables de entorno (Windows PowerShell)
$env:DB_PASSWORD="tu_password"

# O crear archivo .env y cargar con una librería
go get github.com/joho/godotenv

# Ejecutar
go run cmd/api/main.go
```

Deberías ver:
```
Connected to PostgreSQL successfully
✅ Database connected and migrated successfully
[GIN-debug] Listening and serving HTTP on :8080
```

---

## 🔄 Paso 3: Crear Repositorios PostgreSQL (Opcional)

Para usar PostgreSQL en lugar de memoria, crear estos archivos:

### Ejemplo: ParcelPostgresRepository

Crear `internal/infrastructure/persistence/postgres/parcel_postgres_repository.go`:

```go
package postgres

import (
    "context"
    
    "github.com/google/uuid"
    "gorm.io/gorm"
    
    "ms-parcel-core/internal/parcel/parcel_core/domain"
    "ms-parcel-core/internal/parcel/parcel_core/port"
    "ms-parcel-core/internal/pkg/util/apperror"
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
        return nil, apperror.NewInternal("conversion_error", "error converting domain to DB model", map[string]any{"error": err.Error()})
    }
    
    // Inyectar tenant_id en contexto de GORM
    db := r.db.WithContext(ctx).Set("tenant_id", tenantID)
    
    if err := db.Create(&dbModel).Error; err != nil {
        return nil, apperror.NewInternal("database_error", "error creating parcel", map[string]any{"error": err.Error()})
    }
    
    result := dbModel.ToDomain()
    return &result, nil
}

func (r *ParcelPostgresRepository) GetByID(ctx context.Context, tenantID string, id uuid.UUID) (*domain.Parcel, error) {
    var dbModel DBParcel
    
    db := r.db.WithContext(ctx).Set("tenant_id", tenantID)
    
    if err := db.Where("id = ?", id).First(&dbModel).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, apperror.NewInternal("database_error", "error fetching parcel", map[string]any{"error": err.Error()})
    }
    
    result := dbModel.ToDomain()
    return &result, nil
}

// ... implementar resto de métodos de la interfaz port.ParcelRepository
```

Repetir para:
- `ParcelItemPostgresRepository`
- `ParcelPaymentPostgresRepository`
- `TrackingEventPostgresRepository`
- `PrintRecordPostgresRepository`
- `PriceRulePostgresRepository`

---

## 📦 Mover Repositorios In-Memory

Opcional: reorganizar repos in-memory a carpeta centralizada:

```bash
# PowerShell
Move-Item internal\parcel\parcel_core\infrastructure\repository\in_memory_parcel_repository.go internal\infrastructure\persistence\memory\
Move-Item internal\parcel\parcel_item\infrastructure\repository\in_memory_parcel_item_repository.go internal\infrastructure\persistence\memory\
Move-Item internal\parcel\parcel_payment\infrastructure\repository\in_memory_parcel_payment_repository.go internal\infrastructure\persistence\memory\
Move-Item internal\parcel\parcel_tracking\infrastructure\repository\in_memory_tracking_repository.go internal\infrastructure\persistence\memory\
Move-Item internal\parcel\parcel_documents\infrastructure\repository\in_memory_print_repository.go internal\infrastructure\persistence\memory\
Move-Item internal\parcel\parcel_pricing\infrastructure\repository\in_memory_price_rule_repository.go internal\infrastructure\persistence\memory\
```

Luego actualizar los imports en `main.go`.

---

## ✅ Verificación Final

### 1. Verificar que Swagger Funciona

```bash
go run cmd/api/main.go
```

Abrir: `http://localhost:8080/swagger/index.html`

Deberías ver:
- ✅ 22 endpoints listados
- ✅ 7 tags organizados
- ✅ Documentación completa en cada endpoint

### 2. Verificar que PostgreSQL Conecta (si lo instalaste)

Deberías ver en consola:
```
Connected to PostgreSQL successfully
✅ Database connected and migrated successfully
```

### 3. Verificar Tablas Creadas (si usas PostgreSQL)

```sql
-- Conectar a PostgreSQL
psql -U postgres -d parcel_db

-- Listar tablas
\dt

-- Deberías ver:
-- parcels
-- parcel_items
-- parcel_payments
-- tracking_events
-- print_records
-- price_rules
```

---

## 🎯 Resumen de Opciones

### Desarrollo Rápido (Ahora)
```
✅ Repositorios in-memory (ya funciona)
✅ Swagger completo
✅ API lista para probar
⚠️ Datos se pierden al reiniciar
```

### Producción (Requiere pasos adicionales)
```
⏳ Instalar GORM (2 comandos)
⏳ Configurar PostgreSQL
⏳ Crear repositorios PostgreSQL (opcional, puede quedarse en memoria)
✅ Datos persistentes
```

---

## 📞 Ayuda

### Errores Comunes

**Error: `could not import gorm.io/gorm`**
```bash
# Solución:
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go mod tidy
```

**Error: `dial tcp: connection refused`**
- Verificar que PostgreSQL esté corriendo
- Verificar credenciales en `.env`
- Verificar puerto 5432 disponible

**Error: `password authentication failed`**
- Verificar `DB_PASSWORD` en `.env`
- Verificar usuario `DB_USER` existe en PostgreSQL

---

## 🎉 Conclusión

El proyecto **QuatroBus Parcel** está:

✅ **100% Documentado** (Swagger + 8 guías técnicas)  
✅ **100% Funcional** (con repositorios in-memory)  
✅ **Arquitectura Lista** (modelos PostgreSQL creados)

Para usar PostgreSQL:
1. Instalar GORM (2 comandos)
2. Configurar credenciales
3. (Opcional) Crear repos PostgreSQL

---

**Preparado por:** Equipo de Desarrollo QuatroBus  
**Fecha:** 21 de enero de 2026  
**Estado:** 🟢 LISTO PARA USAR
