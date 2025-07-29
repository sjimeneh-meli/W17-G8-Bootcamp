# Makefile para generar coverage de archivos testeados
# Coverage para handlers, services y repositories que tienen tests

# Variables
COVERAGE_DIR := coverage
COVERAGE_FILE := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html

# Paquetes específicos que tienen tests
# Solo incluir archivos que tienen tests correspondientes
TESTED_HANDLERS := \
	./internal/handlers/buyerHandler.go \
	./internal/handlers/employeeHandler.go \
	./internal/handlers/productHandler.go \
	./internal/handlers/sectionHandler.go \
	./internal/handlers/sellerHandler.go \
	./internal/handlers/warehouse_handler.go

TESTED_SERVICES := \
	./internal/services/buyerService.go \
	./internal/services/employeeService.go \
	./internal/services/warehouse_service.go

TESTED_REPOSITORIES := \
	./internal/repositories/buyerRepository.go \
	./internal/repositories/employeeRepository.go \
	./internal/repositories/productRepository.go \
	./internal/repositories/sectionRepository.go \
	./internal/repositories/sellerRepository.go \
	./internal/repositories/warehouse_repository.go

# Variables de archivos testeados definidas arriba - no necesitamos COVER_PACKAGES

.PHONY: coverage coverage-html coverage-report clean-coverage help

# Target principal para generar coverage
coverage: clean-coverage
	@echo "🧪 Generando coverage para archivos testeados..."
	@mkdir -p $(COVERAGE_DIR)
	
	# Ejecutar todos los tests y luego filtrar el coverage
	@echo "📁 Ejecutando tests de handlers..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/handlers.out -coverpkg="./internal/handlers" ./internal/handlers/... || true
	@echo "📁 Ejecutando tests de services..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/services.out -coverpkg="./internal/services" ./internal/services/... || true
	@echo "📁 Ejecutando tests de repositories..."
	@go test -v -coverprofile=$(COVERAGE_DIR)/repositories.out -coverpkg="./internal/repositories" ./internal/repositories/... || true
	
	@echo "🔗 Combinando archivos de coverage..."
	@echo "mode: set" > $(COVERAGE_FILE)
	@if [ -f $(COVERAGE_DIR)/handlers.out ]; then tail -n +2 $(COVERAGE_DIR)/handlers.out >> $(COVERAGE_FILE); fi
	@if [ -f $(COVERAGE_DIR)/services.out ]; then tail -n +2 $(COVERAGE_DIR)/services.out >> $(COVERAGE_FILE); fi
	@if [ -f $(COVERAGE_DIR)/repositories.out ]; then tail -n +2 $(COVERAGE_DIR)/repositories.out >> $(COVERAGE_FILE); fi
	
	@echo "🎯 Filtrando solo archivos con tests..."
	# Filtrar coverage para incluir SOLO archivos que tienen tests correspondientes
	@if [ -f $(COVERAGE_FILE) ]; then \
		(echo "mode: set"; \
		 grep -E "(buyerHandler\.go|employeeHandler\.go|productHandler\.go|sectionHandler\.go|sellerHandler\.go|warehouse_handler\.go|buyerService\.go|employeeService\.go|warehouse_service\.go|buyerRepository\.go|employeeRepository\.go|productRepository\.go|sectionRepository\.go|sellerRepository\.go|warehouse_repository\.go)" $(COVERAGE_FILE) || true) > $(COVERAGE_FILE).filtered; \
		if [ -s $(COVERAGE_FILE).filtered ] && [ "$$(wc -l < $(COVERAGE_FILE).filtered)" -gt 1 ]; then \
			mv $(COVERAGE_FILE).filtered $(COVERAGE_FILE); \
		else \
			echo "⚠️  No se encontró coverage válido para archivos testeados"; \
			rm -f $(COVERAGE_FILE).filtered; \
		fi; \
	fi
	
	@echo "✅ Coverage generado en $(COVERAGE_FILE)"
	@if [ -f $(COVERAGE_FILE) ] && [ -s $(COVERAGE_FILE) ]; then \
		go tool cover -func=$(COVERAGE_FILE); \
	else \
		echo "⚠️  No se generó coverage válido"; \
	fi

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
	@echo "\n🎯 ARCHIVOS TESTEADOS:"
	@echo "Handlers:"
	@echo "  - buyerHandler.go"
	@echo "  - employeeHandler.go" 
	@echo "  - productHandler.go"
	@echo "  - sectionHandler.go"
	@echo "  - sellerHandler.go"
	@echo "  - warehouse_handler.go"
	@echo "\nServices:"
	@echo "  - buyerService.go"
	@echo "  - employeeService.go"
	@echo "  - warehouse_service.go"
	@echo "\nRepositories:"
	@echo "  - buyerRepository.go"
	@echo "  - employeeRepository.go"
	@echo "  - productRepository.go"
	@echo "  - sectionRepository.go"
	@echo "  - sellerRepository.go"
	@echo "  - warehouse_repository.go"
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
	@echo "make coverage        - Generar coverage para archivos testeados"
	@echo "make coverage-html   - Generar coverage + reporte HTML"
	@echo "make coverage-report - Mostrar reporte detallado de coverage"
	@echo "make clean-coverage  - Limpiar archivos de coverage"
	@echo "make help           - Mostrar esta ayuda"
	@echo ""
	@echo "ℹ️  NOTA: Solo se incluyen archivos que tienen tests correspondientes"
	@echo "   Los archivos sin tests son automáticamente ignorados"

# Target por defecto
.DEFAULT_GOAL := coverage 