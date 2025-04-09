package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

var (
	CORS = func(r *ghttp.Request) {
		r.Response.Header().Set("Access-Control-Allow-Origin", "*")
		r.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		r.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
		r.Response.Header().Set("Access-Control-Max-Age", "3600")
		r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
		
		if r.Method == "OPTIONS" {
			r.Response.WriteStatus(200)
			r.Exit()
		}
		
		r.Middleware.Next()
	}
) 