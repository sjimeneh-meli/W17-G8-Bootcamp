# 📊 Coverage Report - Guía de Uso

Esta carpeta contiene los reportes de coverage generados automáticamente por el Makefile del proyecto.

## 🎯 ¿Qué es esto?

El sistema de coverage está configurado para analizar **únicamente** los archivos que tienen tests correspondientes, ignorando automáticamente aquellos archivos sin tests en las capas de:
- **Handlers** 
- **Services**
- **Repositories**

## 🚀 Comandos Disponibles

### 📋 Ver ayuda
```bash
make help
```
Muestra todos los comandos disponibles con descripción.

### 🧪 Generar coverage básico
```bash
make coverage
```
- Ejecuta todos los tests
- Genera archivo de coverage filtrado
- Muestra resumen en consola
- **Resultado**: `coverage/coverage.out`

### 🌐 Generar reporte HTML
```bash
make coverage-html
```
- Ejecuta `make coverage`
- Genera reporte visual en HTML
- **Resultado**: `coverage/coverage.html`
- **Abre el archivo HTML en tu navegador** para ver el reporte interactivo

### 📊 Reporte detallado
```bash
make coverage-report
```
- Ejecuta `make coverage`
- Muestra lista de archivos analizados
- Muestra porcentaje por archivo
- Muestra coverage total del proyecto

### 🧹 Limpiar archivos
```bash
make clean-coverage
```
Elimina archivos de coverage generados **pero preserva este README**.

## 📁 Archivos Generados

### 🔍 `coverage.out`
Archivo de coverage en formato Go. Contiene la información de cobertura que puede ser procesada por herramientas de Go.

### 🌐 `coverage.html`
Reporte visual interactivo que muestra:
- **Archivos analizados** (solo los que tienen tests)
- **Líneas cubiertas** en verde
- **Líneas no cubiertas** en rojo
- **Porcentaje por archivo**
- **Coverage total**

### 📋 Archivos temporales (`.out`)
Archivos intermedios generados durante el proceso:
- `handlers.out`
- `services.out`
- `repositories.out`

## ✅ Archivos Incluidos en el Coverage

### 🎮 Handlers (6 archivos)
- ✅ `buyerHandler.go`
- ✅ `employeeHandler.go`
- ✅ `productHandler.go`
- ✅ `sectionHandler.go`
- ✅ `sellerHandler.go`
- ✅ `warehouse_handler.go`

### ⚙️ Services (3 archivos)
- ✅ `buyerService.go`
- ✅ `employeeService.go`
- ✅ `warehouse_service.go`

### 🗄️ Repositories (6 archivos)
- ✅ `buyerRepository.go`
- ✅ `employeeRepository.go`
- ✅ `productRepository.go`
- ✅ `sectionRepository.go`
- ✅ `sellerRepository.go`
- ✅ `warehouse_repository.go`

## ❌ Archivos Automáticamente Excluidos

Los siguientes archivos **NO aparecen** en el reporte porque no tienen tests:
- ❌ `carry_handler.go`
- ❌ `inboundOrderHandler.go`
- ❌ `localityHandler.go`
- ❌ `productBatchHandler.go`
- ❌ `productRecordHandler.go`
- ❌ `purchaseOrderHandler.go`
- ❌ `localityService.go`
- ❌ `productService.go`
- ❌ `sectionService.go`
- ❌ `carry_services.go`
- ❌ Y otros archivos sin tests correspondientes

## 🔄 Flujo de Trabajo Recomendado

### 1. 📊 Generar reporte inicial
```bash
make coverage-html
```

### 2. 🌐 Ver reporte en navegador
```bash
open coverage/coverage.html
# o en Linux/Windows: abre manualmente el archivo
```

### 3. 📈 Analizar resultados
- Identifica archivos con bajo coverage
- Revisa líneas no cubiertas
- Planifica mejoras en tests

### 4. 🔄 Regenerar después de cambios
```bash
make clean-coverage
make coverage-html
```

## 💡 Tips y Buenas Prácticas

### 🎯 Para mejorar coverage:
1. **Revisa líneas rojas** en el HTML
2. **Agrega tests** para funciones no cubiertas
3. **Considera casos edge** en tus tests
4. **Testa error handlers**

### 🚀 Para CI/CD:
```bash
# Script básico para pipeline
make coverage
if [ $? -eq 0 ]; then
  echo "✅ Coverage generado exitosamente"
else
  echo "❌ Error generando coverage"
  exit 1
fi
```

### 📊 Interpretar resultados:
- **Verde**: Línea ejecutada por tests
- **Rojo**: Línea NO ejecutada por tests
- **Sin color**: Línea no relevante (comentarios, imports, etc.)

## 🆘 Troubleshooting

### ❓ No se genera coverage
```bash
# Verificar que go y make estén instalados
go version
make --version

# Limpiar y regenerar
make clean-coverage
make coverage
```

### ❓ Archivo HTML vacío
```bash
# Verificar que existan tests
ls internal/handlers/*_test.go
ls internal/services/*_test.go
ls internal/repositories/*_test.go

# Ejecutar tests manualmente
go test ./internal/handlers/...
go test ./internal/services/...
go test ./internal/repositories/...
```

### ❓ Coverage muy bajo
- **Normal**: Es esperado si hay muchas funciones nuevas sin tests
- **Solución**: Agregar tests para las funciones identificadas

## 📞 Comandos de Referencia Rápida

```bash
make help              # Ver ayuda
make coverage          # Coverage básico
make coverage-html     # Reporte HTML
make coverage-report   # Reporte detallado
make clean-coverage    # Limpiar archivos (preserva README)
```

---

**📝 Nota**: Este sistema está configurado para mostrar únicamente archivos con tests existentes, proporcionando una visión realista del coverage actual del proyecto.

**🔒 Importante**: El comando `make clean-coverage` ahora preserva este archivo README.md automáticamente. 