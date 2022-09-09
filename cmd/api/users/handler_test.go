package users

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/johan-ag/testing/internal/users"
	"github.com/mercadolibre/fury_go-platform/pkg/dbtest"
	"github.com/stretchr/testify/require"
)

// Unit test of successful POST operation.
// This example receives data through the body of a request. In this case a User with a name and age.
// We use the field excecuteBeforeTest to perform the expected behavior of the mocked service.
// In case our test wont get to excecute the service we can leave this field as nil.
// The body field is used to store the json request as a string and it will be passed down to the NewRequest method as the body of the request.
// The rest of the fields are the necessary data to perform the service call or to validate the results (e.g. expectedCode).
func TestHandlerSave(t *testing.T) {
	tests := []struct {
		name               string
		excecuteBeforeTest func(ctx context.Context, s *users.MockService, u users.User)
		expectedContext    context.Context
		user               users.User
		body               string
		expectedCode       int
	}{
		{
			name: "save user successfully",
			excecuteBeforeTest: func(ctx context.Context, s *users.MockService, u users.User) {
				s.EXPECT().Save(ctx, u.Name, u.Age).Return(u, nil)
			},
			expectedContext: context.Background(),
			user: users.User{
				ID:   1,
				Name: "Jane",
				Age:  30,
			},
			body:         `{"name":"Jane", "age": 30}`,
			expectedCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db, teardown := dbtest.New(t, "root", "root")
			defer teardown()
			db.Load()

			w := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))
			r := httptest.NewRecorder()
			service := users.NewMockService(gomock.NewController(t))
			handler := NewHandler(service)

			if tt.excecuteBeforeTest != nil {
				tt.excecuteBeforeTest(tt.expectedContext, service, tt.user)
			}
			// when
			handler.Save(r, w)
			// then
			require.Equal(t, tt.expectedCode, r.Code)
		})
	}
}

/*func TestIntegrationHandlerSave(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tests := []struct {
		name             string
		createRepository func(dBase *sql.DB) users.Repository
		body             string
		expectedCode     int
	}{
		{
			name: "save handler test successful",
			createRepository: func(dBase *sql.DB) users.Repository {
				return users.NewRepository(dBase)
			},
			body:         `{"name":"name", "age": 30}`,
			expectedCode: http.StatusCreated,
		},
		{
			name: "save handler test bad request",
			createRepository: func(dBase *sql.DB) users.Repository {
				return users.NewRepository(dBase)
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

			ctrl := gomock.NewController(t)
			repository := users.NewMockRepository(ctrl)
			service := users.NewService(repository)
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
