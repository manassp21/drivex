package api

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
	"github.com/manassp21/drivex/backend/routes"
)

func TestRegisterHandler_RejectsShortPassword(t *testing.T) {
    body, _ := json.Marshal(map[string]string{
        "email":    "test@example.com",
        "password": "short",
    })
    req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    routes.RegisterUser(w, req)

    if w.Code != http.StatusBadRequest {
        t.Errorf("expected 400, got %d", w.Code)
    }
}