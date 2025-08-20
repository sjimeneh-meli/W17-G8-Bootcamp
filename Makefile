# Makefile para generar coverage de todo el proyecto
# Coverage completo para todos los paquetes y archivos del proyecto

# Variables
COVERAGE_DIR := coverage
COVERAGE_FILE := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html

# Paquetes del proyecto para coverage
PACKAGES := \
	./internal/handlers/... \
	./internal/services/... \
	./internal/repositories/... \
	./internal/mappers/... \
	./internal/validations/... \
	./internal/models/... \
	./internal/application/... \
	./internal/config/... \
	./internal/container/... \
	./internal/error_message/... \
	./internal/routes/... \
	./pkg/...

.PHONY: coverage coverage-html coverage-report clean-coverage help test-all

# Target principal para generar coverage completo
coverage: clean-coverage
	@echo "🧪 Generando coverage para todo el proyecto..."
	@mkdir -p $(COVERAGE_DIR)
	
	@echo "📁 Ejecutando tests de handlers..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/handlers.out -coverpkg="./internal/handlers" ./internal/handlers/... || true
	@echo "📁 Ejecutando tests de services..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/services.out -coverpkg="./internal/services" ./internal/services/... || true
	@echo "📁 Ejecutando tests de repositories..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/repositories.out -coverpkg="./internal/repositories" ./internal/repositories/... || true
	@echo "📁 Ejecutando tests de mappers..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/mappers.out -coverpkg="./internal/mappers" ./internal/mappers/... || true
	@echo "📁 Ejecutando tests de validations..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/validations.out -coverpkg="./internal/validations" ./internal/validations/... || true
	@echo "📁 Ejecutando tests de models..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/models.out -coverpkg="./internal/models" ./internal/models/... || true
	@echo "📁 Ejecutando tests de application..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/application.out -coverpkg="./internal/application" ./internal/application/... || true
	@echo "📁 Ejecutando tests de config..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/config.out -coverpkg="./internal/config" ./internal/config/... || true
	@echo "📁 Ejecutando tests de container..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/container.out -coverpkg="./internal/container" ./internal/container/... || true
	@echo "📁 Ejecutando tests de error_message..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/error_message.out -coverpkg="./internal/error_message" ./internal/error_message/... || true
	@echo "📁 Ejecutando tests de routes..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/routes.out -coverpkg="./internal/routes" ./internal/routes/... || true
	@echo "📁 Ejecutando tests de pkg..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/pkg.out -coverpkg="./pkg" ./pkg/... || true
	
	@echo "🔗 Combinando archivos de coverage..."
	@echo "mode: set" > $(COVERAGE_FILE)
	@for file in $(COVERAGE_DIR)/*.out; do \
		if [ -f "$$file" ] && [ "$$file" != "$(COVERAGE_FILE)" ]; then \
			tail -n +2 "$$file" >> $(COVERAGE_FILE) 2>/dev/null || true; \
		fi; \
	done
	
	@echo "✅ Coverage generado en $(COVERAGE_FILE)"
	@if [ -f $(COVERAGE_FILE) ] && [ -s $(COVERAGE_FILE) ]; then \
		go tool cover -func=$(COVERAGE_FILE); \
	else \
		echo "⚠️  No se generó coverage válido"; \
	fi

# Target para ejecutar todos los tests sin coverage
test-all:
	@echo "🧪 Ejecutando todos los tests del proyecto..."
	@go test -v $(PACKAGES)

# Generar reporte HTML
coverage-html: coverage
	@echo "🌐 Generando reporte HTML..."
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "✅ Reporte HTML generado en $(COVERAGE_HTML)"
	@echo "🔗 Abrir en navegador: open $(COVERAGE_HTML)"

# Mostrar reporte detallado
coverage-report: coverage
	@echo "\n📊 REPORTE DETALLADO DE COVERAGE"
	@echo "=================================="
	@echo "\n🎯 PAQUETES ANALIZADOS:"
	@echo "  - internal/handlers"
	@echo "  - internal/services"
	@echo "  - internal/repositories"
	@echo "  - internal/mappers"
	@echo "  - internal/validations"
	@echo "  - internal/models"
	@echo "  - internal/application"
	@echo "  - internal/config"
	@echo "  - internal/container"
	@echo "  - internal/error_message"
	@echo "  - internal/routes"
	@echo "  - pkg"
	@echo "\n📈 PORCENTAJE DE COVERAGE POR ARCHIVO:"
	@go tool cover -func=$(COVERAGE_FILE) | grep -E "\.(go):" | grep -v "_test\.go" | sort
	@echo "\n📊 COVERAGE TOTAL:"
	@go tool cover -func=$(COVERAGE_FILE) | grep "total:"



# Limpiar archivos de coverage
clean-coverage:
	@echo "🧹 Limpiando archivos de coverage anteriores..."
	@mkdir -p $(COVERAGE_DIR)
	@rm -f $(COVERAGE_DIR)/*.out $(COVERAGE_DIR)/coverage.html
	@echo "ℹ️  README.md preservado en $(COVERAGE_DIR)/"

# Mostrar ayuda
help:
	@echo "📚 COMANDOS DISPONIBLES:"
	@echo "========================"
	@echo "make coverage        - Generar coverage para todo el proyecto"
	@echo "make coverage-html   - Generar coverage + reporte HTML"
	@echo "make coverage-report - Mostrar reporte detallado de coverage"
	@echo "make test-all        - Ejecutar todos los tests sin coverage"
	@echo "make clean-coverage  - Limpiar archivos de coverage"
	@echo ""

	@echo "make help           - Mostrar esta ayuda"
	@echo ""
	@echo "📦 INCLUYE: Todos los paquetes del proyecto interno y pkg"
	@echo "ℹ️  NOTA: Se analizan todos los archivos Go del proyecto"
	@echo "   Archivos sin tests mostrarán 0% de cobertura"

# Target por defecto
.DEFAULT_GOAL := coverage 