package utils

import (
	"beego-crud-api/conf"
	"beego-crud-api/models"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"regexp"
	"strconv"
	"unicode"
)

var SecretKey = conf.GetEnvConst("ACCESS_SECRET")

const CodeLength = 40

func ValidateMobilePhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

	if re.MatchString(phone) {
		return true
	}
	return false
}

// CanRegistered checks if e-mail is available.
func CanRegisteredOrChanged(email string) (bool, error) {
	cond := orm.NewCondition()
	cond = cond.Or("Email", email)
	var maps []orm.Params
	o := orm.NewOrm()
	n, err := o.QueryTable("users").SetCond(cond).Values(&maps, "Email")

	if err != nil {
		return false, err
	}
	emailCheck := true

	if n > 0 {
		for _, m := range maps {
			if emailCheck && orm.ToStr(m["Email"]) == email {
				emailCheck = false
			}
		}
	}
	return emailCheck, nil
}

// CanRegisteredPhone checks if phone is available.
func CanRegisteredPhone(phone string) (bool, error) {
	cond := orm.NewCondition()
	cond = cond.Or("Phone", phone)
	var maps []orm.Params
	o := orm.NewOrm()
	n, err := o.QueryTable("users").SetCond(cond).Values(&maps, "Phone")

	if err != nil {
		return false, err
	}
	phoneCheck := true

	if n > 0 {
		for _, m := range maps {

			if phoneCheck && orm.ToStr(m["Phone"]) == phone {
				phoneCheck = false
			}
		}
	}
	return phoneCheck, nil
}

func RandomNumberString(n int) string {
	var letters = []rune("1234567890")
	s := make([]rune, n)

	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func RandomString(n int) string {
	var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func IsInt(s string) bool {

	for _, c := range s {

		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,6}$`)
	return Re.MatchString(email)
}

// get a code for email confirmation
func GetEmailConfirmationCode(user *models.Users, startInf interface{}) string {
	data := strconv.Itoa(int(user.ID)) + user.Email
	hexData := strconv.Itoa(int(user.ID)) + "|" + conf.GetEnvConst("ACCESS_SECRET")
	code := CreateEmailConfirmationCode(data)
	// add tail hex user id and secret key
	code += hex.EncodeToString([]byte(hexData))
	return code
}

// create code, it's format: 40 sha1 encoded string
func CreateEmailConfirmationCode(data string) string {
	// create sha1 encode string
	sh := sha1.New()
	sh.Write([]byte(data + SecretKey))
	encoded := hex.EncodeToString(sh.Sum(nil))
	code := fmt.Sprintf("%s", encoded)
	return code
}

// verify code
func VerifyEmailConfirmationCode(data string, code string) bool {

	if len(code) <= 18 {
		return false
	}
	// right active code
	retCode := CreateEmailConfirmationCode(data)

	if retCode == code {
		return true
	}
	return false
}
