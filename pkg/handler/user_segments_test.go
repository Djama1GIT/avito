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

func TestHandler_patchSegment(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUserSegments, input structures.UserSegments)

	tests := []struct {
		name                 string
		inputBody            string
		inputData            structures.UserSegments
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"user_id": 1, "segments_to_add": ["segment1", "segment2"], "segments_to_delete": ["segment3"]}`,
			inputData: structures.UserSegments{
				UserId:           1,
				SegmentsToAdd:    []string{"segment1", "segment2"},
				SegmentsToDelete: []string{"segment3"},
			},
			mockBehavior: func(s *mock_service.MockUserSegments, input structures.UserSegments) {
				s.EXPECT().Patch(input).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user_id":1}`,
		},
		{
			name:      "InvaligSegmentToAdd",
			inputBody: `{"user_id": 1, "segments_to_add": ["segment1-", "segment2"], "segments_to_delete": ["segment3"]}`,
			inputData: structures.UserSegments{
				UserId:           1,
				SegmentsToAdd:    []string{"segment1-", "segment2"},
				SegmentsToDelete: []string{"segment3"},
			},
			mockBehavior: func(s *mock_service.MockUserSegments, input structures.UserSegments) {

			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid slug (segment to add: segment1-)"}`,
		},
		{
			name:      "InvaligSegmentToDelete",
			inputBody: `{"user_id": 1, "segments_to_add": ["segment1", "segment2"], "segments_to_delete": ["_segment3/"]}`,
			inputData: structures.UserSegments{
				UserId:           1,
				SegmentsToAdd:    []string{"segment1", "segment2"},
				SegmentsToDelete: []string{"_segment3/"},
			},
			mockBehavior: func(s *mock_service.MockUserSegments, input structures.UserSegments) {

			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid slug (segment to delete: _segment3/)"}`,
		},
		{
			name:      "InvalidJSON",
			inputBody: `"user_id": 1, "segments_to_add": ["segment1", "segment2"], "segments_to_delete": ["segment3"]}`,
			inputData: structures.UserSegments{
				UserId:           1,
				SegmentsToAdd:    []string{"segment1", "segment2"},
				SegmentsToDelete: []string{"segment3"},
			},
			mockBehavior: func(s *mock_service.MockUserSegments, input structures.UserSegments) {

			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"json: cannot unmarshal string into Go value of type structures.UserSegments"}`,
		},
		{
			name:      "ServiceFail",
			inputBody: `{"user_id": 1, "segments_to_add": ["segment1", "segment2"], "segments_to_delete": ["segment3"]}`,
			inputData: structures.UserSegments{
				UserId:           1,
				SegmentsToAdd:    []string{"segment1", "segment2"},
				SegmentsToDelete: []string{"segment3"},
			},
			mockBehavior: func(s *mock_service.MockUserSegments, input structures.UserSegments) {
				s.EXPECT().Patch(input).Return(0, errors.New("service fail"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service fail"}`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mock := mock_service.NewMockUserSegments(ctl)
			testCase.mockBehavior(mock, testCase.inputData)

			services := &service.Service{UserSegments: mock}
			h := Handler{services}

			r := gin.New()
			r.PATCH("/segments/", h.patchSegment)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/segments/", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getUsersInSegment(t *testing.T) {
	type mockBehavior func(us *mock_service.MockUserSegments, s *mock_service.MockSegment, input structures.User)

	tests := []struct {
		name                 string
		queryParams          map[string]string
		inputBody            string
		inputData            structures.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			queryParams: map[string]string{"user_id": "1"},
			inputBody:   `{"id": 1}`,
			inputData: structures.User{
				Id: 1,
			},
			mockBehavior: func(us *mock_service.MockUserSegments, s *mock_service.MockSegment, input structures.User) {
				us.EXPECT().GetUsersInSegment(input).Return([]string{"segment1", "segment2"}, nil)
				s.EXPECT().GetPercentageSegments().Return(map[string]int{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"segments":["segment1","segment2"],"user_id":1}`,
		},
		{
			name:        "OK",
			queryParams: map[string]string{"user_id": "1"},
			inputBody:   `{"id": 1}`,
			inputData: structures.User{
				Id: 1,
			},
			mockBehavior: func(us *mock_service.MockUserSegments, s *mock_service.MockSegment, input structures.User) {
				us.EXPECT().GetUsersInSegment(input).Return([]string{"segment1", "segment2"}, nil)
				s.EXPECT().GetPercentageSegments().Return(map[string]int{"segment3": 100}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"segments":["segment1","segment2","segment3"],"user_id":1}`,
		},
		{
			name:        "EmptySegments",
			queryParams: map[string]string{"user_id": "1"},
			inputBody:   `{"id": 1}`,
			inputData: structures.User{
				Id: 1,
			},
			mockBehavior: func(us *mock_service.MockUserSegments, s *mock_service.MockSegment, input structures.User) {
				us.EXPECT().GetUsersInSegment(input).Return(nil, nil)
				s.EXPECT().GetPercentageSegments().Return(map[string]int{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"segments":[],"user_id":1}`,
		},
		{
			name:        "InvalidUserID",
			queryParams: map[string]string{"user_id": "invalid"},
			inputBody:   ``,
			inputData:   structures.User{},
			mockBehavior: func(us *mock_service.MockUserSegments, s *mock_service.MockSegment, input structures.User) {

			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"strconv.Atoi: parsing \"invalid\": invalid syntax"}`,
		},
		{
			name:        "InvalidJSON",
			queryParams: map[string]string{},
			inputBody:   `"id": 1}`,
			inputData:   structures.User{},
			mockBehavior: func(us *mock_service.MockUserSegments, s *mock_service.MockSegment, input structures.User) {

			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"json: cannot unmarshal string into Go value of type structures.User"}`,
		},
		{
			name:        "ServiceFail",
			queryParams: map[string]string{"user_id": "1"},
			inputBody:   `{"id": 1}`,
			inputData: structures.User{
				Id: 1,
			},
			mockBehavior: func(us *mock_service.MockUserSegments, s *mock_service.MockSegment, input structures.User) {
				us.EXPECT().GetUsersInSegment(input).Return(nil, errors.New("service fail"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service fail"}`,
		},
		{
			name:        "ServiceFail",
			queryParams: map[string]string{"user_id": "1"},
			inputBody:   `{"id": 1}`,
			inputData: structures.User{
				Id: 1,
			},
			mockBehavior: func(us *mock_service.MockUserSegments, s *mock_service.MockSegment, input structures.User) {
				us.EXPECT().GetUsersInSegment(input).Return(nil, nil)
				s.EXPECT().GetPercentageSegments().Return(nil, errors.New("service fail"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service fail"}`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			usMock := mock_service.NewMockUserSegments(ctl)
			sMock := mock_service.NewMockSegment(ctl)
			testCase.mockBehavior(usMock, sMock, testCase.inputData)

			services := &service.Service{UserSegments: usMock, Segment: sMock}
			h := Handler{services}

			r := gin.New()
			r.GET("/segments/", h.getUsersInSegment)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/segments/", bytes.NewBufferString(testCase.inputBody))

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
