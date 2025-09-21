package core

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
		welcomeHTML := `<!DOCTYPE html>
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
        <div class="version">Go-ney v1.0.1 | Puerto: <no value></div>
    </div>
</body>
</html>`
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
