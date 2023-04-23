package http


import (
  "bytes"
  "encoding/json"
  "log"
  "net/http"
  "fmt"
  "time"
  "io/ioutil"
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


func Board(gameUrl string, authToken string) ([]string, error) {
  c := http.Client{Timeout: time.Duration(1) * time.Second}
  req, err := http.NewRequest("GET", gameUrl + "/board", nil)
  if err != nil {
    fmt.Printf("error %s", err)
    return nil, err
  }
  req.Header.Add("Accept", `application/json`)
  req.Header.Add("X-Auth-Token", authToken)

  resp, err := c.Do(req)
  if err != nil {
    fmt.Printf("Error %s", err)
    return nil, err
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Printf("Error %s", err)
    return nil, err
  }

  m := make(map[string][]string)
  err = json.Unmarshal(body, &m)
  if err != nil {
    fmt.Printf("Error %s", err)
    return nil, err
  }

  shipPositions := m["board"]
  return shipPositions, err
}
