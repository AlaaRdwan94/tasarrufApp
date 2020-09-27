package entities

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ErrorMessage is the standard error message returned
// upon error in api request.
type ErrorMessage struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SendDBError returns a DB error message
func SendDBError(c *gin.Context, message string, err error) {
	log.Error(err)
	c.AbortWithStatusJSON(500, gin.H{
		"error":   fmt.Sprintf("DB Error : %s", err.Error()),
		"message": message,
	})
}

// SendValidationError return a validation error message
func SendValidationError(c *gin.Context, message string, err error) {
	log.Error(err)
	c.AbortWithStatusJSON(400, gin.H{
		"error":   fmt.Sprintf("validation error : %s", err.Error()),
		"message": message,
	})
}

// SendNotFoundError retruns a not found error message
func SendNotFoundError(c *gin.Context, message string, err error) {
	log.Error(err)
	c.AbortWithStatusJSON(404, gin.H{
		"error":   fmt.Sprintf("asset not found : %s", err.Error()),
		"message": message,
	})
}

// SendAuthError returns an authentication error message
func SendAuthError(c *gin.Context, message string, err error) {
	log.Error(err)
	c.AbortWithStatusJSON(401, gin.H{
		"error":   fmt.Sprintf("authorization error : %s", err.Error()),
		"message": message,
	})
}

// SendParsingError return a parsing error message
func SendParsingError(c *gin.Context, message string, err error) {
	log.Error(err)
	c.AbortWithStatusJSON(401, gin.H{
		"error":   fmt.Sprintf("error parsing json body: %s", err.Error()),
		"message": message,
	})
}

// SendServerError return a parsing error message
func SendServerError(c *gin.Context, message string, err error) {
	log.Error(err)
	c.AbortWithStatusJSON(401, gin.H{
		"error":   fmt.Sprintf("server error : %s", err.Error()),
		"message": message,
	})
}
