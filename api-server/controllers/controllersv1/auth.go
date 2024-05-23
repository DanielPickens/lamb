package controllersv1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/danielpickens/yeti-schemas/schemasv1"
	"github.com/danielpickens/yeti/api-server/models"
	"github.com/danielpickens/yeti/api-server/services"
	"github.com/danielpickens/yeti/api-server/transformers/transformersv1"
	"github.com/danielpickens/yeti/common/scookie"
	"github.com/danielpickens/yeti/common/utils"
)

type authController struct {
	
	baseController
}
//create auth controller user services register context schema
var AuthController = authController{}

func (*authController) Register(ctx *gin.Context, schema *schemasv1.RegisterUserSchema) (*schemasv1.UserSchema, error) {
	user, err := services.UserService.Create(ctx, services.CreateUserOption{
		Name:      schema.Name,
		FirstName: schema.FirstName,
		LastName:  schema.LastName,
		Email:     utils.StringPtrWithoutEmpty(schema.Email),
		Password:  schema.Password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "create user")
	}
	err = scookie.SetUsernameToCookie(ctx, user.Name)
	if err != nil {
		return nil, errors.Wrap(err, "set login cookie")
	}
	return transformersv1.ToUserSchema(ctx, user)
}

func (*authController) Login(ctx *gin.Context, schema *schemasv1.LoginUserSchema) (*schemasv1.UserSchema, error) {
	isEmail := strings.Contains(schema.NameOrEmail, "@")
	var err error
	var user *models.User
	if isEmail {
		user, err = services.UserService.GetByEmail(ctx, schema.NameOrEmail)
	} else {
		user, err = services.UserService.GetByName(ctx, schema.NameOrEmail)
	}
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	if user.Email == nil || *user.Email == "" {
		return nil, errors.Errorf("user %s email is empty, it looks like yeti did not complete the setup process", user.Name)
	}
	if err = services.UserService.CheckPassword(ctx, user, schema.Password); err != nil {
		return nil, err
	}
	err = scookie.SetUsernameToCookie(ctx, user.Name)
	if err != nil {
		return nil, errors.Wrap(err, "set login cookie")
	}
	redirectUri := ctx.Query("redirect")
	if redirectUri == "" {
		redirectUri = "/"
	}
	ctx.Redirect(http.StatusSeeOther, redirectUri)
	return transformersv1.ToUserSchema(ctx, user)
}

func (*authController) GetCurrentUser(ctx *gin.Context) (*schemasv1.UserSchema, error) {
	user, err := services.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return transformersv1.ToUserSchema(ctx, user)
}

func (*authController) ResetPassword(ctx *gin.Context, schema *schemasv1.ResetPasswordSchema) (*schemasv1.UserSchema, error) {
	currentUser, err := services.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	user, err := services.UserService.UpdatePassword(ctx, currentUser, schema.CurrentPassword, schema.NewPassword)
	if err != nil {
		return nil, err
	}

	return transformersv1.ToUserSchema(ctx, user)
}
