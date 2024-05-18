package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/redbonzai/user-management-api/internal/domain/user"
	"github.com/redbonzai/user-management-api/internal/infrastructure"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *echo.Echo {
	return infrastructure.NewRouter()
}

func TestGetUsers(t *testing.T) {
	e := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users")

	// Call the handler function directly
	handler := e.Router().Routes()[0].Handler
	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCreateUser(t *testing.T) {
	e := setupRouter()
	u := user.User{Name: "Test User", Email: "test.user@example.com"}
	userJSON, _ := json.Marshal(u)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users")

	// Call the handler function directly
	handler := e.Router().Routes()[1].Handler
	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var resp map[string]int
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NotZero(t, resp["id"])
	}
}

func TestUpdateUser(t *testing.T) {
	e := setupRouter()
	u := user.User{Name: "Updated User", Email: "updated.user@example.com"}
	userJSON, _ := json.Marshal(u)

	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Call the handler function directly
	handler := e.Router().Routes()[2].Handler
	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	e := setupRouter()
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Call the handler function directly
	handler := e.Router().Routes()[3].Handler
	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}
