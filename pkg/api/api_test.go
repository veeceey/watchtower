package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	token = "123123123"
)

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}

var _ = Describe("API", func() {
	api := New(token)

	Describe("RequireToken middleware", func() {
		It("should return 401 Unauthorized when token is not provided", func() {
			handlerFunc := api.RequireToken(testHandler)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/hello", nil)

			handlerFunc(rec, req)

			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return 401 Unauthorized when token is invalid", func() {
			handlerFunc := api.RequireToken(testHandler)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/hello", nil)
			req.Header.Set("Authorization", "Bearer 123")

			handlerFunc(rec, req)

			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return 200 OK when token is valid", func() {
			handlerFunc := api.RequireToken(testHandler)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/hello", nil)
			req.Header.Set("Authorization", "Bearer " + token)

			handlerFunc(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("Custom ServeMux", func() {
		It("should return 404 for unregistered paths", func() {
			api.RegisterFunc("/v1/test", testHandler)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/debug/vars", nil)
			api.mux.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should serve registered paths with valid token", func() {
			api.RegisterFunc("/v1/hello", testHandler)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/v1/hello", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			api.mux.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should require token for registered paths", func() {
			api.RegisterFunc("/v1/secured", testHandler)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/v1/secured", nil)
			api.mux.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})

func testHandler(w http.ResponseWriter, req *http.Request) {
	_, _ = io.WriteString(w, "Hello!")
}
