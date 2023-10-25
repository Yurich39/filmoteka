package sort

import (
	"context"
	"net/http"
	"strings"
)

const (
	ASC                     = "ASC"
	DESC                    = "DESC"
	SortByContextKey    Str = "sort_by"
	SortOrderContextKey Str = "sort_order"
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

		// Если порядок сортировки не указан, то создаем список значений ASC
		if sortOrder == nil {

			sortOrder = []string{}
			for i := 0; i < len(sortBy); i++ {
				sortOrder = append(sortOrder, ASC)
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

		// Наполним контекст запроса новыми парами ключ/значение
		ctx := context.WithValue(r.Context(), SortByContextKey, sortBy)
		ctx = context.WithValue(ctx, SortOrderContextKey, sortOrder)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
