param(
  [string]$ModuleName = "quatrobus-parcel",
  [switch]$InitGoMod = $false
)

$ErrorActionPreference = "Stop"

function Ensure-Dir([string]$Path) {
  if (-not (Test-Path $Path)) { New-Item -ItemType Directory -Path $Path | Out-Null }
}

function Ensure-File([string]$Path, [string]$Content) {
  if (-not (Test-Path $Path)) {
    $parent = Split-Path $Path -Parent
    if ($parent) { Ensure-Dir $parent }
    Set-Content -Path $Path -Value $Content -Encoding UTF8
  }
}

# --- Folders (mínimo útil) ---
$dirs = @(
  ".github",
  ".github/instructions",
  ".github/prompts",
  "scripts",

  "app",
  "internal/config",

  "internal/pkg/util/apperror",
  "internal/pkg/util/resilience",

  "internal/infrastructure/logger",
  "internal/infrastructure/http/router",
  "internal/infrastructure/http/middleware",
  "internal/infrastructure/http/handler",
  "internal/infrastructure/persistence/database",
  "internal/infrastructure/persistence/postgres",
  "internal/infrastructure/clients",

  # Monolito modular de paquetería
  "internal/parcel",

  "internal/parcel/parcel_core/domain",
  "internal/parcel/parcel_core/port",
  "internal/parcel/parcel_core/usecase",

  "internal/parcel/parcel_item/domain",
  "internal/parcel/parcel_item/port",
  "internal/parcel/parcel_item/usecase",

  "internal/parcel/parcel_tracking/domain",
  "internal/parcel/parcel_tracking/port",
  "internal/parcel/parcel_tracking/usecase",

  "internal/parcel/parcel_payment/domain",
  "internal/parcel/parcel_payment/port",
  "internal/parcel/parcel_payment/usecase",

  "internal/parcel/parcel_manifest/domain",
  "internal/parcel/parcel_manifest/port",
  "internal/parcel/parcel_manifest/usecase",

  "internal/parcel/parcel_documents/domain",
  "internal/parcel/parcel_documents/port",
  "internal/parcel/parcel_documents/usecase",

  "internal/parcel/parcel_incident/domain",
  "internal/parcel/parcel_incident/port",
  "internal/parcel/parcel_incident/usecase"
)

$dirs | ForEach-Object { Ensure-Dir $_ }

# --- Copilot instructions (repo-wide) ---
$copilotInstructions = @"
# QuatroBus Parcel - Instrucciones Mandatorias

**NO cambies el estilo existente ni reestructures carpetas.**
- No renombres paquetes/archivos ya creados.
- No “refactors” masivos.
- Cambia SOLO los archivos explícitamente solicitados, o crea nuevos.

## Perfil Go (mandatorio)
- Framework HTTP: Gin (+ gin-contrib/cors).
- Arquitectura tipo Clean/Hexagonal: domain + port + usecase; adaptadores en infrastructure.
- Errores: usar AppError (code, message, details, http_status, timestamp) y respuesta JSON uniforme.
- Middleware-first: request id, auth, logging request/response, error handler global.
- Logging: zap estructurado + lumberjack; loggers por concern (app/http/auth).
- Persistencia: PostgreSQL + GORM; modelos DB* en postgres; migraciones con AutoMigrate.
- DTOs acotados en handlers con tags de binding/validator.
- Auth: Bearer JWT; inyectar tenant_id, user_id, user_name en context y gin keys.

## Parcel (monolito modular)
- Un solo deploy/runtime, módulos internos con límites claros.
- No duplicar lógica de servicios externos: IAM, PERSON, LOCATION, VEHICLE, TRIP/SCHEDULE, PAYMENT, TENANT-CONFIG, NOTIFICATION, TICKETING.
- Consumir externos vía clientes (infrastructure/clients) detrás de interfaces en port.
- Feature flags siempre desde TENANT-CONFIG (no hardcode).

## Convenciones
- Archivos: snake_case.go
- Tipos: XUseCase, XHandler, XRepository, DBX
- Dependencias “hacia adentro”: handler -> usecase -> port; postgres implementa port.
"@

Ensure-File ".github/copilot-instructions.md" $copilotInstructions

# --- Path-specific instructions ---
$goInstructions = @"
---
applyTo: "**/*.go"
---
- Mantén Go moderno/pragmático (Gin). Evita sobre-ingeniería.
- Usa AppError para errores y respuesta JSON uniforme.
- No accedas a Postgres directo desde handlers; usa usecase + port.
- Logging zap estructurado; evita fmt.Println.
- Archivos snake_case.go; paquetes en minúscula.
"@
Ensure-File ".github/instructions/go_profile.instructions.md" $goInstructions

$parcelInstructions = @"
---
applyTo: "internal/parcel/**"
---
- Respeta límites de dominio por módulo (parcel_core, parcel_item, etc.).
- Prohibido: llamar HTTP externo desde domain/usecase sin pasar por port.
- Prohibido: mezclar concerns (documents dentro de core, incident dentro de core).
- Estados del parcel en core; tracking es timeline/historial (no “dueño” del estado).
- Feature flags: siempre por TENANT-CONFIG (cliente externo).
"@
Ensure-File ".github/instructions/parcel_boundaries.instructions.md" $parcelInstructions

# --- Prompt file (ejecutable con / en Copilot Chat) ---
$promptCreateParcel = @"
---
name: parcel_core_create
description: Crear el primer vertical slice: POST /parcels (CREATED) respetando perfil y capas.
argument-hint: "Describe campos mínimos del envío y reglas especiales (si aplica)"
---
Objetivo: implementar un flujo end-to-end mínimo para Paquetería.

**Reglas absolutas**
- NO reestructures el repo ni cambies estilos existentes.
- Usa Gin, middleware, AppError, zap, gorm según .github/copilot-instructions.md.
- Toca SOLO archivos nuevos o los que se indiquen explícitamente en este prompt.

Implementar:
1) Endpoint: POST /api/v1/parcels
2) Handler + DTO request/response (binding tags).
3) Usecase: CreateParcel (estado inicial CREATED).
4) Domain: entidades/enums mínimos en parcel_core/domain (Parcel, ParcelStatus).
5) Port: ParcelRepository (Create) + TenantConfigClient (para flags mínimos, stub).
6) Infrastructure:
   - Router: registrar ruta en internal/infrastructure/http/router.
   - Repo stub in-memory (temporal) o Postgres placeholder, pero sin migraciones aún.
   - Error handling: devolver AppError con http_status adecuado.

Salida esperada:
- Código compilable.
- Respuestas JSON uniformes (success/error).
- Comentarios TODO donde falten integraciones externas.

Input del usuario (argumento):
{{args}}
"@
Ensure-File ".github/prompts/parcel_core_create.prompt.md" $promptCreateParcel

Write-Host "OK. Scaffold + reglas Copilot + prompt file creados (sin sobreescribir existentes)." -ForegroundColor Green
Write-Host "Siguiente: en VS Code, Copilot Chat -> /parcel_core_create <descripcion>" -ForegroundColor Cyan
