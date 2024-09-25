package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

var mySigningKey = []byte(os.Getenv("JWT_SIGNING_KEY"))

func getRefreshInterval() time.Duration {
	durationTime, err := time.ParseDuration("168h")
	if err != nil {
		panic("invalid auth refresh interval")
	}
	return durationTime
}

var ErrTokenExpired = jwt.ErrTokenExpired

func CreateToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Returns the existing token if it's not expired, returns the refreshed token if it's expired and in the refreshed window, returns a jwt.ErrorTokenExpired if expired and not able to be refreshed.
// Caller should store the response.
func ValidateToken(userId string, tokenString string) (newTokenString string, err error) {
	token, err := parseToken(tokenString)
	if err != nil {
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenUserId := claims["user_id"].(string)
		exp := claims["exp"].(float64)
		if tokenUserId == userId && int64(exp) > time.Now().Unix() {
			fmt.Println("Token is valid and user is authorized")
			return tokenString, nil
		} else if tokenUserId == userId && int64(exp) > time.Now().Add(getRefreshInterval()).Unix() {
			fmt.Println("Token is valid and within the refresh window")
			// TODO: review
			return CreateToken(uuid.Must(uuid.FromString(tokenUserId)))
		} else if tokenUserId == userId {
			return "", jwt.ErrTokenExpired
		} else {
			fmt.Println("Token is invalid or user is not authorized")
			return "", jwt.ErrTokenUnverifiable
		}
	} else {
		fmt.Println("Token is invalid")
		return "", jwt.ErrTokenMalformed
	}
}

func parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
