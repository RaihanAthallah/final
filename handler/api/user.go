package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAPI interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserTaskCategory(c *gin.Context)
	GetUserProfile(c *gin.Context)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Register(c *gin.Context) {
	var user model.UserRegister

	fmt.Println("masuk register")

	if err := c.BindJSON(&user); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" || user.Fullname == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("register data is empty"))
		return
	}

	var recordUser = model.User{
		NIK:      user.NIK,
		Fullname: user.Fullname,
		Address:  user.Address,
		Email:    user.Email,
		Password: user.Password,
		IDCard:   user.IDCard,
	}

	fmt.Println(recordUser)

	recordUser, err := u.userService.Register(&recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse("register success"))
}

func (u *userAPI) Login(c *gin.Context) {
	var user model.UserLogin

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("login data is empty"))
		return
	}
	userLogin := model.User{
		Email:    user.Email,
		Password: user.Password,
	}
	token, err := u.userService.Login(&userLogin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}
	expirationTime := 20 * 1000000000

	c.SetCookie("session_token", *token, expirationTime, "", "", false, true)

	c.JSON(http.StatusOK, model.NewSuccessResponse("login success"))
}

func (u *userAPI) GetUserTaskCategory(c *gin.Context) {
	// userID,_ := c.Get("user_id").(int)
	userID := c.GetInt("user_id")
	var taskCategory []model.UserTaskCategory

	taskCategory, err := u.userService.GetUserTaskCategory(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid user id"))
		return
	}
	c.JSON(http.StatusOK, taskCategory)
}

func (u *userAPI) GetUserProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	userProfile, err := u.userService.GetUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid user id"))
		return
	}
	c.JSON(http.StatusOK, userProfile)
}
