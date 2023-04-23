package main


import (
  "fmt"
  "battleship/http"
  "battleship/battleship"
  "github.com/fatih/color"
)


func main() {
  gameUrl := "https://go-pjatk-server.fly.dev/api/game"
  authToken := http.InitGame(gameUrl)
  fmt.Printf("Token %v\n", authToken)
  rawPositions, _ := http.Board(gameUrl, authToken)
  coordinates := battleship.ConvertToCoordinates(rawPositions)

  color.Blue("\nYour board")
  battleship.DisplayBoard(&coordinates)

  color.Red("\nYour oponent's board")
  enemyBoard := make(map[string]battleship.Coordinate)
  battleship.DisplayBoard(&enemyBoard)
}
