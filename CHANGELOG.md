# Changelog

Todos los cambios notables de Go-ney Framework serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
y este proyecto adhiere al [Versionado Semántico](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-09-20

### 🎉 Release Inicial

#### Added
- **CLI Framework completo** con arquitectura MVC
- **Comando `new`** para crear nuevos proyectos
- **Comando `generate`** con múltiples opciones:
  - `crud` - Generar CRUD completo con todas las capas
  - `controller` - Generar controllers HTTP
  - `service` - Generar services de lógica de negocio
  - `repository` - Generar repositories de acceso a datos
  - `guard` - Generar guards de autenticación
  - `interceptor` - Generar interceptors de middleware
  - `microservice` - Generar microservicios (TCP, NATS, gRPC)

#### Features Principales
- **🌐 Página de bienvenida** profesional con logo Go-ney
- **🔧 Puerto configurable** via variables de entorno
- **🌍 Modo Global vs Específico** para DTOs y modelos
- **🐳 Docker support** con Dockerfile y docker-compose
- **⚡ Hot reload** para desarrollo
- **🛡️ Sistema de Guards** (Auth, Role, Throttle)
- **🔄 Interceptors** (Logging, Validation, Cache, Transform)
- **🎨 Decorators** (@Get, @Post, @Cache, @Validate, etc.)

#### CLI Options
- **`--global`** - Usar DTOs y modelos globales
- **`--no-dto`** - No generar DTO específico
- **`--no-model`** - No generar modelo específico

#### Templates Incluidos
- **Controllers** con decoradores y documentación Swagger
- **Services** con inyección de dependencias
- **Repositories** con GORM y operaciones CRUD
- **Models** con timestamps y soft delete
- **DTOs** con validación y transformación
- **Enums** con métodos de validación

#### Microservicios
- **gRPC** - Servidor gRPC con reflection
- **NATS** - Cliente NATS con pub/sub y request/reply
- **TCP** - Servidor TCP con manejo de conexiones

#### DevOps
- **Makefile** con comandos de desarrollo
- **GitHub Actions** CI/CD pipeline
- **Multi-platform builds** (Linux, macOS, Windows)
- **Docker images** automáticos
- **Release automation** con binarios para todas las plataformas

#### Configuración
- **Variables de entorno** con .env support
- **Configuración de base de datos** (PostgreSQL)
- **Configuración CORS** personalizable
- **JWT support** para autenticación
- **Logging configurable** por niveles

### 🏗️ Arquitectura

```
proyecto/
├── controllers/     # HTTP controllers
├── services/        # Business logic
├── repositories/    # Data access
├── models/          # Database models
├── dto/             # Data transfer objects
├── guards/          # Authentication guards
├── interceptors/    # Middleware interceptors
├── decorators/      # Custom decorators
├── enums/           # Enumerations
├── config/          # Configuration
└── pkg/core/        # Framework core
```

### 📦 Instalación

```bash
# Script automático
curl -fsSL https://raw.githubusercontent.com/tu-usuario/goney/main/install.sh | bash

# Manual
make install

# Solo usuario
make install-user
```

### 🚀 Uso

```bash
# Crear proyecto
goney new mi-api

# Generar CRUD completo
goney generate crud Usuario --global

# Iniciar servidor
goney start
```

---

## [Unreleased]

### Planned Features
- [ ] ORM abstraction para múltiples bases de datos
- [ ] Plugin system para extensiones
- [ ] GraphQL support
- [ ] WebSocket support
- [ ] Metrics y monitoring integrado
- [ ] OpenAPI/Swagger generation automático
- [ ] Template engine para vistas
- [ ] Session management
- [ ] Rate limiting avanzado
- [ ] Internacionalización (i18n)