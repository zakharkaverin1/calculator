package tests

import (
	"testing"

	"github.com/zakharkaverin1/calculator/pkg/calculation"
)

func TestCompute(t *testing.T) {
	tests := []struct {
		op        string  
		a, b      float64 
		expected  float64 
		shouldErr bool    
	}{
		{"+", 1, 1, 2, false}, 
		{"-", 3, 3, 0, false}, 
		{"*", 3, 3, 9, false}, 
		{"/", 10, 2, 5, false},
		{"^", 2, 3, 0, true},  
	}
	for _, tc := range tests {
		result, err := calculation.Compute(tc.op, tc.a, tc.b)
		if tc.shouldErr {
			if err == nil {
				t.Errorf("Ожидалась ошибка для операции %s с операндами %f и %f", tc.op, tc.a, tc.b)
			}
		} else {
			if err != nil {
				t.Errorf("Неожиданная ошибка для операции %s с операндами %f и %f: %v", tc.op, tc.a, tc.b, err)
			}
			if result != tc.expected {
				t.Errorf("Compute(%s, %f, %f) = %f; ожидалось %f", tc.op, tc.a, tc.b, result, tc.expected)
			}
		}
	}
}
