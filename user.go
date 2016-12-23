package main

import (
	"crypto/md5"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string
	Email          string
	HashedPassword string
	UserName       string
}

const (
	passwordLength = 8
	hashCost       = 10
	userIDLength   = 16
)

// Create new user and persist to store
func NewUser(username, email, password string) (User, error) {
	user := User{
		Email:    email,
		UserName: username,
	}

	// Check if empty
	if username == "" {
		return user, errNoUserName
	}
	if email == "" {
		return user, errNoEmail
	}
	if password == "" {
		return user, errNoPassword
	}

	// Check if username exists
	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, errUsernameExist
	}
	// Check if email existingUser
	existingEmail, err := globalUserStore.FindByEmail(email)
	if err != nil {
		return user, err
	}
	if existingEmail != nil {
		return user, errEmailExist
	}
	// Check password length
	if len(password) < passwordLength {
		return user, errPasswordTooShort
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)

	user.HashedPassword = string(hashedPassword)
	user.ID = GenerateID("usr", userIDLength)

	// Persist user
	err = globalUserStore.Save(user)
	if err != nil {
		panic(err)
	}
	return user, err
}

func FindUser(username, password string) (*User, error) {
	out := &User{
		UserName: username,
	}

	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return out, err
	}
	if existingUser == nil {
		return out, errCredentialsIncorrect
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.HashedPassword), []byte(password))
	if err != nil {
		return out, errCredentialsIncorrect
	}
	return existingUser, nil
}

func UpdateUser(user *User, email, currentPassword, newPassword string) (User, error) {
	out := *user
	out.Email = email

	// Check if email exist
	existingUser, err := globalUserStore.FindByEmail(email)
	if err != nil {
		return out, err
	}
	if existingUser != nil && existingUser.ID != user.ID {
		return out, errEmailExist
	}

	// At this point, we can update the email address
	// Since this is a pointer, the calling code will see this updated email
	// address
	user.Email = email

	// No current password? Don't try update the password.
	if currentPassword == "" {
		return out, nil
	}

	if bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(currentPassword)) != nil {
		return out, errPasswordIncorrect
	}

	if newPassword == "" {
		return out, errNoPassword
	}

	if len(newPassword) < passwordLength {
		return out, errPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), hashCost)
	user.HashedPassword = string(hashedPassword)

	return out, err
}

func (user *User) AvatarURL() string {
	return fmt.Sprintf("//www.gravatar.com/avatar/%x", md5.Sum([]byte(user.Email)))
}

func (user *User) ImagesRoute() string {
	return "/user/" + user.ID
}
