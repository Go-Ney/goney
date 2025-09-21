#!/bin/bash

# Go-ney Framework - Script de Instalación
# Instala Go-ney Framework en tu sistema

set -e

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="goney"
REPO_URL="https://github.com/Go-Ney/goney"

# Función para imprimir mensajes con colores
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar si Go está instalado
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go no está instalado. Por favor instala Go desde https://golang.org/dl/"
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go $GO_VERSION detectado"
}

# Verificar si Git está instalado
check_git() {
    if ! command -v git &> /dev/null; then
        print_error "Git no está instalado. Por favor instala Git"
        exit 1
    fi
    print_success "Git detectado"
}

# Función para instalar desde código fuente
install_from_source() {
    print_status "Instalando Go-ney desde código fuente..."

    # Crear directorio temporal
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    # Si estamos en el directorio del proyecto, usar el código local
    if [ -f "./go.mod" ] && grep -q "goney" "./go.mod"; then
        print_status "Usando código fuente local..."
        cp -r "$(dirname "$0")" ./goney
        cd goney
    else
        print_status "Clonando repositorio..."
        git clone "$REPO_URL" goney
        cd goney
    fi

    # Construir e instalar
    print_status "Construyendo Go-ney..."
    go mod tidy
    go build -o "$BINARY_NAME" ./cmd/

    # Mover al directorio de instalación
    print_status "Instalando en $INSTALL_DIR..."
    if [ -w "$INSTALL_DIR" ]; then
        mv "$BINARY_NAME" "$INSTALL_DIR/"
    else
        sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi

    # Limpiar
    cd /
    rm -rf "$TEMP_DIR"

    print_success "Go-ney instalado exitosamente en $INSTALL_DIR/$BINARY_NAME"
}

# Función para instalar usando go install
install_with_go_install() {
    print_status "Instalando Go-ney usando 'go install'..."

    if [ -f "./go.mod" ] && grep -q "goney" "./go.mod"; then
        # Estamos en el directorio del proyecto
        go install ./cmd/
    else
        # Instalar desde repositorio remoto
        go install "$REPO_URL/cmd@latest"
    fi

    print_success "Go-ney instalado usando 'go install'"
}

# Verificar instalación
verify_installation() {
    print_status "Verificando instalación..."

    if command -v "$BINARY_NAME" &> /dev/null; then
        VERSION=$($BINARY_NAME --version 2>/dev/null || echo "unknown")
        print_success "Go-ney está instalado correctamente"
        print_status "Versión: $VERSION"
        print_status "Ubicación: $(which $BINARY_NAME)"
    else
        print_error "Go-ney no se pudo instalar correctamente"
        exit 1
    fi
}

# Mostrar ayuda de uso
show_usage() {
    echo -e "${GREEN}"
    echo "   ___           _  __"
    echo "  / _ \___      / |/ /__ __ __"
    echo " / (_ / _ \_   /    / -_) // /"
    echo " \___/\___(_) /_/|_/\__/\_, /"
    echo "                       /___/"
    echo -e "${NC}"
    echo -e "${GREEN}🚀 Go-ney Framework instalado exitosamente!${NC}"
    echo ""
    echo -e "${YELLOW}Comandos disponibles:${NC}"
    echo "  goney new <proyecto>           - Crear nuevo proyecto"
    echo "  goney generate crud <nombre>   - Generar CRUD completo"
    echo "  goney generate controller <nombre> - Generar controller"
    echo "  goney generate service <nombre>    - Generar service"
    echo "  goney start                    - Iniciar servidor de desarrollo"
    echo ""
    echo -e "${YELLOW}Ejemplo de uso:${NC}"
    echo "  goney new mi-api"
    echo "  cd mi-api"
    echo "  goney generate crud usuarios"
    echo "  goney start"
    echo ""
    echo -e "${BLUE}Para más información: goney --help${NC}"
    echo -e "${BLUE}Documentación: https://github.com/Go-Ney/goney${NC}"
}

# Función principal
main() {
    echo -e "${GREEN}"
    echo "   ___           _  __"
    echo "  / _ \___      / |/ /__ __ __"
    echo " / (_ / _ \_   /    / -_) // /"
    echo " \___/\___(_) /_/|_/\__/\_, /"
    echo "                       /___/"
    echo -e "${NC}"
    echo -e "${GREEN}🚀 Instalador de Go-ney Framework${NC}"
    echo "====================================="
    echo ""

    # Verificar prerrequisitos
    check_go
    check_git

    # Determinar método de instalación
    if [ "$1" = "--source" ]; then
        install_from_source
    else
        # Intentar go install primero, luego source como fallback
        if install_with_go_install 2>/dev/null; then
            print_success "Instalación completada con 'go install'"
        else
            print_warning "'go install' falló, intentando instalación desde código fuente..."
            install_from_source
        fi
    fi

    # Verificar instalación
    verify_installation

    # Mostrar ayuda de uso
    show_usage
}

# Manejar opciones de línea de comandos
case "$1" in
    --help|-h)
        echo "Instalador de Go-ney Framework"
        echo ""
        echo "Uso: $0 [opción]"
        echo ""
        echo "Opciones:"
        echo "  --source    Instalar desde código fuente"
        echo "  --help      Mostrar esta ayuda"
        echo ""
        exit 0
        ;;
    --uninstall)
        print_status "Desinstalando Go-ney..."
        sudo rm -f "$INSTALL_DIR/$BINARY_NAME"
        print_success "Go-ney desinstalado"
        exit 0
        ;;
    *)
        main "$@"
        ;;
esac