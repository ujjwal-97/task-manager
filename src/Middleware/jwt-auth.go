package Middleware

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"../Service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var UserID primitive.ObjectID

// AuthorizeJWT validates the token from the http request, returning a 401 if it's not valid
func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			log.Println("Unauthorized Access")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]

		token, err := Service.NewJWTService().ValidateToken(tokenString)

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claims[Uid]: ", claims["Uid"])
			useridString := claims["Uid"].(string)
			if uid, err := primitive.ObjectIDFromHex(useridString); err != nil {
				log.Println(err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			} else {
				UserID = uid
			}
			//log.Println("Claims[Issuer]: ", claims["iss"])
			//log.Println("Claims[IssuedAt]: ", claims["iat"])
			//log.Println("Claims[ExpiresAt]: ", claims["exp"])
		} else {
			log.Println(err.Error())
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}