package entities

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

var jwtKey []byte

// EncryptPassword returns the hash of the given password
func EncryptPassword(password string) ([]byte, error) {
	hpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hpass, nil
}

// HashMatch returns true if the given text hashes to the given value
func HashMatch(hash []byte, text string) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(text)) == nil
}

// VerifyPassword compares the given password with the password of the given user
func (user *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)) == nil
}

// CheckVerificationCode compares the given code with the hashed value of it
func (user *User) CheckVerificationCode(code string) bool {
	return bcrypt.CompareHashAndPassword(user.HashedVerificationCode, []byte(code)) == nil
}

// VerifyPassword compares the given code with the hashed value of it
func (otp *OTP) VerifyPassword(code string) bool {
	return bcrypt.CompareHashAndPassword(otp.HashedPassword, []byte(code)) == nil
}

// Claims defines a Clames struct for the id value
type Claims struct {
	ID uint `json:"id"`
	jwt.StandardClaims
}

// GenerateAuthToken generates a login token with jwt
// The login token expires in 7 days.
// The token is the encrypted value of user ID.
func GenerateAuthToken(id uint) (string, error) {
	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

	expirationTime := time.Now().Add(7 * 24 * 60 * time.Minute)
	claim := &Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken gets the user ID out the token
func ParseToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		log.Error(err)
		return 0, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, errors.New("could not retrieve id from token")
	}
	return claims.ID, nil
}
