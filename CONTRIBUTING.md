# Contribuyendo a Go-ney Framework 🚀

¡Gracias por tu interés en contribuir a Go-ney! Este documento te guiará sobre cómo participar en el desarrollo del framework.

## 🔒 Política de Contribuciones

### Repositorio Público con Control de Acceso
- ✅ **Repositorio público**: Cualquiera puede ver y usar el código
- 🛡️ **Ramas protegidas**: Solo maintainers autorizados pueden hacer push directo
- 🔄 **Pull Requests**: Todas las contribuciones deben pasar por PR y revisión

## 👥 Roles y Permisos

### 🔑 **Maintainers (Acceso completo)**
- Pueden hacer push directo a `main` y `develop`
- Aprueban y hacen merge de PRs
- Gestionan releases y versiones
- Administran permisos del repositorio

### 🌟 **Contributors (Contribuciones via PR)**
- Fork del repositorio
- Crear branch para features/fixes
- Enviar Pull Request para revisión
- Seguir guidelines de desarrollo

### 📖 **Users (Solo lectura)**
- Clonar y usar el framework
- Reportar issues y bugs
- Proponer nuevas funcionalidades

## 🚀 Proceso de Contribución

### Para Maintainers Autorizados:

```bash
# 1. Clonar repositorio principal
git clone https://github.com/tu-usuario/goney-framework.git
cd goney-framework

# 2. Crear rama para nueva funcionalidad
git checkout -b feature/nueva-funcionalidad

# 3. Desarrollar y probar
make dev
make test

# 4. Commit y push
git add .
git commit -m "feat: agregar nueva funcionalidad"
git push origin feature/nueva-funcionalidad

# 5. Crear PR hacia develop
# (Revisión opcional pero recomendada)
```

### Para Contributors Externos:

```bash
# 1. Fork del repositorio en GitHub
# 2. Clonar tu fork
git clone https://github.com/tu-usuario/goney-framework.git
cd goney-framework

# 3. Configurar upstream
git remote add upstream https://github.com/usuario-oficial/goney-framework.git

# 4. Crear rama para tu contribución
git checkout -b feature/mi-contribucion

# 5. Desarrollar siguiendo las guidelines
make dev
make test

# 6. Commit con mensaje descriptivo
git add .
git commit -m "feat: descripción de la funcionalidad"

# 7. Push a tu fork
git push origin feature/mi-contribucion

# 8. Crear Pull Request desde GitHub
# Será revisado por maintainers antes del merge
```

## 📋 Guidelines de Desarrollo

### 🔧 **Configuración del Entorno**

```bash
# Instalar dependencias
make dev

# Ejecutar tests
make test

# Verificar linting
golangci-lint run

# Probar CLI
./goney --help
./goney new test-project
```

### 📝 **Estándares de Código**

1. **Go Conventions**: Seguir estándares oficiales de Go
2. **Testing**: Toda nueva funcionalidad debe incluir tests
3. **Documentación**: Actualizar README y docs cuando sea necesario
4. **Commit Messages**: Usar formato conventional commits

```bash
# Ejemplos de commits válidos:
feat: agregar comando para generar middleware
fix: corregir error en generación de DTOs
docs: actualizar documentación de instalación
test: agregar tests para generador de CRUD
refactor: mejorar estructura de templates
```

### 🧪 **Testing Requirements**

- ✅ Tests unitarios para nuevas funcionalidades
- ✅ Tests de integración para CLI
- ✅ Coverage mínimo del 80%
- ✅ Todos los tests deben pasar en CI/CD

## 🛡️ Configuración de Protección de Ramas

### En GitHub Repository Settings:

```yaml
# Protección para rama 'main'
Branch Protection Rules:
  - Require pull request reviews before merging: ✅
  - Required number of reviewers: 2
  - Dismiss stale reviews: ✅
  - Require status checks: ✅
    - CI Tests (Go 1.21, 1.22, 1.23)
    - Linting
    - Build verification
  - Require branches to be up to date: ✅
  - Restrict pushes that create files larger than 100MB: ✅
  - Do not allow bypassing the above settings: ✅

# Protección para rama 'develop'
Branch Protection Rules:
  - Require pull request reviews: ✅
  - Required reviewers: 1
  - Require status checks: ✅
```

## 📦 Release Process

### Para Maintainers:

```bash
# 1. Merge develop → main
git checkout main
git merge develop

# 2. Crear tag de versión
git tag -a v1.0.0 -m "Release v1.0.0: Initial stable release"
git push origin v1.0.0

# 3. Crear GitHub Release
# CI/CD automáticamente construirá binarios para todas las plataformas
```

## 🐛 Reportar Issues

### Template para Issues:

```markdown
## Descripción del Problema
Descripción clara del problema o funcionalidad solicitada.

## Pasos para Reproducir
1. Ejecutar comando...
2. Ver error...

## Comportamiento Esperado
Lo que debería suceder.

## Comportamiento Actual
Lo que está sucediendo.

## Entorno
- OS: [e.g. macOS, Linux, Windows]
- Go Version: [e.g. 1.23]
- Go-ney Version: [e.g. 1.0.0]
```

## 🤝 Código de Conducta

- 🤝 Ser respetuoso con todos los contributors
- 💬 Usar lenguaje inclusivo y profesional
- 🎯 Enfocarse en el código, no en las personas
- 🧠 Estar abierto a diferentes perspectivas
- 📚 Ayudar a nuevos contributors

## 🏆 Reconocimientos

Los contributors serán reconocidos en:
- 📝 CONTRIBUTORS.md
- 🎉 Release notes
- 🌟 GitHub contributors section

## 📞 Contacto

- 📧 **Issues**: Para bugs y feature requests
- 💬 **Discussions**: Para preguntas generales
- 📱 **Email**: maintainers@goney.dev (para temas sensibles)

---

¡Gracias por ayudar a hacer Go-ney mejor! 🚀✨