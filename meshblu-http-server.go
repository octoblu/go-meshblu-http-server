package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

  "github.com/garyburd/redigo/redis"
	"github.com/gocraft/web"
)

var globalRedisPool *redis.Pool

// AuthContext stores the uuid and token
type AuthContext struct {
	uuid,token string
	redis *redis.Pool
}

// Healthcheck responds with online true
func Healthcheck(response web.ResponseWriter, request *web.Request) {
	response.WriteHeader(http.StatusOK)
	fmt.Fprint(response, `{"online": true}`)
}

// AttachRedis attaches redis connection
func (context *AuthContext) AttachRedis(response web.ResponseWriter, request *web.Request, next web.NextMiddlewareFunc){
	context.redis = globalRedisPool
	next(response,request)
}

// MeshbluAuth checks auth headers and puts them into the context
func (context *AuthContext) MeshbluAuth(response web.ResponseWriter, request *web.Request, next web.NextMiddlewareFunc) {
	if request.URL.Path == "/healthcheck" {
		next(response, request)
		return
	}
	uuid,token,ok := request.BasicAuth()
	if !ok {
		response.WriteHeader(http.StatusForbidden)
		fmt.Fprint(response, `{"error": "Not Authorized"}`)
		return
	}
	context.uuid = uuid
	context.token = token
	next(response,request)
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			// if _, err := c.Do("AUTH", password); err != nil {
			// 	c.Close()
			// 	return nil, err
			// }
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func getRedisServerName() string {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	return fmt.Sprintf("%s:%s", host, port)
}

func main() {
	var err error

	globalRedisPool = newPool(getRedisServerName())

	router := web.New(AuthContext{}).
	  Middleware((*AuthContext).AttachRedis).
	  Middleware((*AuthContext).MeshbluAuth).
		Post("/messages", (*AuthContext).CreateMessage).
		Get("/healthcheck", Healthcheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	host := fmt.Sprintf("0.0.0.0:%v", port)
	log.Printf("Listening on: %s", host)
	err = http.ListenAndServe(host, router) // Start the server!
	if err != nil {
		log.Panicf("FATAL ERRROR: %v", err.Error())
	}
}
