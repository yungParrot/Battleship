package main


import (
  "fmt"
  "battleship/http"
)


func main() {
  gameUrl := "https://go-pjatk-server.fly.dev/api/game"
  res := http.InitGame(gameUrl)
  fmt.Printf("Token %v\n", res)
}
