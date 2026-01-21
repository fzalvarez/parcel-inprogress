# Índice de Documentación - QuatroBus Parcel

**Última Actualización:** 21 de enero de 2026  
**Versión del Proyecto:** 2.0

---

## 📚 Documentación Disponible

### 🚀 Inicio Rápido

| Documento | Propósito | Para Quién |
|-----------|-----------|-----------|
| [README.md](../../README.md) | Descripción general del proyecto | Todos |
| [SETUP.md](../../SETUP.md) | Instalación y configuración | DevOps, Desarrolladores |
| [scripts/init_parcel.ps1](../../scripts/init_parcel.ps1) | Script de inicialización | DevOps |

---

### 📖 Arquitectura y Diseño

| Documento | Propósito | Para Quién |
|-----------|-----------|-----------|
| [.github/instructions/go_profile.instructions.md](../../.github/instructions/go_profile.instructions.md) | Perfil de Go: framework, pattern, librerías | Desarrolladores Go |
| [.github/instructions/parcel_boundaries.instructions.md](../../.github/instructions/parcel_boundaries.instructions.md) | Límites de dominio y módulos de Parcel | Arquitectos, Desarrolladores |
| [docs/architecture_diagram.md](./architecture_diagram.md) | 🆕 Diagrama visual completo del sistema | Arquitectos, Desarrolladores, PMs |
| [docs/persistence_architecture.md](./persistence_architecture.md) | 🆕 Arquitectura de persistencia (PostgreSQL + memoria) | Arquitectos, Desarrolladores |

---

### 🎯 Pricing (Motor de Precios)

| Documento | Propósito | Para Quién |
|-----------|-----------|-----------|
| [docs/pricing_rules_guide.md](./pricing_rules_guide.md) | Guía completa del motor de precios jerárquico | Desarrolladores, Product Managers |
| **Temas Cubiertos:** | | |
| - Arquitectura del motor | Sistema de prioridad y comodines | Desarrolladores |
| - Cálculo de peso facturable | Peso real vs volumétrico | Desarrolladores |
| - Reglas específicas y comodines | Ejemplos prácticos | Desarrolladores, PMs |
| - Configuración por tenant | Feature flags | Desarrolladores |
| - Integración con API | DTOs y endpoints | Desarrolladores |

---

### 🔗 API y Swagger

| Documento | Propósito | Para Quién |
|-----------|-----------|-----------|
| [docs/swagger_endpoints_reference.md](./swagger_endpoints_reference.md) | Referencia completa de todos los endpoints | API Consumers, Desarrolladores |
| [docs/swagger_maintenance_guide.md](./swagger_maintenance_guide.md) | Cómo mantener Swagger actualizado | Desarrolladores |
| [docs/SWAGGER_UPDATE_SUMMARY.md](./SWAGGER_UPDATE_SUMMARY.md) | Resumen de actualización Swagger (Enero 2026) | Desarrolladores, Revisores |

**Endpoints Documentados:**
- Parcels (Envíos): 8 endpoints
- ParcelItems (Artículos): 3 endpoints
- ParcelPayments (Pagos): 3 endpoints
- ParcelTracking (Historial): 1 endpoint
- Manifests (Manifiestos): 2 endpoints
- Pricing (Precios): 3 endpoints
- ParcelDocuments (Documentos): 2 endpoints

**Total: 28 endpoints documentados**

---

## 📁 Estructura del Proyecto

```
parcel-inprogress/
├── cmd/
│   └── api/
│       ├── main.go                    # Entrada de la aplicación
│       └── swagger_meta.go            # Metadata Swagger
├── internal/
│   ├── infrastructure/
│   │   ├── http/
│   │   │   ├── dto/                   # Data Transfer Objects
│   │   │   ├── handler/               # Handlers HTTP (🎯 CON SWAGGER)
│   │   │   ├── middleware/            # Middlewares (auth, error handling)
│   │   │   └── router/                # Rutas
│   │   └── clients/                   # Clientes de servicios externos
│   └── parcel/
│       ├── parcel_core/               # Módulo core: estados y transiciones
│       ├── parcel_item/               # Módulo items: artículos/bultos
│       ├── parcel_payment/            # Módulo payment: información de pago
│       ├── parcel_tracking/           # Módulo tracking: historial de eventos
│       ├── parcel_manifest/           # Módulo manifest: manifiesto virtual
│       ├── parcel_pricing/            # Módulo pricing: motor de precios (🎯)
│       └── parcel_documents/          # Módulo documents: impresión y docs
├── .github/
│   ├── docs/                          # 📚 DOCUMENTACIÓN PRINCIPAL
│   │   ├── pricing_rules_guide.md
│   │   ├── swagger_endpoints_reference.md
│   │   ├── swagger_maintenance_guide.md
│   │   ├── SWAGGER_UPDATE_SUMMARY.md
│   │   └── INDEX.md (este archivo)
│   └── instructions/
│       ├── go_profile.instructions.md
│       └── parcel_boundaries.instructions.md
├── docs/                              # Swagger JSON generado (gitignore)
│   ├── swagger.json
│   └── swagger.yaml
├── scripts/
│   └── init_parcel.ps1
├── go.mod
├── go.sum
└── README.md
```

---

## 🔑 Conceptos Clave

### 1. **Motor de Precios** (parcel_pricing)
- Sistema jerárquico de búsqueda de reglas
- Soporte de comodines (`*`) para mayor flexibilidad
- Prioridad (0-100) para control de orden de evaluación
- Cálculo automático de peso facturable (real vs volumétrico)
- Ver: [pricing_rules_guide.md](./pricing_rules_guide.md)

### 2. **Estados del Envío** (parcel_core)
```
CREATED → REGISTERED → BOARDED → EN_ROUTE → ARRIVED → DELIVERED
```

### 3. **Módulos de Dominio**
- **parcel_core**: Gestión de envíos y ciclo de vida
- **parcel_item**: Artículos/bultos dentro de envíos
- **parcel_payment**: Información y transacciones de pago
- **parcel_tracking**: Historial temporal de eventos
- **parcel_pricing**: Motor de cálculo de tarifas (🎯)
- **parcel_manifest**: Manifiesto virtual de cargas
- **parcel_documents**: Impresión y documentación

### 4. **Arquitectura Clean/Hexagonal**
```
Handler → UseCase → Port (Interface)
                ↓
            Repository/Client (Implementation)
```

Ver: [.github/instructions/parcel_boundaries.instructions.md](../../.github/instructions/parcel_boundaries.instructions.md)

---

## 🎯 Casos de Uso Comunes

### Crear un Envío
```
1. POST /parcels → Estado: CREATED
2. POST /parcels/{id}/items → Agregar artículos
3. POST /parcels/{id}/register → Estado: REGISTERED
```

### Embarcar y Entregar
```
1. POST /parcels/{id}/board → Estado: BOARDED (asignar vehículo)
2. POST /parcels/{id}/depart → Estado: EN_ROUTE
3. POST /parcels/{id}/arrive → Estado: ARRIVED
4. POST /parcels/{id}/deliver → Estado: DELIVERED
```

### Gestionar Pagos
```
1. PUT /parcels/{id}/payment → Crear/actualizar pago
2. GET /parcels/{id}/payment → Obtener detalles
3. POST /parcels/{id}/payment/mark-paid → Marcar como PAID
```

### Configurar Precios
```
1. POST /pricing/rules → Crear regla (soporta comodines)
2. PUT /pricing/rules/{id} → Actualizar regla
3. GET /pricing/rules → Listar todas las reglas
```

---

## 🛠️ Herramientas y Tecnologías

| Componente | Tecnología | Versión |
|-----------|-----------|---------|
| **Framework HTTP** | Gin | v1.9+ |
| **Base de Datos** | PostgreSQL + GORM | - |
| **Logging** | Zap + Lumberjack | - |
| **Documentación** | Swagger/OpenAPI | 2.0 |
| **Generador Swagger** | swag | v1.16+ |

### Dependencias de Go
- `github.com/gin-gonic/gin` - Framework HTTP
- `github.com/gin-contrib/cors` - CORS middleware
- `gorm.io/gorm` - ORM
- `gorm.io/driver/postgres` - Driver PostgreSQL
- `go.uber.org/zap` - Logging estructurado
- `gopkg.in/natefinch/lumberjack.v2` - Log rotation
- `github.com/google/uuid` - UUID generation
- `github.com/swaggo/swag` - Swagger generation
- `github.com/swaggo/gin-swagger` - Swagger UI en Gin

---

## 📋 Checklist de Développador

### Al Empezar a Trabajar
- [ ] Leer [go_profile.instructions.md](../../.github/instructions/go_profile.instructions.md)
- [ ] Leer [parcel_boundaries.instructions.md](../../.github/instructions/parcel_boundaries.instructions.md)
- [ ] Entender la arquitectura Clean/Hexagonal
- [ ] Familiarizarse con el motor de precios

### Al Crear un Endpoint
- [ ] Crear handler en `infrastructure/http/handler/`
- [ ] Crear usecase correspondiente
- [ ] Implementar port (interface) si es necesario
- [ ] Crear DTOs en `infrastructure/http/dto/`
- [ ] Agregar Swagger comments completos
- [ ] Validar con `go build`
- [ ] Generar Swagger: `swag init -g cmd/api/main.go -o docs`
- [ ] Probar en Swagger UI

### Al Hacer Commit
- [ ] Correr linter/formatter
- [ ] Validar que compila: `go build ./cmd/api`
- [ ] Generar Swagger: `swag init -g cmd/api/main.go -o docs`
- [ ] Validar en Swagger UI: `go run cmd/api/main.go`
- [ ] Agregar documentación si aplica
- [ ] Commit message en inglés, clara y concisa

---

## 📞 Contacto y Soporte

| Rol | Contacto |
|-----|----------|
| **Tech Lead** | - |
| **Product Manager** | - |
| **DevOps** | - |

---

## 📅 Historial de Cambios

### Enero 20, 2026
- ✅ Actualización completa de comentarios Swagger en todos los handlers
- ✅ Creación de guía de referencia de endpoints (28 endpoints)
- ✅ Creación de guía de mantenimiento de Swagger
- ✅ Creación de resumen de actualización

### Enero 13, 2026
- ✅ Implementación del motor de precios jerárquico con comodines
- ✅ Cálculo de peso volumétrico y facturable
- ✅ Creación de guía de reglas de precios

---

## 📚 Lectura Recomendada (Orden)

### Para Nuevos Desarrolladores
1. [README.md](../../README.md)
2. [docs/SWAGGER_UPDATE_SUMMARY.md](./SWAGGER_UPDATE_SUMMARY.md) - Visión general del proyecto
3. [.github/instructions/go_profile.instructions.md](../../.github/instructions/go_profile.instructions.md) - Perfil técnico
4. [.github/instructions/parcel_boundaries.instructions.md](../../.github/instructions/parcel_boundaries.instructions.md) - Arquitectura
5. [docs/pricing_rules_guide.md](./pricing_rules_guide.md) - Lógica de negocios
6. [docs/swagger_endpoints_reference.md](./swagger_endpoints_reference.md) - API disponible

### Para Desarrolladores Backend
1. [.github/instructions/go_profile.instructions.md](../../.github/instructions/go_profile.instructions.md)
2. [.github/instructions/parcel_boundaries.instructions.md](../../.github/instructions/parcel_boundaries.instructions.md)
3. [docs/pricing_rules_guide.md](./pricing_rules_guide.md)
4. [docs/swagger_maintenance_guide.md](./swagger_maintenance_guide.md)

### Para Integradores/API Consumers
1. [docs/swagger_endpoints_reference.md](./swagger_endpoints_reference.md)
2. Swagger UI en `/swagger/index.html`
3. [docs/pricing_rules_guide.md](./pricing_rules_guide.md) - Especialmente sección de Pricing

### Para Product Managers
1. [docs/pricing_rules_guide.md](./pricing_rules_guide.md)
2. [docs/swagger_endpoints_reference.md](./swagger_endpoints_reference.md) - Listar "Casos de Uso"

---

## 🔗 Enlaces Útiles

- **Local Swagger UI:** `http://localhost:8080/swagger/index.html`
- **Swagger Petstore Demo:** https://petstore.swagger.io/
- **OpenAPI 3.0 Spec:** https://spec.openapis.org/oas/v3.0.0
- **Swag Repository:** https://github.com/swaggo/swag
- **Gin Repository:** https://github.com/gin-gonic/gin

---

## ✅ Estado Actual del Proyecto

| Componente | Estado | Progreso |
|-----------|--------|----------|
| **Arquitectura Core** | ✅ Completado | 100% |
| **Motor de Precios** | ✅ Completado | 100% |
| **Endpoints Principales** | ✅ Completado | 100% |
| **Documentación Swagger** | ✅ Completado | 100% |
| **Guías de Desarrollo** | ✅ Completado | 100% |
| **Testing Unitario** | 🟡 Pendiente | 0% |
| **Testing Integración** | 🟡 Pendiente | 0% |
| **Deployment** | 🟡 Pendiente | 0% |

---

**Próxima Actualización Esperada:** 27 de enero de 2026

---

**Responsable de Documentación:** Equipo de Desarrollo  
**Última Revisión:** 20 de enero de 2026  
**Versión de Documentación:** 1.0
