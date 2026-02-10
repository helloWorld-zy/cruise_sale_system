package middleware

import (
	"log"
	"net/http"

	"cruise_booking_system/internal/data"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

var Enforcer *casbin.Enforcer

func InitCasbin() error {
	adapter, err := gormadapter.NewAdapterByDB(data.DB)
	if err != nil {
		return err
	}

	// Load model configuration
    text := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
    m, err := model.NewModelFromString(text)
	if err != nil {
		return err
	}

	Enforcer, err = casbin.NewEnforcer(m, adapter)
	if err != nil {
		return err
	}

	return Enforcer.LoadPolicy()
}

func CasbinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			// If no role, maybe guest or unauthorized. 
            // For now, strict RBAC: require role.
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No role found in context"})
			return
		}

		obj := c.Request.URL.Path
		act := c.Request.Method
        
        sub := role.(string)

		ok, err := Enforcer.Enforce(sub, obj, act)
		if err != nil {
            log.Printf("Casbin enforce error: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Authorization error"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}

		c.Next()
	}
}
