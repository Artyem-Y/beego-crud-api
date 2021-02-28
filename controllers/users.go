package controllers

import (
	"beego-crud-api/conf"
	"beego-crud-api/models"
	"beego-crud-api/services/mailgun"
	"beego-crud-api/utils"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/gommon/log"
)

// Create a struct to read the email from the request body
type UserEmailData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type uEmail struct {
	Email string `json:"email"`
}

// Create a struct to read password from the request body
type UserPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ValidationCode struct {
	ValidationCode string `json:"code"`
}

type CurrentUser struct {
	ID             int64  `orm:"column(id);pk;auto" json:"id"`
	Name           string `orm:"size(128)" json:"name"`
	Email          string `orm:"size(128);unique" json:"email"`
	Phone          string `orm:"size(128);unique" json:"phone"`
	PinCode        string `orm:"size(128)" json:"pin_code"`
	EmailConfirmed bool   `orm:"size(128)" json:"email_confirmed"`
}

//  UsersController operations for Users
type UsersController struct {
	BaseController
}

// URLMapping ...
func (c *UsersController) URLMapping() {
	c.Mapping("Put", c.PutEmail)
	c.Mapping("Put", c.PutPassword)
	c.Mapping("Get", c.CheckEmail)
	c.Mapping("Get", c.GetCurrent)
	c.Mapping("Get", c.ValidateEmail)
	c.Mapping("Get", c.GetUsers)
}

func (c *UsersController) PutEmail() {
	var err error
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	s := string(c.Ctx.Input.RequestBody)
	var user *models.Users
	var userEmail UserEmailData

	if user, err = models.GetUsersById(id); err != nil {
		log.Error(err)
		c.Response(http.StatusInternalServerError, nil, err)
	}

	if err = json.Unmarshal([]byte(s), &userEmail); err != nil {
		log.Error(err)
		c.Response(http.StatusBadRequest, nil, err)
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userEmail.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		err := errors.New("wrong password, please enter the correct password")
		c.Response(http.StatusUnauthorized, nil, err)
	}
	var canChanged, _ = utils.CanRegisteredOrChanged(userEmail.Email)

	var uEmail uEmail

	if canChanged && user.Email != userEmail.Email {
		user.Email = userEmail.Email

		if err = models.UpdateUsersById(user); err != nil {
			log.Error(err)
			c.Response(http.StatusInternalServerError, nil, err)
		}
		uEmail.Email = user.Email
		c.Response(http.StatusOK, uEmail, nil)

	} else if user.Email == userEmail.Email {
		uEmail.Email = user.Email
		c.Response(http.StatusOK, uEmail, nil)

	} else {
		err := errors.New("such email already exists")
		c.Response(http.StatusConflict, nil, err)
	}
}

func (c *UsersController) PutPassword() {
	var err error
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	s := string(c.Ctx.Input.RequestBody)
	var user *models.Users
	var password UserPassword

	if user, err = models.GetUsersById(id); err != nil {
		log.Error(err)
		c.Response(http.StatusInternalServerError, nil, err)
	}

	if err = json.Unmarshal([]byte(s), &password); err != nil {
		log.Error(err)
		c.Response(http.StatusBadRequest, nil, err)
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password.OldPassword)); err != nil {
		// If the two passwords don't match, return a 401 status
		err := errors.New("wrong password, please enter the correct password")
		c.Response(http.StatusUnauthorized, nil, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), 8)

	if err != nil {
		log.Error(err)
		c.Response(http.StatusBadRequest, nil, err)
	}
	user.Password = string(hashedPassword)

	if err = models.UpdateUsersById(user); err != nil {
		log.Error(err)
		c.Response(http.StatusInternalServerError, nil, err)
	}
	c.Response(http.StatusOK, "password updated", nil)
}

func (c *UsersController) CheckEmail() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	var user *models.Users
	var err error

	if user, err = models.GetUsersById(id); err != nil {
		log.Error(err)
		c.Response(http.StatusInternalServerError, nil, err)
	}

	if user.EmailConfirmed == true {
		c.Response(http.StatusOK, true, nil)
	}
	c.Response(http.StatusBadRequest, false, nil)
}

func (c *UsersController) GetCurrent() {
	var err error
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	var user *models.Users
	var currentUser CurrentUser

	if user, err = models.GetUsersById(id); err != nil {
		log.Error(err)
		c.Response(http.StatusInternalServerError, nil, err)
	}
	currentUser.ID = user.ID
	currentUser.Name = user.Name
	currentUser.Email = user.Email
	currentUser.EmailConfirmed = user.EmailConfirmed
	currentUser.Phone = user.Phone
	c.Response(http.StatusOK, currentUser, nil)
}

func (c *UsersController) ValidateEmail() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	var user *models.Users
	var err error
	var emailConfirmationCode string

	if user, err = models.GetUsersById(id); err != nil {
		log.Error(err)
		c.Response(http.StatusInternalServerError, nil, err)
	}
	emailConfirmationCode = utils.GetEmailConfirmationCode(user, nil)
	url := conf.GetEnvConst("APP_URL") + "/active/" + emailConfirmationCode

	// send Email to forward user email
	_, err = mailgun.SendMail(
		conf.GetEnvConst("NOTIFICATION_EMAIL"),
		user.Email,
		"Email validation code",
		url,
	)

	if err != nil {
		log.Error(err)
		c.Response(http.StatusInternalServerError, nil, err)
	}
	c.Response(http.StatusOK, "Email validation url is sent", nil)
}

// GetAll ...
// @Title Get All
// @Description get Users
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Users
// @Failure 403
// @router / [get]
func (c *UsersController) GetUsers() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllUsers(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}
