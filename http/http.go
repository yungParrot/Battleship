package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)


type GetGameInfoDTO struct {
  Desc string `json:"desc"`
  GameStatus string `json:"game_status"`
  LastGameStatus string `json:"last_game_status"`
  Nick string `json:"nick"`
  OppDesc string `json:"opp_desc"`
  OppShots []string `json:"opp_shots"`
  Opponent string `json:"opponent"`
  ShouldFire bool `json:"should_fire"`
  Timer int `json:"timer"`
}


func SendGetRequest(URL string, authToken string) ([]byte, error) {
  client := &http.Client{}
  req, err := http.NewRequest("GET", URL, nil)
  if err != nil {
    fmt.Printf("Error: %s", err)
    return nil, err
  }
  req.Header.Add("Accept", `application/json`)
  req.Header.Add("X-Auth-Token", authToken)
  resp, err := client.Do(req); 
  if err != nil {
    fmt.Printf("Error: %s", err)
    return nil, err
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Printf("Error: %s", err)
    return nil, err
  }
  return body, nil
}


func InitGame(gameURL string) string {
  data := map[string]bool{"wpbot": true}
  postBody, _ := json.Marshal(data)
  responseBody := bytes.NewBuffer(postBody)
  res, err := http.Post(gameURL, "application/json", responseBody)
  if err != nil {
    log.Fatalf("An Error Occured %v", err)
  }
  token := res.Header.Get("x-auth-token")
  return token
}


func Board(gameURL string, authToken string) ([]string, error) {
  c := &http.Client{}
  req, err := http.NewRequest("GET", gameURL + "/board", nil)
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


func GetGameStatus(gameURL string, authToken string) (GetGameInfoDTO, error) {
  data, err := SendGetRequest(gameURL, authToken)
  var gameDesc GetGameInfoDTO
  if err = json.Unmarshal(data, &gameDesc); err != nil {
    fmt.Printf("GetGameStatus error: %s\n", err)
    return GetGameInfoDTO{}, err
  }
  return gameDesc, err
}

func GetGameDescription(gameURL string, authToken string) (GetGameInfoDTO, error) {
  data, err := SendGetRequest(gameURL + "/desc", authToken)
  var gameDesc GetGameInfoDTO
  if err := json.Unmarshal(data, &gameDesc); err != nil {
    fmt.Printf("GetGameDescription error: %s\n", err)
    return GetGameInfoDTO{}, err
  }
  return gameDesc, err
}
