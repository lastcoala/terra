package domain

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	emailRegex     = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	upperCaseRegex = regexp.MustCompile(`[A-Z]`)
	lowerCaseRegex = regexp.MustCompile(`[a-z]`)
	digitRegex     = regexp.MustCompile(`[0-9]`)
)

// ValidateEmail returns an error if the email address is not valid.
func (u User) ValidateEmail() error {
	if !emailRegex.MatchString(u.Email) {
		return u.error(fmt.Errorf("invalid email: %q", u.Email), "ValidateEmail")
	}
	return nil
}

// ValidateGender returns an error if gender is not "male" or "female".
func (u User) ValidateGender() error {
	switch strings.ToLower(u.Gender) {
	case "male", "female":
		return nil
	default:
		return u.error(fmt.Errorf("invalid gender: %q, must be \"male\" or \"female\"", u.Gender), "ValidateGender")
	}
}

// ValidatePassword returns an error if the password does not meet complexity rules:
// at least 8 characters, one uppercase letter, one lowercase letter, and one digit.
func (u User) ValidatePassword() error {
	if len(u.Password) < 8 {
		return u.error(fmt.Errorf("password must be at least 8 characters long"), "ValidatePassword")
	}
	if !upperCaseRegex.MatchString(u.Password) {
		return u.error(fmt.Errorf("password must contain at least one uppercase letter"), "ValidatePassword")
	}
	if !lowerCaseRegex.MatchString(u.Password) {
		return u.error(fmt.Errorf("password must contain at least one lowercase letter"), "ValidatePassword")
	}
	if !digitRegex.MatchString(u.Password) {
		return u.error(fmt.Errorf("password must contain at least one digit"), "ValidatePassword")
	}
	return nil
}

// EncryptPassword hashes u.Password with bcrypt and stores the result back
// into u.Password. Call this on a pointer receiver so the change persists.
func (u *User) EncryptPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return u.error(err, "EncryptPassword")
	}
	u.Password = string(hashed)
	return nil
}

// Validate returns an error if the user is not valid.
func (u User) Validate(checkPassword bool) error {
	if err := u.ValidateEmail(); err != nil {
		return u.error(err, "Validate")
	}
	if err := u.ValidateGender(); err != nil {
		return u.error(err, "Validate")
	}
	if checkPassword {
		if err := u.ValidatePassword(); err != nil {
			return u.error(err, "Validate")
		}
	}
	return nil
}

// CheckPassword returns an error if the password does not match the stored hash.
func (u User) CheckPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return u.error(err, "CheckPassword")
	}
	return nil
}

func (u User) error(err error, method string, params ...any) error {
	return fmt.Errorf("User.(%v)(%v) %w", method, params, err)
}
