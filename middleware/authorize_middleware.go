package middleware

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		// tokenString, err := c.Cookie("token")
		authentication := c.Request.Header["Authorization"]
		if len(authentication) < 1 {
			c.Set("isLogin", false)
			return
			c.Next()
		}
		tokenString := strings.Split(authentication[0], " ")[1]
		// if err != nil {
		// 	c.Set("isLogin", false)
		// 	return
		// 	c.Next()
		// } 
			
		var mySigningKey = []byte("pa$$w0rd")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Error in parsing")
			}
			return mySigningKey, nil
		})
	
		if err != nil {
			c.Set("isLogin", false)
			return
			c.Next()
		} else {
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				c.Set("isLogin", true)
				c.Set("email", claims["email"])
				c.Set("userID", claims["userID"])
				c.Set("role", claims["role"])
				return
				c.Next()
					
			}
		}
		
	}
}