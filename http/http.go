package http


import (
   "bytes"
   "encoding/json"
   "log"
   "net/http"
)


func InitGame(gameUrl string) string {
  postBody, _ := json.Marshal(map[string]string{})
  responseBody := bytes.NewBuffer(postBody)
  res, err := http.Post(gameUrl, "application/json", responseBody)
  if err != nil {
    log.Fatalf("An Error Occured %v", err)
  }
  token := res.Header.Get("x-auth-token")
  return token
}
