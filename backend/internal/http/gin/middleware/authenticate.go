package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
)

func (m *Middleware) authenticateWithBeared(ctx context.Context, c *gin.Context, bearer string) (bool, error) {
	apiKey := strings.TrimPrefix(bearer, "Bearer ")

	users, _, err := m.UserService.FindUsers(ctx, musichub.UserFilter{APIKey: &apiKey})
	if err != nil {
		return false, err
	} else if len(users) == 0 {
		return false, &musichub.Error{Code: musichub.EUNAUTHORIZED, Message: "Invalid API"}
	}

	c.Request = c.Request.WithContext(musichub.NewContextWithUser(ctx, users[0]))
	return true, nil
}

func (m *Middleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if bearear := c.GetHeader("Authorization"); strings.HasPrefix(bearear, "Bearer ") {
			if ok, err := m.authenticateWithBeared(ctx, c, bearear); err != nil {
				log.Printf("cannot authenticate with bearer: %s", err)
			} else if ok {
				c.Next()
				return
			}
		}

		session, _ := m.sessionManager.Session(c)

		if session.UserID != 0 {
			if user, err := m.UserService.FindUserByID(ctx, session.UserID); err != nil {
				log.Printf("cannot find session user: id=%d err=%s", session.UserID, err)
			} else {
				c.Request = c.Request.WithContext(musichub.NewContextWithUser(ctx, user))
			}
		}
		c.Next()
	}
}
