package pkg

import "crypto/sha512"

func HashPassword(password string, salt string) string {
	sum := sha512.Sum512([]byte(password + salt))
	return string(sum[:])
}
