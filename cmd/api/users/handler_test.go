package users

/*import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/johan-ag/testing/internal/platform/database"
	"github.com/johan-ag/testing/internal/users"
	"github.com/mercadolibre/fury_go-platform/pkg/dbtest"
)

func TestHandlerSave(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tests := []struct {
		name             string
		createRepository func(queries *database.Queries) users.Repository
		body             string
		expectedCode     int
	}{
		{
			name: "save handler test successful",
			createRepository: func(queries *database.Queries) users.Repository {
				return users.NewRepository(queries)
			},
			body:         `{"name":"name", "age": 30}`,
			expectedCode: http.StatusCreated,
		},
		{
			name: "save handler test bad request",
			createRepository: func(queries *database.Queries) users.Repository {
				return users.NewRepository(queries)
			},
			body:         ``,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db, teardown := dbtest.New(t, "root", "root")
			defer teardown()

			db.Load()

			queries := database.New(db.DB)

			ctrl := gomock.NewController(t)

			qkvs := kvs.NewMockQueryableClient(ctrl)

			repository := tt.createRepository(queries)
			service := users.NewService(repository, qkvs)

			handler := NewHandler(service)

			req := httptest.NewRequest("POST", "/api/users?siteId=Soysite", bytes.NewReader([]byte(tt.body)))
			rr := httptest.NewRecorder()

			// when
			handler.Save(rr, req)

			// then
			if rr.Code != tt.expectedCode {
				t.Log(rr.Code)
				t.Log(rr)
				t.Fail()
			}

		})
	}
}*/
