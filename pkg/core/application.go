package core

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Application struct {
	engine *gin.Engine
	config *Config
}

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

func NewApplication(config *Config) *Application {
	app := &Application{
		engine: gin.Default(),
		config: config,
	}

	app.setupMiddleware()
	app.setupRoutes()

	return app
}

func (a *Application) setupMiddleware() {
	a.engine.Use(gin.Logger())
	a.engine.Use(gin.Recovery())

	a.engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

func (a *Application) setupRoutes() {
	// Ruta principal con página de bienvenida
	a.engine.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, getWelcomeHTML())
	})

	// API routes
	api := a.engine.Group("/api/v1")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Go-ney server is running",
			"port":    a.config.Port,
			"version": "1.0.0",
		})
	})
}

func getWelcomeHTML() string {
	return `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go-ney Framework</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            color: #333;
        }

        .container {
            text-align: center;
            background: white;
            padding: 3rem;
            border-radius: 20px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            max-width: 600px;
            margin: 2rem;
        }

        .logo {
            width: 120px;
            height: 120px;
            margin: 0 auto 2rem;
            background: #FFC107;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 3rem;
            font-weight: bold;
            color: #333;
            box-shadow: 0 10px 30px rgba(255, 193, 7, 0.3);
            position: relative;
        }

        .logo::before {
            content: "🐝";
            font-size: 4rem;
        }

        .logo::after {
            content: "GO";
            position: absolute;
            font-size: 1.2rem;
            bottom: 10px;
            right: 15px;
            background: #333;
            color: white;
            padding: 2px 6px;
            border-radius: 4px;
            font-weight: bold;
        }

        h1 {
            font-size: 3rem;
            margin-bottom: 1rem;
            background: linear-gradient(45deg, #667eea, #764ba2);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }

        .subtitle {
            font-size: 1.3rem;
            color: #666;
            margin-bottom: 2rem;
            line-height: 1.6;
        }

        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 1rem;
            margin: 2rem 0;
        }

        .feature {
            background: #f8f9fa;
            padding: 1rem;
            border-radius: 10px;
            border-left: 4px solid #FFC107;
        }

        .feature-title {
            font-weight: bold;
            color: #333;
            margin-bottom: 0.5rem;
        }

        .feature-desc {
            font-size: 0.9rem;
            color: #666;
        }

        .status {
            background: #d4edda;
            color: #155724;
            padding: 1rem;
            border-radius: 10px;
            margin: 2rem 0;
            border: 1px solid #c3e6cb;
        }

        .links {
            display: flex;
            gap: 1rem;
            justify-content: center;
            flex-wrap: wrap;
            margin-top: 2rem;
        }

        .link {
            background: #667eea;
            color: white;
            padding: 0.8rem 1.5rem;
            text-decoration: none;
            border-radius: 8px;
            transition: all 0.3s ease;
            font-weight: 500;
        }

        .link:hover {
            background: #5a6fd8;
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.3);
        }

        .version {
            margin-top: 2rem;
            padding-top: 2rem;
            border-top: 1px solid #eee;
            color: #999;
            font-size: 0.9rem;
        }

        @media (max-width: 600px) {
            .container {
                margin: 1rem;
                padding: 2rem;
            }

            h1 {
                font-size: 2rem;
            }

            .logo {
                width: 100px;
                height: 100px;
            }

            .features {
                grid-template-columns: 1fr;
            }

            .links {
                flex-direction: column;
                align-items: center;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo"></div>

        <h1>Go-ney Framework</h1>
        <p class="subtitle">
            Framework CLI para Go inspirado en NestJS<br>
            <strong>¡Apoyando el ecosistema Go! 🚀</strong>
        </p>

        <div class="status">
            <strong>✅ Servidor funcionando correctamente</strong><br>
            Tu aplicación Go-ney está lista para el desarrollo
        </div>

        <div class="features">
            <div class="feature">
                <div class="feature-title">🏗️ Arquitectura MVC</div>
                <div class="feature-desc">Controllers, Services, Repositories</div>
            </div>
            <div class="feature">
                <div class="feature-title">⚡ CLI Generativa</div>
                <div class="feature-desc">Genera código automáticamente</div>
            </div>
            <div class="feature">
                <div class="feature-title">🔗 Microservicios</div>
                <div class="feature-desc">TCP, NATS y gRPC</div>
            </div>
            <div class="feature">
                <div class="feature-title">🛡️ Guards & Interceptors</div>
                <div class="feature-desc">Seguridad y middleware</div>
            </div>
        </div>

        <div class="links">
            <a href="/api/v1/health" class="link">🩺 Health Check</a>
            <a href="#" onclick="showCommands()" class="link">📋 Comandos CLI</a>
            <a href="https://github.com/tu-usuario/goney" class="link" target="_blank">📖 Documentación</a>
        </div>

        <div class="version">
            Go-ney Framework v1.0.0 | Hecho con ❤️ para la comunidad Go
        </div>
    </div>

    <script>
        function showCommands() {
            alert("🚀 Comandos Go-ney disponibles:\n\n" +
                  "📦 Crear proyecto:\n" +
                  "   goney new mi-proyecto\n\n" +
                  "⚡ Generar CRUD:\n" +
                  "   goney generate crud Usuario\n" +
                  "   goney generate crud Producto --global\n\n" +
                  "🔧 Generar componentes:\n" +
                  "   goney generate controller Usuario\n" +
                  "   goney generate service ProductoService\n" +
                  "   goney generate repository ClienteRepo\n\n" +
                  "🚀 Iniciar servidor:\n" +
                  "   goney start\n\n" +
                  "📚 Más información:\n" +
                  "   goney --help");
        }

        // Animación sutil del logo
        const logo = document.querySelector('.logo');
        setInterval(() => {
            logo.style.transform = 'scale(1.05)';
            setTimeout(() => {
                logo.style.transform = 'scale(1)';
            }, 200);
        }, 3000);
    </script>
</body>
</html>`
}

func (a *Application) Listen(addr string) error {
	fmt.Printf("🚀 Go-ney server starting on %s\n", addr)
	return a.engine.Run(addr)
}

func (a *Application) RegisterController(path string, controller interface{}) {
	// Implementation for registering controllers
}

func (a *Application) Use(middleware gin.HandlerFunc) {
	a.engine.Use(middleware)
}