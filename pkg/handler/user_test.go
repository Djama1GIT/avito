package handler

import (
	"avito/pkg/service"
	mock_service "avito/pkg/service/mocks"
	"avito/pkg/structures"
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestHandler_getUserHistory(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, userHistory structures.UserHistory)

	tests := []struct {
		name                 string
		queryParams          map[string]string
		inputBody            string
		userHistory          structures.UserHistory
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"user_id": 1, "year_month": "2023-08"}`,
			userHistory: structures.UserHistory{
				Id:        1,
				YearMonth: "2023-08",
			},
			mockBehavior: func(s *mock_service.MockUser, userHistory structures.UserHistory) {
				s.EXPECT().GetUserHistory(userHistory).Return("example", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"report":"http://localhost:8000/files/example","user_id":1}`,
		},
		{
			name:        "OK",
			queryParams: map[string]string{"user_id": "1", "year_month": "2023-08"},
			userHistory: structures.UserHistory{
				Id:        1,
				YearMonth: "2023-08",
			},
			mockBehavior: func(s *mock_service.MockUser, userHistory structures.UserHistory) {
				s.EXPECT().GetUserHistory(userHistory).Return("example", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"report":"http://localhost:8000/files/example","user_id":1}`,
		},
		{
			name:        "InvalidData",
			queryParams: map[string]string{"user_id": "1"},
			userHistory: structures.UserHistory{
				Id: 1,
			},
			mockBehavior: func(s *mock_service.MockUser, userHistory structures.UserHistory) {

			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Key: 'UserHistory.YearMonth' Error:Field validation for 'YearMonth' failed on the 'required' tag"}`,
		},
		{
			name:      "InvalidData",
			inputBody: `{"user_id": 1}`,
			userHistory: structures.UserHistory{
				Id: 1,
			},
			mockBehavior: func(s *mock_service.MockUser, userHistory structures.UserHistory) {

			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Key: 'UserHistory.YearMonth' Error:Field validation for 'YearMonth' failed on the 'required' tag"}`,
		},
		{
			name:        "InvalidData",
			queryParams: map[string]string{"user_id": "one", "year_month": "2023-08"},
			userHistory: structures.UserHistory{
				Id:        1,
				YearMonth: "2023-08",
			},
			mockBehavior: func(s *mock_service.MockUser, userHistory structures.UserHistory) {

			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"strconv.Atoi: parsing \"one\": invalid syntax"}`,
		},
		{
			name:      "ServiceFail",
			inputBody: `{"user_id": 1, "year_month": "2023-08"}`,
			userHistory: structures.UserHistory{
				Id:        1,
				YearMonth: "2023-08",
			},
			mockBehavior: func(s *mock_service.MockUser, userHistory structures.UserHistory) {
				s.EXPECT().GetUserHistory(userHistory).Return("example", errors.New("service fail"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service fail"}`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mock := mock_service.NewMockUser(ctl)
			testCase.mockBehavior(mock, testCase.userHistory)

			services := &service.Service{User: mock}
			h := Handler{services}

			r := gin.New()
			r.GET("/history/", h.getUserHistory)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/history/", bytes.NewBufferString(testCase.inputBody))

			q := req.URL.Query()
			for key, value := range testCase.queryParams {
				q.Add(key, value)
			}

			req.URL.RawQuery = q.Encode()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
