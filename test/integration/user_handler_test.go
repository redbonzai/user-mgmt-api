package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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
		userService *mocks.MockService
	)

	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
		mockCtrl = gomock.NewController(GinkgoT())
		userService = mocks.NewMockService(mockCtrl)
		userHandler = handler.NewUserHandler(userService)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetUsers", func() {
		It("should return all users", func() {
			users := []interfaces.User{
				{ID: 1, Name: "User One", Email: "user1@example.com"},
				{ID: 2, Name: "User Two", Email: "user2@example.com"},
			}

			userService.EXPECT().GetUsers().Return(users, nil)

			req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
			ctx := e.NewContext(req, rec)

			err := userHandler.GetUsers(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("User One"))
			Expect(rec.Body.String()).To(ContainSubstring("User Two"))
		})
	})

	Describe("GetUser", func() {
		It("should return a user by ID", func() {
			user := interfaces.User{ID: 1, Name: "User One", Email: "user1@example.com"}

			userService.EXPECT().GetUserByID(1).Return(user, nil)

			req := httptest.NewRequest(http.MethodGet, "/v1/users/1", nil)
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := userHandler.GetUser(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("User One"))
		})
	})

	Describe("CreateUser", func() {
		It("should create a new user", func() {
			user := interfaces.User{Username: "newuser", Password: "newpass", Name: "New User", Email: "new@example.com"}

			userService.EXPECT().CreateUser(gomock.Any()).Return(user, nil)

			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(`{"username":"newuser","password":"newpass","name":"New User","email":"new@example.com"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := userHandler.CreateUser(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusCreated))
			Expect(rec.Body.String()).To(ContainSubstring(`"username":"newuser"`))
		})
	})

	Describe("UpdateUser", func() {
		It("should update a user", func() {
			existingUser := interfaces.User{ID: 1, Username: "existinguser", Password: "existingpass", Name: "Existing User", Email: "existing@example.com"}
			updatedUser := interfaces.User{ID: 1, Username: "updateduser", Password: "updatedpass", Name: "Updated User", Email: "updated@example.com"}

			userService.EXPECT().GetUserByID(1).Return(existingUser, nil)
			userService.EXPECT().UpdateUser(gomock.Any()).Return(updatedUser, nil)

			req := httptest.NewRequest(http.MethodPut, "/v1/users/1", strings.NewReader(`{"username":"updateduser","password":"updatedpass","name":"Updated User","email":"updated@example.com"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := userHandler.UpdateUser(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring(`"username":"updateduser"`))
		})
	})

	Describe("DeleteUser", func() {
		It("should delete a user", func() {
			user := interfaces.User{ID: 1, Name: "User One", Email: "user1@example.com"}

			userService.EXPECT().DeleteUser(1).Return(user, nil)

			req := httptest.NewRequest(http.MethodDelete, "/v1/users/1", nil)
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := userHandler.DeleteUser(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("User One"))
		})
	})

	Describe("Login", func() {
		It("should login a user and return a token", func() {
			user := interfaces.User{ID: 1, Username: "testuser", Password: "$2a$10$7.qGVUb5v4PQcK/n1Ub0RODnpDFnx/38TF/1ntCR3IUmY/ma1DLG2", Name: "Test User", Email: "test@example.com"} // hashed password for "testpass"

			userService.EXPECT().GetUserByUsername("testuser").Return(user, nil)
			userService.EXPECT().HashPassword(gomock.Any()).Return(user.Password, nil)

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"testuser","password":"testpass"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := userHandler.Login(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("token"))
		})
	})

	Describe("Logout", func() {
		It("should logout a user and blacklist the token", func() {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": "testuser",
				"exp":      time.Now().Add(time.Hour * 1).Unix(),
			})
			tokenString, _ := token.SignedString([]byte("secret"))

			req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(""))
			req.Header.Set("Authorization", "Bearer "+tokenString)
			ctx := e.NewContext(req, rec)
			ctx.Set("user", token)

			expTime := time.Unix(token.Claims.(jwt.MapClaims)["exp"].(int64), 0)
			userService.EXPECT().Logout(tokenString, expTime).Return(nil)

			err := userHandler.Logout(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("logged out successfully"))
		})
	})

	Describe("Register", func() {
		It("should register a new user", func() {
			user := interfaces.User{Username: "newuser", Password: "$2a$10$7.qGVUb5v4PQcK/n1Ub0RODnpDFnx/38TF/1ntCR3IUmY/ma1DLG2", Name: "New User", Email: "new@example.com"} // hashed password for "newpass"
			registerRequest := interfaces.RegisterRequest{Username: "newuser", Password: "newpass", Name: "New User", Email: "new@example.com"}

			userService.EXPECT().IsUsernameUnique(registerRequest.Username).Return(true, nil)
			userService.EXPECT().HashPassword(registerRequest.Password).Return(user.Password, nil)
			userService.EXPECT().CreateUser(gomock.Any()).Return(user, nil)

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"newuser","password":"newpass","name":"New User","email":"new@example.com"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := userHandler.Register(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusCreated))
			Expect(rec.Body.String()).To(ContainSubstring(`"username":"newuser"`))
		})
	})

	Describe("GetAuthenticatedUser", func() {
		It("should return the authenticated user", func() {
			user := interfaces.User{ID: 1, Username: "testuser", Name: "Test User", Email: "test@example.com"}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": "testuser",
				"exp":      time.Now().Add(time.Hour * 1).Unix(),
			})
			tokenString, _ := token.SignedString([]byte("secret"))

			userService.EXPECT().GetUserByUsername("testuser").Return(user, nil)

			req := httptest.NewRequest(http.MethodGet, "/v1/current-user", nil)
			req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
			ctx := e.NewContext(req, rec)
			ctx.Set("user", token.Claims.(jwt.MapClaims))

			err := userHandler.GetAuthenticatedUser(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("Test User"))
		})
	})
})
