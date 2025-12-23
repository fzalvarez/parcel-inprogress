---
applyTo: "**/*.go"
---
- Mantén Go moderno/pragmático (Gin). Evita sobre-ingeniería.
- Usa AppError para errores y respuesta JSON uniforme.
- No accedas a Postgres directo desde handlers; usa usecase + port.
- Logging zap estructurado; evita fmt.Println.
- Archivos snake_case.go; paquetes en minúscula.
