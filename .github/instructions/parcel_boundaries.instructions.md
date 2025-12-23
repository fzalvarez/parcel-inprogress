---
applyTo: "internal/parcel/**"
---
- Respeta límites de dominio por módulo (parcel_core, parcel_item, etc.).
- Prohibido: llamar HTTP externo desde domain/usecase sin pasar por port.
- Prohibido: mezclar concerns (documents dentro de core, incident dentro de core).
- Estados del parcel en core; tracking es timeline/historial (no “dueño” del estado).
- Feature flags: siempre por TENANT-CONFIG (cliente externo).
