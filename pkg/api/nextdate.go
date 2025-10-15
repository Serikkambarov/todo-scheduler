package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

// NextDate вычисляет следующую дату задачи по правилам d<N> или y
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга даты: %w", err)
	}

	if repeat == "" {
		return "", fmt.Errorf("repeat пустой")
	}

	parts := strings.Split(repeat, " ")
	mainRule := parts[0]

	// --- d <число>
	if mainRule == "d" && len(parts) == 2 {
		n, err := strconv.Atoi(parts[1])
		if err != nil || n <= 0 || n > 400 {
			return "", fmt.Errorf("repeat: неверный интервал %v", parts[1])
		}
		next := start
		for {
			next = next.AddDate(0, 0, n)
			if afterNow(next, now) {
				break
			}
		}
		return next.Format(dateFormat), nil
	}

	// --- y (ежегодно)
	if mainRule == "y" && len(parts) == 1 {
		next := start
		for {
			year := next.Year() + 1
			month := next.Month()
			day := next.Day()
			// Если 29 февраля и следующий год невисокосный → 1 марта
			if month == time.February && day == 29 && !isLeap(year) {
				next = time.Date(year, time.March, 1, 0, 0, 0, 0, time.UTC)
			} else {
				next = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
			}
			if afterNow(next, now) {
				break
			}
		}
		return next.Format(dateFormat), nil
	}

	return "", fmt.Errorf("repeat в неподдерживаемом формате: %s", repeat)
}

// afterNow возвращает true, если date > now, игнорируя время суток
func afterNow(date, now time.Time) bool {
	return date.Format(dateFormat) > now.Format(dateFormat)
}

// isLeap проверяет, високосный ли год
func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// nextDateHandler обрабатывает GET /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
    // Проверяем метод
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    nowStr := r.FormValue("now")
    dstart := r.FormValue("date")
    repeat := r.FormValue("repeat")

    if dstart == "" || repeat == "" {
        http.Error(w, "missing date or repeat parameter", http.StatusBadRequest)
        return
    }

    var now time.Time
    var err error
    if nowStr == "" {
        now = time.Now()
    } else {
        now, err = time.Parse(dateFormat, nowStr)
        if err != nil {
            http.Error(w, "invalid now date", http.StatusBadRequest)
            return
        }
    }

    next, err := NextDate(now, dstart, repeat)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Write([]byte(next))
}

