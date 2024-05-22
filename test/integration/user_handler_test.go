package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/redbonzai/user-management-api/internal/domain/user"
	"github.com/redbonzai/user-management-api/internal/interfaces/handler"
	"github.com/redbonzai/user-management-api/pkg/logger"
)

func TestUserHandler(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "UserHandler Suite")
}

var _ = ginkgo.Describe("UserHandler", func() {
	var (
		e           *echo.Echo
		rec         *httptest.ResponseRecorder
		ctx         echo.Context
		userHandler *handler.UserHandler
		mockService *handler.MockUserService
	)

	ginkgo.BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
		logger.InitLogger()

		mockService = &handler.MockUserService{}
		userHandler = handler.NewUserHandler(mockService)
	})

	ginkgo.Describe("GetUsers", func() {
		ginkgo.It("should return all users", func() {
			// Mock data
			users := []user.User{
				{ID: 1, Name: "John Doe", Email: "john.doe@example.com"},
				{ID: 2, Name: "Jane Doe", Email: "jane.doe@example.com"},
			}
			mockService.On("GetUsers").Return(users, nil)

			req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
			ctx = e.NewContext(req, rec)

			err := userHandler.GetUsers(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))

			var response []user.User
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response).To(gomega.Equal(users))
		})
	})

	ginkgo.Describe("GetUser", func() {
		ginkgo.It("should return a user by ID", func() {
			// Mock data
			retrievedUser := user.User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}
			mockService.On("GetUserByID", 1).Return(retrievedUser, nil)

			req := httptest.NewRequest(http.MethodGet, "/v1/users/1", nil)
			ctx = e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := userHandler.GetUser(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))

			var response user.User
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response).To(gomega.Equal(retrievedUser))
		})
	})

	ginkgo.Describe("CreateUser", func() {
		ginkgo.It("should create a new user", func() {
			// Mock data
			newUser := user.User{Name: "John Doe", Email: "john.doe@example.com"}
			mockService.On("CreateUser", newUser).Return(newUser, nil)

			userJSON, _ := json.Marshal(newUser)
			req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(userJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx = e.NewContext(req, rec)

			err := userHandler.CreateUser(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusCreated))

			var response user.User
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response).To(gomega.Equal(newUser))
		})
	})

	ginkgo.Describe("UpdateUser", func() {
		ginkgo.It("should update an existing user", func() {
			// Mock data
			updatedUser := user.User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}
			mockService.On("UpdateUser", updatedUser).Return(updatedUser, nil)

			userJSON, _ := json.Marshal(updatedUser)
			req := httptest.NewRequest(http.MethodPut, "/v1/users/1", bytes.NewReader(userJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx = e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := userHandler.UpdateUser(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))

			var response user.User
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response).To(gomega.Equal(updatedUser))
		})
	})

	ginkgo.Describe("DeleteUser", func() {
		ginkgo.It("should delete an existing user", func() {
			// Mock data
			deletedUser := user.User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}
			mockService.On("DeleteUser", 1).Return(deletedUser, nil)

			req := httptest.NewRequest(http.MethodDelete, "/v1/users/1", nil)
			ctx = e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := userHandler.DeleteUser(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))

			var response user.User
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response).To(gomega.Equal(deletedUser))
		})
	})
})
