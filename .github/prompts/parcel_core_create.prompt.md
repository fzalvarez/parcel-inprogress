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
