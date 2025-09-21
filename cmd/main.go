package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "goney",
	Version: "1.0.0",
	Short:   "Go-ney - Framework MVC para Go inspirado en NestJS",
	Long: `Go-ney es un framework CLI inspirado en NestJS para crear aplicaciones Go
con arquitectura MVC modular y soporte para microservicios TCP, NAT y gRPC.`,
}

var newCmd = &cobra.Command{
	Use:   "new [nombre-proyecto]",
	Short: "Crear un nuevo proyecto Go-ney",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		fmt.Printf("Creando proyecto Go-ney: %s\n", projectName)
		createNewProject(projectName)
	},
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g"},
	Short:   "Generar componentes MVC",
}

var controllerCmd = &cobra.Command{
	Use:   "controller [nombre]",
	Short: "Generar un controller",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("Generando controller: %s\n", name)
		generateController(name)
	},
}

var serviceCmd = &cobra.Command{
	Use:   "service [nombre]",
	Short: "Generar un service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("Generando service: %s\n", name)
		generateService(name)
	},
}

var repositoryCmd = &cobra.Command{
	Use:   "repository [nombre]",
	Short: "Generar un repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("Generando repository: %s\n", name)
		generateRepository(name)
	},
}

var microserviceCmd = &cobra.Command{
	Use:   "microservice [tipo] [nombre]",
	Short: "Generar un microservicio (tcp, nats, grpc)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		serviceType := args[0]
		name := args[1]
		fmt.Printf("Generando microservicio %s: %s\n", serviceType, name)
		generateMicroservice(serviceType, name)
	},
}

var guardCmd = &cobra.Command{
	Use:   "guard [nombre]",
	Short: "Generar un guard",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("Generando guard: %s\n", name)
		generateGuard(name)
	},
}

var interceptorCmd = &cobra.Command{
	Use:   "interceptor [nombre]",
	Short: "Generar un interceptor",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("Generando interceptor: %s\n", name)
		generateInterceptor(name)
	},
}

var crudCmd = &cobra.Command{
	Use:   "crud [nombre-modulo]",
	Short: "Generar módulo CRUD completo con tests incluidos",
	Long: `Generar módulo CRUD completo con estructura modular como NestJS.

Ejemplos:
  goney generate crud users                     # Módulo users/ con toda la estructura + tests
  goney generate crud products --global        # Módulo con DTOs y modelos globales
  goney generate crud orders --no-dto          # Módulo sin generar DTO
  goney generate crud clients --no-model       # Módulo sin generar modelo`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName := args[0]
		global, _ := cmd.Flags().GetBool("global")
		noDto, _ := cmd.Flags().GetBool("no-dto")
		noModel, _ := cmd.Flags().GetBool("no-model")

		fmt.Printf("Generando módulo CRUD: %s\n", moduleName)
		generateModuleCRUD(moduleName, global, noDto, noModel)
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Iniciar el servidor de desarrollo Go-ney",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Iniciando servidor de desarrollo Go-ney...")
		startDevServer()
	},
}

func init() {
	// Flags para el comando CRUD
	crudCmd.Flags().Bool("global", false, "Usar DTOs y modelos globales (no genera archivos específicos)")
	crudCmd.Flags().Bool("no-dto", false, "No generar DTO específico")
	crudCmd.Flags().Bool("no-model", false, "No generar modelo específico")

	generateCmd.AddCommand(controllerCmd)
	generateCmd.AddCommand(serviceCmd)
	generateCmd.AddCommand(repositoryCmd)
	generateCmd.AddCommand(microserviceCmd)
	generateCmd.AddCommand(guardCmd)
	generateCmd.AddCommand(interceptorCmd)
	generateCmd.AddCommand(crudCmd)

	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(startCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}