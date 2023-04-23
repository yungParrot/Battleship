package battleship

import (
	"battleship/http"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)


type Coordinate struct {
  column string 
  row string
}


func getColumns() []string {
  return []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "K", "J"}
}

func getRows() []string { return []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
}


func ConvertToCoordinates(rawCoordinates []string) map[string]Coordinate {
  coordinates := make(map[string]Coordinate)
  for _, rawCoordinate := range rawCoordinates {
    column, row, _ := strings.Cut(rawCoordinate, "")
    coordinates[rawCoordinate] = Coordinate{column, row}
  }
  return coordinates
}


func DisplayBoard(yourCoords *map[string]Coordinate, opponentsCoords *map[string]Coordinate) {
  rows := getRows()
  columns := getColumns()
  for _, row := range rows {
    color.Set(color.BgBlue)
    fmt.Printf("%2s", row)
    color.Unset()
    for _, column := range columns {
      position := column + row
      if _, found := (*yourCoords)[position]; found {
        fmt.Printf("%2s", "X")
        continue
      }
      fmt.Printf("%2s", "-")
    }

    fmt.Printf("%2s", "\t|\t")

    color.Set(color.BgRed)
    fmt.Printf("%2s", row)
    color.Unset()
    for _, column := range columns {
      position := column + row
      if _, found := (*opponentsCoords)[position]; found {
        fmt.Printf("%2s", "X")
        continue
      }
      fmt.Printf("%2s", "-")
    }
    fmt.Println()
  }
  color.Set(color.BgBlue)
  fmt.Printf("%3s", "")
  for _, column := range columns {
    fmt.Printf("%s ", column)
  }
  color.Unset()

  fmt.Printf("%2s", "\t|\t")

  color.Set(color.BgRed)
  fmt.Printf("%3s", "")
  for _, column := range columns {
    fmt.Printf("%s ", column)
  }
  color.Unset()

  fmt.Println()
}


func Game(gameURL string) {
  authToken := http.InitGame(gameURL)
  fmt.Printf("Token %v\n", authToken)
  time.Sleep(time.Second)
  acceptedGameStatus := "game_in_progress"
  var gameInfo http.GetGameInfoDTO
  for gameInfo.GameStatus != acceptedGameStatus {
    gameInfo, _ = http.GetGameStatus(gameURL, authToken)   
    fmt.Print("\033[H\033[2J")
    fmt.Println("Waiting for game...")
    time.Sleep(time.Second)
  }
  fmt.Print("\033[H\033[2J")
  fmt.Printf("Token %v\n", authToken)

  fmt.Printf("You: %-25s", gameInfo.Nick)
  fmt.Print("\t|\t")
  fmt.Printf("Opp: %s\n", gameInfo.Opponent)

  rawPositions, _ := http.Board(gameURL, authToken)
  yourCoords := ConvertToCoordinates(rawPositions)
  opponentsCoords := ConvertToCoordinates(gameInfo.OppShots)

  DisplayBoard(&yourCoords, &opponentsCoords)
  for {
    gamePrompt := promptui.Select{
      Label: "Select an option\n",
      Items: []string{  
        "h",
        "i",
        "q",
      },
    }
    _, option, err := gamePrompt.Run()
    fmt.Print("\033[H\033[2J")
    if err != nil {
      fmt.Printf("Prompt failed %v\n", err)
      return
    }

    switch option {
    case "h":
      fmt.Println("Help:")
      fmt.Println("\th - display this help")
      fmt.Println("\ti - game info")
      fmt.Println("\tq - quit the app")
    case "i":
      gameDesc, _ := http.GetGameDescription(gameURL, authToken)   
      fmt.Printf("You: %s\n", gameDesc.Nick)
      fmt.Printf("%s\n", gameDesc.Desc)
      fmt.Printf("Opp: %s\n", gameDesc.Opponent)
      fmt.Printf("%s\n", gameDesc.OppDesc)
    case "q":
      fmt.Println("Leaving the game...")
      return
    default:
      fmt.Println("Unknown option %s", option)
    }
    DisplayBoard(&yourCoords, &opponentsCoords)
  }
}
