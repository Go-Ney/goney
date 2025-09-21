package guards

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type Interceptor interface {
	Before(ctx *gin.Context) error
	After(ctx *gin.Context, response interface{}) error
}

type LoggingInterceptor struct {
	logger *log.Logger
}

func NewLoggingInterceptor(logger *log.Logger) *LoggingInterceptor {
	return &LoggingInterceptor{logger: logger}
}

func (i *LoggingInterceptor) Before(ctx *gin.Context) error {
	start := time.Now()
	ctx.Set("request_start_time", start)

	i.logger.Printf("[REQUEST] %s %s - IP: %s",
		ctx.Request.Method,
		ctx.Request.URL.Path,
		ctx.ClientIP())

	return nil
}

func (i *LoggingInterceptor) After(ctx *gin.Context, response interface{}) error {
	startTime, exists := ctx.Get("request_start_time")
	if !exists {
		return nil
	}

	duration := time.Since(startTime.(time.Time))

	i.logger.Printf("[RESPONSE] %s %s - Status: %d - Duration: %v",
		ctx.Request.Method,
		ctx.Request.URL.Path,
		ctx.Writer.Status(),
		duration)

	return nil
}

type ValidationInterceptor struct{}

func NewValidationInterceptor() *ValidationInterceptor {
	return &ValidationInterceptor{}
}

func (i *ValidationInterceptor) Before(ctx *gin.Context) error {
	if ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			return err
		}

		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			return err
		}

		ctx.Set("request_body", data)
	}
	return nil
}

func (i *ValidationInterceptor) After(ctx *gin.Context, response interface{}) error {
	return nil
}

type TransformInterceptor struct {
	transformers map[string]func(interface{}) interface{}
}

func NewTransformInterceptor() *TransformInterceptor {
	return &TransformInterceptor{
		transformers: make(map[string]func(interface{}) interface{}),
	}
}

func (i *TransformInterceptor) AddTransformer(path string, transformer func(interface{}) interface{}) {
	i.transformers[path] = transformer
}

func (i *TransformInterceptor) Before(ctx *gin.Context) error {
	return nil
}

func (i *TransformInterceptor) After(ctx *gin.Context, response interface{}) error {
	transformer, exists := i.transformers[ctx.Request.URL.Path]
	if exists {
		transformed := transformer(response)
		ctx.Set("transformed_response", transformed)
	}
	return nil
}

func InterceptorMiddleware(interceptors ...Interceptor) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, interceptor := range interceptors {
			if err := interceptor.Before(ctx); err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				ctx.Abort()
				return
			}
		}

		ctx.Next()

		for _, interceptor := range interceptors {
			interceptor.After(ctx, nil)
		}
	}
}

type CacheInterceptor struct {
	cache map[string]interface{}
	ttl   map[string]time.Time
}

func NewCacheInterceptor() *CacheInterceptor {
	return &CacheInterceptor{
		cache: make(map[string]interface{}),
		ttl:   make(map[string]time.Time),
	}
}

func (i *CacheInterceptor) Before(ctx *gin.Context) error {
	if ctx.Request.Method == "GET" {
		key := ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery

		if cachedResponse, exists := i.cache[key]; exists {
			if time.Now().Before(i.ttl[key]) {
				ctx.JSON(200, cachedResponse)
				ctx.Abort()
				return nil
			} else {
				delete(i.cache, key)
				delete(i.ttl, key)
			}
		}
	}
	return nil
}

func (i *CacheInterceptor) After(ctx *gin.Context, response interface{}) error {
	if ctx.Request.Method == "GET" && ctx.Writer.Status() == 200 {
		key := ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery
		i.cache[key] = response
		i.ttl[key] = time.Now().Add(5 * time.Minute)
	}
	return nil
}