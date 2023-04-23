package main

import (
	"battleship/battleship"
	"battleship/http"
	"fmt"
	"github.com/manifoldco/promptui"
)


type Option struct {
  Name string
  Description string
}


func main() {
  for {
    mainPrompt := promptui.Select{
      Label: "Select an option",
      Items: []string{
        "h",
        "n",
        "q",
      },
    }
    _, option, err := mainPrompt.Run()
    if err != nil {
      fmt.Printf("Prompt failed %v\n", err)
      return
    }

    switch option {
    case "h":
      fmt.Println("Help:")
      fmt.Println("\th - display this help")
      fmt.Println("\tn - create a new game")
      fmt.Println("\tq - quit the app")
    case "n":
      gameUrl := "https://go-pjatk-server.fly.dev/api/game"
      authToken := http.InitGame(gameUrl)
      fmt.Printf("Token %v\n", authToken)
      rawPositions, _ := http.Board(gameUrl, authToken)
      yourCoords := battleship.ConvertToCoordinates(rawPositions)
      opponentsCoords := make(map[string]battleship.Coordinate)
      battleship.DisplayBoard(&yourCoords, &opponentsCoords)
    case "q":
      fmt.Println("Leaving the game...")
      return
    default:
      fmt.Println("Unknown option %s", option)
    }
  }
}
