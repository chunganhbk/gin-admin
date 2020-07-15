package util

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// MD5Hash MD5
func MD5Hash(b []byte) string {
	h := md5.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// MD5HashString MD5
func MD5HashString(s string) string {
	return MD5Hash([]byte(s))
}

func bcryptPwd(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
func BcryptPwd(pwd string) string {
	return bcryptPwd([]byte(pwd))
}
func comparePasswords(plainPwd []byte, hashedPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
func ComparePasswords(plainPwd string, hashedPwd string) bool {
	return comparePasswords([]byte(plainPwd), hashedPwd)
}
