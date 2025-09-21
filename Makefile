# Go-ney Framework Makefile

# Variables
BINARY_NAME=goney
BUILD_DIR=bin
INSTALL_PATH=/usr/local/bin

# Colores para output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: build install uninstall clean test dev help

# Construir el binario
build:
	@echo "$(YELLOW)🔨 Construyendo Go-ney...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/
	@echo "$(GREEN)✅ Go-ney construido en $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Instalar globalmente
install: build
	@echo "$(YELLOW)📦 Instalando Go-ney globalmente...$(NC)"
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/
	@sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "$(GREEN)✅ Go-ney instalado en $(INSTALL_PATH)/$(BINARY_NAME)$(NC)"
	@echo "$(GREEN)🚀 Ahora puedes usar 'goney' desde cualquier directorio$(NC)"

# Instalar para el usuario actual (sin sudo)
install-user: build
	@echo "$(YELLOW)📦 Instalando Go-ney para el usuario actual...$(NC)"
	@mkdir -p ~/bin
	@cp $(BUILD_DIR)/$(BINARY_NAME) ~/bin/
	@chmod +x ~/bin/$(BINARY_NAME)
	@echo "$(GREEN)✅ Go-ney instalado en ~/bin/$(BINARY_NAME)$(NC)"
	@echo "$(YELLOW)⚠️  Asegúrate de que ~/bin esté en tu PATH$(NC)"
	@echo "$(YELLOW)   Agrega esto a tu ~/.bashrc o ~/.zshrc:$(NC)"
	@echo "$(YELLOW)   export PATH=\"$$HOME/bin:$$PATH\"$(NC)"

# Desinstalar
uninstall:
	@echo "$(YELLOW)🗑️  Desinstalando Go-ney...$(NC)"
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@rm -f ~/bin/$(BINARY_NAME)
	@echo "$(GREEN)✅ Go-ney desinstalado$(NC)"

# Limpiar archivos generados
clean:
	@echo "$(YELLOW)🧹 Limpiando archivos generados...$(NC)"
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "$(GREEN)✅ Limpieza completada$(NC)"

# Ejecutar tests
test:
	@echo "$(YELLOW)🧪 Ejecutando tests...$(NC)"
	@go test ./...
	@echo "$(GREEN)✅ Tests completados$(NC)"

# Modo desarrollo (instalar dependencias y construir)
dev:
	@echo "$(YELLOW)🔧 Configurando entorno de desarrollo...$(NC)"
	@go mod tidy
	@go mod download
	@$(MAKE) build
	@echo "$(GREEN)✅ Entorno de desarrollo listo$(NC)"

# Actualizar dependencias
update:
	@echo "$(YELLOW)📦 Actualizando dependencias...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)✅ Dependencias actualizadas$(NC)"

# Verificar instalación
check:
	@echo "$(YELLOW)🔍 Verificando instalación...$(NC)"
	@which $(BINARY_NAME) || echo "$(RED)❌ Go-ney no está instalado globalmente$(NC)"
	@$(BINARY_NAME) --version 2>/dev/null || echo "$(RED)❌ Error ejecutando Go-ney$(NC)"
	@echo "$(GREEN)✅ Verificación completada$(NC)"

# Release (construir para múltiples plataformas)
release: clean
	@echo "$(YELLOW)🚀 Creando release para múltiples plataformas...$(NC)"
	@mkdir -p $(BUILD_DIR)/release

	# Linux AMD64
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-amd64 ./cmd/

	# Linux ARM64
	@GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-arm64 ./cmd/

	# macOS AMD64
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-amd64 ./cmd/

	# macOS ARM64 (Apple Silicon)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-arm64 ./cmd/

	# Windows AMD64
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-windows-amd64.exe ./cmd/

	@echo "$(GREEN)✅ Release creado en $(BUILD_DIR)/release/$(NC)"

# Ayuda
help:
	@echo "$(GREEN)🐝 Go-ney Framework - Comandos disponibles:$(NC)"
	@echo ""
	@echo "$(YELLOW)Construcción e instalación:$(NC)"
	@echo "  build        - Construir el binario"
	@echo "  install      - Instalar globalmente (requiere sudo)"
	@echo "  install-user - Instalar para el usuario actual"
	@echo "  uninstall    - Desinstalar Go-ney"
	@echo ""
	@echo "$(YELLOW)Desarrollo:$(NC)"
	@echo "  dev          - Configurar entorno de desarrollo"
	@echo "  test         - Ejecutar tests"
	@echo "  clean        - Limpiar archivos generados"
	@echo "  update       - Actualizar dependencias"
	@echo ""
	@echo "$(YELLOW)Utilidades:$(NC)"
	@echo "  check        - Verificar instalación"
	@echo "  release      - Crear release multiplataforma"
	@echo "  help         - Mostrar esta ayuda"
	@echo ""
	@echo "$(GREEN)Ejemplo de uso:$(NC)"
	@echo "  make install     # Instalar Go-ney"
	@echo "  goney new mi-app # Crear nuevo proyecto"
	@echo "  cd mi-app        # Entrar al proyecto"
	@echo "  goney generate crud Usuario  # Generar CRUD completo"
	@echo "  goney start     # Iniciar servidor"