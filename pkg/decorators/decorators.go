package decorators

import (
	"reflect"
	"runtime"
	"strings"
)

type DecoratorFunc func(interface{}) interface{}

type DecoratorRegistry struct {
	decorators map[string][]DecoratorFunc
}

func NewDecoratorRegistry() *DecoratorRegistry {
	return &DecoratorRegistry{
		decorators: make(map[string][]DecoratorFunc),
	}
}

func (r *DecoratorRegistry) AddDecorator(target string, decorator DecoratorFunc) {
	r.decorators[target] = append(r.decorators[target], decorator)
}

func (r *DecoratorRegistry) ApplyDecorators(target interface{}) interface{} {
	targetName := getTypeName(target)
	decorators, exists := r.decorators[targetName]
	if !exists {
		return target
	}

	result := target
	for _, decorator := range decorators {
		result = decorator(result)
	}
	return result
}

func getTypeName(obj interface{}) string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

type RouteDecorator struct {
	Method string
	Path   string
	Guards []string
}

func Route(method, path string, guards ...string) RouteDecorator {
	return RouteDecorator{
		Method: method,
		Path:   path,
		Guards: guards,
	}
}

func Get(path string, guards ...string) RouteDecorator {
	return Route("GET", path, guards...)
}

func Post(path string, guards ...string) RouteDecorator {
	return Route("POST", path, guards...)
}

func Put(path string, guards ...string) RouteDecorator {
	return Route("PUT", path, guards...)
}

func Delete(path string, guards ...string) RouteDecorator {
	return Route("DELETE", path, guards...)
}

type ValidateDecorator struct {
	Rules map[string]interface{}
}

func Validate(rules map[string]interface{}) ValidateDecorator {
	return ValidateDecorator{Rules: rules}
}

type CacheDecorator struct {
	TTL     int
	Key     string
	Enabled bool
}

func Cache(ttl int, key string) CacheDecorator {
	return CacheDecorator{
		TTL:     ttl,
		Key:     key,
		Enabled: true,
	}
}

type LogDecorator struct {
	Level   string
	Message string
}

func Log(level, message string) LogDecorator {
	return LogDecorator{
		Level:   level,
		Message: message,
	}
}

type RateLimitDecorator struct {
	Limit  int
	Window int
}

func RateLimit(limit, window int) RateLimitDecorator {
	return RateLimitDecorator{
		Limit:  limit,
		Window: window,
	}
}

type TransactionDecorator struct {
	Isolation string
}

func Transaction(isolation string) TransactionDecorator {
	return TransactionDecorator{Isolation: isolation}
}

type TimeoutDecorator struct {
	Duration int
}

func Timeout(duration int) TimeoutDecorator {
	return TimeoutDecorator{Duration: duration}
}

type RetryDecorator struct {
	MaxAttempts int
	Delay       int
}

func Retry(maxAttempts, delay int) RetryDecorator {
	return RetryDecorator{
		MaxAttempts: maxAttempts,
		Delay:       delay,
	}
}

type SecurityDecorator struct {
	Permissions []string
	Roles       []string
}

func RequirePermissions(permissions ...string) SecurityDecorator {
	return SecurityDecorator{Permissions: permissions}
}

func RequireRoles(roles ...string) SecurityDecorator {
	return SecurityDecorator{Roles: roles}
}

func getFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func extractPackageName(fullName string) string {
	parts := strings.Split(fullName, ".")
	if len(parts) > 1 {
		return parts[len(parts)-2]
	}
	return fullName
}

type ComponentDecorator struct {
	Name        string
	Singleton   bool
	Dependencies []string
}

func Component(name string, singleton bool, dependencies ...string) ComponentDecorator {
	return ComponentDecorator{
		Name:         name,
		Singleton:    singleton,
		Dependencies: dependencies,
	}
}

type ServiceDecorator struct {
	Name      string
	Scope     string
	Transient bool
}

func Service(name, scope string) ServiceDecorator {
	return ServiceDecorator{
		Name:      name,
		Scope:     scope,
		Transient: false,
	}
}

type ControllerDecorator struct {
	BasePath    string
	Version     string
	Middlewares []string
}

func Controller(basePath, version string, middlewares ...string) ControllerDecorator {
	return ControllerDecorator{
		BasePath:    basePath,
		Version:     version,
		Middlewares: middlewares,
	}
}