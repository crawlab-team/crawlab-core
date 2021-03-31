package services

import (
	"errors"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

func InitUserService() error {
	_ = CreateNewUser("admin", "admin", constants.RoleAdmin, "")
	return nil
}

func MakeToken(user *models.User) (tokenStr string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"nbf":      time.Now().Unix(),
	})

	return token.SignedString([]byte(viper.GetString("server.secret")))

}

//func GetToken(username string) (tokenStr string, err error) {
//	user, err := model.GetUserByUsername(username)
//	if err != nil {
//		log.Errorf(err.Error())
//		debug.PrintStack()
//		return
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"id":       user.Id,
//		"username": user.Username,
//		"nbf":      time.Now().Unix(),
//	})
//
//	tokenStr, err = token.SignedString([]byte(viper.GetString("server.secret")))
//	if err != nil {
//		return
//	}
//	return
//}

func SecretFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("server.secret")), nil
	}
}

func CheckToken(tokenStr string) (user models.User, err error) {
	token, err := jwt.Parse(tokenStr, SecretFunc())
	if err != nil {
		return
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return
	}

	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return
	}

	id, err := primitive.ObjectIDFromHex(claim["id"].(string))
	if err != nil {
		return user, err
	}
	username := claim["username"].(string)
	user, err = models.UserService.GetById(id)
	if err != nil {
		err = errors.New("cannot get user")
		return
	}

	if username != user.Username {
		err = errors.New("username does not match")
		return
	}

	return
}

func CreateNewUser(username string, password string, role string, email string) error {
	user := models.User{
		Username: strings.ToLower(username),
		Password: utils.EncryptPassword(password),
		Role:     role,
		Email:    email,
		Setting: models.UserSetting{
			NotificationTrigger: constants.NotificationTriggerNever,
			EnabledNotifications: []string{
				constants.NotificationTypeMail,
				constants.NotificationTypeDingTalk,
				constants.NotificationTypeWechat,
			},
		},
	}
	if err := user.Add(); err != nil {
		return err
	}
	return nil
}

func GetCurrentUser(c *gin.Context) *models.User {
	data, _ := c.Get(constants.ContextUser)
	if data == nil {
		return &models.User{}
	}
	return data.(*models.User)
}

func GetCurrentUserId(c *gin.Context) primitive.ObjectID {
	return GetCurrentUser(c).Id
}

func GetAdminUser() (user *models.User, err error) {
	u, err := models.UserService.Get(bson.M{"username": "admin"}, nil)
	if err != nil {
		return user, err
	}
	return &u, nil
}
