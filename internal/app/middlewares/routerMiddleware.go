package middlewares

import (
	"context"
	"net/http"
)

type Router struct {
	routes map[string]*http.HandlerFunc
}

func NewRouter(pathsHandlersMap map[string]*http.HandlerFunc) *Router {
	router := &Router{
		routes: pathsHandlersMap,
	}

	return router
}

func (r *Router) ProcessRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			path := req.URL.Path

			var params map[string]string
			var ok bool

			for route, handler := range r.routes {
				if params, ok = r.parsePath(route, path); ok {
					// Добавляем параметры в контекст
					ctx := context.WithValue(req.Context(), "params", params)
					(*handler)(rw, req.WithContext(ctx))

					break
				}
			}

			if !ok {
				http.Error(rw, "Path does not match the pattern", http.StatusNotFound)
				return
			}

			if next != nil {
				next.ServeHTTP(rw, req)
			}
		},
	)
}

// Функция для парсинга пути с использованием двух указателей
func (r *Router) parsePath(pattern, path string) (map[string]string, bool) {
	params := make(map[string]string)
	patternLen := len(pattern)
	pathLen := len(path)
	patternIndex, pathIndex := 0, 0

	for patternIndex < patternLen && pathIndex < pathLen {
		if pattern[patternIndex] == '{' {
			// Начало параметра
			patternIndex++
			paramStart := patternIndex

			// Ищем конец параметра (закрывающую скобку)
			for patternIndex < patternLen && pattern[patternIndex] != '}' {
				patternIndex++
			}
			if patternIndex >= patternLen {
				return nil, false // Некорректный шаблон
			}

			paramName := pattern[paramStart:patternIndex]
			patternIndex++

			// Ищем конец значения параметра (следующий '/' или конец строки)
			paramValueStart := pathIndex
			for pathIndex < pathLen && (patternIndex >= patternLen || path[pathIndex] != '/') {
				pathIndex++
			}

			params[paramName] = path[paramValueStart:pathIndex]
		} else if pattern[patternIndex] == path[pathIndex] {
			// Символы совпадают, продолжаем
			patternIndex++
			pathIndex++
		} else {
			// Символы не совпадают, путь не соответствует шаблону
			return nil, false
		}
	}

	// Проверяем, что оба индекса достигли конца строк
	if patternIndex != patternLen || pathIndex != pathLen {
		return nil, false
	}

	return params, true
}
