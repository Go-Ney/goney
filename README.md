# Go-ney Framework 🚀

**Go-ney** es un framework CLI para Go inspirado en NestJS que permite crear aplicaciones con arquitectura MVC modular y soporte para microservicios TCP, NATS y gRPC.

## 🚀 Características

- **🏗️ Arquitectura Modular**: Estructura por módulos como NestJS
- **⚡ CLI Generativa**: Genera código automáticamente
- **🔗 Microservicios**: Soporte para TCP, NATS y gRPC
- **🛡️ Guards & Interceptors**: Sistema de autenticación y middleware
- **🎨 Decorators**: Sistema de anotaciones similar a NestJS
- **⚙️ Configuración**: Sistema de configuración flexible
- **🗄️ Base de Datos**: Integración con GORM
- **📁 Estructura Modular**: src/modules/ como en NestJS

## 📦 Instalación

### Opción 1: Script de instalación automática (Recomendado)
```bash
curl -fsSL https://raw.githubusercontent.com/Go-Ney/goney/main/install.sh | bash
```

### Opción 2: Usando Makefile
```bash
git clone https://github.com/Go-Ney/goney
cd goney
make install
```

### Opción 3: Instalación manual
```bash
git clone https://github.com/Go-Ney/goney
cd goney
go build -o goney ./cmd/
sudo cp goney /usr/local/bin/
```

### Opción 4: Solo para usuario actual
```bash
git clone https://github.com/Go-Ney/goney
cd goney
make install-user
```

## 🛠️ Comandos CLI

### Crear nuevo proyecto
```bash
goney new mi-proyecto
```

### Generar componentes

#### ⚡ Módulos CRUD (Como NestJS)
```bash
# Generar módulo completo con estructura modular
goney generate crud users        # Módulo src/modules/users/
goney g crud products           # Módulo src/modules/products/

# Con DTOs y modelos globales
goney generate crud orders --global

# Sin generar DTO específico
goney generate crud clients --no-dto

# Sin generar modelo específico
goney generate crud inventory --no-model

# Combinaciones
goney generate crud notifications --global --no-dto
```

#### 📁 **Nueva Estructura Modular**
```
src/
├── modules/
│   ├── users/              # 🆕 Módulo Users
│   │   ├── controllers/    # users.controller.go
│   │   ├── services/       # users.service.go
│   │   ├── repositories/   # users.repository.go
│   │   ├── dto/           # users.dto.go
│   │   ├── models/        # users.model.go
│   │   └── users.module.go # Configuración del módulo
│   └── products/          # 🆕 Módulo Products
│       ├── controllers/
│       ├── services/
│       └── ...
├── common/                # Código compartido
│   ├── dto/              # DTOs globales
│   ├── guards/           # Guards reutilizables
│   └── interceptors/     # Interceptors globales
└── config/               # Configuración
```

#### Componentes individuales
```bash
# Generar controller
goney generate controller Usuario
goney g controller Product

# Generar service
goney generate service Usuario
goney g service Product

# Generar repository
goney generate repository Usuario
goney g repository Product

# Generar microservicio
goney generate microservice grpc UserService
goney generate microservice nats NotificationService
goney generate microservice tcp ChatService

# Generar guards e interceptors
goney generate guard AuthGuard
goney generate interceptor LoggingInterceptor
```

### Iniciar proyecto
```bash
# Iniciar servidor de desarrollo
goney start
```

## 🏗️ Estructura de Proyecto

```
mi-proyecto/
├── controllers/          # Controladores HTTP
├── services/            # Lógica de negocio
├── repositories/        # Acceso a datos
├── models/              # Modelos de datos
├── dto/                 # Data Transfer Objects
├── config/              # Configuración
├── guards/              # Guards de autenticación
├── interceptors/        # Interceptors
├── decorators/          # Decorators personalizados
├── enums/               # Enumeraciones
├── middleware/          # Middleware personalizado
└── main.go              # Punto de entrada
```

## 📝 Ejemplo de Uso

### 1. Controller con Decorators
```go
package controllers

import (
    "github.com/gin-gonic/gin"
    "mi-proyecto/pkg/decorators"
    "mi-proyecto/pkg/guards"
)

type UsuarioController struct {
    usuarioService *services.UsuarioService
}

// @Controller("/api/v1/usuarios")
// @RequireAuth()
func NewUsuarioController(service *services.UsuarioService) *UsuarioController {
    return &UsuarioController{usuarioService: service}
}

// @Get("/")
// @Cache(300, "usuarios")
func (c *UsuarioController) GetAll(ctx *gin.Context) {
    // Implementación
}

// @Post("/")
// @Validate({"name": "required", "email": "required|email"})
func (c *UsuarioController) Create(ctx *gin.Context) {
    // Implementación
}
```

### 2. Service con Inyección de Dependencias
```go
package services

type UsuarioService struct {
    usuarioRepo *repositories.UsuarioRepository
}

func NewUsuarioService(repo *repositories.UsuarioRepository) *UsuarioService {
    return &UsuarioService{usuarioRepo: repo}
}

func (s *UsuarioService) CreateUser(req *dto.CreateUsuarioRequest) (*dto.UsuarioResponse, error) {
    // Lógica de negocio
}
```

### 3. Microservicio gRPC
```go
// Generado automáticamente
goney generate microservice grpc UserService
```

### 4. Guards y Middleware
```go
package main

import (
    "mi-proyecto/pkg/guards"
)

func setupRoutes(app *core.Application) {
    authGuard := guards.NewAuthGuard("secret-key")
    roleGuard := guards.NewRoleGuard("admin", "user")

    app.Use(guards.GuardMiddleware(authGuard, roleGuard))
}
```

## 🔧 Configuración

```go
// config/config.go
type Config struct {
    Port     string
    Database DatabaseConfig
    Grpc     GrpcConfig
    Nats     NatsConfig
}
```

## 🔐 Guards Disponibles

- **AuthGuard**: Autenticación por token
- **RoleGuard**: Autorización por roles
- **ThrottleGuard**: Limitación de velocidad

## 🔄 Interceptors Disponibles

- **LoggingInterceptor**: Logging de requests/responses
- **ValidationInterceptor**: Validación de datos
- **CacheInterceptor**: Cache de respuestas
- **TransformInterceptor**: Transformación de datos

## 🎯 Decorators Disponibles

- **@Get/@Post/@Put/@Delete**: Rutas HTTP
- **@Cache**: Cache de respuestas
- **@Validate**: Validación de datos
- **@RateLimit**: Limitación de velocidad
- **@Transaction**: Transacciones de BD
- **@RequirePermissions**: Permisos requeridos

## 🚀 Microservicios

### gRPC
```bash
goney generate microservice grpc UserService
```

### NATS
```bash
goney generate microservice nats EventService
```

### TCP
```bash
goney generate microservice tcp SocketService
```

## 🚀 Inicio Rápido

```bash
# 1. Instalar Go-ney
curl -fsSL https://raw.githubusercontent.com/Go-Ney/goney/main/install.sh | bash

# 2. Crear nuevo proyecto
goney new mi-api
cd mi-api

# 3. Generar módulos CRUD
goney generate crud users          # Módulo users/ completo
goney generate crud products --global  # Con DTOs globales

# 4. Iniciar servidor
goney start
```

## 🌐 Modo Global vs Específico

### Modo Específico (por defecto)
Genera archivos individuales para cada entidad:
- `dto/usuario_dto.go` - DTOs específicos
- `models/usuario_model.go` - Modelo específico

### Modo Global (--global)
Usa DTOs y modelos base reutilizables:
- `dto/global_dto.go` - DTOs base (BaseResponse, BaseCreateRequest, etc.)
- `models/global_model.go` - Modelos base (BaseModel, NamedModel, etc.)

**Ventajas del modo global:**
- ✅ Menos archivos duplicados
- ✅ Consistencia en toda la aplicación
- ✅ Fácil mantenimiento
- ✅ Ideal para proyectos grandes

## 📚 Ejemplos

Consulta la carpeta `examples/` para ver ejemplos completos de:
- API REST con autenticación
- Microservicio gRPC
- Cliente NATS
- Servidor TCP

### Comandos útiles del Makefile

```bash
make help          # Ver todos los comandos disponibles
make install       # Instalar globalmente
make install-user  # Instalar solo para usuario actual
make dev           # Configurar entorno de desarrollo
make clean         # Limpiar archivos generados
make test          # Ejecutar tests
make release       # Crear release multiplataforma
```

## 🤝 Contribuir

1. Fork del proyecto
2. Crear rama feature (`git checkout -b feature/AmazingFeature`)
3. Commit cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## 🌟 Características Avanzadas

- **Dependency Injection**: Sistema de inyección automática
- **Middleware Pipeline**: Pipeline de middleware configurable
- **Health Checks**: Endpoints de salud automáticos
- **Metrics**: Métricas integradas con Prometheus
- **Tracing**: Trazabilidad con OpenTelemetry
- **Hot Reload**: Recarga automática en desarrollo

---

Hecho con ❤️ para la comunidad Go