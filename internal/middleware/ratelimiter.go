package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"router-cloud-platform/internal/utils"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

func init() {
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

func getClient(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	if c, ok := clients[ip]; ok {
		c.lastSeen = time.Now()
		return c.limiter
	}
	limiter := rate.NewLimiter(10, 20)
	clients[ip] = &client{limiter: limiter, lastSeen: time.Now()}
	return limiter
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !getClient(ip).Allow() {
			utils.Error(c, http.StatusTooManyRequests, "Too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}