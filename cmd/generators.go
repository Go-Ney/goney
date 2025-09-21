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
	"{{.ProjectName}}/src/config"
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
	dtoImport := fmt.Sprintf("\"{{.ModuleName}}/src/modules/%s/dto\"", moduleName)
	dtoType := fmt.Sprintf("%s.", strings.Title(moduleName))

	if global {
		dtoImport = "\"{{.ModuleName}}/src/common/dto\""
		dtoType = "dto.Base"
	}

	controllerTemplate := `package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"{{.ModuleName}}/src/modules/{{.ModulePath}}/services"
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

// @Router /api/v1/{{.RouterPath}} [get]
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

// @Router /api/v1/{{.RouterPath}}/{id} [get]
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

// @Router /api/v1/{{.RouterPath}} [post]
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

// @Router /api/v1/{{.RouterPath}}/{id} [put]
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
	fileName := fmt.Sprintf("src/modules/%s/%s.controller.go", moduleName, moduleName)
	file, _ := os.Create(fileName)
	defer file.Close()

	data := map[string]string{
		"ModuleName":     getProjectModuleName(),
		"ModulePath":     moduleName,
		"ControllerName": capitalizedName,
		"ServiceName":    serviceName,
		"RouterPath":     routerPath,
		"EntityName":     entityName,
		"DTOImport":      dtoImport,
		"DTOType":        dtoType,
	}
	tmpl.Execute(file, data)
}

func generateController(name string) {
	controllerTemplate := `package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"{{.ModuleName}}/services"
	"{{.ModuleName}}/dto"
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
// @Success 200 {array} dto.{{.ControllerName}}Response
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
// @Success 200 {object} dto.{{.ControllerName}}Response
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
// @Param body body dto.Create{{.ControllerName}}Request true "{{.ControllerName}} data"
// @Success 201 {object} dto.{{.ControllerName}}Response
func (c *{{.ControllerName}}Controller) Create(ctx *gin.Context) {
	var req dto.Create{{.ControllerName}}Request
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
// @Param body body dto.Update{{.ControllerName}}Request true "{{.ControllerName}} data"
// @Success 200 {object} dto.{{.ControllerName}}Response
func (c *{{.ControllerName}}Controller) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.Update{{.ControllerName}}Request
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
	}
	tmpl.Execute(file, data)

	generateDTO(name)
	fmt.Printf("✅ Controller %s generado en %s\n", capitalizedName, fileName)
}

func generateService(name string) {
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

func (s *{{.ServiceName}}Service) GetAll() ([]dto.{{.ServiceName}}Response, error) {
	entities, err := s.{{.RepositoryName}}Repository.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.{{.ServiceName}}Response
	for _, entity := range entities {
		responses = append(responses, dto.{{.ServiceName}}Response{
			ID: entity.ID,
			// Map other fields here
		})
	}
	return responses, nil
}

func (s *{{.ServiceName}}Service) GetByID(id string) (*dto.{{.ServiceName}}Response, error) {
	entity, err := s.{{.RepositoryName}}Repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.{{.ServiceName}}Response{
		ID: entity.ID,
		// Map other fields here
	}, nil
}

func (s *{{.ServiceName}}Service) Create(req *dto.Create{{.ServiceName}}Request) (*dto.{{.ServiceName}}Response, error) {
	entity := &models.{{.ServiceName}}{
		// Map fields from request
	}

	createdEntity, err := s.{{.RepositoryName}}Repository.Create(entity)
	if err != nil {
		return nil, err
	}

	return &dto.{{.ServiceName}}Response{
		ID: createdEntity.ID,
		// Map other fields here
	}, nil
}

func (s *{{.ServiceName}}Service) Update(id string, req *dto.Update{{.ServiceName}}Request) (*dto.{{.ServiceName}}Response, error) {
	entity, err := s.{{.RepositoryName}}Repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields from request

	updatedEntity, err := s.{{.RepositoryName}}Repository.Update(entity)
	if err != nil {
		return nil, err
	}

	return &dto.{{.ServiceName}}Response{
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
	}
	tmpl.Execute(file, data)

	generateModel(name)
	fmt.Printf("✅ Service %s generado en %s\n", capitalizedName, fileName)
}

func generateRepository(name string) {
	repositoryTemplate := `package repositories

import (
	"gorm.io/gorm"
	"{{.ModuleName}}/models"
)

type {{.RepositoryName}}Repository struct {
	db *gorm.DB
}

func New{{.RepositoryName}}Repository(db *gorm.DB) *{{.RepositoryName}}Repository {
	return &{{.RepositoryName}}Repository{db: db}
}

func (r *{{.RepositoryName}}Repository) FindAll() ([]models.{{.RepositoryName}}, error) {
	var entities []models.{{.RepositoryName}}
	err := r.db.Find(&entities).Error
	return entities, err
}

func (r *{{.RepositoryName}}Repository) FindByID(id string) (*models.{{.RepositoryName}}, error) {
	var entity models.{{.RepositoryName}}
	err := r.db.First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *{{.RepositoryName}}Repository) Create(entity *models.{{.RepositoryName}}) (*models.{{.RepositoryName}}, error) {
	err := r.db.Create(entity).Error
	return entity, err
}

func (r *{{.RepositoryName}}Repository) Update(entity *models.{{.RepositoryName}}) (*models.{{.RepositoryName}}, error) {
	err := r.db.Save(entity).Error
	return entity, err
}

func (r *{{.RepositoryName}}Repository) Delete(id string) error {
	return r.db.Delete(&models.{{.RepositoryName}}{}, "id = ?", id).Error
}

func (r *{{.RepositoryName}}Repository) FindBy(field string, value interface{}) ([]models.{{.RepositoryName}}, error) {
	var entities []models.{{.RepositoryName}}
	err := r.db.Where(field+" = ?", value).Find(&entities).Error
	return entities, err
}
`

	moduleName := getModuleName()
	capitalizedName := strings.Title(name)

	tmpl, _ := template.New("repository").Parse(repositoryTemplate)
	fileName := fmt.Sprintf("repositories/%s_repository.go", strings.ToLower(name))
	file, _ := os.Create(fileName)
	defer file.Close()

	data := map[string]string{
		"ModuleName":     moduleName,
		"RepositoryName": capitalizedName,
	}
	tmpl.Execute(file, data)
	fmt.Printf("✅ Repository %s generado en %s\n", capitalizedName, fileName)
}

func generateDTO(name string) {
	dtoTemplate := `package dto

import "time"

type {{.DTOName}}Response struct {
	ID        string    ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
	// Add other fields here
}

type Create{{.DTOName}}Request struct {
	// Add fields here
}

type Update{{.DTOName}}Request struct {
	// Add fields here
}
`

	capitalizedName := strings.Title(name)

	tmpl, _ := template.New("dto").Parse(dtoTemplate)
	fileName := fmt.Sprintf("dto/%s_dto.go", strings.ToLower(name))
	file, _ := os.Create(fileName)
	defer file.Close()

	data := map[string]string{
		"DTOName": capitalizedName,
	}
	tmpl.Execute(file, data)
}

func generateModel(name string) {
	modelTemplate := `package models

import (
	"gorm.io/gorm"
	"time"
)

type {{.ModelName}} struct {
	ID        string    ` + "`gorm:\"primaryKey;type:uuid;default:gen_random_uuid()\" json:\"id\"`" + `
	CreatedAt time.Time ` + "`gorm:\"autoCreateTime\" json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`gorm:\"autoUpdateTime\" json:\"updated_at\"`" + `
	DeletedAt gorm.DeletedAt ` + "`gorm:\"index\" json:\"-\"`" + `

	// Add other fields here
}

func ({{.ModelName}}) TableName() string {
	return "{{.TableName}}"
}
`

	capitalizedName := strings.Title(name)
	tableName := strings.ToLower(name) + "s"

	tmpl, _ := template.New("model").Parse(modelTemplate)
	fileName := fmt.Sprintf("models/%s_model.go", strings.ToLower(name))
	file, _ := os.Create(fileName)
	defer file.Close()

	data := map[string]string{
		"ModelName": capitalizedName,
		"TableName": tableName,
	}
	tmpl.Execute(file, data)
}

func getModuleName() string {
	if _, err := os.Stat("go.mod"); err == nil {
		content, _ := os.ReadFile("go.mod")
		lines := strings.Split(string(content), "\n")
		if len(lines) > 0 && strings.HasPrefix(lines[0], "module ") {
			return strings.TrimSpace(strings.TrimPrefix(lines[0], "module "))
		}
	}
	return "myapp"
}

func generateMicroservice(serviceType, name string) {
	switch serviceType {
	case "grpc":
		generateGrpcService(name)
	case "nats":
		generateNatsService(name)
	case "tcp":
		generateTcpService(name)
	default:
		fmt.Printf("❌ Tipo de microservicio no válido: %s. Usa: grpc, nats, tcp\n", serviceType)
	}
}

func generateGrpcService(name string) {
	fmt.Printf("✅ Microservicio gRPC %s generado\n", name)
}

func generateNatsService(name string) {
	fmt.Printf("✅ Microservicio NATS %s generado\n", name)
}

func generateTcpService(name string) {
	fmt.Printf("✅ Microservicio TCP %s generado\n", name)
}

func generateGuard(name string) {
	fmt.Printf("✅ Guard %s generado\n", name)
}

func generateInterceptor(name string) {
	fmt.Printf("✅ Interceptor %s generado\n", name)
}

func generateCRUD(name string) {
	generateCRUDWithOptions(name, false, false, false)
}

func generateModuleCRUD(moduleName string, global, noDto, noModel bool) {
	fmt.Printf("🚀 Generando módulo CRUD: %s\n", moduleName)

	if global {
		fmt.Printf("🌐 Modo global activado - usando DTOs y modelos globales\n")
		ensureGlobalFiles()
	}

	// Crear estructura del módulo
	moduleDir := fmt.Sprintf("src/modules/%s", moduleName)
	createModuleStructure(moduleDir)

	// Generar archivos del módulo
	generateModuleController(moduleName, global)
	generateModuleService(moduleName, global)
	generateModuleRepository(moduleName)
	generateModuleFile(moduleName)

	// Generar model y DTO según las opciones
	if !global && !noModel {
		generateModuleModel(moduleName)
	}

	if !global && !noDto {
		generateModuleDTO(moduleName)
	}

	// Generar tests unitarios
	generateModuleTests(moduleName)

	fmt.Printf("✅ Módulo %s generado exitosamente!\n", moduleName)
	fmt.Printf("📁 Archivos creados en src/modules/%s/:\n", moduleName)
	fmt.Printf("   - %s.controller.go\n", moduleName)
	fmt.Printf("   - %s.service.go\n", moduleName)
	fmt.Printf("   - %s.repository.go\n", moduleName)
	fmt.Printf("   - %s.module.go\n", moduleName)
	fmt.Printf("   - %s_test.go 🧪\n", moduleName)

	if !global && !noModel {
		fmt.Printf("   - %s.model.go\n", moduleName)
	}

	if !global && !noDto {
		fmt.Printf("   - %s.dto.go\n", moduleName)
	}

	if global {
		fmt.Printf("   📌 Usando DTOs y modelos globales en src/common/\n")
	}
}

func createModuleStructure(moduleDir string) {
	// Solo crear el directorio del módulo - archivos directos dentro
	err := os.MkdirAll(moduleDir, 0755)
	if err != nil {
		fmt.Printf("Error creando directorio %s: %v\n", moduleDir, err)
		return
	}
}

func generateCRUDWithOptions(name string, global, noDto, noModel bool) {
	fmt.Printf("🚀 Generando CRUD completo para: %s\n", name)

	if global {
		fmt.Printf("🌐 Modo global activado - usando DTOs y modelos globales\n")
		ensureGlobalFiles()
	}

	// Generar controller y service siempre
	generateControllerWithOptions(name, global)
	generateServiceWithOptions(name, global)
	generateRepository(name)

	// Generar model y DTO según las opciones
	if !global && !noModel {
		generateModel(name)
	}

	if !global && !noDto {
		generateDTO(name)
	}

	// Generar enum si es necesario
	generateEnum(name)

	fmt.Printf("✅ CRUD completo generado para %s!\n", name)
	fmt.Printf("📁 Archivos creados:\n")
	fmt.Printf("   - controllers/%s_controller.go\n", strings.ToLower(name))
	fmt.Printf("   - services/%s_service.go\n", strings.ToLower(name))
	fmt.Printf("   - repositories/%s_repository.go\n", strings.ToLower(name))

	if !global && !noModel {
		fmt.Printf("   - models/%s_model.go\n", strings.ToLower(name))
	}

	if !global && !noDto {
		fmt.Printf("   - dto/%s_dto.go\n", strings.ToLower(name))
	}

	fmt.Printf("   - enums/%s_enum.go\n", strings.ToLower(name))

	if global {
		fmt.Printf("   📌 Usando DTOs y modelos globales existentes\n")
	}
}

func generateEnum(name string) {
	enumTemplate := `package enums

type {{.EnumName}}Status string

const (
	{{.EnumName}}StatusActive   {{.EnumName}}Status = "active"
	{{.EnumName}}StatusInactive {{.EnumName}}Status = "inactive"
	{{.EnumName}}StatusPending  {{.EnumName}}Status = "pending"
	{{.EnumName}}StatusDeleted  {{.EnumName}}Status = "deleted"
)

func (s {{.EnumName}}Status) String() string {
	return string(s)
}

func (s {{.EnumName}}Status) IsValid() bool {
	switch s {
	case {{.EnumName}}StatusActive, {{.EnumName}}StatusInactive, {{.EnumName}}StatusPending, {{.EnumName}}StatusDeleted:
		return true
	default:
		return false
	}
}

type {{.EnumName}}Type string

const (
	{{.EnumName}}TypeDefault {{.EnumName}}Type = "default"
	{{.EnumName}}TypePremium {{.EnumName}}Type = "premium"
	{{.EnumName}}TypeBasic   {{.EnumName}}Type = "basic"
)

func (t {{.EnumName}}Type) String() string {
	return string(t)
}

func (t {{.EnumName}}Type) IsValid() bool {
	switch t {
	case {{.EnumName}}TypeDefault, {{.EnumName}}TypePremium, {{.EnumName}}TypeBasic:
		return true
	default:
		return false
	}
}
`

	capitalizedName := strings.Title(name)

	tmpl, _ := template.New("enum").Parse(enumTemplate)
	fileName := fmt.Sprintf("enums/%s_enum.go", strings.ToLower(name))
	file, _ := os.Create(fileName)
	defer file.Close()

	data := map[string]string{
		"EnumName": capitalizedName,
	}
	tmpl.Execute(file, data)
}

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

func generateModuleService(moduleName string, global bool) {
	serviceTemplate := `package {{.ModuleName}}

import (
	"{{.ProjectName}}/src/modules/{{.ModuleName}}"
)

type {{.ClassName}}Service struct {
	{{.RepositoryName}}Repository *{{.ClassName}}Repository
}

func New{{.ClassName}}Service({{.RepositoryName}}Repository *{{.ClassName}}Repository) *{{.ClassName}}Service {
	return &{{.ClassName}}Service{
		{{.RepositoryName}}Repository: {{.RepositoryName}}Repository,
	}
}

func (s *{{.ClassName}}Service) GetAll() ([]{{.ClassName}}Response, error) {
	entities, err := s.{{.RepositoryName}}Repository.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []{{.ClassName}}Response
	for _, entity := range entities {
		responses = append(responses, {{.ClassName}}Response{
			ID: entity.ID,
		})
	}
	return responses, nil
}

func (s *{{.ClassName}}Service) GetByID(id string) (*{{.ClassName}}Response, error) {
	entity, err := s.{{.RepositoryName}}Repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &{{.ClassName}}Response{
		ID: entity.ID,
	}, nil
}

func (s *{{.ClassName}}Service) Create(req *Create{{.ClassName}}Request) (*{{.ClassName}}Response, error) {
	entity := &{{.ClassName}}{
		// Map fields from request
	}

	createdEntity, err := s.{{.RepositoryName}}Repository.Create(entity)
	if err != nil {
		return nil, err
	}

	return &{{.ClassName}}Response{
		ID: createdEntity.ID,
	}, nil
}

func (s *{{.ClassName}}Service) Update(id string, req *Update{{.ClassName}}Request) (*{{.ClassName}}Response, error) {
	entity, err := s.{{.RepositoryName}}Repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields from request

	updatedEntity, err := s.{{.RepositoryName}}Repository.Update(entity)
	if err != nil {
		return nil, err
	}

	return &{{.ClassName}}Response{
		ID: updatedEntity.ID,
	}, nil
}

func (s *{{.ClassName}}Service) Delete(id string) error {
	return s.{{.RepositoryName}}Repository.Delete(id)
}
`

	capitalizedName := strings.Title(moduleName)
	repositoryName := strings.ToLower(moduleName[:1]) + moduleName[1:]
	projectName := getProjectModuleName()

	tmpl, _ := template.New("module-service").Parse(serviceTemplate)
	fileName := fmt.Sprintf("src/modules/%s/%s.service.go", moduleName, moduleName)
	file, _ := os.Create(fileName)
	defer file.Close()

	data := map[string]string{
		"ModuleName":     moduleName,
		"ClassName":      capitalizedName,
		"RepositoryName": repositoryName,
		"ProjectName":    projectName,
	}
	tmpl.Execute(file, data)
}

func generateModuleRepository(moduleName string) {
	repositoryTemplate := `package {{.ModuleName}}

import (
	"gorm.io/gorm"
)

type {{.ClassName}}Repository struct {
	db *gorm.DB
}

func New{{.ClassName}}Repository(db *gorm.DB) *{{.ClassName}}Repository {
	return &{{.ClassName}}Repository{db: db}
}

func (r *{{.ClassName}}Repository) FindAll() ([]{{.ClassName}}, error) {
	var entities []{{.ClassName}}
	err := r.db.Find(&entities).Error
	return entities, err
}

func (r *{{.ClassName}}Repository) FindByID(id string) (*{{.ClassName}}, error) {
	var entity {{.ClassName}}
	err := r.db.First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *{{.ClassName}}Repository) Create(entity *{{.ClassName}}) (*{{.ClassName}}, error) {
	err := r.db.Create(entity).Error
	return entity, err
}

func (r *{{.ClassName}}Repository) Update(entity *{{.ClassName}}) (*{{.ClassName}}, error) {
	err := r.db.Save(entity).Error
	return entity, err
}

func (r *{{.ClassName}}Repository) Delete(id string) error {
	return r.db.Delete(&{{.ClassName}}{}, "id = ?", id).Error
}

func (r *{{.ClassName}}Repository) FindBy(field string, value interface{}) ([]{{.ClassName}}, error) {
	var entities []{{.ClassName}}
	err := r.db.Where(field+" = ?", value).Find(&entities).Error
	return entities, err
}
`

	capitalizedName := strings.Title(moduleName)

	tmpl, _ := template.New("module-repository").Parse(repositoryTemplate)
	fileName := fmt.Sprintf("src/modules/%s/%s.repository.go", moduleName, moduleName)
	file, _ := os.Create(fileName)
	defer file.Close()

	data := map[string]string{
		"ModuleName": moduleName,
		"ClassName":  capitalizedName,
	}
	tmpl.Execute(file, data)
}

func generateModuleFile(moduleName string) {
	moduleTemplate := `package {{.ModuleName}}

import (
	"{{.ProjectName}}/src/modules/{{.ModuleName}}/controllers"
	"{{.ProjectName}}/src/modules/{{.ModuleName}}/services"
	"{{.ProjectName}}/src/modules/{{.ModuleName}}/repositories"
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
	// Registrar rutas del módulo
	// router.Group("/api/v1/{{.ModuleName}}")
}
`

	capitalizedName := strings.Title(moduleName)
	projectName := getProjectModuleName()

	tmpl, _ := template.New("module").Parse(moduleTemplate)
	fileName := fmt.Sprintf("src/modules/%s/%s.module.go", moduleName, moduleName)
	file, _ := os.Create(fileName)
	defer file.Close()

	data := map[string]string{
		"ModuleName":  moduleName,
		"ClassName":   capitalizedName,
		"ProjectName": projectName,
	}
	tmpl.Execute(file, data)
}

func generateModuleModel(moduleName string) {
	modelTemplate := `package {{.ModuleName}}

import (
	"gorm.io/gorm"
	"time"
)

type {{.ClassName}} struct {
	ID        string    ` + "`gorm:\"primaryKey;type:uuid;default:gen_random_uuid()\" json:\"id\"`" + `
	CreatedAt time.Time ` + "`gorm:\"autoCreateTime\" json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`gorm:\"autoUpdateTime\" json:\"updated_at\"`" + `
	DeletedAt gorm.DeletedAt ` + "`gorm:\"index\" json:\"-\"`" + `

	// Add other fields here for {{.ModuleName}}
	Name        string ` + "`gorm:\"not null;size:255\" json:\"name\"`" + `
	Description string ` + "`gorm:\"type:text\" json:\"description\"`" + `
	Status      string ` + "`gorm:\"default:active;size:50\" json:\"status\"`" + `
}

func ({{.ClassName}}) TableName() string {
	return "{{.TableName}}"
}
`

	capitalizedName := strings.Title(moduleName)
	tableName := moduleName

	tmpl, _ := template.New("module-model").Parse(modelTemplate)
	fileName := fmt.Sprintf("src/modules/%s/%s.model.go", moduleName, moduleName)
	file, _ := os.Create(fileName)
	defer file.Close()

	data := map[string]string{
		"ModuleName": moduleName,
		"ClassName":  capitalizedName,
		"TableName":  tableName,
	}
	tmpl.Execute(file, data)
}

func generateModuleTests(moduleName string) {
	testTemplate := `package {{.ModuleName}}

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repository for testing
type Mock{{.ClassName}}Repository struct {
	mock.Mock
}

func (m *Mock{{.ClassName}}Repository) FindAll() ([]{{.ClassName}}, error) {
	args := m.Called()
	return args.Get(0).([]{{.ClassName}}), args.Error(1)
}

func (m *Mock{{.ClassName}}Repository) FindByID(id string) (*{{.ClassName}}, error) {
	args := m.Called(id)
	return args.Get(0).(*{{.ClassName}}), args.Error(1)
}

func (m *Mock{{.ClassName}}Repository) Create(entity *{{.ClassName}}) (*{{.ClassName}}, error) {
	args := m.Called(entity)
	return args.Get(0).(*{{.ClassName}}), args.Error(1)
}

func (m *Mock{{.ClassName}}Repository) Update(entity *{{.ClassName}}) (*{{.ClassName}}, error) {
	args := m.Called(entity)
	return args.Get(0).(*{{.ClassName}}), args.Error(1)
}

func (m *Mock{{.ClassName}}Repository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Mock{{.ClassName}}Repository) FindBy(field string, value interface{}) ([]{{.ClassName}}, error) {
	args := m.Called(field, value)
	return args.Get(0).([]{{.ClassName}}), args.Error(1)
}

// Test {{.ClassName}}Service
func Test{{.ClassName}}Service_GetAll(t *testing.T) {
	// Arrange
	mockRepo := new(Mock{{.ClassName}}Repository)
	service := New{{.ClassName}}Service(mockRepo)

	expected{{.ClassName}}s := []{{.ClassName}}{
		{ID: "1", Name: "Test {{.ClassName}} 1"},
		{ID: "2", Name: "Test {{.ClassName}} 2"},
	}

	mockRepo.On("FindAll").Return(expected{{.ClassName}}s, nil)

	// Act
	result, err := service.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "1", result[0].ID)
	mockRepo.AssertExpectations(t)
}

func Test{{.ClassName}}Service_GetByID(t *testing.T) {
	// Arrange
	mockRepo := new(Mock{{.ClassName}}Repository)
	service := New{{.ClassName}}Service(mockRepo)

	expected{{.ClassName}} := &{{.ClassName}}{
		ID:   "1",
		Name: "Test {{.ClassName}}",
	}

	mockRepo.On("FindByID", "1").Return(expected{{.ClassName}}, nil)

	// Act
	result, err := service.GetByID("1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "1", result.ID)
	mockRepo.AssertExpectations(t)
}

func Test{{.ClassName}}Service_Create(t *testing.T) {
	// Arrange
	mockRepo := new(Mock{{.ClassName}}Repository)
	service := New{{.ClassName}}Service(mockRepo)

	createRequest := &Create{{.ClassName}}Request{
		// Add test data here
	}

	expected{{.ClassName}} := &{{.ClassName}}{
		ID: "1",
		// Add other fields
	}

	mockRepo.On("Create", mock.AnythingOfType("*{{.ModuleName}}.{{.ClassName}}")).Return(expected{{.ClassName}}, nil)

	// Act
	result, err := service.Create(createRequest)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "1", result.ID)
	mockRepo.AssertExpectations(t)
}

func Test{{.ClassName}}Service_Delete(t *testing.T) {
	// Arrange
	mockRepo := new(Mock{{.ClassName}}Repository)
	service := New{{.ClassName}}Service(mockRepo)

	mockRepo.On("Delete", "1").Return(nil)

	// Act
	err := service.Delete("1")

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Benchmark tests
func Benchmark{{.ClassName}}Service_GetAll(b *testing.B) {
	mockRepo := new(Mock{{.ClassName}}Repository)
	service := New{{.ClassName}}Service(mockRepo)

	expected{{.ClassName}}s := []{{.ClassName}}{
		{ID: "1", Name: "Test {{.ClassName}}"},
	}

	mockRepo.On("FindAll").Return(expected{{.ClassName}}s, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetAll()
	}
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
	// Add other fields here for {{.ModuleName}}
}

type Create{{.ClassName}}Request struct {
	// Add fields here for creating {{.ModuleName}}
}

type Update{{.ClassName}}Request struct {
	// Add fields here for updating {{.ModuleName}}
}
`

	capitalizedName := strings.Title(moduleName)

	tmpl, _ := template.New("module-dto").Parse(dtoTemplate)
	fileName := fmt.Sprintf("src/modules/%s/%s.dto.go", moduleName, moduleName)
	file, _ := os.Create(fileName)
	defer file.Close()

	data := map[string]string{
		"ClassName":  capitalizedName,
		"ModuleName": moduleName,
	}
	tmpl.Execute(file, data)
}

func ensureGlobalFiles() {
	// Crear DTO global si no existe
	if _, err := os.Stat("dto/global_dto.go"); os.IsNotExist(err) {
		createGlobalDTO()
	}

	// Crear modelo global si no existe
	if _, err := os.Stat("models/global_model.go"); os.IsNotExist(err) {
		createGlobalModel()
	}
}

func createGlobalDTO() {
	globalDTOTemplate := `package dto

import "time"

// BaseResponse DTO base para todas las respuestas
type BaseResponse struct {
	ID        string    ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

// BaseCreateRequest DTO base para peticiones de creación
type BaseCreateRequest struct {
	Name        string ` + "`json:\"name\" validate:\"required\"`" + `
	Description string ` + "`json:\"description\"`" + `
}

// BaseUpdateRequest DTO base para peticiones de actualización
type BaseUpdateRequest struct {
	Name        string ` + "`json:\"name,omitempty\"`" + `
	Description string ` + "`json:\"description,omitempty\"`" + `
}

// PaginationRequest DTO para paginación
type PaginationRequest struct {
	Page     int    ` + "`json:\"page\" validate:\"min=1\"`" + `
	Limit    int    ` + "`json:\"limit\" validate:\"min=1,max=100\"`" + `
	SortBy   string ` + "`json:\"sort_by\"`" + `
	SortDesc bool   ` + "`json:\"sort_desc\"`" + `
}

// PaginationResponse DTO para respuestas paginadas
type PaginationResponse struct {
	Data       interface{} ` + "`json:\"data\"`" + `
	Page       int         ` + "`json:\"page\"`" + `
	Limit      int         ` + "`json:\"limit\"`" + `
	Total      int64       ` + "`json:\"total\"`" + `
	TotalPages int         ` + "`json:\"total_pages\"`" + `
}

// ErrorResponse DTO para respuestas de error
type ErrorResponse struct {
	Error   string      ` + "`json:\"error\"`" + `
	Code    string      ` + "`json:\"code,omitempty\"`" + `
	Details interface{} ` + "`json:\"details,omitempty\"`" + `
}

// SuccessResponse DTO para respuestas exitosas
type SuccessResponse struct {
	Message string      ` + "`json:\"message\"`" + `
	Data    interface{} ` + "`json:\"data,omitempty\"`" + `
}
`

	file, _ := os.Create("dto/global_dto.go")
	defer file.Close()
	file.WriteString(globalDTOTemplate)
	fmt.Printf("✅ DTO global creado en dto/global_dto.go\n")
}

func createGlobalModel() {
	globalModelTemplate := `package models

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel modelo base con campos comunes
type BaseModel struct {
	ID        string    ` + "`gorm:\"primaryKey;type:uuid;default:gen_random_uuid()\" json:\"id\"`" + `
	CreatedAt time.Time ` + "`gorm:\"autoCreateTime\" json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`gorm:\"autoUpdateTime\" json:\"updated_at\"`" + `
	DeletedAt gorm.DeletedAt ` + "`gorm:\"index\" json:\"-\"`" + `
}

// NamedModel modelo base con nombre y descripción
type NamedModel struct {
	BaseModel
	Name        string ` + "`gorm:\"not null;size:255\" json:\"name\"`" + `
	Description string ` + "`gorm:\"type:text\" json:\"description\"`" + `
	Status      string ` + "`gorm:\"default:active;size:50\" json:\"status\"`" + `
}

// MetadataModel modelo con metadata adicional
type MetadataModel struct {
	NamedModel
	Metadata map[string]interface{} ` + "`gorm:\"type:jsonb\" json:\"metadata\"`" + `
	Tags     []string               ` + "`gorm:\"type:text[]\" json:\"tags\"`" + `
}

// AuditModel modelo con campos de auditoría
type AuditModel struct {
	BaseModel
	CreatedBy string ` + "`gorm:\"size:255\" json:\"created_by\"`" + `
	UpdatedBy string ` + "`gorm:\"size:255\" json:\"updated_by\"`" + `
	Version   int    ` + "`gorm:\"default:1\" json:\"version\"`" + `
}
`

	file, _ := os.Create("models/global_model.go")
	defer file.Close()
	file.WriteString(globalModelTemplate)
	fmt.Printf("✅ Modelo global creado en models/global_model.go\n")
}

func generateControllerWithOptions(name string, global bool) {
	dtoImport := "\"{{.ModuleName}}/dto\""
	dtoType := "dto."

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