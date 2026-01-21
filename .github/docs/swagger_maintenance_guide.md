# Guía de Actualización y Mantenimiento de Swagger - QuatroBus Parcel

**Versión:** 1.0  
**Fecha:** 20 de enero de 2026  
**Propósito:** Guía para desarrolladores sobre cómo mantener la documentación Swagger actualizada.

---

## Tabla de Contenidos

1. [Introducción](#introducción)
2. [Estructura de Comentarios Swagger](#estructura-de-comentarios-swagger)
3. [Convenciones Obligatorias](#convenciones-obligatorias)
4. [Ejemplos Prácticos](#ejemplos-prácticos)
5. [Cómo Generar el Archivo Swagger](#cómo-generar-el-archivo-swagger)
6. [Validación y Testing](#validación-y-testing)
7. [Checklist de Actualización](#checklist-de-actualización)

---

## Introducción

Todos los endpoints en QuatroBus Parcel deben tener documentación **completa y actualizada** en Swagger/OpenAPI. Esto permite:

- ✅ Generar documentación interactiva automáticamente
- ✅ Facilitar testing manual con Swagger UI
- ✅ Generar clientes/SDKs automáticamente
- ✅ Mantener sincronización código ↔ documentación

**Herramienta utilizada:** [swag](https://github.com/swaggo/swag) + [gin-swagger](https://github.com/swaggo/gin-swagger)

---

## Estructura de Comentarios Swagger

### Anatomía de un Endpoint

```go
// MyHandler godoc
// @Summary Descripción corta (máx 80 caracteres)
// @Description Descripción larga con contexto y comportamiento esperado
// @Tags TagName
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "Descripción del parámetro" Format(uuid)
// @Param query_param query int false "Descripción" minimum(0) maximum(100)
// @Param payload body RequestDTO true "Descripción del request body"
// @Success 200 {object} ResponseDTO "Descripción del éxito"
// @Success 201 {object} ResponseDTO "Para POST/PUT"
// @Failure 400 {object} ErrorResponse "Validación fallida"
// @Failure 401 {object} ErrorResponse "No autorizado"
// @Failure 404 {object} ErrorResponse "Recurso no encontrado"
// @Failure 409 {object} ErrorResponse "Conflicto"
// @Failure 500 {object} ErrorResponse "Error interno"
// @Router /path [post|get|put|delete]
func (h *Handler) MyHandler(c *gin.Context) {
    // implementación
}
```

---

## Convenciones Obligatorias

### 1. Summary (Resumen)

**Obligatorio.** Máximo 80 caracteres.

✅ **Bueno:**
```
// @Summary Crear nuevo envío
// @Summary Obtener detalles del envío
// @Summary Listar artículos del envío
```

❌ **Malo:**
```
// @Summary Crear
// @Summary GET endpoint para obtener información detallada completa del envío con todos los campos
```

---

### 2. Description (Descripción)

**Obligatorio.** Proporciona contexto, comportamiento, transiciones de estado, etc.

✅ **Bueno:**
```
// @Description Crea un nuevo envío en estado CREATED. Requiere tipos de envío, 
// oficinas origen/destino, personas (remitente/destinatario) y opcionales notes. 
// El package_key permite proteger operaciones posteriores con confirmación.
```

❌ **Malo:**
```
// @Description Crea un envío
// @Description Endpoint para crear un envío (consultar documentación)
```

---

### 3. Tags

**Obligatorio.** Agrupa endpoints por dominio.

**Tags Válidos:**
- `Parcels` - Gestión de envíos y estados
- `ParcelItems` - Artículos/bultos
- `ParcelPayments` - Información de pago
- `ParcelTracking` - Historial de eventos
- `Manifests` - Manifiesto virtual
- `Pricing` - Reglas y cálculo de precios
- `ParcelDocuments` - Impresión y documentación

```go
// @Tags Parcels
// @Tags ParcelItems
```

---

### 4. Parámetros

**Format de Path Parameter:**
```go
// @Param id path string true "UUID del envío" Format(uuid)
// @Param item_id path string true "UUID del item" Format(uuid)
```

**Format de Query Parameter:**
```go
// @Param limit query int false "Límite (default: 50)" minimum(1) maximum(200)
// @Param offset query int false "Desplazamiento (default: 0)" minimum(0)
// @Param status query string false "Estado (CREATED, REGISTERED, BOARDED, EN_ROUTE, ARRIVED, DELIVERED)"
```

**Format de Body Parameter:**
```go
// @Param payload body CreateParcelRequest true "Solicitud con datos del envío"
```

**Regla:** Siempre incluir descripción significativa.

---

### 5. Respuestas (Success y Failure)

**Success:**
```go
// @Success 200 {object} ResponseDTO "Descripción clara"
// @Success 201 {object} ResponseDTO "Recurso creado exitosamente"
```

**Failure (ordenadas lógicamente):**
```go
// @Failure 400 {object} ErrorResponse "Validación fallida: payload malformado o valores inválidos"
// @Failure 401 {object} ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 404 {object} ErrorResponse "Recurso no encontrado"
// @Failure 409 {object} ErrorResponse "Conflicto: estado incompatible o operación no permitida"
// @Failure 500 {object} ErrorResponse "Error interno del servidor"
```

---

## Ejemplos Prácticos

### Ejemplo 1: Crear Recurso (POST)

```go
// CreateParcelItem godoc
// @Summary Agregar artículo al envío
// @Description Agrega un bulto/artículo al envío con cálculo automático de peso 
// facturable y precio. Soporta dimensiones opcionales para peso volumétrico.
// @Tags ParcelItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID del envío" Format(uuid)
// @Param payload body CreateParcelItemRequest true "Datos del item (descripción y peso requeridos)"
// @Success 201 {object} ParcelItemResponse "Item creado exitosamente"
// @Failure 400 {object} ErrorResponse "Validación fallida: payload malformado"
// @Failure 401 {object} ErrorResponse "No autorizado"
// @Failure 404 {object} ErrorResponse "Envío no encontrado"
// @Failure 409 {object} ErrorResponse "Regla de precios no encontrada"
// @Failure 500 {object} ErrorResponse "Error interno"
// @Router /parcels/{id}/items [post]
func (h *ParcelItemHandler) Add(c *gin.Context) {
    // implementación
}
```

---

### Ejemplo 2: Listar con Paginación (GET)

```go
// ListParcels godoc
// @Summary Listar envíos
// @Description Lista envíos del tenant con filtros y paginación. Soporta búsqueda 
// por query, filtros por estado, offices, personas y rango de fechas.
// @Tags Parcels
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param q query string false "Búsqueda por código o ID"
// @Param status query string false "Estado (CREATED, REGISTERED, BOARDED, EN_ROUTE, ARRIVED, DELIVERED)"
// @Param origin_office_id query string false "Filtrar por oficina origen (UUID)" Format(uuid)
// @Param destination_office_id query string false "Filtrar por oficina destino (UUID)" Format(uuid)
// @Param limit query int false "Límite (default: 50, max: 200)" minimum(1) maximum(200)
// @Param offset query int false "Desplazamiento (default: 0)" minimum(0)
// @Param from_created_at query string false "Desde (RFC3339)" Format(date-time)
// @Param to_created_at query string false "Hasta (RFC3339)" Format(date-time)
// @Success 200 {object} ParcelListResponse "Lista paginada de envíos"
// @Failure 400 {object} ErrorResponse "Parámetros inválidos"
// @Failure 401 {object} ErrorResponse "No autorizado"
// @Failure 500 {object} ErrorResponse "Error interno"
// @Router /parcels [get]
func (h *ParcelHandler) List(c *gin.Context) {
    // implementación
}
```

---

### Ejemplo 3: Transición de Estado (POST sin Body)

```go
// RegisterParcel godoc
// @Summary Registrar envío
// @Description Transiciona el envío de estado CREATED a REGISTERED. Marca el envío 
// como listo para ser embarcado en un vehículo.
// @Tags Parcels
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID del envío" Format(uuid)
// @Success 200 {object} ParcelResponse "Envío registrado (estado: REGISTERED)"
// @Failure 400 {object} ErrorResponse "ID inválido"
// @Failure 401 {object} ErrorResponse "No autorizado"
// @Failure 404 {object} ErrorResponse "Envío no encontrado"
// @Failure 409 {object} ErrorResponse "Transición no permitida (estado actual no es CREATED)"
// @Failure 500 {object} ErrorResponse "Error interno"
// @Router /parcels/{id}/register [post]
func (h *ParcelHandler) Register(c *gin.Context) {
    // implementación
}
```

---

### Ejemplo 4: Operación con Body (PUT/POST con Body)

```go
// BoardParcel godoc
// @Summary Embarcar envío
// @Description Transiciona el envío a estado BOARDED. Asigna el envío a un 
// vehículo y opcionalmente a un viaje.
// @Tags Parcels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID del envío" Format(uuid)
// @Param payload body BoardParcelRequest true "Solicitud con UUID vehículo (requerido)"
// @Success 200 {object} ParcelResponse "Envío embarcado (estado: BOARDED)"
// @Failure 400 {object} ErrorResponse "Validación fallida"
// @Failure 401 {object} ErrorResponse "No autorizado"
// @Failure 404 {object} ErrorResponse "Envío no encontrado"
// @Failure 409 {object} ErrorResponse "Estado incompatible o vehículo inválido"
// @Failure 500 {object} ErrorResponse "Error interno"
// @Router /parcels/{id}/board [post]
func (h *ParcelHandler) Board(c *gin.Context) {
    // implementación
}
```

---

## Cómo Generar el Archivo Swagger

### Paso 1: Instalar swag

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Paso 2: Generar Swagger JSON

Navega a la raíz del proyecto y ejecuta:

```bash
swag init -g cmd/api/main.go -o docs
```

Esto generará:
- `docs/swagger.json` - Especificación OpenAPI
- `docs/swagger.yaml` - Especificación en YAML (opcional)

### Paso 3: Verificar Documentación

Inicia el servidor:

```bash
go run cmd/api/main.go
```

Accede a Swagger UI en:

```
http://localhost:8080/swagger/index.html
```

---

## Validación y Testing

### 1. Validación Manual en Swagger UI

1. Abre `http://localhost:8080/swagger/index.html`
2. Haz clic en "Authorize" e ingresa tu token Bearer
3. Expande cada endpoint y verifica:
   - ✅ Summary es claro y conciso
   - ✅ Description está completo
   - ✅ Parámetros están documentados
   - ✅ Respuestas tienen ejemplos claros
   - ✅ Errores están correctamente mapeados

### 2. Validación de JSON

```bash
# Validate swagger.json
python -m json.tool docs/swagger.json > /dev/null && echo "Valid"
```

### 3. Testing con curl

```bash
# Test endpoint con token Bearer
curl -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"field":"value"}' \
     http://localhost:8080/parcels
```

---

## Checklist de Actualización

Cuando **crees o modificas** un endpoint, asegúrate de:

- [ ] Agregar comentario `godoc` con formato `FunctionName godoc`
- [ ] Summary (máx 80 caracteres)
- [ ] Description (detallado, con contexto)
- [ ] @Tags con valor válido
- [ ] @Accept json (si aplica)
- [ ] @Produce json
- [ ] @Security BearerAuth
- [ ] @Param Authorization header (si es privado)
- [ ] @Param para cada path parameter (con Format y descripción)
- [ ] @Param para cada query parameter (con min/max si aplica)
- [ ] @Param para body (si aplica)
- [ ] @Success 200 o 201 con objeto correcto
- [ ] @Failure 400 con descripción de validación
- [ ] @Failure 401 con mensaje de autenticación
- [ ] @Failure 404 si aplica
- [ ] @Failure 409 para conflictos de negocio
- [ ] @Failure 500 para errores internos
- [ ] @Router con método HTTP correcto

---

### Antes de Hacer Commit

```bash
# 1. Generar Swagger
swag init -g cmd/api/main.go -o docs

# 2. Validar compilación
go build ./cmd/api

# 3. Verificar en Swagger UI
go run cmd/api/main.go &
# Abre http://localhost:8080/swagger/index.html

# 4. Commit
git add .
git commit -m "feat: add/update endpoint documentation in Swagger"
```

---

## Errores Comunes

### ❌ Error: "Multiple Failure Responses"

```go
// MALO: Demasiados failure codes
// @Failure 200 // No tiene sentido
// @Failure 201 // Usar Success en su lugar
// @Failure 300 // No es estándar
```

✅ **Usar solo códigos HTTP estándar:**
- `200` Success (GET, PUT, DELETE)
- `201` Created (POST)
- `204` No Content (DELETE sin body)
- `400` Bad Request
- `401` Unauthorized
- `404` Not Found
- `409` Conflict
- `500` Internal Server Error

---

### ❌ Error: "Parameter format not recognized"

```go
// MALO
// @Param id path string "ID"

// BUENO
// @Param id path string true "UUID del envío" Format(uuid)
```

Formatos válidos:
- `Format(uuid)` - UUID
- `Format(date-time)` - RFC3339
- `Format(email)` - Email
- `minimum(0)`, `maximum(100)` - Validación numérica
- `oneof=VALUE1,VALUE2` - Enum

---

### ❌ Error: "Response model not found"

```go
// MALO
// @Success 200 {object} UnknownDTO

// BUENO (debe existir la struct)
// @Success 200 {object} ParcelResponse
```

Asegúrate de que el DTO esté definido en el mismo paquete o importado.

---

## Tips y Buenas Prácticas

### 1. Mantener Consistencia

Usa las mismas descripciones de error para operaciones similares:

```go
// Consistente para todas las validaciones
// @Failure 400 {object} ErrorResponse "Validación fallida: payload malformado o valores inválidos"

// Consistente para autenticación
// @Failure 401 {object} ErrorResponse "No autorizado: token inválido o credenciales faltantes"

// Consistente para recursos no encontrados
// @Failure 404 {object} ErrorResponse "Recurso no encontrado"
```

---

### 2. Documentar Comodines y Reglas

Para endpoints con lógica especial:

```go
// @Description Crea una regla de precios. Soporta comodines (*) en 
// ShipmentType, OriginOfficeID y DestinationOfficeID. La prioridad (0-100) 
// determina el orden en búsquedas jerárquicas: específicas primero, luego comodines.
```

---

### 3. Incluir Ejemplos en Description

```go
// @Description Agrega un item con cálculo automático de peso facturable. 
// Ejemplos: weight_kg=5.5 (real), length_cm=55, width_cm=25, height_cm=5 
// (volumétrico: 6.875). El mayor se usa como facturable.
```

---

## Recursos Adicionales

- [Swag Documentation](https://github.com/swaggo/swag)
- [OpenAPI 3.0 Spec](https://spec.openapis.org/oas/v3.0.0)
- [Gin Swagger](https://github.com/swaggo/gin-swagger)
- [Swagger UI Demo](https://petstore.swagger.io/)

---

**Última Actualización:** 20 de enero de 2026  
**Responsable:** Equipo de Desarrollo - QuatroBus  
**Próxima Revisión:** 27 de enero de 2026
