package login

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/img21326/fb_chat/structure/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/img21326/fb_chat/mock"
)

func TestRegisterNotSetGender(t *testing.T) {
	c := gomock.NewController(t)

	authUsecase := mock.NewMockAuthUsecaseInterFace(c)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	NewLoginController(r, authUsecase)

	req, _ := http.NewRequest("GET", "/register", nil)
	r.ServeHTTP(w, req)

	body := w.Body.Bytes()

	assert.Equal(t, 410, w.Code)
	assert.Equal(t, `{"error":"should add params with gender"}`, string(body[:]))
}

func TestRegisterWithSuccess(t *testing.T) {
	c := gomock.NewController(t)

	u := &user.User{
		Gender: "male",
	}

	uid, _ := uuid.Parse("1cad3ac9-b72f-4b96-8c20-0acb2debec49")

	authUsecase := mock.NewMockAuthUsecaseInterFace(c)
	authUsecase.EXPECT().GenerateToken(gomock.Any(), u).DoAndReturn(func(ctx context.Context, u *user.User) (string, error) {
		u.UUID = uid
		return "jwt_token", nil
	})

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	NewLoginController(r, authUsecase)

	req, _ := http.NewRequest("GET", "/register?gender=male", nil)
	r.ServeHTTP(w, req)

	body := w.Body.Bytes()

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"token":"jwt_token","uuid":"1cad3ac9-b72f-4b96-8c20-0acb2debec49"}`, string(body[:]))
}
