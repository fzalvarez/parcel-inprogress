# Arquitectura QuatroBus Parcel - Diagrama Visual

## Vista General del Sistema

```
┌─────────────────────────────────────────────────────────────────┐
│                     QUATROBUS PARCEL API                         │
│                    (Monolito Modular - Go/Gin)                   │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                      CAPA HTTP (Gin)                             │
├─────────────────────────────────────────────────────────────────┤
│  Handlers:                                                       │
│  ├─ ParcelHandler          (CRUD + transiciones de estado)     │
│  ├─ ParcelItemHandler      (gestión de bultos/artículos)       │
│  ├─ ParcelPaymentHandler   (pagos y marcado como pagado)       │
│  ├─ ParcelTrackingHandler  (historial de eventos)              │
│  ├─ ParcelDocumentsHandler (impresiones)                        │
│  ├─ ManifestHandler        (preview de manifiestos)             │
│  └─ PriceRuleHandler       (reglas de precios)                  │
│                                                                  │
│  Middleware:                                                     │
│  ├─ ErrorMiddleware        (manejo global de errores)           │
│  ├─ DevClaimsMiddleware    (inyección de claims en desarrollo) │
│  └─ (auth, request-id, logging - pendientes)                   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    CAPA DE CASOS DE USO                          │
├─────────────────────────────────────────────────────────────────┤
│  Parcel Core:                                                    │
│  ├─ CreateParcelUseCase                                          │
│  ├─ RegisterParcelUseCase                                        │
│  ├─ BoardParcelUseCase                                           │
│  ├─ DepartParcelUseCase                                          │
│  ├─ ArriveParcelUseCase                                          │
│  ├─ DeliverParcelUseCase                                         │
│  ├─ GetParcelUseCase                                             │
│  └─ ListParcelsUseCase                                           │
│                                                                  │
│  Parcel Item:                                                    │
│  ├─ AddParcelItemUseCase    (✨ calcula peso volumétrico)      │
│  ├─ ListParcelItemsUseCase                                       │
│  └─ DeleteParcelItemUseCase                                      │
│                                                                  │
│  Parcel Payment:                                                 │
│  ├─ CreateOrUpdateParcelPaymentUseCase                           │
│  ├─ GetParcelPaymentUseCase                                      │
│  └─ MarkPaidParcelPaymentUseCase                                 │
│                                                                  │
│  Parcel Tracking:                                                │
│  └─ RecordTrackingEventUseCase                                   │
│                                                                  │
│  Parcel Documents:                                               │
│  └─ RegisterPrintUseCase                                         │
│                                                                  │
│  Parcel Manifest:                                                │
│  └─ BuildManifestPreviewUseCase                                  │
│                                                                  │
│  Parcel Pricing:                                                 │
│  ├─ CreatePriceRuleUseCase  (✨ soporta comodines *)           │
│  ├─ UpdatePriceRuleUseCase                                       │
│  └─ ListPriceRulesUseCase                                        │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    CAPA DE PUERTOS (Interfaces)                  │
├─────────────────────────────────────────────────────────────────┤
│  Repositories (out):                                             │
│  ├─ ParcelRepository                                             │
│  ├─ ParcelItemRepository                                         │
│  ├─ ParcelPaymentRepository                                      │
│  ├─ TrackingRepository                                           │
│  ├─ PrintRepository                                              │
│  └─ PriceRuleRepository     (✨ búsqueda jerárquica)           │
│                                                                  │
│  Clients (out):                                                  │
│  ├─ TenantConfigClient      (feature flags, opciones)            │
│  ├─ CashboxClient           (integración con cajas)              │
│  └─ QRGenerator             (generación de QR codes)             │
│                                                                  │
│  Readers (out):                                                  │
│  └─ ParcelReader            (queries optimizadas)                │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│              CAPA DE INFRAESTRUCTURA (Adaptadores)               │
├─────────────────────────────────────────────────────────────────┤
│  Persistence:                                                    │
│  ├─ database/                                                    │
│  │   ├─ connect.go          (🆕 conexión PostgreSQL)           │
│  │   └─ migrate.go          (🆕 migraciones automáticas)       │
│  │                                                               │
│  ├─ postgres/                                                    │
│  │   ├─ tenant_scope.go     (🆕 multi-tenancy automático)      │
│  │   ├─ parcel_model.go     (🆕 DBParcel)                      │
│  │   ├─ parcel_item_model.go (🆕 DBParcelItem)                 │
│  │   ├─ parcel_payment_model.go (🆕 DBParcelPayment)           │
│  │   ├─ tracking_event_model.go (🆕 DBTrackingEvent)           │
│  │   ├─ print_record_model.go (🆕 DBPrintRecord)               │
│  │   ├─ price_rule_model.go (🆕 DBPriceRule)                   │
│  │   └─ *_postgres_repository.go (⏳ pendientes)               │
│  │                                                               │
│  └─ memory/                  (🔄 mover repos in-memory aquí)    │
│      ├─ in_memory_parcel_repository.go                           │
│      ├─ in_memory_parcel_item_repository.go                      │
│      ├─ in_memory_parcel_payment_repository.go                   │
│      ├─ in_memory_tracking_repository.go                         │
│      ├─ in_memory_print_repository.go                            │
│      └─ in_memory_price_rule_repository.go (✨ jerarquía OK)   │
│                                                                  │
│  Clients:                                                        │
│  └─ TenantConfigStubClient   (stub para desarrollo)              │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    CAPA DE DOMINIO                               │
├─────────────────────────────────────────────────────────────────┤
│  Entidades:                                                      │
│  ├─ Parcel                   (estados, tracking_code, etc.)      │
│  ├─ ParcelItem               (✨ peso volumétrico, facturable) │
│  ├─ ParcelPayment            (tipos de pago, estados)            │
│  ├─ TrackingEvent            (historial de eventos)              │
│  ├─ PrintRecord              (registros de impresión)            │
│  ├─ PriceRule                (✨ comodines, prioridad)          │
│  └─ ManifestPreview          (vista previa de manifiesto)        │
│                                                                  │
│  Value Objects:                                                  │
│  ├─ TenantOptions            (✨ volumetric_enabled, divisor)  │
│  ├─ ParcelStatus             (CREATED → DELIVERED)               │
│  ├─ ShipmentType             (BUS, CARGUERO)                     │
│  ├─ PaymentType              (CASH, FOB, CARD, etc.)             │
│  ├─ PaymentStatus            (PENDING, PAID)                     │
│  ├─ DocumentType             (LABEL, RECEIPT, MANIFEST, GUIDE)   │
│  └─ PriceUnit                (PER_KG, PER_ITEM)                  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                  SERVICIOS EXTERNOS (No implementados)           │
├─────────────────────────────────────────────────────────────────┤
│  ├─ IAM Service              (autenticación, autorización)       │
│  ├─ PERSON Service           (gestión de personas)               │
│  ├─ LOCATION Service         (oficinas, zonas)                   │
│  ├─ VEHICLE Service          (vehículos, viajes)                 │
│  ├─ TRIP/SCHEDULE Service    (horarios, rutas)                   │
│  ├─ PAYMENT Service          (procesamiento de pagos)            │
│  ├─ TENANT-CONFIG Service    (✨ feature flags, opciones)       │
│  ├─ NOTIFICATION Service     (alertas, notificaciones)           │
│  └─ TICKETING Service        (boletos)                           │
└─────────────────────────────────────────────────────────────────┘
```

---

## Flujo de Datos - Ejemplo: Agregar Item

```
1. HTTP Request
   POST /parcels/{id}/items
   Body: { description, quantity, weight_kg, length_cm, width_cm, height_cm }
   ↓
   
2. Handler (ParcelItemHandler.Add)
   - Valida JWT y extrae tenant_id
   - Parsea UUID del path param
   - Valida request body (binding)
   ↓
   
3. UseCase (AddParcelItemUseCase.Execute)
   - Consulta TenantOptions desde TenantConfigClient
   - Calcula peso volumétrico si aplica:
     volumetric_weight = (L × W × H) / divisor
   - Determina peso facturable:
     billable_weight = max(weight_kg, volumetric_weight)
   - Busca precio en PriceRuleRepository:
     * Intenta regla específica: origin → destination
     * Si no existe, busca: origin → *
     * Si no existe, busca: * → destination
     * Si no existe, busca: * → *
     * Si sigue sin encontrar y allow_manual_price=true → OK
     * Si no, retorna error con sugerencia de crear regla
   - Calcula precio total
   - Crea ParcelItem
   - Persiste en ParcelItemRepository
   ↓
   
4. Repository (InMemoryParcelItemRepository)
   - Genera UUID
   - Guarda en map[tenantID]map[itemID]ParcelItem
   - Retorna ParcelItem creado
   ↓
   
5. Handler (respuesta)
   - Convierte domain → DTO
   - Retorna JSON 201 Created
```

---

## Estados del Parcel

```
CREATED
   ↓ Register
REGISTERED
   ↓ Board (asignar vehículo)
BOARDED
   ↓ Depart (salir de oficina)
EN_ROUTE
   ↓ Arrive (llegar a destino)
ARRIVED
   ↓ Deliver (entregar al destinatario)
DELIVERED
```

---

## Motor de Pricing - Búsqueda Jerárquica

```
Configuración tenant: volumetric_enabled=true, divisor=5000

Ejemplo: Agregar item de Lima → Arequipa, 10kg, 50×40×30 cm

Paso 1: Calcular peso volumétrico
  volumetric_weight = (50 × 40 × 30) / 5000 = 12 kg

Paso 2: Determinar peso facturable
  billable_weight = max(10, 12) = 12 kg

Paso 3: Buscar regla de precio (jerarquía de especificidad)
  ┌─────────────────────┬─────────────┬──────────┐
  │ Búsqueda            │ Especifidad │ Prioridad│
  ├─────────────────────┼─────────────┼──────────┤
  │ Lima → Arequipa     │ 20 (exact)  │ Alta     │
  │ Lima → *            │ 11 (hybrid) │ Media    │
  │ * → Arequipa        │ 11 (hybrid) │ Media    │
  │ * → *               │ 2 (wildcard)│ Baja     │
  └─────────────────────┴─────────────┴──────────┘

  Score = origin_score + dest_score
  - Exact match: 10 puntos
  - Wildcard: 1 punto

Paso 4: Aplicar precio
  Si encuentra regla: unit_price × billable_weight
  Si no encuentra y allow_manual_price=true: usar precio manual del request
  Si no: error "no_price_rule_found" con sugerencia
```

---

## Documentación Swagger

Todos los endpoints están documentados con:
- `@Summary` - Descripción corta
- `@Description` - Descripción detallada de funcionalidad
- `@Tags` - Agrupación por módulo
- `@Param` - Especificación de parámetros
- `@Success` - Respuestas exitosas
- `@Failure` - Códigos de error posibles

Ver: [swagger_endpoints_reference.md](./swagger_endpoints_reference.md)

---

## Leyenda

- ✅ Implementado y funcionando
- ✨ Implementado con mejoras recientes
- 🆕 Creado recientemente (arquitectura nueva)
- 🔄 Requiere reorganización
- ⏳ Pendiente de implementación
- ⚠️ Requiere ajustes

---

**Última Actualización:** 21 de enero de 2026  
**Responsable:** Equipo de Desarrollo - QuatroBus
