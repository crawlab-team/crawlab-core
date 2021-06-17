package user

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"time"
)

type Service struct {
	// settings variables
	jwtSecret        string
	jwtSigningMethod jwt.SigningMethod

	// dependencies
	modelSvc service.ModelService
}

func (svc *Service) Init() (err error) {
	_, err = svc.modelSvc.GetUserByUsername(constants.DefaultAdminUsername, nil)
	if err == nil {
		return nil
	}
	if err.Error() != mongo.ErrNoDocuments.Error() {
		return err
	}
	return svc.Create(&interfaces.UserCreateOptions{
		Username: constants.DefaultAdminUsername,
		Password: constants.DefaultAdminPassword,
		Role:     constants.RoleAdmin,
	})
}

func (svc *Service) SetJwtSecret(secret string) {
	svc.jwtSecret = secret
}

func (svc *Service) SetJwtSigningMethod(method jwt.SigningMethod) {
	svc.jwtSigningMethod = method
}

func (svc *Service) Create(opts *interfaces.UserCreateOptions) (err error) {
	if opts.Username == "" || opts.Password == "" {
		return trace.TraceError(errors.ErrorUserMissingRequiredFields)
	}
	if opts.Role == "" {
		opts.Role = constants.RoleNormal
	}
	if _, err := svc.modelSvc.GetUserByUsername(opts.Username, nil); err == nil {
		return trace.TraceError(errors.ErrorUserAlreadyExists)
	}
	u := &models.User{
		Username: opts.Username,
		Password: utils.EncryptPassword(opts.Password),
		Role:     opts.Role,
		Email:    opts.Email,
	}
	if err := delegate.NewModelDelegate(u).Add(); err != nil {
		return err
	}
	return nil
}

func (svc *Service) Login(opts *interfaces.UserLoginOptions) (token string, u interfaces.User, err error) {
	u, err = svc.modelSvc.GetUserByUsername(opts.Username, nil)
	if err != nil {
		return "", nil, err
	}
	if u.GetPassword() != utils.EncryptPassword(opts.Password) {
		return "", nil, errors.ErrorUserMismatch
	}
	token, err = svc.makeToken(u)
	if err != nil {
		return "", nil, err
	}
	return token, u, nil
}

func (svc *Service) CheckToken(tokenStr string) (u interfaces.User, err error) {
	return svc.checkToken(tokenStr)
}

func (svc *Service) makeToken(user interfaces.User) (tokenStr string, err error) {
	token := jwt.NewWithClaims(svc.jwtSigningMethod, jwt.MapClaims{
		"id":       user.GetId(),
		"username": user.GetUsername(),
		"nbf":      time.Now().Unix(),
	})
	return token.SignedString([]byte(svc.jwtSecret))
}

func (svc *Service) checkToken(tokenStr string) (user interfaces.User, err error) {
	token, err := jwt.Parse(tokenStr, svc.getSecretFunc())
	if err != nil {
		return
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.ErrorUserInvalidType
		return
	}

	if !token.Valid {
		err = errors.ErrorUserInvalidToken
		return
	}

	id, err := primitive.ObjectIDFromHex(claim["id"].(string))
	if err != nil {
		return user, err
	}
	username := claim["username"].(string)
	user, err = svc.modelSvc.GetUserById(id)
	if err != nil {
		err = errors.ErrorUserNotExists
		return
	}

	if username != user.GetUsername() {
		err = errors.ErrorUserMismatch
		return
	}

	return
}

func (svc *Service) getSecretFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(svc.jwtSecret), nil
	}
}

func NewUserService(opts ...Option) (svc2 interfaces.UserService, err error) {
	// service
	svc := &Service{
		jwtSecret:        "crawlab",
		jwtSigningMethod: jwt.SigningMethodHS256,
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(modelSvc service.ModelService) {
		svc.modelSvc = modelSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// initialize
	if err := svc.Init(); err != nil {
		return nil, trace.TraceError(err)
	}

	return svc, nil
}

func ProvideUserService(opts ...Option) func() (svc interfaces.UserService, err error) {
	return func() (svc interfaces.UserService, err error) {
		return NewUserService(opts...)
	}
}

var userSvc interfaces.UserService

func GetUserService(opts ...Option) (svc interfaces.UserService, err error) {
	if userSvc != nil {
		return userSvc, nil
	}
	svc, err = NewUserService(opts...)
	if err != nil {
		return nil, err
	}
	userSvc = svc
	return svc, nil
}

func ProvideGetUserService(opts ...Option) func() (svr interfaces.UserService, err error) {
	return func() (svr interfaces.UserService, err error) {
		return GetUserService(opts...)
	}
}
