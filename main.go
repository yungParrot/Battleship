package main

import (
	"battleship/battleship"
	"fmt"
	"github.com/manifoldco/promptui"
  "github.com/grupawp/warships-lightgui/v2"
)


type Option struct {
  Name string
  Description string
}


const gameURL string = "https://go-pjatk-server.fly.dev/api/game"


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
    fmt.Print("\033[H\033[2J")
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
      cfg := board.NewConfig()
      battleship.Game(gameURL, cfg)
      return
    case "q":
      fmt.Println("Leaving the game...")
      return
    default:
      fmt.Printf("Unknown option %s\n", option)
    }
  }
}
