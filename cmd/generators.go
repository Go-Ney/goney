package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

func createNewProject(projectName string) {
	dirs := []string{
		"src/modules",
		"src/common/dto",
		"src/common/guards",
		"src/common/interceptors",
		"src/common/decorators",
		"src/common/enums",
		"src/common/middleware",
		"src/common/models",
		"src/config",
		"pkg/core",
		"docs",
		"tests",
	}


	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(projectName, dir), 0755)
		if err != nil {
			fmt.Printf("Error creando directorio %s: %v\n", dir, err)
			return
		}
	}

	createMainFile(projectName)
	createConfigFile(projectName)
	createCoreFile(projectName)
	createGoMod(projectName)
	createEnvFile(projectName)
	createDockerfile(projectName)
	createAppModule(projectName)
	fmt.Printf("✅ Proyecto Go-ney %s creado exitosamente!\n", projectName)
	fmt.Printf("🌐 Para iniciar: cd %s && goney start\n", projectName)
	fmt.Printf("🔧 Puerto por defecto: 8080\n")
	fmt.Printf("📝 Variables de entorno: .env\n")
	fmt.Printf("📁 Estructura modular como NestJS creada en src/modules/\n")
}

func createMainFile(projectName string) {
	mainTemplate := `package main

import (
	"fmt"
	"log"
	"os"
	"{{.ProjectName}}/config"
	"{{.ProjectName}}/pkg/core"
)

func main() {
	cfg := config.Load()
	app := core.NewApplication(cfg)

	port := cfg.Port
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("🚀 Servidor Go-ney iniciado en puerto %s\n", port)
	fmt.Printf("🌐 Visita: http://localhost:%s\n", port)
	fmt.Printf("🩺 Health: http://localhost:%s/api/v1/health\n", port)

	log.Fatal(app.Listen(":" + port))
}
`
	tmpl, _ := template.New("main").Parse(mainTemplate)
	file, _ := os.Create(filepath.Join(projectName, "main.go"))
	defer file.Close()
	tmpl.Execute(file, map[string]string{"ProjectName": projectName})
}

func createConfigFile(projectName string) {
	configTemplate := `package config

import (
	"os"
)

type Config struct {
	Port     string
	Database DatabaseConfig
	Grpc     GrpcConfig
	Nats     NatsConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type GrpcConfig struct {
	Port string
}

type NatsConfig struct {
	URL string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "{{.ProjectName}}"),
		},
		Grpc: GrpcConfig{
			Port: getEnv("GRPC_PORT", "50051"),
		},
		Nats: NatsConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
`
	tmpl, _ := template.New("config").Parse(configTemplate)
	file, _ := os.Create(filepath.Join(projectName, "config", "config.go"))
	defer file.Close()
	tmpl.Execute(file, map[string]string{"ProjectName": projectName})
}

func createCoreFile(projectName string) {
	coreTemplate := `package core

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port string
}

type Application struct {
	Router *gin.Engine
	Config *Config
}

func NewApplication(cfg interface{}) *Application {
	router := gin.Default()

	// Página de bienvenida con logo ASCII
	router.GET("/", func(c *gin.Context) {
		welcomeHTML := ` + "`" + `<!DOCTYPE html>
<html>
<head>
    <title>Go-ney Framework</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, sans-serif;
            margin: 0;
            padding: 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
        }
        .container {
            text-align: center;
            max-width: 800px;
            padding: 2rem;
        }
        .logo {
            font-family: monospace;
            font-size: 2rem;
            margin-bottom: 2rem;
            white-space: pre-line;
            color: #ffd700;
        }
        .title { font-size: 3rem; margin-bottom: 1rem; }
        .subtitle { font-size: 1.2rem; opacity: 0.9; }
        .version { margin-top: 2rem; opacity: 0.7; }
        .links { margin-top: 2rem; }
        .links a {
            color: #ffd700;
            text-decoration: none;
            margin: 0 1rem;
            padding: 0.5rem 1rem;
            border: 2px solid #ffd700;
            border-radius: 5px;
            transition: all 0.3s ease;
        }
        .links a:hover {
            background: #ffd700;
            color: #764ba2;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
   ___           _  __
  / _ \___      / |/ /__ __ __
 / (_ / _ \_   /    / -_) // /
 \___/\___(_) /_/|_/\__/\_, /
                       /___/
        </div>
        <h1 class="title">¡Bienvenido a Go-ney!</h1>
        <p class="subtitle">Framework MVC para Go inspirado en NestJS</p>
        <div class="links">
            <a href="/api/v1/health">🩺 Health Check</a>
            <a href="https://github.com/tu-usuario/go-ney" target="_blank">📚 Documentación</a>
        </div>
        <div class="version">Go-ney v1.0.1 | Puerto: {{.Port}}</div>
    </div>
</body>
</html>` + "`" + `
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, welcomeHTML)
	})

	// Health check endpoint
	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Go-ney API está funcionando correctamente",
			"version": "1.0.1",
			"port":    cfg,
		})
	})

	return &Application{
		Router: router,
		Config: &Config{},
	}
}

func (app *Application) Listen(addr string) error {
	fmt.Printf("🚀 Go-ney iniciado en %s\n", addr)
	return app.Router.Run(addr)
}
`
	tmpl, _ := template.New("core").Parse(coreTemplate)
	file, _ := os.Create(filepath.Join(projectName, "pkg", "core", "application.go"))
	defer file.Close()
	tmpl.Execute(file, map[string]string{"ProjectName": projectName})
}

func createGoMod(projectName string) {
	goModTemplate := `module {{.ProjectName}}

go 1.23

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/spf13/viper v1.17.0
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
	github.com/nats-io/nats.go v1.31.0
	gorm.io/gorm v1.25.5
	gorm.io/driver/postgres v1.5.4
)
`
	tmpl, _ := template.New("gomod").Parse(goModTemplate)
	file, _ := os.Create(filepath.Join(projectName, "go.mod"))
	defer file.Close()
	tmpl.Execute(file, map[string]string{"ProjectName": projectName})
}

func createEnvFile(projectName string) {
	envTemplate := `# Go-ney Configuration
# Puerto del servidor (por defecto: 8080)
PORT=8080

# Base de datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME={{.ProjectName}}

# gRPC
GRPC_PORT=50051

# NATS
NATS_URL=nats://localhost:4222

# Configuración de la aplicación
APP_ENV=development
APP_DEBUG=true
APP_LOG_LEVEL=info

# JWT Secret (cambiar en producción)
JWT_SECRET=tu-jwt-secret-super-seguro

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization
`

	tmpl, _ := template.New("env").Parse(envTemplate)
	file, _ := os.Create(filepath.Join(projectName, ".env"))
	defer file.Close()
	tmpl.Execute(file, map[string]string{"ProjectName": projectName})

	// Crear también .env.example
	exampleFile, _ := os.Create(filepath.Join(projectName, ".env.example"))
	defer exampleFile.Close()
	tmpl.Execute(exampleFile, map[string]string{"ProjectName": projectName + "_example"})
}

func createDockerfile(projectName string) {
	dockerTemplate := `# Go-ney Dockerfile
FROM golang:1.23-alpine AS builder

# Instalar dependencias del sistema
RUN apk --no-cache add ca-certificates git

WORKDIR /app

# Copiar go mod y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Imagen final
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binario desde builder
COPY --from=builder /app/main .

# Crear usuario no-root
RUN adduser -D -s /bin/sh goney
USER goney

# Exponer puerto
EXPOSE 8080

# Comando por defecto
CMD ["./main"]
`

	file, _ := os.Create(filepath.Join(projectName, "Dockerfile"))
	defer file.Close()
	file.WriteString(dockerTemplate)

	// Crear docker-compose.yml
	dockerComposeTemplate := `version: '3.8'

services:
  {{.ProjectName}}:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=password
      - DB_NAME={{.ProjectName}}
    depends_on:
      - postgres
      - redis
      - nats
    volumes:
      - .:/app
    restart: unless-stopped

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: {{.ProjectName}}
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    restart: unless-stopped

  nats:
    image: nats:2.10-alpine
    ports:
      - "4222:4222"
      - "8222:8222"
    restart: unless-stopped

volumes:
  postgres_data:
`

	tmpl, _ := template.New("docker-compose").Parse(dockerComposeTemplate)
	composeFile, _ := os.Create(filepath.Join(projectName, "docker-compose.yml"))
	defer composeFile.Close()
	tmpl.Execute(composeFile, map[string]string{"ProjectName": projectName})
}

func createAppModule(projectName string) {
	appModuleTemplate := `package main

import (
	"{{.ProjectName}}/config"
	"{{.ProjectName}}/pkg/core"
)

type AppModule struct {
	Config *config.Config
	Core   *core.Application
}

func NewAppModule() *AppModule {
	cfg := config.Load()

	coreConfig := &core.Config{
		Port: cfg.Port,
		Database: core.DatabaseConfig{
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			Name:     cfg.Database.Name,
		},
		Grpc: core.GrpcConfig{
			Port: cfg.Grpc.Port,
		},
		Nats: core.NatsConfig{
			URL: cfg.Nats.URL,
		},
	}

	app := core.NewApplication(coreConfig)

	return &AppModule{
		Config: cfg,
		Core:   app,
	}
}

func (app *AppModule) Bootstrap() error {
	// Registrar módulos aquí
	// app.RegisterModule(&users.UsersModule{})
	// app.RegisterModule(&products.ProductsModule{})

	return nil
}
`

	tmpl, _ := template.New("app-module").Parse(appModuleTemplate)
	file, _ := os.Create(filepath.Join(projectName, "src", "app.module.go"))
	defer file.Close()
	tmpl.Execute(file, map[string]string{"ProjectName": projectName})
}

func generateModuleController(moduleName string, global bool) {
    project := getProjectModuleName()
    dtoImport := fmt.Sprintf("\"%s/src/modules/%s/dto\"", project, moduleName)
    dtoTypePrefix := "dto."
    if global {
        dtoImport = fmt.Sprintf("\"%s/src/common/dto\"", project)
        dtoTypePrefix = "dto.Base"
    }

    controllerTemplate := `package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "{{.Project}}/src/modules/{{.ModulePath}}/services"
    {{.DTOImport}}
)

type {{.ControllerName}}Controller struct {
    {{.ServiceName}}Service *services.{{.ControllerName}}Service
}

func New{{.ControllerName}}Controller({{.ServiceName}}Service *services.{{.ControllerName}}Service) *{{.ControllerName}}Controller {
    return &{{.ControllerName}}Controller{
        {{.ServiceName}}Service: {{.ServiceName}}Service,
    }
}

// @Router /api/v1/{{.RouterPath}} [get]
// @Summary Get all {{.EntityName}}
// @Tags {{.ControllerName}}
// @Accept json
// @Produce json
// @Success 200 {array} {{.DTOTypePrefix}}Response
func (c *{{.ControllerName}}Controller) GetAll(ctx *gin.Context) {
    result, err := c.{{.ServiceName}}Service.GetAll()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, result)
}

// @Router /api/v1/{{.RouterPath}}/{id} [get]
// @Summary Get {{.EntityName}} by ID
// @Tags {{.ControllerName}}
// @Accept json
// @Produce json
// @Param id path string true "{{.ControllerName}} ID"
// @Success 200 {object} {{.DTOTypePrefix}}Response
func (c *{{.ControllerName}}Controller) GetByID(ctx *gin.Context) {
    id := ctx.Param("id")
    result, err := c.{{.ServiceName}}Service.GetByID(id)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, result)
}

// @Router /api/v1/{{.RouterPath}} [post]
// @Summary Create {{.EntityName}}
// @Tags {{.ControllerName}}
// @Accept json
// @Produce json
// @Param body body {{.DTOTypePrefix}}CreateRequest true "{{.ControllerName}} data"
// @Success 201 {object} {{.DTOTypePrefix}}Response
func (c *{{.ControllerName}}Controller) Create(ctx *gin.Context) {
    var req {{.DTOTypePrefix}}CreateRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := c.{{.ServiceName}}Service.Create(&req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, result)
}

// @Router /api/v1/{{.RouterPath}}/{id} [put]
// @Summary Update {{.EntityName}}
// @Tags {{.ControllerName}}
// @Accept json
// @Produce json
// @Param id path string true "{{.ControllerName}} ID"
// @Param body body {{.DTOTypePrefix}}UpdateRequest true "{{.ControllerName}} data"
// @Success 200 {object} {{.DTOTypePrefix}}Response
func (c *{{.ControllerName}}Controller) Update(ctx *gin.Context) {
    id := ctx.Param("id")
    var req {{.DTOTypePrefix}}UpdateRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := c.{{.ServiceName}}Service.Update(id, &req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, result)
}

// @Router /api/v1/{{.RouterPath}}/{id} [delete]
// @Summary Delete {{.EntityName}}
// @Tags {{.ControllerName}}
// @Accept json
// @Produce json
// @Param id path string true "{{.ControllerName}} ID"
// @Success 204
func (c *{{.ControllerName}}Controller) Delete(ctx *gin.Context) {
    id := ctx.Param("id")
    err := c.{{.ServiceName}}Service.Delete(id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.Status(http.StatusNoContent)
}
`

    capitalizedName := strings.Title(moduleName)
    serviceName := strings.ToLower(moduleName[:1]) + moduleName[1:]
    routerPath := moduleName
    entityName := strings.ToLower(moduleName)

    tmpl, _ := template.New("module-controller").Parse(controllerTemplate)
    fileName := fmt.Sprintf("src/modules/%s/controllers/%s.controller.go", moduleName, moduleName)
    file, _ := os.Create(fileName)
    defer file.Close()

    data := map[string]string{
        "Project":        project,
        "ModulePath":     moduleName,
        "ControllerName": capitalizedName,
        "ServiceName":    serviceName,
        "RouterPath":     routerPath,
        "EntityName":     entityName,
        "DTOImport":      dtoImport,
        "DTOTypePrefix":  dtoTypePrefix,
    }
    tmpl.Execute(file, data)
}

func createModuleStructure(moduleDir string) {
    // Crear solo el directorio del módulo (estructura plana)
    if err := os.MkdirAll(moduleDir, 0755); err != nil {
        fmt.Printf("Error creando directorio %s: %v\n", moduleDir, err)
        return
    }
}

func generateModuleService(moduleName string, global bool) {
    project := getProjectModuleName()
    dtoImport := fmt.Sprintf("\"%s/src/modules/%s/dto\"", project, moduleName)
    modelImport := fmt.Sprintf("\"%s/src/modules/%s/models\"", project, moduleName)
    dtoTypePrefix := "dto."
    modelTypePrefix := "models."
    if global {
        dtoImport = fmt.Sprintf("\"%s/src/common/dto\"", project)
        modelImport = fmt.Sprintf("\"%s/src/common/models\"", project)
        dtoTypePrefix = "dto.Base"
        modelTypePrefix = "models.Base"
    }

    serviceTemplate := `package services

import (
    "{{.Project}}/src/modules/{{.ModuleName}}/repositories"
    {{.DTOImport}}
    {{.ModelImport}}
)

type {{.ClassName}}Service struct {
    {{.RepositoryVar}}Repository *repositories.{{.ClassName}}Repository
}

func New{{.ClassName}}Service({{.RepositoryVar}}Repository *repositories.{{.ClassName}}Repository) *{{.ClassName}}Service {
    return &{{.ClassName}}Service{
        {{.RepositoryVar}}Repository: {{.RepositoryVar}}Repository,
    }
}

func (s *{{.ClassName}}Service) GetAll() ([]{{.DTOTypePrefix}}Response, error) {
    entities, err := s.{{.RepositoryVar}}Repository.FindAll()
    if err != nil {
        return nil, err
    }

    var responses []{{.DTOTypePrefix}}Response
    for _, entity := range entities {
        responses = append(responses, {{.DTOTypePrefix}}Response{
            ID: entity.ID,
        })
    }
    return responses, nil
}

func (s *{{.ClassName}}Service) GetByID(id string) (*{{.DTOTypePrefix}}Response, error) {
    entity, err := s.{{.RepositoryVar}}Repository.FindByID(id)
    if err != nil {
        return nil, err
    }

    return &{{.DTOTypePrefix}}Response{
        ID: entity.ID,
    }, nil
}

func (s *{{.ClassName}}Service) Create(req *{{.DTOTypePrefix}}CreateRequest) (*{{.DTOTypePrefix}}Response, error) {
    entity := &{{.ModelTypePrefix}}{{.ClassName}}{}
    createdEntity, err := s.{{.RepositoryVar}}Repository.Create(entity)
    if err != nil {
        return nil, err
    }
    return &{{.DTOTypePrefix}}Response{ID: createdEntity.ID}, nil
}

func (s *{{.ClassName}}Service) Update(id string, req *{{.DTOTypePrefix}}UpdateRequest) (*{{.DTOTypePrefix}}Response, error) {
    entity, err := s.{{.RepositoryVar}}Repository.FindByID(id)
    if err != nil {
        return nil, err
    }
    updatedEntity, err := s.{{.RepositoryVar}}Repository.Update(entity)
    if err != nil {
        return nil, err
    }
    return &{{.DTOTypePrefix}}Response{ID: updatedEntity.ID}, nil
}

func (s *{{.ClassName}}Service) Delete(id string) error {
    return s.{{.RepositoryVar}}Repository.Delete(id)
}
`

    capitalizedName := strings.Title(moduleName)
    repositoryVar := strings.ToLower(moduleName[:1]) + moduleName[1:]

    tmpl, _ := template.New("module-service").Parse(serviceTemplate)
    fileName := fmt.Sprintf("src/modules/%s/services/%s.service.go", moduleName, moduleName)
    file, _ := os.Create(fileName)
    defer file.Close()

    data := map[string]string{
        "Project":         project,
        "ModuleName":      moduleName,
        "ClassName":       capitalizedName,
        "RepositoryVar":   repositoryVar,
        "DTOImport":       dtoImport,
        "ModelImport":     modelImport,
        "DTOTypePrefix":   dtoTypePrefix,
        "ModelTypePrefix": modelTypePrefix,
    }
    tmpl.Execute(file, data)
}

func generateModuleRepository(moduleName string) {
    repositoryTemplate := `package repositories

// NOTE: Implementa persistencia real según tu proyecto (GORM, SQL, etc.)
type {{.ClassName}}Repository struct{}

func New{{.ClassName}}Repository() *{{.ClassName}}Repository { return &{{.ClassName}}Repository{} }

func (r *{{.ClassName}}Repository) FindAll() ([]{{.ClassName}}, error) { return []{{.ClassName}}{}, nil }
func (r *{{.ClassName}}Repository) FindByID(id string) (*{{.ClassName}}, error) { return &{{.ClassName}}{ID: id}, nil }
func (r *{{.ClassName}}Repository) Create(entity *{{.ClassName}}) (*{{.ClassName}}, error) { return entity, nil }
func (r *{{.ClassName}}Repository) Update(entity *{{.ClassName}}) (*{{.ClassName}}, error) { return entity, nil }
func (r *{{.ClassName}}Repository) Delete(id string) error { return nil }
func (r *{{.ClassName}}Repository) FindBy(field string, value interface{}) ([]{{.ClassName}}, error) { return []{{.ClassName}}{}, nil }
`

    capitalizedName := strings.Title(moduleName)

    tmpl, _ := template.New("module-repository").Parse(repositoryTemplate)
    fileName := fmt.Sprintf("src/modules/%s/repositories/%s.repository.go", moduleName, moduleName)
    file, _ := os.Create(fileName)
    defer file.Close()

    data := map[string]string{
        "ClassName": capitalizedName,
    }
    tmpl.Execute(file, data)
}

func generateModuleFile(moduleName string) {
    moduleTemplate := `package {{.ModuleName}}

import (
    "{{.Project}}/src/modules/{{.ModuleName}}/controllers"
    "{{.Project}}/src/modules/{{.ModuleName}}/services"
    "{{.Project}}/src/modules/{{.ModuleName}}/repositories"
)

type {{.ClassName}}Module struct {
    Controller *controllers.{{.ClassName}}Controller
    Service    *services.{{.ClassName}}Service
    Repository *repositories.{{.ClassName}}Repository
}

func New{{.ClassName}}Module() *{{.ClassName}}Module {
    repository := repositories.New{{.ClassName}}Repository()
    service := services.New{{.ClassName}}Service(repository)
    controller := controllers.New{{.ClassName}}Controller(service)

    return &{{.ClassName}}Module{
        Controller: controller,
        Service:    service,
        Repository: repository,
    }
}

func (m *{{.ClassName}}Module) RegisterRoutes(router interface{}) {
    // TODO: Registrar rutas del módulo en tu router
}
`

    capitalizedName := strings.Title(moduleName)
    project := getProjectModuleName()

    tmpl, _ := template.New("module").Parse(moduleTemplate)
    fileName := fmt.Sprintf("src/modules/%s/%s.module.go", moduleName, moduleName)
    file, _ := os.Create(fileName)
    defer file.Close()

    data := map[string]string{
        "ModuleName": moduleName,
        "ClassName":  capitalizedName,
        "Project":    project,
    }
    tmpl.Execute(file, data)
}

func generateModuleModel(moduleName string) {
    modelTemplate := `package models

import (
    "time"
)

type {{.ClassName}} struct {
    ID        string    ` + "`json:\"id\"`" + `
    CreatedAt time.Time ` + "`json:\"created_at\"`" + `
    UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
    Name        string ` + "`json:\"name\"`" + `
    Description string ` + "`json:\"description\"`" + `
    Status      string ` + "`json:\"status\"`" + `
}

func ({{.ClassName}}) TableName() string {
    return "{{.TableName}}"
}
`

    capitalizedName := strings.Title(moduleName)
    tableName := moduleName

    tmpl, _ := template.New("module-model").Parse(modelTemplate)
    fileName := fmt.Sprintf("src/modules/%s/models/%s.model.go", moduleName, moduleName)
    file, _ := os.Create(fileName)
    defer file.Close()

    data := map[string]string{
        "ClassName":  capitalizedName,
        "TableName":  tableName,
    }
    tmpl.Execute(file, data)
}

func generateModuleTests(moduleName string) {
    testTemplate := `package {{.ModuleName}}

import "testing"

func Test{{.ClassName}}Module_Basic(t *testing.T) {
    t.Run("placeholder", func(t *testing.T) {
        // TODO: implement tests
    })
}
`

    capitalizedName := strings.Title(moduleName)

    tmpl, _ := template.New("module-tests").Parse(testTemplate)
    fileName := fmt.Sprintf("src/modules/%s/%s_test.go", moduleName, moduleName)
    file, _ := os.Create(fileName)
    defer file.Close()

    data := map[string]string{
        "ModuleName": moduleName,
        "ClassName":  capitalizedName,
    }
    tmpl.Execute(file, data)
}

func generateModuleDTO(moduleName string) {
    dtoTemplate := `package dto

import "time"

type {{.ClassName}}Response struct {
    ID        string    ` + "`json:\"id\"`" + `
    CreatedAt time.Time ` + "`json:\"created_at\"`" + `
    UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

type Create{{.ClassName}}Request struct {
    Name        string ` + "`json:\"name\"`" + `
    Description string ` + "`json:\"description\"`" + `
}

type Update{{.ClassName}}Request struct {
    Name        string ` + "`json:\"name,omitempty\"`" + `
    Description string ` + "`json:\"description,omitempty\"`" + `
}
`

    capitalizedName := strings.Title(moduleName)

    tmpl, _ := template.New("module-dto").Parse(dtoTemplate)
    fileName := fmt.Sprintf("src/modules/%s/dto/%s.dto.go", moduleName, moduleName)
    file, _ := os.Create(fileName)
    defer file.Close()

    data := map[string]string{
        "ClassName": capitalizedName,
    }
    tmpl.Execute(file, data)
}

func ensureGlobalFiles() {
    // Crear DTO global si no existe en src/common/dto
    if _, err := os.Stat("src/common/dto"); os.IsNotExist(err) {
        os.MkdirAll("src/common/dto", 0755)
    }
    if _, err := os.Stat("src/common/models"); os.IsNotExist(err) {
        os.MkdirAll("src/common/models", 0755)
    }
    if _, err := os.Stat("src/common/dto/base.go"); os.IsNotExist(err) {
        createGlobalDTO()
    }
    if _, err := os.Stat("src/common/models/base.go"); os.IsNotExist(err) {
        createGlobalModel()
    }
}

func createGlobalDTO() {
    globalDTOTemplate := `package dto

import "time"

type BaseResponse struct {
    ID        string    ` + "`json:\"id\"`" + `
    CreatedAt time.Time ` + "`json:\"created_at\"`" + `
    UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

type BaseCreateRequest struct {
    Name        string ` + "`json:\"name\"`" + `
    Description string ` + "`json:\"description\"`" + `
}

type BaseUpdateRequest struct {
    Name        string ` + "`json:\"name,omitempty\"`" + `
    Description string ` + "`json:\"description,omitempty\"`" + `
}
`

    file, _ := os.Create("src/common/dto/base.go")
    defer file.Close()
    file.WriteString(globalDTOTemplate)
    fmt.Printf("✅ DTO global creado en src/common/dto/base.go\n")
}

func createGlobalModel() {
    globalModelTemplate := `package models

import "time"

type Base struct {
    ID        string    ` + "`json:\"id\"`" + `
    CreatedAt time.Time ` + "`json:\"created_at\"`" + `
    UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

type Named struct {
    Base
    Name        string ` + "`json:\"name\"`" + `
    Description string ` + "`json:\"description\"`" + `
    Status      string ` + "`json:\"status\"`" + `
}
`

    file, _ := os.Create("src/common/models/base.go")
    defer file.Close()
    file.WriteString(globalModelTemplate)
    fmt.Printf("✅ Modelo global creado en src/common/models/base.go\n")
}

func generateControllerWithOptions(name string, global bool) {
    dtoImport := "\"{{.ModuleName}}/dto\""
    dtoType := "dto."
    if global {
        dtoType = "dto.Base"
    }

	if global {
		dtoType = "dto.Base"
	}

	controllerTemplate := `package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"{{.ModuleName}}/services"
	{{.DTOImport}}
)

type {{.ControllerName}}Controller struct {
	{{.ServiceName}}Service *services.{{.ServiceName}}Service
}

func New{{.ControllerName}}Controller({{.ServiceName}}Service *services.{{.ServiceName}}Service) *{{.ControllerName}}Controller {
	return &{{.ControllerName}}Controller{
		{{.ServiceName}}Service: {{.ServiceName}}Service,
	}
}

// @Router /{{.RouterPath}} [get]
// @Summary Get all {{.EntityName}}
// @Tags {{.ControllerName}}
// @Accept json
// @Produce json
// @Success 200 {array} {{.DTOType}}Response
func (c *{{.ControllerName}}Controller) GetAll(ctx *gin.Context) {
	result, err := c.{{.ServiceName}}Service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// @Router /{{.RouterPath}}/{id} [get]
// @Summary Get {{.EntityName}} by ID
// @Tags {{.ControllerName}}
// @Accept json
// @Produce json
// @Param id path string true "{{.ControllerName}} ID"
// @Success 200 {object} {{.DTOType}}Response
func (c *{{.ControllerName}}Controller) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.{{.ServiceName}}Service.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// @Router /{{.RouterPath}} [post]
// @Summary Create {{.EntityName}}
// @Tags {{.ControllerName}}
// @Accept json
// @Produce json
// @Param body body {{.DTOType}}CreateRequest true "{{.ControllerName}} data"
// @Success 201 {object} {{.DTOType}}Response
func (c *{{.ControllerName}}Controller) Create(ctx *gin.Context) {
	var req {{.DTOType}}CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.{{.ServiceName}}Service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, result)
}

// @Router /{{.RouterPath}}/{id} [put]
// @Summary Update {{.EntityName}}
// @Tags {{.ControllerName}}
// @Accept json
// @Produce json
// @Param id path string true "{{.ControllerName}} ID"
// @Param body body {{.DTOType}}UpdateRequest true "{{.ControllerName}} data"
// @Success 200 {object} {{.DTOType}}Response
func (c *{{.ControllerName}}Controller) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req {{.DTOType}}UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.{{.ServiceName}}Service.Update(id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// @Router /{{.RouterPath}}/{id} [delete]
// @Summary Delete {{.EntityName}}
// @Tags {{.ControllerName}}
// @Accept json
// @Produce json
// @Param id path string true "{{.ControllerName}} ID"
// @Success 204
func (c *{{.ControllerName}}Controller) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.{{.ServiceName}}Service.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
`

	moduleName := getModuleName()
	capitalizedName := strings.Title(name)
	serviceName := strings.ToLower(name[:1]) + name[1:]
	routerPath := strings.ToLower(name) + "s"
	entityName := strings.ToLower(name)

	tmpl, _ := template.New("controller").Parse(controllerTemplate)
	fileName := fmt.Sprintf("controllers/%s_controller.go", strings.ToLower(name))
	file, _ := os.Create(fileName)
	defer file.Close()

	data := map[string]string{
		"ModuleName":     moduleName,
		"ControllerName": capitalizedName,
		"ServiceName":    serviceName,
		"RouterPath":     routerPath,
		"EntityName":     entityName,
		"DTOImport":      dtoImport,
		"DTOType":        dtoType,
	}
	tmpl.Execute(file, data)

	if !global {
		generateDTO(name)
	}
	fmt.Printf("✅ Controller %s generado en %s\n", capitalizedName, fileName)
}

func generateServiceWithOptions(name string, global bool) {
	dtoType := "dto."
	modelType := "models."

	if global {
		dtoType = "dto.Base"
		modelType = "models.Named"
	} else {
		dtoType = "dto."
		modelType = "models."
	}

	serviceTemplate := `package services

import (
	"{{.ModuleName}}/repositories"
	"{{.ModuleName}}/dto"
	"{{.ModuleName}}/models"
)

type {{.ServiceName}}Service struct {
	{{.RepositoryName}}Repository *repositories.{{.ServiceName}}Repository
}

func New{{.ServiceName}}Service({{.RepositoryName}}Repository *repositories.{{.ServiceName}}Repository) *{{.ServiceName}}Service {
	return &{{.ServiceName}}Service{
		{{.RepositoryName}}Repository: {{.RepositoryName}}Repository,
	}
}

func (s *{{.ServiceName}}Service) GetAll() ([]{{.DTOType}}Response, error) {
	entities, err := s.{{.RepositoryName}}Repository.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []{{.DTOType}}Response
	for _, entity := range entities {
		responses = append(responses, {{.DTOType}}Response{
			ID: entity.ID,
			// Map other fields here
		})
	}
	return responses, nil
}

func (s *{{.ServiceName}}Service) GetByID(id string) (*{{.DTOType}}Response, error) {
	entity, err := s.{{.RepositoryName}}Repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &{{.DTOType}}Response{
		ID: entity.ID,
		// Map other fields here
	}, nil
}

func (s *{{.ServiceName}}Service) Create(req *{{.DTOType}}CreateRequest) (*{{.DTOType}}Response, error) {
	entity := &{{.ModelType}}{{.ServiceName}}{
		// Map fields from request
	}

	createdEntity, err := s.{{.RepositoryName}}Repository.Create(entity)
	if err != nil {
		return nil, err
	}

	return &{{.DTOType}}Response{
		ID: createdEntity.ID,
		// Map other fields here
	}, nil
}

func (s *{{.ServiceName}}Service) Update(id string, req *{{.DTOType}}UpdateRequest) (*{{.DTOType}}Response, error) {
	entity, err := s.{{.RepositoryName}}Repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields from request

	updatedEntity, err := s.{{.RepositoryName}}Repository.Update(entity)
	if err != nil {
		return nil, err
	}

	return &{{.DTOType}}Response{
		ID: updatedEntity.ID,
		// Map other fields here
	}, nil
}

func (s *{{.ServiceName}}Service) Delete(id string) error {
	return s.{{.RepositoryName}}Repository.Delete(id)
}
`

	moduleName := getModuleName()
	capitalizedName := strings.Title(name)
	repositoryName := strings.ToLower(name[:1]) + name[1:]

	tmpl, _ := template.New("service").Parse(serviceTemplate)
	fileName := fmt.Sprintf("services/%s_service.go", strings.ToLower(name))
	file, _ := os.Create(fileName)
	defer file.Close()

	data := map[string]string{
		"ModuleName":     moduleName,
		"ServiceName":    capitalizedName,
		"RepositoryName": repositoryName,
		"DTOType":        dtoType,
		"ModelType":      modelType,
	}
	tmpl.Execute(file, data)

	if !global {
		generateModel(name)
	}
	fmt.Printf("✅ Service %s generado en %s\n", capitalizedName, fileName)
}

func generateModule(moduleName string, crud, global, noDto, noModel bool) {
	fmt.Printf("🚀 Generando módulo: %s\n", moduleName)

	if global {
		fmt.Printf("🌐 Modo global activado - usando DTOs y modelos globales\n")
		ensureGlobalFiles()
	}

	// Crear estructura del módulo
	moduleDir := fmt.Sprintf("src/modules/%s", moduleName)
	createModuleStructure(moduleDir)

	// Generar un único archivo plano del módulo + test
	writeSingleFileModule(moduleName, global, noDto, noModel)

	// Generar model y DTO según las opciones
	// (Ya incluidos dentro del archivo único si no se usan global/no-dto/no-model)

	// Si se pidió CRUD, también generamos tests
	if crud {
		// ya se generó un test básico; se puede extender aquí si se requiere
	}

	fmt.Printf("✅ Módulo %s generado exitosamente!\n", moduleName)
	fmt.Printf("📁 Archivos creados en src/modules/%s/:\n", moduleName)
	fmt.Printf("   - %s.go\n", moduleName)
	fmt.Printf("   - %s_test.go\n", moduleName)

	if global {
		fmt.Printf("   📌 Usando DTOs y modelos globales en src/common/\n")
	}

	fmt.Printf("\n💡 Para generar un módulo con CRUD completo, usa: goney generate module %s --crud\n", moduleName)
}

// writeSingleFileModule crea un archivo plano con Controller, Service, Repository,
// DTO y Model en el mismo paquete del módulo, sin subcarpetas.
func writeSingleFileModule(moduleName string, global, noDto, noModel bool) {
    pkg := moduleName

    // Construcción condicional de secciones DTO/Model
    dtoSection := ""
    modelSection := ""
    if global {
        dtoSection = "// Usando DTOs globales desde src/common/dto (import manual en tus handlers)\n"
        modelSection = "// Usando modelos globales desde src/common/models (import manual si aplica)\n"
    } else {
        if !noDto {
            dtoSection = "type " + strings.Title(moduleName) + "Response struct {\n\tID string `json:\"id\"`\n}\n\n" +
                "type Create" + strings.Title(moduleName) + "Request struct {\n\tName string `json:\"name\"`\n}\n\n" +
                "type Update" + strings.Title(moduleName) + "Request struct {\n\tName string `json:\"name,omitempty\"`\n}\n\n"
        }
        if !noModel {
            modelSection = "type " + strings.Title(moduleName) + " struct {\n\tID string `json:\"id\"`\n\tName string `json:\"name\"`\n}\n\n"
        }
    }

    code := "package " + pkg + "\n\n" +
        "// Archivo único generado por Goney.\n" +
        dtoSection +
        modelSection +
        "// Repository\n" +
        "type " + strings.Title(moduleName) + "Repository struct{}\n" +
        "func New" + strings.Title(moduleName) + "Repository() *" + strings.Title(moduleName) + "Repository { return &" + strings.Title(moduleName) + "Repository{} }\n" +
        "func (r *" + strings.Title(moduleName) + "Repository) FindAll() ([]" + strings.Title(moduleName) + ", error) { return []" + strings.Title(moduleName) + "{}, nil }\n" +
        "func (r *" + strings.Title(moduleName) + "Repository) FindByID(id string) (*" + strings.Title(moduleName) + ", error) { return &" + strings.Title(moduleName) + "{ID: id}, nil }\n" +
        "func (r *" + strings.Title(moduleName) + "Repository) Create(e *" + strings.Title(moduleName) + ") (*" + strings.Title(moduleName) + ", error) { return e, nil }\n" +
        "func (r *" + strings.Title(moduleName) + "Repository) Update(e *" + strings.Title(moduleName) + ") (*" + strings.Title(moduleName) + ", error) { return e, nil }\n" +
        "func (r *" + strings.Title(moduleName) + "Repository) Delete(id string) error { return nil }\n\n" +
        "// Service\n" +
        "type " + strings.Title(moduleName) + "Service struct { repo *" + strings.Title(moduleName) + "Repository }\n" +
        "func New" + strings.Title(moduleName) + "Service(repo *" + strings.Title(moduleName) + "Repository) *" + strings.Title(moduleName) + "Service { return &" + strings.Title(moduleName) + "Service{repo: repo} }\n\n" +
        "// Controller (placeholder)\n" +
        "type " + strings.Title(moduleName) + "Controller struct { svc *" + strings.Title(moduleName) + "Service }\n" +
        "func New" + strings.Title(moduleName) + "Controller(svc *" + strings.Title(moduleName) + "Service) *" + strings.Title(moduleName) + "Controller { return &" + strings.Title(moduleName) + "Controller{svc: svc} }\n\n" +
        "// Module wiring\n" +
        "type " + strings.Title(moduleName) + "Module struct { Controller *" + strings.Title(moduleName) + "Controller; Service *" + strings.Title(moduleName) + "Service; Repository *" + strings.Title(moduleName) + "Repository }\n" +
        "func New" + strings.Title(moduleName) + "Module() *" + strings.Title(moduleName) + "Module {\n" +
        "\trepo := New" + strings.Title(moduleName) + "Repository()\n" +
        "\tsvc := New" + strings.Title(moduleName) + "Service(repo)\n" +
        "\tctrl := New" + strings.Title(moduleName) + "Controller(svc)\n" +
        "\treturn &" + strings.Title(moduleName) + "Module{Controller: ctrl, Service: svc, Repository: repo}\n}" + "\n"

    // Escribir archivo principal del módulo
    filePath := fmt.Sprintf("src/modules/%s/%s.go", moduleName, moduleName)
    _ = os.WriteFile(filePath, []byte(code), 0644)

    // Escribir test básico
    testCode := "package " + pkg + "\n\nimport \"testing\"\n\nfunc Test" + strings.Title(moduleName) + "Module_Bootstrap(t *testing.T) { _ = New" + strings.Title(moduleName) + "Module() }\n"
    testPath := fmt.Sprintf("src/modules/%s/%s_test.go", moduleName, moduleName)
    _ = os.WriteFile(testPath, []byte(testCode), 0644)
}

func getProjectModuleName() string {
	if _, err := os.Stat("go.mod"); err == nil {
		content, _ := os.ReadFile("go.mod")
		lines := strings.Split(string(content), "\n")
		if len(lines) > 0 && strings.HasPrefix(lines[0], "module ") {
			return strings.TrimSpace(strings.TrimPrefix(lines[0], "module "))
		}
	}
	return "myapp"
}

func getModuleName() string { return getProjectModuleName() }

func startDevServer() {
	if _, err := os.Stat("main.go"); err != nil {
		fmt.Println("❌ No se encontró main.go. Asegúrate de estar en un proyecto Go-ney.")
		return
	}

	fmt.Println("📦 Instalando dependencias...")
	if err := runCommand("go", "mod", "tidy"); err != nil {
		fmt.Printf("❌ Error instalando dependencias: %v\n", err)
		return
	}

	fmt.Println("🔧 Compilando proyecto...")
	if err := runCommand("go", "build", "-o", "app", "."); err != nil {
		fmt.Printf("❌ Error compilando: %v\n", err)
		return
	}

	fmt.Println("🚀 Iniciando servidor...")
	if err := runCommand("./app"); err != nil {
		fmt.Printf("❌ Error iniciando servidor: %v\n", err)
		return
	}
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// --- Stubs de compatibilidad para comandos legacy ---
func generateController(name string) { ensureFlatModuleDir(name); generateFlatController(name); fmt.Printf("✅ Controller %s generado\n", name) }
func generateService(name string)    { ensureFlatModuleDir(name); generateFlatService(name); fmt.Printf("✅ Service %s generado\n", name) }
func generateRepository(name string) { ensureFlatModuleDir(name); generateFlatRepository(name); fmt.Printf("✅ Repository %s generado\n", name) }
func generateDTO(name string)        { ensureFlatModuleDir(name); generateFlatDTO(name); fmt.Printf("✅ DTO %s generado\n", name) }
func generateModel(name string)      { ensureFlatModuleDir(name); generateFlatModel(name); fmt.Printf("✅ Model %s generado\n", name) }

func generateMicroservice(serviceType, name string) {
    fmt.Printf("✅ Microservicio %s %s (stub). Implementa tu plantilla según necesidades.\n", serviceType, name)
}

func generateGuard(name string) {
    fmt.Printf("✅ Guard %s (stub). Implementa tu plantilla según necesidades.\n", name)
}

func generateInterceptor(name string) {
    fmt.Printf("✅ Interceptor %s (stub). Implementa tu plantilla según necesidades.\n", name)
}

func generateModuleCRUD(moduleName string, global, noDto, noModel bool) {
    // Reutiliza el generador de módulo con flag crud para incluir tests
    generateModule(moduleName, true, global, noDto, noModel)
}

func ensureFlatModuleDir(name string) {
    moduleDir := filepath.Join("src", "modules", name)
    _ = os.MkdirAll(moduleDir, 0755)
}

func generateFlatController(name string) {
    class := strings.Title(name)
    code := "package " + name + "\n\n" +
        "// Controller generado\n" +
        "type " + class + "Controller struct { svc *" + class + "Service }\n" +
        "func New" + class + "Controller(svc *" + class + "Service) *" + class + "Controller { return &" + class + "Controller{svc: svc} }\n"
    path := filepath.Join("src", "modules", name, name+".controller.go")
    _ = os.WriteFile(path, []byte(code), 0644)
}

func generateFlatService(name string) {
    class := strings.Title(name)
    code := "package " + name + "\n\n" +
        "// Service generado\n" +
        "type " + class + "Service struct { repo *" + class + "Repository }\n" +
        "func New" + class + "Service(repo *" + class + "Repository) *" + class + "Service { return &" + class + "Service{repo: repo} }\n"
    path := filepath.Join("src", "modules", name, name+".service.go")
    _ = os.WriteFile(path, []byte(code), 0644)
}

func generateFlatRepository(name string) {
    class := strings.Title(name)
    code := "package " + name + "\n\n" +
        "// Repository generado (stub)\n" +
        "type " + class + "Repository struct{}\n" +
        "func New" + class + "Repository() *" + class + "Repository { return &" + class + "Repository{} }\n"
    path := filepath.Join("src", "modules", name, name+".repository.go")
    _ = os.WriteFile(path, []byte(code), 0644)
}

func generateFlatDTO(name string) {
    class := strings.Title(name)
    code := "package " + name + "\n\n" +
        "// DTOs generados\n" +
        "type " + class + "Response struct { ID string `json:\"id\"` }\n" +
        "type Create" + class + "Request struct { Name string `json:\"name\"` }\n" +
        "type Update" + class + "Request struct { Name string `json:\"name,omitempty\"` }\n"
    path := filepath.Join("src", "modules", name, name+".dto.go")
    _ = os.WriteFile(path, []byte(code), 0644)
}

func generateFlatModel(name string) {
    class := strings.Title(name)
    code := "package " + name + "\n\n" +
        "// Model generado\n" +
        "type " + class + " struct { ID string `json:\"id\"`; Name string `json:\"name\"` }\n"
    path := filepath.Join("src", "modules", name, name+".model.go")
    _ = os.WriteFile(path, []byte(code), 0644)
}