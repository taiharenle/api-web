package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"time"
)

func GinJWTMiddlewareInit() (authMiddleware *jwt.GinJWTMiddleware) {
	authMiddleware, _ = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "web",
		Key:         []byte("secret key"),
		Timeout:     3 * 24 * time.Hour,
		MaxRefresh:  3 * 24 * time.Hour,
		IdentityKey: "session",
		//PayloadFunc: func(data interface{}) jwt.MapClaims {
		//	if v, ok := data.(*model.Session); ok {
		//		return jwt.MapClaims{
		//			"id":        v.ID,
		//			"type":      v.Type,
		//			"accountid": v.AccountID,
		//		}
		//	}
		//	return jwt.MapClaims{}
		//},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		TokenLookup:   "header: Authorization, query: token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	return
}
