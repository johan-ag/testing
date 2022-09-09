package users

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/johan-ag/testing/internal/users"
	"github.com/mercadolibre/fury_go-core/pkg/web"
	"github.com/mercadolibre/fury_go-platform/pkg/dbtest"
	"github.com/stretchr/testify/require"
)

// grabQueryParam grabs the query param from a request so it can be used in the test.
func grabQueryParam(key, value string, r *http.Request) {
	param := r.URL.Query()
	param.Set(key, value)
	r.URL.RawQuery = param.Encode()
}

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
			db, teardown := dbtest.New(t, "root", "T1m1t1*root")
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

// Unit test of a failed POST operation.
// This example receives data through the body of a request. In this case a User with a name and age.
// We use the field excecuteBeforeTest to perform the expected behavior of the mocked service.
// In the first case our test wont get to use it because it fails before excecuting the service function, so we leave the field as nil.
// The body field is used to store the json request as a string and it will be passed down to the NewRequest method as the body of the request.
// The rest of the fields are the necessary data to perform the service call or to validate the results (e.g. expectedCode).
func TestHandlerSaveFail(t *testing.T) {
	tests := []struct {
		name               string
		excecuteBeforeTest func(ctx context.Context, s *users.MockService, u users.User)
		expectedContext    context.Context
		user               users.User
		body               string
		expectedCode       int
	}{
		{
			name:               "fail to decode",
			excecuteBeforeTest: nil,
			expectedContext:    context.Background(),
			user: users.User{
				ID:   1,
				Name: "Jane",
				Age:  30,
			},
			body:         `{"name":"Jane", "age": 30,}`,
			expectedCode: http.StatusBadGateway,
		},
		{
			name: "service fail",
			excecuteBeforeTest: func(ctx context.Context, s *users.MockService, u users.User) {
				s.EXPECT().Save(ctx, gomock.Eq(u.Name), gomock.Eq(u.Age)).Return(u, errors.New("service fail"))
			},
			expectedContext: context.Background(),
			user: users.User{
				Name: "Jane",
				Age:  30,
			},
			body:         `{"name":"Jane", "age": 30}`,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db, teardown := dbtest.New(t, "root", "T1m1t1*root")
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
			err := handler.Save(r, w)
			// then
			var wErr *web.Error
			b := errors.As(err, &wErr)
			require.True(t, b)
			require.Equal(t, tt.expectedCode, wErr.Status)
		})
	}
}

func TestHandlerFind(t *testing.T) {
	tests := []struct {
		name               string
		excecuteBeforeTest func(ctx context.Context, s *users.MockService, u users.User)
		expectedContext    context.Context
		user               users.User
		userID             string
		expectedCode       int
	}{
		{
			name: "find user successfully",
			excecuteBeforeTest: func(ctx context.Context, s *users.MockService, u users.User) {
				s.EXPECT().Find(ctx, uint(u.ID)).Return(u, nil)
			},
			expectedContext: context.Background(),
			user: users.User{
				ID:   1,
				Name: "Jane",
				Age:  30,
			},
			userID:       "1",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db, teardown := dbtest.New(t, "root", "T1m1t1*root")
			defer teardown()
			db.Load()

			tt.expectedContext = web.WithParams(tt.expectedContext, web.URIParams{"id": tt.userID})
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r = r.WithContext(tt.expectedContext)

			service := users.NewMockService(gomock.NewController(t))
			handler := NewHandler(service)

			if tt.excecuteBeforeTest != nil {
				tt.excecuteBeforeTest(tt.expectedContext, service, tt.user)
			}
			// when
			handler.Find(w, r)
			// then
			require.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestHandlerFindFail(t *testing.T) {
	tests := []struct {
		name               string
		excecuteBeforeTest func(ctx context.Context, s *users.MockService, u users.User)
		expectedContext    context.Context
		user               users.User
		userID             string
		expectedCode       int
	}{
		{
			name:               "empty id",
			excecuteBeforeTest: nil,
			expectedContext:    context.Background(),
			user:               users.User{},
			expectedCode:       http.StatusBadRequest,
		},
		{
			name:               "user not found",
			excecuteBeforeTest: nil,
			expectedContext:    context.Background(),
			user:               users.User{},
			userID:             "7",
			expectedCode:       http.StatusNotFound,
		},
		{
			name: "service fail",
			excecuteBeforeTest: func(ctx context.Context, s *users.MockService, u users.User) {
				s.EXPECT().Find(ctx, uint(u.ID)).Return(u, errors.New("service fail"))
			},
			expectedContext: context.Background(),
			user: users.User{
				ID:   1,
				Name: "Jane",
				Age:  30,
			},
			userID:       "1",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db, teardown := dbtest.New(t, "root", "T1m1t1*root")
			defer teardown()
			db.Load()

			tt.expectedContext = web.WithParams(tt.expectedContext, web.URIParams{"id": tt.userID})
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r = r.WithContext(tt.expectedContext)

			service := users.NewMockService(gomock.NewController(t))
			handler := NewHandler(service)

			if tt.excecuteBeforeTest != nil {
				tt.excecuteBeforeTest(tt.expectedContext, service, tt.user)
			}
			// when
			err := handler.Find(w, r)
			// then
			var wErr *web.Error
			b := errors.As(err, &wErr)
			require.True(t, b)
			require.Equal(t, tt.expectedCode, wErr.Status)
		})
	}
}

func TestHandlerFindByNameAndAge(t *testing.T) {
	tests := []struct {
		name               string
		excecuteBeforeTest func(ctx context.Context, s *users.MockService, u users.User)
		expectedContext    context.Context
		user               users.User
		ageParam           string
		expectedCode       int
	}{
		{
			name: "find user successfully",
			excecuteBeforeTest: func(ctx context.Context, s *users.MockService, u users.User) {
				s.EXPECT().FindByParams(ctx, gomock.Eq(u.Name), gomock.Eq(int32(u.Age))).Return([]users.User{u}, nil)
			},
			expectedContext: context.Background(),
			user: users.User{
				ID:   1,
				Name: "Jane",
				Age:  30,
			},
			ageParam:     "30",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db, teardown := dbtest.New(t, "root", "T1m1t1*root")
			defer teardown()
			db.Load()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			// Grab query params from the request
			grabQueryParam("name", tt.user.Name, r)
			grabQueryParam("age", tt.ageParam, r)

			service := users.NewMockService(gomock.NewController(t))
			handler := NewHandler(service)

			if tt.excecuteBeforeTest != nil {
				tt.excecuteBeforeTest(tt.expectedContext, service, tt.user)
			}
			// when
			handler.FindByNameAndAge(w, r)
			// then
			require.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestHandlerFindByNameAndAgeFail(t *testing.T) {
	tests := []struct {
		name               string
		excecuteBeforeTest func(ctx context.Context, s *users.MockService, u users.User)
		expectedContext    context.Context
		user               users.User
		ageParam           string
		expectedCode       int
	}{
		{
			name:               "invalid age param",
			excecuteBeforeTest: nil,
			expectedContext:    context.Background(),
			user: users.User{
				ID:   1,
				Name: "Jane",
				Age:  30,
			},
			ageParam:     "AA",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "service fail",
			excecuteBeforeTest: func(ctx context.Context, s *users.MockService, u users.User) {
				s.EXPECT().FindByParams(ctx, u.Name, int32(u.Age)).Return([]users.User{u}, errors.New("service fail"))
			},
			expectedContext: context.Background(),
			user: users.User{
				ID:   1,
				Name: "Jane",
				Age:  30,
			},
			ageParam:     "30",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db, teardown := dbtest.New(t, "root", "T1m1t1*root")
			defer teardown()
			db.Load()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			// Grab query params from the request
			grabQueryParam("name", tt.user.Name, r)
			grabQueryParam("age", tt.ageParam, r)

			service := users.NewMockService(gomock.NewController(t))
			handler := NewHandler(service)

			if tt.excecuteBeforeTest != nil {
				tt.excecuteBeforeTest(tt.expectedContext, service, tt.user)
			}
			// when
			err := handler.FindByNameAndAge(w, r)
			// then
			var wErr *web.Error
			b := errors.As(err, &wErr)
			require.True(t, b)
			require.Equal(t, tt.expectedCode, wErr.Status)
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
