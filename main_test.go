package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenOK(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// запрос сформирован корректно и код ответа 200
	status := responseRecorder.Code
	require.Equal(t, http.StatusOK, status)

	// тело ответа не пустое
	body := responseRecorder.Body.String()
	assert.NotEmpty(t, body)
}

func TestMainHandlerWhenWrongCity(t *testing.T) {
	// город, который передаётся в параметре city, не поддерживается
	req := httptest.NewRequest("GET", "/cafe?count=10&city=ryazan", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа
	status := responseRecorder.Code
	assert.Equal(t, http.StatusBadRequest, status)

	expected := `wrong city value`
	body := responseRecorder.Body.String()
	assert.Equal(t, expected, body)
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	// если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")

	assert.Len(t, list, totalCount)
}
