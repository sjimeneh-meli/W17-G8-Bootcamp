PKG=./internal
EMPLOYEE_FILES=handlers/employeeHandler.go services/employeeService.go repositories/employeeRepository.go validations/employeeValidatioin.go mappers/employeeMapper.go
EMPLOYEE_COVERAGE_FILE=employee-coverage-out.out

# Colores ANSI
YELLOW=\033[1;33m
GREEN=\033[1;32m
BLUE=\033[1;34m
RESET=\033[0m

.PHONY: coverage-employee coverage-employee-html tests 

coverage-employee:
	@go test $(PKG)/... -coverprofile=$(EMPLOYEE_COVERAGE_FILE) -covermode=atomic > /dev/null 2>&1

	@echo ""
	@echo -e "$(BLUE)==== COVERAGE POR ARCHIVO (EMPLOYEE) ====$(RESET)"
	@for file in $(EMPLOYEE_FILES); do \
		echo ""; \
		echo -e "$(GREEN)Archivo: $$file$(RESET)"; \
		go tool cover -func=$(EMPLOYEE_COVERAGE_FILE) | grep "$$file" || echo -e "$(YELLOW)No hay coverage para $$file$(RESET)"; \
	done


tests:
	@echo -e "$(YELLOW)Ejecutando todos los tests del proyecto...$(RESET)"
	go test ./...


coverage-html:
	@echo -e "$(YELLOW)Generando coverage HTML para todo el proyecto...$(RESET)"
	@go test $(PKG)/... -coverprofile=coverage-out.out -covermode=atomic > /dev/null 2>&1
	@go tool cover -html=coverage-out.out -o coverage-out.html
	@echo -e "$(GREEN)Coverage HTML generado: coverage-out.html$(RESET)" 