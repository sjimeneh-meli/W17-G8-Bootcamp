PKG=./internal
BUYER_FILES=handlers/buyerHandler.go services/buyerService.go repositories/buyerRepository.go validations/buyerValidations.go mappers/buyerMapper.go
BUYERS_COVERAGE_FILE=buyers-coverage-out.out
COVERAGE_HTML=coverage-out.html

# Colores ANSI
YELLOW=\033[1;33m
GREEN=\033[1;32m
BLUE=\033[1;34m
RESET=\033[0m

.PHONY: coverage-buyer coverage-buyer-html tests clean

coverage-buyers: clean
	@go test $(PKG)/... -coverprofile=$(BUYERS_COVERAGE_FILE) -covermode=atomic > /dev/null 2>&1

	@echo ""
	@echo -e "$(BLUE)==== COVERAGE POR ARCHIVO (BUYERS) ====$(RESET)"
	@for file in $(BUYER_FILES); do \
		echo ""; \
		echo -e "$(GREEN)Archivo: $$file$(RESET)"; \
		go tool cover -func=$(BUYERS_COVERAGE_FILE) | grep "$$file" || echo -e "$(YELLOW)No hay coverage para $$file$(RESET)"; \
	done


tests:
	@echo -e "$(YELLOW)Ejecutando todos los tests del proyecto...$(RESET)"
	go test ./...

clean:
	rm -f $(BUYERS_COVERAGE_FILE) $(COVERAGE_HTML)


coverage-html: clean
	@echo -e "$(YELLOW)Generando coverage HTML para todo el proyecto...$(RESET)"
	@go test $(PKG)/... -coverprofile=coverage-out.out -covermode=atomic > /dev/null 2>&1
	@go tool cover -html=coverage-out.out -o coverage-out.html
	@echo -e "$(GREEN)Coverage HTML generado: coverage-out.html$(RESET)"