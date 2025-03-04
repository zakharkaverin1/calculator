package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zakharkaverin1/calculator/internal/application"
)

func TestCalculateHandler(t *testing.T) {
	o := application.NewOrchestrator()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o.CreateHandler(w, r)
	})
	reqBody := `{"expression": "(1+2)*3"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusCreated, rr.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}

	id, ok := resp["id"]
	if !ok || id == "" {
		t.Errorf("Невалидный айди, получено: %v", resp)
	}
}
