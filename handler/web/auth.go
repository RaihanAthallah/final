package web

import (
	"a21hc3NpZ25tZW50/client"
	"a21hc3NpZ25tZW50/service"
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/gin-gonic/gin"
)

type AuthWeb interface {
	Login(c *gin.Context)
	LoginProcess(c *gin.Context)
	Register(c *gin.Context)
	RegisterProcess(c *gin.Context)
	Logout(c *gin.Context)
}

type authWeb struct {
	userClient     client.UserClient
	sessionService service.SessionService
	embed          embed.FS
}

func NewAuthWeb(userClient client.UserClient, sessionService service.SessionService, embed embed.FS) *authWeb {
	return &authWeb{userClient, sessionService, embed}
}

func (a *authWeb) Login(c *gin.Context) {
	var filepath = path.Join("views", "auth", "login.html")
	var header = path.Join("views", "general", "header.html")

	var tmpl, err = template.ParseFS(a.embed, filepath, header)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
	}
}

func (a *authWeb) LoginProcess(c *gin.Context) {
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")

	status, err := a.userClient.Login(email, password)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	session, err := a.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	if status == 200 {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:   "session_token",
			Value:  session.Token,
			Path:   "/",
			MaxAge: 31536000,
			Domain: "",
		})

		c.Redirect(http.StatusSeeOther, "/client/dashboard")
	} else {
		c.Redirect(http.StatusSeeOther, "/client/login")
	}
}

func (a *authWeb) Register(c *gin.Context) {
	var header = path.Join("views", "general", "header.html")
	var filepath = path.Join("views", "auth", "register.html")

	var tmpl, err = template.ParseFS(a.embed, filepath, header)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
	}
}

func (a *authWeb) RegisterProcess(c *gin.Context) {

	// fmt.Printf("address: %+v\n", address)
	if err := c.Request.ParseMultipartForm(1024); err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	nik := c.Request.FormValue("nik")
	fullname := c.Request.FormValue("fullname")
	address := c.Request.FormValue("address")
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	idCard, handler, err := c.Request.FormFile("id_card")

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	allowedExtensions := []string{".jpg", ".jpeg", ".png"}
	ext := filepath.Ext(handler.Filename)
	validExtension := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			validExtension = true
			break
		}
	}

	if !validExtension {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Invalid file extension")
		return
	}

	uploadFolder := "uploads"
	if err := os.MkdirAll(uploadFolder, os.ModePerm); err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Failed to create upload folder")
		return
	}

	imageFilename := handler.Filename

	fmt.Printf("imageFilename: %+v\n", imageFilename)

	// Create a new file in the specified folder
	imagePath := filepath.Join(uploadFolder, imageFilename)
	fmt.Printf("imagePath: %+v\n", imagePath)
	imageFile, err := os.Create(imagePath)
	// encrypt.
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}
	defer imageFile.Close()

	// Copy the uploaded image to the new file
	_, err = io.Copy(imageFile, idCard)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	status, err := a.userClient.Register(nik, fullname, address, email, password, imageFilename)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	if status == 201 {
		c.Redirect(http.StatusSeeOther, "/client/login")
	} else {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Register Failed!")
	}
}

func (a *authWeb) Logout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "", false, false)
	c.Redirect(http.StatusSeeOther, "/client/login")
}
