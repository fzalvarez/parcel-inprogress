# Guía del Motor de Precios (Pricing Engine)

## Resumen
El sistema de precios de QuatroBus Parcel soporta **cálculo automático** basado en reglas jerárquicas que combinan peso volumétrico, zonas geográficas y precios manuales.

---

## 1. Componentes del Motor

### A. Peso Facturable (Billable Weight)
El sistema calcula automáticamente el "peso facturable" usando la fórmula de la industria logística:

```
VolumetricWeight = (Largo_cm × Ancho_cm × Alto_cm) / Divisor
BillableWeight = Max(PesoReal_kg, VolumetricWeight)
```

**Configuración (Tenant Options):**
- `UseVolumetricWeight`: Activa/desactiva el cálculo volumétrico.
- `VolumetricDivisor`: Factor de conversión (default: 6000, estándar IATA).

**Ejemplo:**
- Paquete: 50×40×30 cm, 5 kg real.
- Volumétrico: (50×40×30)/6000 = 10 kg.
- **Facturable: 10 kg** (se cobra por el volumétrico).

---

## 2. Reglas de Precios (Price Rules)

### A. Estructura de una Regla
```json
{
  "shipment_type": "STANDARD",
  "origin_office_id": "uuid-oficina-lima",
  "destination_office_id": "uuid-oficina-cusco",
  "unit": "PER_KG",
  "price": 2.50,
  "currency": "PEN",
  "priority": 10,
  "active": true
}
```

### B. Comodines (Wildcards)
Usa `"*"` para crear reglas que apliquen a múltiples rutas:

| Origin | Destination | Descripción |
|--------|-------------|-------------|
| `uuid-lima` | `uuid-cusco` | Precio específico Lima → Cusco |
| `uuid-lima` | `*` | Precio desde Lima a cualquier destino |
| `*` | `uuid-cusco` | Precio desde cualquier origen a Cusco |
| `*` | `*` | **Precio global (fallback)** |

### C. Sistema de Prioridad
Si existen múltiples reglas coincidentes, el sistema elige usando:
1. **Especificidad** (mayor puntaje = más específica):
   - Origen exacto + Destino exacto = 20 puntos
   - Origen exacto + Destino comodín = 11 puntos
   - Origen comodín + Destino exacto = 11 puntos
   - Comodín + Comodín = 2 puntos
2. **Campo Priority** (en caso de empate, gana el mayor valor).

---

## 3. Flujo de Cálculo al Agregar Items

```
POST /api/v1/parcels/:id/items
{
  "description": "Laptop",
  "quantity": 1,
  "weight_kg": 2.5,
  "length_cm": 40,
  "width_cm": 30,
  "height_cm": 10,
  "unit_price": 0  // Opcional, 0 = calcular automático
}
```

**Proceso:**
1. Calcula `BillableWeight` (si `UseVolumetricWeight=true`).
2. Busca regla de precio jerárquicamente.
3. Calcula precio sugerido:
   - `PER_ITEM`: `Price × Quantity`
   - `PER_KG`: `Price × BillableWeight`
4. Aplica precio final:
   - Si enviaste `unit_price > 0` y `AllowOverridePriceTable=true`: usa el manual.
   - Si no: usa el sugerido.

---

## 4. Casos de Uso Prácticos

### Caso 1: Tenant con rutas limitadas (1-5 oficinas)
**Solución:** Crear reglas específicas.
```bash
# Lima → Arequipa
POST /pricing/rules { origin: "uuid-lima", destination: "uuid-arequipa", price: 3.0 }
# Lima → Cusco
POST /pricing/rules { origin: "uuid-lima", destination: "uuid-cusco", price: 4.5 }
```

### Caso 2: Tenant con muchas oficinas (red nacional)
**Solución:** Usar regla global + excepciones.
```bash
# Regla base nacional
POST /pricing/rules { origin: "*", destination: "*", price: 2.0, priority: 1 }
# Excepción para rutas largas
POST /pricing/rules { origin: "uuid-lima", destination: "uuid-iquitos", price: 8.0, priority: 10 }
```

### Caso 3: Permitir precio manual en casos especiales
**Configuración:** `AllowManualPrice=true` + `AllowOverridePriceTable=true`.
El operador puede enviar un precio personalizado que sobrescribe el calculado.

---

## 5. Mensajes de Error

### `price_rule_not_found`
**Causa:** No existe regla para la ruta Origen→Destino.
**Solución:**
1. Crear regla específica: `POST /pricing/rules`.
2. O crear regla global: `{ origin: "*", destination: "*", price: X }`.

**Hint en el error:**
```json
{
  "code": "price_rule_not_found",
  "message": "regla de precios no encontrada para esta ruta. Defina una regla específica o use comodín (*)",
  "details": {
    "hint": "Puede crear una regla global usando '*' como origin_office_id o destination_office_id"
  }
}
```

---

## 6. Mejores Prácticas

1. **Siempre define una regla global (`*→*`)** antes de activar `UsePriceTable=true`.
2. Usa `Priority` para ordenar reglas cuando haya múltiples comodines.
3. Activa `UseVolumetricWeight` si transportas bultos ligeros pero voluminosos.
4. Configura `VolumetricDivisor` según tu tipo de transporte:
   - Aéreo: 6000 (estándar IATA)
   - Terrestre: 3000-5000 (más tolerante)

---

## 7. Roadmap (Próximas Mejoras)

- [ ] Soporte de **zonas geográficas** (agrupar oficinas por región).
- [ ] **Multiplicadores de distancia** (precio base × factor de km).
- [ ] **Registro masivo** de reglas (bulk create para evitar N requests).
- [ ] Integración con servicio **ms-routes** para cálculo de distancia real.
