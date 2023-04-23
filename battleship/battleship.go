package battleship


import (
  "fmt"
  "strings"
  "github.com/fatih/color"
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
