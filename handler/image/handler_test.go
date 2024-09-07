package image

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func setTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *Handler) {
	t.Helper()

	// setting stub db
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while opening stub database connection: %v", err)
	}

	// setting up test gin server
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	stubStore := NewStore(db)
	mockService := NewService(stubStore)
	mockHandler := NewHandler(mockService)
	mockHandler.RegisterRoutes(router.Group(""))

	return db, mock, mockHandler
}

func TestGetImages(t *testing.T) {
	t.Run("No images found", func(t *testing.T) {
		db, mock, handler := setTest(t)
		defer db.Close()

		expOut := `{"error":"no images found","status":false}`

		// filling expected values
		rows := sqlmock.NewRows([]string{"id", "name", "path"})
		mock.ExpectQuery("SELECT id, name, path FROM mst_images").WillReturnRows(rows)

		// calling the handler
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// c := gin.CreateTestContextOnly(w, router)
		c.Request = httptest.NewRequest("GET", "/getImages", nil)
		handler.GetImages(c)

		// asserting the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expOut, w.Body.String())
	})

	t.Run("Images found", func(t *testing.T) {
		db, mock, handler := setTest(t)
		defer db.Close()

		expOut := `{"data":[{"id":1,"name":"test-name","path":"test-title"}],"status":true}`

		// filling expected values
		rows := sqlmock.NewRows([]string{"id", "name", "path"}).AddRow(1, "test-name", "test-title")
		mock.ExpectQuery("SELECT id, name, path FROM mst_images").WillReturnRows(rows)

		// calling the handler
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// c := gin.CreateTestContextOnly(w, router)
		c.Request = httptest.NewRequest("GET", "/getImages", nil)
		handler.GetImages(c)

		// asserting the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expOut, w.Body.String())
	})
}