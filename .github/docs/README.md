# 📚 Documentación QuatroBus Parcel

Bienvenido a la documentación técnica del proyecto **QuatroBus Parcel**.

---

## 🚀 Inicio Rápido

### Para Ver la API (Swagger)
```bash
# 1. Iniciar el servidor
go run cmd/api/main.go

# 2. Abrir Swagger UI en navegador
http://localhost:8080/swagger/index.html
```

### Para Regenerar Swagger
```bash
swag init -g cmd/api/main.go
```

---

## 📖 Documentos Disponibles

### 🎯 Documentos Principales

| Documento | Para Quién | Descripción |
|-----------|-----------|-------------|
| [**FINAL_STATUS.md**](./FINAL_STATUS.md) | 👨‍💼 Todos | ✅ Estado final del proyecto completado |
| [**INDEX.md**](./INDEX.md) | 👨‍💼 Todos | 📋 Índice maestro de toda la documentación |
| [**architecture_diagram.md**](./architecture_diagram.md) | 🏗️ Arquitectos | 📐 Diagrama visual completo del sistema |

### 🔌 API y Endpoints

| Documento | Para Quién | Descripción |
|-----------|-----------|-------------|
| [**swagger_endpoints_reference.md**](./swagger_endpoints_reference.md) | 💻 Desarrolladores | 📚 Referencia completa de 22 endpoints |
| [**swagger_maintenance_guide.md**](./swagger_maintenance_guide.md) | 💻 Desarrolladores | 🔧 Guía de mantenimiento de Swagger |
| [**SWAGGER_UPDATE_SUMMARY.md**](./SWAGGER_UPDATE_SUMMARY.md) | 📝 Revisores | 📋 Resumen de actualización Swagger |

### 🏛️ Arquitectura

| Documento | Para Quién | Descripción |
|-----------|-----------|-------------|
| [**persistence_architecture.md**](./persistence_architecture.md) | 🏗️ Arquitectos | 💾 Arquitectura de persistencia (PostgreSQL + memoria) |
| [**WORK_COMPLETED_SUMMARY.md**](./WORK_COMPLETED_SUMMARY.md) | 📊 Project Managers | 📈 Resumen completo de trabajo realizado |

### 💰 Módulos de Negocio

| Documento | Para Quién | Descripción |
|-----------|-----------|-------------|
| [**pricing_rules_guide.md**](./pricing_rules_guide.md) | 💻 Desarrolladores<br>💼 Product Managers | 💵 Guía completa del motor de precios |

---

## 🎓 Rutas de Aprendizaje

### 👨‍💻 Soy Nuevo en el Proyecto

```
1. 📋 Leer INDEX.md (vista general)
   ↓
2. 📐 Ver architecture_diagram.md (entender el flujo)
   ↓
3. 🔌 Explorar Swagger UI (probar la API)
   ↓
4. 💵 Revisar pricing_rules_guide.md (motor de precios)
```

### 🏗️ Soy Arquitecto/Tech Lead

```
1. ✅ Leer FINAL_STATUS.md (estado del proyecto)
   ↓
2. 📐 Revisar architecture_diagram.md (diseño completo)
   ↓
3. 💾 Analizar persistence_architecture.md (decisiones de BD)
   ↓
4. 📊 Ver WORK_COMPLETED_SUMMARY.md (trabajo realizado)
```

### 💼 Soy Product Manager

```
1. ✅ Leer FINAL_STATUS.md (qué está listo)
   ↓
2. 📚 Explorar swagger_endpoints_reference.md (capacidades)
   ↓
3. 💵 Entender pricing_rules_guide.md (sistema de precios)
   ↓
4. 🔌 Probar Swagger UI (endpoints en vivo)
```

### 🔧 Voy a Mantener el Código

```
1. 📋 Leer INDEX.md (navegación general)
   ↓
2. 🔧 Revisar swagger_maintenance_guide.md (mantener Swagger)
   ↓
3. 📐 Entender architecture_diagram.md (capas y flujos)
   ↓
4. 💻 Seguir go_profile.instructions.md (convenciones)
```

---

## 📊 Estado del Proyecto

```
✅ Swagger:           100% documentado (22 endpoints)
✅ Arquitectura:      Reorganizada y documentada
✅ Motor de Pricing:  Implementado con jerarquía
✅ PostgreSQL:        Modelos creados, listo para usar
✅ Documentación:     8 guías técnicas completas
```

---

## 🔗 Links Rápidos

### Documentación de Código
- [Perfil Go](../instructions/go_profile.instructions.md)
- [Límites de Parcel](../instructions/parcel_boundaries.instructions.md)

### Herramientas
- Swagger UI: `http://localhost:8080/swagger/index.html`
- Generar Swagger: `swag init -g cmd/api/main.go`

---

## 📞 Ayuda

¿No encuentras lo que buscas?

1. **Empieza por:** [INDEX.md](./INDEX.md)
2. **Busca en:** La tabla de contenidos de cada documento
3. **Pregunta al equipo:** Desarrollo QuatroBus

---

## 🗂️ Estructura de Esta Carpeta

```
.github/docs/
├── README.md                      # 👈 Este archivo
├── INDEX.md                       # Índice maestro
├── FINAL_STATUS.md                # Estado final completado
├── WORK_COMPLETED_SUMMARY.md      # Resumen de trabajo
├── architecture_diagram.md        # Diagrama de arquitectura
├── persistence_architecture.md    # Arquitectura de persistencia
├── pricing_rules_guide.md         # Guía de pricing
├── swagger_endpoints_reference.md # Referencia de endpoints
├── swagger_maintenance_guide.md   # Mantenimiento de Swagger
└── SWAGGER_UPDATE_SUMMARY.md      # Resumen de actualización
```

---

**Última Actualización:** 21 de enero de 2026  
**Mantenido por:** Equipo de Desarrollo QuatroBus
