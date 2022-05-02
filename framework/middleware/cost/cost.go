package cost

import (
	"log"
	"time"

	"github.com/Pangjiping/eehe/framework/gin"
)

type Cost struct{}

// Func define cost middleware.
func (c *Cost) Func() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Printf("api uri start: %v", c.Request.RequestURI)

		c.Next()

		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri end: %v, cost: %v", c.Request.RequestURI, cost.Seconds())
	}
}
