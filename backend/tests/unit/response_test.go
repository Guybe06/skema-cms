package unit_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"skema-api/core/response"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// Crée un contexte Gin factice pour les tests de réponses HTTP.
func newTestContext(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, nil)
	return c, w
}

// Vérifie que OK retourne un statut 200 avec success à true.
func TestResponse_OK(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	response.OK(c, "tout va bien", gin.H{"key": "value"})

	if w.Code != http.StatusOK {
		t.Fatalf("statut attendu 200, obtenu %d", w.Code)
	}
	var body response.Body
	json.Unmarshal(w.Body.Bytes(), &body)
	if !body.Success {
		t.Fatal("success doit être true pour une réponse OK")
	}
}

// Vérifie que Created retourne un statut 201.
func TestResponse_Created(t *testing.T) {
	c, w := newTestContext(http.MethodPost, "/")
	response.Created(c, "créé", nil)

	if w.Code != http.StatusCreated {
		t.Fatalf("statut attendu 201, obtenu %d", w.Code)
	}
}

// Vérifie que BadRequest retourne un statut 400 avec success à false.
func TestResponse_BadRequest(t *testing.T) {
	c, w := newTestContext(http.MethodPost, "/")
	response.BadRequest(c, "données invalides")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("statut attendu 400, obtenu %d", w.Code)
	}
	var body response.ErrorBody
	json.Unmarshal(w.Body.Bytes(), &body)
	if body.Success {
		t.Fatal("success doit être false pour une réponse d'erreur")
	}
}

// Vérifie que Unauthorized retourne un statut 401.
func TestResponse_Unauthorized(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	response.Unauthorized(c, "non autorisé")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("statut attendu 401, obtenu %d", w.Code)
	}
}

// Vérifie que Forbidden retourne un statut 403.
func TestResponse_Forbidden(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	response.Forbidden(c, "accès refusé")

	if w.Code != http.StatusForbidden {
		t.Fatalf("statut attendu 403, obtenu %d", w.Code)
	}
}

// Vérifie que NotFound retourne un statut 404.
func TestResponse_NotFound(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	response.NotFound(c, "introuvable")

	if w.Code != http.StatusNotFound {
		t.Fatalf("statut attendu 404, obtenu %d", w.Code)
	}
}

// Vérifie que NoContent ne retourne aucun corps JSON.
func TestResponse_NoContent(t *testing.T) {
	c, w := newTestContext(http.MethodDelete, "/")
	response.NoContent(c)

	if w.Body.Len() != 0 {
		t.Fatalf("la réponse NoContent ne doit pas avoir de corps, obtenu : %q", w.Body.String())
	}
}

// Vérifie que List retourne les métadonnées de pagination correctement calculées.
func TestResponse_List_Pagination(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	response.List(c, "liste", []string{"a", "b"}, 25, 1, 10)

	if w.Code != http.StatusOK {
		t.Fatalf("statut attendu 200, obtenu %d", w.Code)
	}
	var body response.ListBody
	json.Unmarshal(w.Body.Bytes(), &body)
	if body.Pagination.TotalPages != 3 {
		t.Fatalf("total_pages attendu 3, obtenu %d", body.Pagination.TotalPages)
	}
}
