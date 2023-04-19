package main


import (
  "fmt"
  "battleship/http"
  "battleship/battleship"
)


func main() {
  gameUrl := "https://go-pjatk-server.fly.dev/api/game"
  authToken := http.InitGame(gameUrl)
  fmt.Printf("Token %v\n", authToken)
  rawPositions, _ := http.Board(gameUrl, authToken)
  coordinates := battleship.ConvertToCoordinates(rawPositions)
  battleship.DisplayBoard(coordinates)
}
