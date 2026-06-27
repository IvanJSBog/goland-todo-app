package core_http_response

import (
	"net/http"
)

var (
	StatusCodeUninitialized = -1
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	// Перехватываем вызов ДО того, как он уйдет в стандартную библиотеку.
	// Если статус не был задан явно, значит сейчас Go отправит 200 OK.
	if rw.statusCode == StatusCodeUninitialized {
		rw.statusCode = http.StatusOK // Фиксируем этот факт у себя
	}
	// Передаем данные дальше внутреннему родителю для реальной отправки
	return rw.ResponseWriter.Write(b)
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, statusCode: StatusCodeUninitialized}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) GetStatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUninitialized {
		panic("no status code set")
	}
	return rw.statusCode
}
