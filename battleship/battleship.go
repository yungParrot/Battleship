package battleship

import (
	"battleship/http"
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
  "github.com/grupawp/warships-lightgui/v2"
)


func getColumns() []string {
  return []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
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


func createBorder(gameBoard *board.Board, shipCoords []string) {

  gameBoard.CreateBorder()

}


func DisplayBoard(
  gameBoard *board.Board,
  yourCoords []string,
  yourShots []string,
  opponentsCoords []string,
  opponentsShots []string,
) {
  for _, coord := range yourCoords {
    gameBoard.Set(board.Left, coord, board.Ship)
  }
  for _, coord := range opponentsCoords {
    gameBoard.Set(board.Right, coord, board.Ship)
  }
  for _, coord := range yourShots {
    state, _ := gameBoard.HitOrMiss(board.Right, coord)
    gameBoard.Set(board.Right, coord, state)
  }
  for _, coord := range opponentsShots {
    state, _ := gameBoard.HitOrMiss(board.Left, coord)
    gameBoard.Set(board.Left, coord, state)
  }
  gameBoard.Display()
}


func DisplayPlayers(status *http.GetGameInfoDTO) {
  fmt.Printf("You: %-25s", status.Nick)
  fmt.Printf("Opp: %s\n", status.Opponent)
  if status.ShouldFire {
    color.Set(color.BgBlue)
    fmt.Printf("Your turn")
    color.Unset()
    fmt.Printf("%-21s", "")
    fmt.Println()
  } else {
    fmt.Printf("%-30s", "")
    color.Set(color.BgRed)
    fmt.Printf("Opp's turn\n")
    color.Unset()
  }
}

func Game(gameURL string, config *board.Config) {
  authToken := http.InitGame(gameURL)
  fmt.Printf("Token %v\n", authToken)
  time.Sleep(time.Second)

  acceptedGameStatus := "game_in_progress"
  var gameInfo http.GetGameInfoDTO
  var oppCoords []string
  var yourShots []string
  for gameInfo.GameStatus != acceptedGameStatus {
    gameInfo, _ = http.GetGameStatus(gameURL, authToken)   
    fmt.Print("\033[H\033[2J")
    fmt.Println("Waiting for game...")
    time.Sleep(time.Second)
  }
  fmt.Print("\033[H\033[2J")
  fmt.Printf("Token %v\n", authToken)

  gameBoard := board.New(config)

  yourCoords, _ := http.Board(gameURL, authToken)
  oppShots := gameInfo.OppShots

  DisplayBoard(gameBoard, yourCoords, yourShots, oppCoords, oppShots)
  DisplayPlayers(&gameInfo)
  for {
    if !gameInfo.ShouldFire {
      gameInfo, _ = http.GetGameStatus(gameURL, authToken)   
      time.Sleep(time.Second)
      fmt.Print("\033[H\033[2J")
      oppShots := gameInfo.OppShots
      DisplayBoard(gameBoard, yourCoords, yourShots, oppCoords, oppShots)
      DisplayPlayers(&gameInfo)
      fmt.Printf("Opp shot at %q", gameInfo.OppShots)
      continue
    }
    gamePrompt := promptui.Select{
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
      // TODO: this should be displayed later
      gameDesc, _ := http.GetGameDescription(gameURL, authToken)   
      fmt.Printf("You: %s\n", gameDesc.Nick)
      fmt.Printf("%s\n", gameDesc.Desc)
      fmt.Printf("Opp: %s\n", gameDesc.Opponent)
      fmt.Printf("%s\n", gameDesc.OppDesc)
    case "f":
      stillShooting := true
      for stillShooting {
        oppShots = gameInfo.OppShots
        DisplayBoard(gameBoard, yourCoords, yourShots, oppCoords, oppShots)
        DisplayPlayers(&gameInfo)
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
        yourShots = append(yourShots, coord)
        fmt.Printf("%s\n", result.Result)
        time.Sleep(time.Second)
        hit := result.Result == "hit"
        sunk := result.Result == "sunk"
        stillShooting = hit || sunk
        if stillShooting {
          oppCoords = append(oppCoords, coord)
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
    oppShots = gameInfo.OppShots
    DisplayBoard(gameBoard, yourCoords, yourShots, oppCoords, oppShots)
    DisplayPlayers(&gameInfo)
  }
}
