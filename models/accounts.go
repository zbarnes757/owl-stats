package models

import (
	"fmt"
	"os"
	u "owl-stats/utils"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Token is a JWT claims struct
type Token struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

// Account represents the user account
type Account struct {
	Base
	ID       uuid.UUID `gorm:"type:uuid;primary_key;" jsonapi:"primary,accounts"`
	Email    string    `jsonapi:"attr,email"`
	Password string    `jsonapi:"attr,password,omitempty"`
	Token    string    `jsonapi:"attr,token,omitempty" sql:"-"`
}

// Validate the provided user information
func (account *Account) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	// Email must be unique
	temp := &Account{}

	// check for errors and duplicate emails
	err := db.Table("accounts").Where("email = ?", account.Email).First(temp).Error
	fmt.Println(err)
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

// Create a new user
func (account *Account) Create() {
	if _, ok := account.Validate(); !ok {
		return
	}

	now := time.Now()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	account.ID = uuid.NewV4()
	account.CreatedAt = now
	account.UpdatedAt = now

	db.Create(account)

	// Create new JWT token for the newly registered account
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	// delete password
	account.Password = ""
}

// Login the user
func Login(email, password string) map[string]interface{} {
	account := &Account{}
	err := db.Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))

	// Password does not match!
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials. Please try again")
	}

	// Worked! Logged In
	account.Password = ""

	// Create JWT token
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	// Store the token in the response
	account.Token = tokenString

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

// GetUser returns a user
func GetUser(u uint) *Account {

	acc := &Account{}
	db.Table("accounts").Where("id = ?", u).First(acc)

	// User not found!
	if acc.Email == "" {
		return nil
	}

	acc.Password = ""
	return acc
}
