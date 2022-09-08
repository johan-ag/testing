package users

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestServiceSave(t *testing.T) {
	tests := []struct {
		name              string
		executeBeforeTest func(ctx context.Context, r *MockRepository, expectedName string, expectedAge uint)
		expectedContext   context.Context
		expectedName      string
		expectedAge       uint
		withError         bool
		expectedError     error
	}{
		{
			name: "save service test successful",
			executeBeforeTest: func(ctx context.Context, r *MockRepository, expectedName string, expectedAge uint) {
				r.
					EXPECT().
					Save(gomock.Eq(ctx), gomock.Eq(expectedName), gomock.Eq(expectedAge)).
					Return(uint(1), nil)
			},
			expectedContext: context.Background(),
			expectedName:    "name",
			expectedAge:     43,
			withError:       false,
			expectedError:   nil,
		},
		{
			name: "save service test failure",
			executeBeforeTest: func(ctx context.Context, r *MockRepository, expectedName string, expectedAge uint) {
				r.
					EXPECT().
					Save(gomock.Eq(ctx), gomock.Eq(expectedName), gomock.Eq(expectedAge)).
					Return(uint(0), ErrorSavingToDB)
			},
			expectedContext: context.Background(),
			expectedName:    "name",
			expectedAge:     43,
			withError:       true,
			expectedError:   ErrorSavingToDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			repository := NewMockRepository(ctrl)

			tt.executeBeforeTest(tt.expectedContext, repository, tt.expectedName, tt.expectedAge)

			service := NewService(repository)

			// when
			_, err := service.Save(tt.expectedContext, tt.expectedName, tt.expectedAge)

			// then

			if err != tt.expectedError {
				t.Fail()
			}
		})
	}

}
