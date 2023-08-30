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

func TestHandler_createSegment(t *testing.T) {
	type mockBehavior func(s *mock_service.MockSegment, segment structures.Segment)

	tests := []struct {
		name                 string
		inputBody            string
		inputSegment         structures.Segment
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"slug": "example-slug"}`,
			inputSegment: structures.Segment{
				Slug: "example-slug",
			},
			mockBehavior: func(s *mock_service.MockSegment, segment structures.Segment) {
				s.EXPECT().Create(segment).Return("example-slug", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"slug":"example-slug"}`,
		},
		{
			name:      "EmptySlug",
			inputBody: `{"slug": ""}`,
			inputSegment: structures.Segment{
				Slug: "",
			},
			mockBehavior: func(s *mock_service.MockSegment, segment structures.Segment) {

			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Key: 'Segment.Slug' Error:Field validation for 'Slug' failed on the 'required' tag"}`,
		},
		{
			name:      "InvalidSlug",
			inputBody: `{"slug": "iueefiuwq-"}`,
			inputSegment: structures.Segment{
				Slug: "iueefiuwq-",
			},
			mockBehavior: func(s *mock_service.MockSegment, segment structures.Segment) {

			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid slug"}`,
		},
		{
			name:      "InvalidJSON",
			inputBody: `{"slug":example-slug"}`,
			inputSegment: structures.Segment{
				Slug: "example-slug",
			},
			mockBehavior: func(s *mock_service.MockSegment, segment structures.Segment) {

			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid character 'e' looking for beginning of value"}`,
		},
		{
			name:      "ServiceFail",
			inputBody: `{"slug":"example-slug"}`,
			inputSegment: structures.Segment{
				Slug: "example-slug",
			},
			mockBehavior: func(s *mock_service.MockSegment, segment structures.Segment) {
				s.EXPECT().Create(segment).Return("example-slug", errors.New("service fail"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service fail"}`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mock := mock_service.NewMockSegment(ctl)
			testCase.mockBehavior(mock, testCase.inputSegment)

			services := &service.Service{Segment: mock}
			h := Handler{services}

			r := gin.New()
			r.POST("/segments/", h.createSegment)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/segments/", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_deleteSegment(t *testing.T) {
	type mockBehavior func(s *mock_service.MockSegment, segment structures.Segment)

	tests := []struct {
		name                string
		inputBody           string
		inputSegment        structures.Segment
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"slug": "example-slug"}`,
			inputSegment: structures.Segment{
				Slug: "example-slug",
			},
			mockBehavior: func(s *mock_service.MockSegment, segment structures.Segment) {
				s.EXPECT().Delete(segment).Return("example-slug", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"slug":"example-slug"}`,
		},
		{
			name:      "EmptySlug",
			inputBody: `{"slug": ""}`,
			inputSegment: structures.Segment{
				Slug: "",
			},
			mockBehavior: func(s *mock_service.MockSegment, segment structures.Segment) {

			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"Key: 'Segment.Slug' Error:Field validation for 'Slug' failed on the 'required' tag"}`,
		},
		{
			name:      "InvalidSlug",
			inputBody: `{"slug": "iueefiuwq-"}`,
			inputSegment: structures.Segment{
				Slug: "iueefiuwq-",
			},
			mockBehavior: func(s *mock_service.MockSegment, segment structures.Segment) {

			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid slug"}`,
		},
		{
			name:      "InvalidJSON",
			inputBody: `{"slug":example-slug"}`,
			inputSegment: structures.Segment{
				Slug: "example-slug",
			},
			mockBehavior: func(s *mock_service.MockSegment, segment structures.Segment) {

			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid character 'e' looking for beginning of value"}`,
		},
		{
			name:      "ServiceFail",
			inputBody: `{"slug":"example-slug"}`,
			inputSegment: structures.Segment{
				Slug: "example-slug",
			},
			mockBehavior: func(s *mock_service.MockSegment, segment structures.Segment) {
				s.EXPECT().Delete(segment).Return("example-slug", errors.New("service fail"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service fail"}`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mock := mock_service.NewMockSegment(ctl)
			testCase.mockBehavior(mock, testCase.inputSegment)

			services := &service.Service{Segment: mock}
			h := Handler{services}

			r := gin.New()
			r.DELETE("/segments/", h.deleteSegment)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/segments/", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
