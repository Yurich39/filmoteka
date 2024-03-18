package sort

import (
	"context"
	"net/http"
	"strings"
)

const (
	ASC                       = "ASC"
	DESC                      = "DESC"
	SortOptionsContextKey Str = "sort_options"
)

type Str string

// The following Middleware injects sorting options into request context.
// In case options can't be found, we stop here and return error response.

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Так как параметров сортировки может быть несколько, получим их списком
		sortBy := r.URL.Query()["sort_by"] // Получаем список значений

		// Порядок сортировки ASC или DESC
		sortOrder := r.URL.Query()["sort_order"] // Получаем список значений

		// Если сортировка не требуется
		if sortBy == nil {
			next.ServeHTTP(w, r)
			return
		}

		// Если порядок сортировки не указан, то создаем список значений DESC
		if sortOrder == nil {

			sortOrder = []string{}
			for i := 0; i < len(sortBy); i++ {
				sortOrder = append(sortOrder, DESC)
			}

		} else {

			// Дополняем список sortOrder
			if len(sortOrder) != len(sortBy) {
				for len(sortOrder) < len(sortBy) {
					sortOrder = append(sortOrder, ASC)
				}
			}

			// Переводим значения в upper case
			for i := range sortOrder {
				sortOrder[i] = strings.ToUpper(sortOrder[i])
			}

			// Проверяем направление сортировки
			for _, val := range sortOrder {

				if val != ASC && val != DESC {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("incorrect sort order"))

					next.ServeHTTP(w, r)
				}

			}

		}
		options := make(map[string]string)
		for i := 0; i < len(sortOrder); i++ {
			options[sortBy[i]] = sortOrder[i]
		}

		// Наполним контекст запроса новой парой ключ/значение
		ctx := context.WithValue(r.Context(), SortOptionsContextKey, options)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
