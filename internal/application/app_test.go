package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/zakharkaverin1/http_calc_go/pkg/calculation"
)

func TestCalcHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    Request
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Верное",
			requestBody:    Request{Expression: "3 + 2"},
			expectedStatus: http.StatusOK,
			expectedBody:   Result{Result: "5.000000"},
		},
		{
			name:           "Неверное",
			requestBody:    Request{Expression: "3 + a"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   Error{Error: calculation.ErrInvalidExpression.Error()},
		},
		{
			name:           "Деление на ноль",
			requestBody:    Request{Expression: "3 / 0"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   Error{Error: calculation.ErrDivisionByZero.Error()},
		},
		{
			name:           "Пустое",
			requestBody:    Request{Expression: ""},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   Error{Error: calculation.ErrExpressionNotString.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			CalcHandler(w, req)

			res := w.Result()
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			var responseBody interface{}
			if tt.expectedStatus == http.StatusOK {
				responseBody = Result{}
			} else {
				responseBody = Error{}
			}
			json.NewDecoder(res.Body).Decode(&responseBody)

			if !reflect.DeepEqual(responseBody, tt.expectedBody) {
				t.Errorf("expected body %+v, got %+v", tt.expectedBody, responseBody)
			}
		})
	}
}
