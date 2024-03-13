package passwd

import (
	"crypto/sha512"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/anaskhan96/go-password-encoder"
)

var options = &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}

// Encode options := &password.Options{SaltLen: 10, Iterations: 10000, KeyLen: 50, HashFunction: passwd.New}
func Encode(rawPassword string) string {
	salt, encodedPwd := password.Encode(rawPassword, options)
	dbPassword := fmt.Sprintf("$sha512$%s$%s", salt, encodedPwd)
	// fmt.Println(dbPassword)
	// fmt.Println(encodedPwd)
	return dbPassword
}

func Verify(rawPassword string, dbPassword string) bool {
	splits := strings.Split(dbPassword, "$")
	if len(splits) != 4 {
		zap.L().Error("Invalid password format")
		return false
	}
	// [ sha512 WI56n0Wmte5Ul0ui f89fdc8f8c7dc87c220d2331007b48d53e5eb2a9f64d2a31f0cedb4dd3f7c874]
	// fmt.Println(splits) // splits[0]是空字符 splits[2]才是盐
	valid := password.Verify(rawPassword, splits[2], splits[3], options)
	// fmt.Println(valid) // true
	return valid
}
