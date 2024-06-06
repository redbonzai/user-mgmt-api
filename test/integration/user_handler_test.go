package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/redbonzai/user-management-api/internal/interfaces"
	"github.com/redbonzai/user-management-api/internal/interfaces/handler"
	"github.com/redbonzai/user-management-api/internal/services/mocks"
)

var _ = Describe("UserHandler", func() {
	var (
		e           *echo.Echo
		rec         *httptest.ResponseRecorder
		userHandler *handler.UserHandler
		mockCtrl    *gomock.Controller
		service     *mocks.MockService
	)

	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
		mockCtrl = gomock.NewController(GinkgoT())
		service = mocks.NewMockService(mockCtrl)
		userHandler = handler.NewUserHandler(service)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetUsers", func() {
		activeStatus := "active"
		inactiveStatus := "inactive"
		It("should return all users", func() {
			users := []interfaces.User{
				{ID: 1, Name: "User One", Email: "user1@example.com", Status: &activeStatus},
				{ID: 2, Name: "User Two", Email: "user2@example.com", Status: &inactiveStatus},
			}

			service.EXPECT().GetUsers().Return(users, nil)

			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			ctx := e.NewContext(req, rec)

			err := userHandler.GetUsers(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("User One"))
			Expect(rec.Body.String()).To(ContainSubstring("User Two"))
		})
	})

	Describe("GetUser", func() {
		activeStatus := "active"

		It("should return a user by ID", func() {
			user := interfaces.User{ID: 1, Name: "User One", Email: "user1@example.com", Status: &activeStatus}

			service.EXPECT().GetUserByID(1).Return(user, nil)

			req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := userHandler.GetUser(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("User One"))
		})

		It("should return an error if the user ID is invalid", func() {
			req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("invalid")

			err := userHandler.GetUser(ctx)
			Expect(err).To(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
			Expect(rec.Body.String()).To(ContainSubstring("Invalid ID"))
		})
	})

	Describe("CreateUser", func() {
		activeStatus := "active"
		It("should create a new user", func() {
			user := interfaces.User{Username: "newuser", Password: "newpass", Name: "New User", Email: "new@example.com", Status: &activeStatus}

			service.EXPECT().CreateUser(gomock.Any()).Return(user, nil)

			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"username":"newuser","password":"newpass","name":"New User","email":"new@example.com","status":"active"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := userHandler.CreateUser(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusCreated))
			Expect(rec.Body.String()).To(ContainSubstring(`"username":"newuser"`))
		})

		It("should return an error for invalid input", func() {
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"username":"","password":"","name":"","email":"","status":""}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := userHandler.CreateUser(ctx)
			Expect(err).To(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
			Expect(rec.Body.String()).To(ContainSubstring("Invalid input"))
		})
	})

	Describe("UpdateUser", func() {
		activeStatus := "active"
		It("should update a user", func() {
			user := interfaces.User{ID: 1, Username: "updateduser", Password: "updatedpass", Name: "Updated User", Email: "updated@example.com", Status: &activeStatus}

			service.EXPECT().UpdateUser(gomock.Any()).Return(user, nil)

			req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(`{"username":"updateduser","password":"updatedpass","name":"Updated User","email":"updated@example.com","status":"active"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := userHandler.UpdateUser(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring(`"username":"updateduser"`))
		})

		It("should return an error for invalid input", func() {
			req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(`{"username":"","password":"","name":"","email":"","status":""}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := userHandler.UpdateUser(ctx)
			Expect(err).To(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
			Expect(rec.Body.String()).To(ContainSubstring("Invalid input"))
		})
	})

	Describe("DeleteUser", func() {
		activeStatus := "active"
		It("should delete a user", func() {
			user := interfaces.User{ID: 1, Name: "User One", Email: "user1@example.com", Status: &activeStatus}

			service.EXPECT().DeleteUser(1).Return(user, nil)

			req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := userHandler.DeleteUser(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("User One"))
		})

		It("should return an error if the user ID is invalid", func() {
			req := httptest.NewRequest(http.MethodDelete, "/users/invalid", nil)
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("invalid")

			err := userHandler.DeleteUser(ctx)
			Expect(err).To(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
			Expect(rec.Body.String()).To(ContainSubstring("Invalid ID"))
		})
	})

	Describe("Login", func() {
		activeStatus := "active"
		It("should login a user and return a token", func() {
			req := httptest.NewRequest(http.MethodPost, "/v1/login", strings.NewReader(`{"username":"testuser","password":"testpass"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			service.EXPECT().GetUserByUsername("testuser").Return(interfaces.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Status:   &activeStatus,
				Username: "testuser",
				Password: "testpass",
			}, nil)

			err := userHandler.Login(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("token"))
		})

		It("should return an error for invalid credentials", func() {
			req := httptest.NewRequest(http.MethodPost, "/v1/login", strings.NewReader(`{"username":"testuser","password":"wrongpass"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			service.EXPECT().GetUserByUsername("testuser").Return(interfaces.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Status:   &activeStatus,
				Username: "testuser",
				Password: "testpass",
			}, nil)

			err := userHandler.Login(ctx)
			Expect(err).To(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
			Expect(rec.Body.String()).To(ContainSubstring("Invalid username or password"))
		})
	})

	Describe("Register", func() {
		activeStatus := "active"
		It("should register a new user", func() {
			req := httptest.NewRequest(http.MethodPost, "/v1/register", strings.NewReader(`{"username":"newuser","password":"newpass", "name":"New User", "email":"new@example.com", "status":"active"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			service.EXPECT().CreateUser(gomock.Any()).Return(interfaces.User{
				ID:       1,
				Name:     "New User",
				Email:    "new@example.com",
				Status:   &activeStatus,
				Username: "newuser",
				Password: "newpass",
			}, nil)

			err := userHandler.Register(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusCreated))
			Expect(rec.Body.String()).To(ContainSubstring(`"username":"newuser"`))
		})

		It("should return an error for invalid input", func() {
			req := httptest.NewRequest(http.MethodPost, "/v1/register", strings.NewReader(`{"username":"","password":"","name":"","email":"","status":""}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := userHandler.Register(ctx)
			Expect(err).To(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
			Expect(rec.Body.String()).To(ContainSubstring("Invalid input"))
		})
	})
})
