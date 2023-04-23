package battleship

import (
	"battleship/http"
	"bufio"
	"fmt"
	"os"
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


func ValidateCoords(rawCoordinate string) bool {
  valid := false
  var availableOptions []string
  for _, column := range getColumns() {
    for _, row := range getRows() {
      availableOptions = append(availableOptions, column + row)
    }
  }
  for _, option := range availableOptions {
    if rawCoordinate == option {
      valid = true
      break
    }
  }
  return valid
}


func ConvertToCoordinates(rawCoordinates []string) map[string]Coordinate {
  coordinates := make(map[string]Coordinate)
  for _, rawCoordinate := range rawCoordinates {
    column, row, _ := strings.Cut(rawCoordinate, "")
    coordinates[rawCoordinate] = Coordinate{column, row}
  }
  return coordinates
}


func DisplayBoard(
  yourCoords *map[string]Coordinate,
  yourShots *map[string]Coordinate,
  opponentsCoords *map[string]Coordinate,
  opponentsShots *map[string]Coordinate,
) {
  rows := getRows()
  columns := getColumns()
  for _, row := range rows {
    color.Set(color.BgBlue)
    fmt.Printf("%2s", row)
    color.Unset()
    for _, column := range columns {
      position := column + row
      _, yoursFound := (*yourCoords)[position] 
      _, oppsFound := (*opponentsShots)[position] 
      if yoursFound && oppsFound {
        color.Set(color.FgRed)
        fmt.Printf("%2s", "X")
        color.Unset()
        continue
      }
      if yoursFound {
        fmt.Printf("%2s", "X")
        continue
      }
      if oppsFound {
        fmt.Printf("%2s", "o")
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
      _, yoursFound := (*yourShots)[position] 
      _, oppsFound := (*opponentsCoords)[position] 
      if yoursFound && oppsFound {
        color.Set(color.FgGreen)
        fmt.Printf("%2s", "X")
        color.Unset()
        continue
      }
      if yoursFound {
        fmt.Printf("%2s", "o")
        continue
      }
      if oppsFound {
        fmt.Printf("%2s", "O")
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


func DisplayPlayers(status *http.GetGameInfoDTO) {
  fmt.Printf("You: %-25s", status.Nick)
  fmt.Print("\t|\t")
  fmt.Printf("Opp: %s\n", status.Opponent)
  if status.ShouldFire {
    color.Set(color.BgBlue)
    fmt.Printf("Your turn")
    color.Unset()
    fmt.Printf("%-21s", "")
    fmt.Print("\t|\t\n")
  } else {
    fmt.Printf("%-30s", "")
    fmt.Print("\t|\t")
    color.Set(color.BgRed)
    fmt.Printf("Opp's turn\n")
    color.Unset()
  }
}

func Game(gameURL string) {
  authToken := http.InitGame(gameURL)
  fmt.Printf("Token %v\n", authToken)
  time.Sleep(time.Second)

  acceptedGameStatus := "game_in_progress"
  var gameInfo http.GetGameInfoDTO
  var oppRawCoords []string
  var yourRawShots []string
  for gameInfo.GameStatus != acceptedGameStatus {
    gameInfo, _ = http.GetGameStatus(gameURL, authToken)   
    fmt.Print("\033[H\033[2J")
    fmt.Println("Waiting for game...")
    time.Sleep(time.Second)
  }
  fmt.Print("\033[H\033[2J")
  fmt.Printf("Token %v\n", authToken)

  DisplayPlayers(&gameInfo)

  yourRawCoords, _ := http.Board(gameURL, authToken)
  yourCoords := ConvertToCoordinates(yourRawCoords)
  yourShots := ConvertToCoordinates(yourRawShots)
  opponentsCoords := ConvertToCoordinates(oppRawCoords)
  oppShots := ConvertToCoordinates(gameInfo.OppShots)

  DisplayBoard(&yourCoords, &yourShots, &opponentsCoords, &oppShots)
  for {
    if !gameInfo.ShouldFire {
      gameInfo, _ = http.GetGameStatus(gameURL, authToken)   
      time.Sleep(time.Second)
      fmt.Print("\033[H\033[2J")
      DisplayPlayers(&gameInfo)
      oppShots := ConvertToCoordinates(gameInfo.OppShots)
      DisplayBoard(&yourCoords, &yourShots, &opponentsCoords, &oppShots)
      fmt.Printf("Opp shot at %q", gameInfo.OppShots)
      continue
    }
    gamePrompt := promptui.Selec{
      Label: "Select an option\n",
      Items: []string{  
        "h",
        "i",
        "f",
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
      fmt.Println("\tf - fire at the opp")
      fmt.Println("\tq - quit the app")
    case "i":
      gameDesc, _ := http.GetGameDescription(gameURL, authToken)   
      fmt.Printf("You: %s\n", gameDesc.Nick)
      fmt.Printf("%s\n", gameDesc.Desc)
      fmt.Printf("Opp: %s\n", gameDesc.Opponent)
      fmt.Printf("%s\n", gameDesc.OppDesc)
    case "f":
      stillShooting := true
      for stillShooting {
        DisplayPlayers(&gameInfo)
        yourShots = ConvertToCoordinates(yourRawShots)
        opponentsCoords = ConvertToCoordinates(oppRawCoords)
        oppShots = ConvertToCoordinates(gameInfo.OppShots)
        DisplayBoard(&yourCoords, &yourShots, &opponentsCoords, &oppShots)
        fmt.Print("Coords [A-J][1-10] ↩️: ")
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        coord := scanner.Text()
        if !ValidateCoords(coord) {
          fmt.Printf("Coordinate not available: %q\n", coord)
          time.Sleep(time.Millisecond * 1500)
          fmt.Print("\033[H\033[2J")
          continue
        }
        fmt.Printf("Firing at %q...\t", coord)
        result, _ := http.FireAtOpp(gameURL, authToken, coord)
        yourRawShots = append(yourRawShots, coord)
        fmt.Printf("%s\n", result.Result)
        time.Sleep(time.Second)
        stillShooting = result.Result == "hit"
        if stillShooting {
          oppRawCoords = append(oppRawCoords, coord)
        }
        fmt.Print("\033[H\033[2J")
      }
    case "q":
      fmt.Println("Leaving the game...")
      return
    default:
      fmt.Printf("Unknown option %s\n", option)
    }
    gameInfo, _ = http.GetGameStatus(gameURL, authToken)   

    DisplayPlayers(&gameInfo)
    yourShots = ConvertToCoordinates(yourRawShots)
    opponentsCoords = ConvertToCoordinates(oppRawCoords)
    oppShots = ConvertToCoordinates(gameInfo.OppShots)
    DisplayBoard(&yourCoords, &yourShots, &opponentsCoords, &oppShots)
  }
}
