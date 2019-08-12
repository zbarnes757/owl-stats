package models

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// Token is a JWT claims struct
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// Account represents the user account
type Account struct {
	Base
	ID       uint   `gorm:"primary_key,AUTO_INCREMENT" jsonapi:"primary,accounts"`
	Email    string `gorm:"UNIQUE_INDEX;not null" jsonapi:"attr,email"`
	Password string `jsonapi:"attr,password,omitempty"`
	Token    string `jsonapi:"attr,token,omitempty" sql:"-"`
}

// Validate the provided user information
func (account *Account) Validate() (string, bool) {
	if !strings.Contains(account.Email, "@") {
		return "Email address is required", false
	}

	if len(account.Password) < 6 {
		return "Password is required", false
	}

	// Email must be unique
	temp := &Account{}

	// check for errors and duplicate emails
	err := db.Table("accounts").Where("email = ?", account.Email).First(temp).Error
	fmt.Println(err)
	if err != nil && err != gorm.ErrRecordNotFound {
		return "Connection error. Please retry", false
	}
	if temp.Email != "" {
		return "Email address already in use by another user.", false
	}

	return "Requirement passed", true
}

// Create a new user
func (account *Account) Create() error {
	if reason, ok := account.Validate(); !ok {
		return errors.New(reason)
	}

	now := time.Now()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
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
	return nil
}

// Login the user
func (account *Account) Login() error {
	tempPass := account.Password
	err := db.Table("accounts").Where("email = ?", account.Email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("Email address not found")
		}
		return errors.New("Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(tempPass))

	// Password does not match!
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return errors.New("Invalid login credentials. Please try again")
	}

	// Worked! Logged In
	account.Password = ""

	// Create JWT token
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	// Store the token in the response
	account.Token = tokenString

	return nil
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
