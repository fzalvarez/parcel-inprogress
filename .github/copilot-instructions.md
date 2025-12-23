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
