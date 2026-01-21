# Referencia Completa de Endpoints Swagger - QuatroBus Parcel

**Versión:** 1.0  
**Fecha de Actualización:** 20 de enero de 2026  
**Descripción:** Documentación completa de todos los endpoints de la API de Parcel con tags, parámetros, respuestas y ejemplos de error.

---

## Tabla de Contenidos

1. [Introducción](#introducción)
2. [Autenticación](#autenticación)
3. [Tags y Categorías](#tags-y-categorías)
4. [Endpoints de Parcels (Envíos)](#endpoints-de-parcels-envíos)
5. [Endpoints de Parcel Items (Artículos)](#endpoints-de-parcel-items-artículos)
6. [Endpoints de Parcel Payments (Pagos)](#endpoints-de-parcel-payments-pagos)
7. [Endpoints de Parcel Tracking (Historial)](#endpoints-de-parcel-tracking-historial)
8. [Endpoints de Manifests (Manifiestos)](#endpoints-de-manifests-manifiestos)
9. [Endpoints de Pricing (Precios)](#endpoints-de-pricing-precios)
10. [Endpoints de Documents (Documentos)](#endpoints-de-documents-documentos)

---

## Introducción

Los endpoints de QuatroBus Parcel están organizados por **dominio modular**:
- **Parcels**: Gestión de envíos y estados
- **ParcelItems**: Artículos/bultos dentro de envíos
- **ParcelPayments**: Información de pago
- **ParcelTracking**: Historial de eventos
- **Manifests**: Manifiesto virtual de cargas
- **Pricing**: Reglas y cálculo de precios
- **ParcelDocuments**: Impresión y documentación

Todos los endpoints requieren **token Bearer JWT** en header `Authorization`.

---

## Autenticación

### Header Requerido

```http
Authorization: Bearer <jwt_token>
```

### Claims Esperados

El token JWT debe contener:
- `tenant_id`: ID del tenant (obligatorio)
- `user_id`: ID del usuario (obligatorio)
- `user_name`: Nombre del usuario (obligatorio)

### Respuesta de Falta de Autenticación

```json
{
  "success": false,
  "error": {
    "code": "unauthorized",
    "message": "credenciales inválidas",
    "http_status": 401,
    "timestamp": "2026-01-20T10:30:00Z"
  }
}
```

---

## Tags y Categorías

| Tag | Descripción | Responsable |
|-----|-------------|-------------|
| **Parcels** | Creación, listado y gestión del ciclo de vida de envíos | parcel_core |
| **ParcelItems** | Gestión de artículos/bultos dentro de envíos | parcel_item |
| **ParcelPayments** | Información de pago y transiciones de estado | parcel_payment |
| **ParcelTracking** | Historial completo de eventos y cambios de estado | parcel_tracking |
| **Manifests** | Manifiesto virtual y preview de cargas | parcel_manifest |
| **Pricing** | Reglas de precios y cálculo de tarifas | parcel_pricing |
| **ParcelDocuments** | Impresión y documentación de envíos | parcel_documents |

---

## Endpoints de Parcels (Envíos)

### 1. Crear Nuevo Envío

```
POST /parcels
```

**Tag:** `Parcels`

**Descripción:**  
Crea un nuevo envío en estado `CREATED`. Requiere tipos de envío, oficinas origen/destino, personas (remitente/destinatario) y opcionales notes. El `package_key` permite proteger operaciones posteriores con confirmación.

**Request Body:**
```json
{
  "shipment_type": "STANDARD",
  "origin_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "destination_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d480",
  "sender_person_id": "f47ac10b-58cc-4372-a567-0e02b2c3d481",
  "recipient_person_id": "f47ac10b-58cc-4372-a567-0e02b2c3d482",
  "notes": "Frágil - Manejar con cuidado",
  "package_key": "ABC123",
  "package_key_confirm": "ABC123"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "data": {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d483",
    "status": "CREATED",
    "shipment_type": "STANDARD",
    "origin_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "destination_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d480",
    "sender_person_id": "f47ac10b-58cc-4372-a567-0e02b2c3d481",
    "recipient_person_id": "f47ac10b-58cc-4372-a567-0e02b2c3d482",
    "notes": "Frágil - Manejar con cuidado",
    "created_at": "2026-01-20T10:30:00Z"
  }
}
```

**Error Responses:**
- `400`: Payload malformado o valores inválidos
- `401`: No autorizado
- `409`: Offices o personas no válidas
- `500`: Error interno

---

### 2. Listar Envíos

```
GET /parcels
```

**Tag:** `Parcels`

**Descripción:**  
Lista envíos del tenant con filtros y paginación. Soporta búsqueda por query, filtros por estado, offices, personas, vehículo y rango de fechas.

**Query Parameters:**
- `q` (string): Búsqueda por código o ID
- `status` (string): Estado (CREATED, REGISTERED, BOARDED, EN_ROUTE, ARRIVED, DELIVERED)
- `origin_office_id` (string, UUID): Filtrar por oficina origen
- `destination_office_id` (string, UUID): Filtrar por oficina destino
- `sender_person_id` (string, UUID): Filtrar por remitente
- `recipient_person_id` (string, UUID): Filtrar por destinatario
- `vehicle_id` (string, UUID): Filtrar por vehículo
- `from_created_at` (string, RFC3339): Desde
- `to_created_at` (string, RFC3339): Hasta
- `limit` (int): Límite (default: 50, max: 200)
- `offset` (int): Desplazamiento (default: 0)

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "items": [...],
    "pagination": {
      "limit": 50,
      "offset": 0,
      "count": 125
    }
  }
}
```

---

### 3. Obtener Detalles del Envío

```
GET /parcels/{id}
```

**Tag:** `Parcels`

**Descripción:**  
Devuelve información detallada de un envío incluyendo estado actual, timeline de transiciones, asignación de vehículo/viaje, y auditoría.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d483",
    "status": "BOARDED",
    "shipment_type": "STANDARD",
    "origin_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "destination_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d480",
    "created_at": "2026-01-20T10:30:00Z",
    "registered_at": "2026-01-20T10:35:00Z",
    "boarded_vehicle_id": "f47ac10b-58cc-4372-a567-0e02b2c3d484",
    "boarded_trip_id": "f47ac10b-58cc-4372-a567-0e02b2c3d485",
    "boarded_at": "2026-01-20T10:40:00Z",
    "boarded_by_user_id": "f47ac10b-58cc-4372-a567-0e02b2c3d486"
  }
}
```

---

### 4. Registrar Envío

```
POST /parcels/{id}/register
```

**Tag:** `Parcels`

**Descripción:**  
Transiciona el envío de estado `CREATED` a `REGISTERED`. Marca el envío como listo para ser embarcado.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d483",
    "status": "REGISTERED",
    "registered_at": "2026-01-20T10:35:00Z"
  }
}
```

**Error Responses:**
- `404`: Envío no encontrado
- `409`: Transición no permitida (estado actual no es CREATED)

---

### 5. Embarcar Envío

```
POST /parcels/{id}/board
```

**Tag:** `Parcels`

**Descripción:**  
Transiciona el envío de estado `REGISTERED` a `BOARDED`. Asigna el envío a un vehículo específico y opcionalmente a un viaje.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

**Request Body:**
```json
{
  "vehicle_id": "f47ac10b-58cc-4372-a567-0e02b2c3d484",
  "trip_id": "f47ac10b-58cc-4372-a567-0e02b2c3d485",
  "departure_at": "2026-01-20T11:00:00Z",
  "origin_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
}
```

---

### 6. Registrar Salida del Envío

```
POST /parcels/{id}/depart
```

**Tag:** `Parcels`

**Descripción:**  
Transiciona el envío de estado `BOARDED` a `EN_ROUTE`. Confirma la partida real del vehículo.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

**Request Body:**
```json
{
  "departure_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "vehicle_id": "f47ac10b-58cc-4372-a567-0e02b2c3d484",
  "departed_at": "2026-01-20T11:15:00Z"
}
```

---

### 7. Registrar Llegada del Envío

```
POST /parcels/{id}/arrive
```

**Tag:** `Parcels`

**Descripción:**  
Transiciona el envío de estado `EN_ROUTE` a `ARRIVED`. Marca la llegada a destino.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

**Request Body:**
```json
{
  "destination_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d480"
}
```

---

### 8. Entregar Envío

```
POST /parcels/{id}/deliver
```

**Tag:** `Parcels`

**Descripción:**  
Transiciona el envío de estado `ARRIVED` a `DELIVERED`. Finaliza la cadena de custodia.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

**Request Body:**
```json
{
  "package_key": "ABC123"
}
```

---

## Endpoints de Parcel Items (Artículos)

### 1. Agregar Artículo

```
POST /parcels/{id}/items
```

**Tag:** `ParcelItems`

**Descripción:**  
Agrega un bulto/artículo al envío con cálculo automático de peso facturable y precio. Soporta dimensiones opcionales para cálculo de peso volumétrico.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

**Request Body:**
```json
{
  "description": "Monitor de 24 pulgadas",
  "quantity": 2,
  "weight_kg": 5.5,
  "length_cm": 55.0,
  "width_cm": 25.0,
  "height_cm": 5.0,
  "unit_price": 299.99,
  "content_type": "ELECTRONICS",
  "notes": "Frágil"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "data": {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d500",
    "parcel_id": "f47ac10b-58cc-4372-a567-0e02b2c3d483",
    "description": "Monitor de 24 pulgadas",
    "quantity": 2,
    "weight_kg": 5.5,
    "length_cm": 55.0,
    "width_cm": 25.0,
    "height_cm": 5.0,
    "volumetric_weight": 7.656,
    "billable_weight": 7.656,
    "unit_price": 299.99,
    "created_at": "2026-01-20T10:40:00Z"
  }
}
```

**Error Responses:**
- `400`: Payload malformado
- `404`: Envío no encontrado
- `409`: Regla de precios no encontrada (sugerencia: use comodines *)

---

### 2. Listar Artículos

```
GET /parcels/{id}/items
```

**Tag:** `ParcelItems`

**Descripción:**  
Lista todos los artículos agregados a un envío.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

**Success Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": "f47ac10b-58cc-4372-a567-0e02b2c3d500",
      "parcel_id": "f47ac10b-58cc-4372-a567-0e02b2c3d483",
      "description": "Monitor de 24 pulgadas",
      "quantity": 2,
      "weight_kg": 5.5,
      "billable_weight": 7.656,
      "unit_price": 299.99,
      "created_at": "2026-01-20T10:40:00Z"
    }
  ]
}
```

---

### 3. Eliminar Artículo

```
DELETE /parcels/{id}/items/{item_id}
```

**Tag:** `ParcelItems`

**Descripción:**  
Elimina un artículo específico del envío.

**Path Parameters:**
- `id` (string, UUID): UUID del envío
- `item_id` (string, UUID): UUID del item

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "deleted": true
  }
}
```

---

## Endpoints de Parcel Payments (Pagos)

### 1. Crear o Actualizar Pago

```
PUT /parcels/{id}/payment
```

**Tag:** `ParcelPayments`

**Descripción:**  
Crea o actualiza la información de pago del envío. Soporta múltiples formas de pago.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

**Request Body:**
```json
{
  "payment_type": "CASH",
  "amount": 150.00,
  "currency": "PEN",
  "channel": "COUNTER",
  "office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "cashbox_id": "CAJA-001",
  "seller_user_id": "f47ac10b-58cc-4372-a567-0e02b2c3d487",
  "notes": "Pago en efectivo"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d510",
    "parcel_id": "f47ac10b-58cc-4372-a567-0e02b2c3d483",
    "payment_type": "CASH",
    "amount": 150.00,
    "currency": "PEN",
    "status": "PENDING",
    "channel": "COUNTER",
    "created_at": "2026-01-20T10:45:00Z"
  }
}
```

---

### 2. Obtener Información de Pago

```
GET /parcels/{id}/payment
```

**Tag:** `ParcelPayments`

**Descripción:**  
Devuelve los detalles completos del pago registrado para un envío.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

---

### 3. Marcar Pago como Realizado

```
POST /parcels/{id}/payment/mark-paid
```

**Tag:** `ParcelPayments`

**Descripción:**  
Transiciona el pago a estado `PAID`. Integración con servicio de caja.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d510",
    "status": "PAID",
    "paid_at": "2026-01-20T10:50:00Z",
    "paid_by_user_id": "f47ac10b-58cc-4372-a567-0e02b2c3d488"
  }
}
```

---

## Endpoints de Parcel Tracking (Historial)

### 1. Listar Historial de Tracking

```
GET /parcels/{id}/tracking
```

**Tag:** `ParcelTracking`

**Descripción:**  
Lista todos los eventos y cambios de estado del envío en orden cronológico.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

**Success Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": "f47ac10b-58cc-4372-a567-0e02b2c3d520",
      "parcel_id": "f47ac10b-58cc-4372-a567-0e02b2c3d483",
      "event_type": "PARCEL_CREATED",
      "occurred_at": "2026-01-20T10:30:00Z",
      "user_id": "f47ac10b-58cc-4372-a567-0e02b2c3d486",
      "user_name": "Juan Pérez",
      "metadata": {
        "status": "CREATED",
        "notes": "Envío creado"
      }
    },
    {
      "id": "f47ac10b-58cc-4372-a567-0e02b2c3d521",
      "parcel_id": "f47ac10b-58cc-4372-a567-0e02b2c3d483",
      "event_type": "PARCEL_REGISTERED",
      "occurred_at": "2026-01-20T10:35:00Z",
      "user_id": "f47ac10b-58cc-4372-a567-0e02b2c3d486",
      "user_name": "Juan Pérez",
      "metadata": {
        "status": "REGISTERED"
      }
    }
  ]
}
```

---

## Endpoints de Manifests (Manifiestos)

### 1. Construir Preview de Manifiesto (POST)

```
POST /manifests/preview
```

**Tag:** `Manifests`

**Descripción:**  
Construye un manifiesto virtual basado en envíos pendientes entre oficinas.

**Request Body:**
```json
{
  "vehicle_id": "f47ac10b-58cc-4372-a567-0e02b2c3d484",
  "origin_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "destination_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d480"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "parcels": [...],
    "totals": {
      "count": 45,
      "weight_kg": 350.5,
      "volume_cm3": 125000
    }
  }
}
```

---

### 2. Construir Preview de Manifiesto (GET)

```
GET /manifests/preview?vehicle_id=...&origin_office_id=...&destination_office_id=...
```

**Tag:** `Manifests`

**Descripción:**  
Construye un manifiesto virtual basado en parámetros de query.

**Query Parameters:**
- `vehicle_id` (string, UUID): UUID del vehículo
- `origin_office_id` (string, UUID): UUID de la oficina origen
- `destination_office_id` (string, UUID): UUID de la oficina destino

---

## Endpoints de Pricing (Precios)

### 1. Crear Regla de Precios

```
POST /pricing/rules
```

**Tag:** `Pricing`

**Descripción:**  
Crea una nueva regla de precios para el tenant. Soporta comodines (`*`) en campos de ruta.

**Request Body:**
```json
{
  "shipment_type": "STANDARD",
  "origin_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "destination_office_id": "*",
  "unit": "PER_KG",
  "price": 2.50,
  "currency": "PEN",
  "priority": 50,
  "active": true
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d530",
    "shipment_type": "STANDARD",
    "origin_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "destination_office_id": "*",
    "unit": "PER_KG",
    "price": 2.50,
    "currency": "PEN",
    "priority": 50,
    "active": true,
    "created_at": "2026-01-20T10:55:00Z"
  }
}
```

**Error Responses:**
- `400`: Valores inválidos
- `409`: Regla duplicada

---

### 2. Actualizar Regla de Precios

```
PUT /pricing/rules/{id}
```

**Tag:** `Pricing`

**Descripción:**  
Actualiza una regla de precios existente.

**Path Parameters:**
- `id` (string, UUID): UUID de la regla

---

### 3. Listar Reglas de Precios

```
GET /pricing/rules
```

**Tag:** `Pricing`

**Descripción:**  
Lista todas las reglas de precios activas del tenant.

**Success Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": "f47ac10b-58cc-4372-a567-0e02b2c3d530",
      "shipment_type": "STANDARD",
      "origin_office_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
      "destination_office_id": "*",
      "unit": "PER_KG",
      "price": 2.50,
      "priority": 50,
      "active": true
    }
  ]
}
```

---

## Endpoints de Documents (Documentos)

### 1. Registrar Impresión de Documento

```
POST /parcels/{id}/documents/print
```

**Tag:** `ParcelDocuments`

**Descripción:**  
Registra un evento de impresión de documento para un envío.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

**Request Body:**
```json
{
  "document_type": "LABEL"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "record": {
      "id": "f47ac10b-58cc-4372-a567-0e02b2c3d540",
      "parcel_id": "f47ac10b-58cc-4372-a567-0e02b2c3d483",
      "document_type": "LABEL",
      "printed_at": "2026-01-20T11:00:00Z",
      "printed_by_user_id": "f47ac10b-58cc-4372-a567-0e02b2c3d488"
    },
    "meta": {
      "qr_code": "...",
      "barcode": "..."
    }
  }
}
```

---

### 2. Listar Impresiones de Envío

```
GET /parcels/{id}/documents/prints
```

**Tag:** `ParcelDocuments`

**Descripción:**  
Lista todos los registros de impresión asociados a un envío.

**Path Parameters:**
- `id` (string, UUID): UUID del envío

---

## Resumen de Códigos de Error HTTP

| Código | Descripción | Causa Común |
|--------|-------------|------------|
| **400** | Bad Request | Validación fallida, payload malformado, parámetros inválidos |
| **401** | Unauthorized | Token inválido, credenciales faltantes |
| **404** | Not Found | Recurso no encontrado (parcel, item, etc.) |
| **409** | Conflict | Conflicto de negocio, transición no permitida, duplicado |
| **500** | Internal Server Error | Error interno del servidor |

---

## Modelo de Respuesta Uniforme

### Success Response
```json
{
  "success": true,
  "data": {
    // Contenido específico del endpoint
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "error_code",
    "message": "Descripción del error",
    "details": {
      "field": "valor_adicional"
    },
    "http_status": 400,
    "timestamp": "2026-01-20T11:00:00Z"
  }
}
```

---

## Tips para Usar Swagger UI

1. **Generar Swagger JSON:**
   ```bash
   swag init -g cmd/api/main.go
   ```

2. **Acceder a Swagger UI:**
   - Inicia el servidor: `go run cmd/api/main.go`
   - Navega a: `http://localhost:8080/swagger/index.html`

3. **Autorización:**
   - En Swagger UI, haz clic en "Authorize"
   - Pega tu token JWT Bearer

4. **Try It Out:**
   - Haz clic en "Try it out" para probar endpoints
   - Rellena parámetros y body
   - Haz clic en "Execute"

---

## Documentación Adicional

- **Guía de Reglas de Precios:** Ver [pricing_rules_guide.md](./pricing_rules_guide.md)
- **Arquitectura:** Ver [../instructions/parcel_boundaries.instructions.md](../instructions/parcel_boundaries.instructions.md)
- **Perfil Go:** Ver [../instructions/go_profile.instructions.md](../instructions/go_profile.instructions.md)

---

**Última Actualización:** 20 de enero de 2026  
**Responsable:** Equipo de Desarrollo - QuatroBus
