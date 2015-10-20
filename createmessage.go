package main

import (
  "bytes"
  "fmt"
	"net/http"
	"time"
  "github.com/gocraft/web"
)

// CreateMessage creates a message
func (context *AuthContext) CreateMessage(response web.ResponseWriter, request *web.Request) {
  buf := new(bytes.Buffer)
  buf.ReadFrom(request.Body)
  body := buf.String()

  currentTime := time.Now().UnixNano() / (1000 * 1000)
  job := fmt.Sprintf(`{"auth":{"uuid":"%v","token":"%v"},"http":%v,"message":%v}`, context.uuid, context.token, currentTime, body)
  conn := context.redis.Get()
  _,err := conn.Do("LPUSH", "meshblu-messages", job)
	if err != nil {
    response.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(response, `{"error": "%v"}`, err.Error())
		return
	}
  conn.Close()

  response.WriteHeader(http.StatusOK)
  fmt.Fprint(response, body)
}
