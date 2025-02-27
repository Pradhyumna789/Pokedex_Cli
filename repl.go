package main

import (
  "fmt"
  "strings"
  "bufio"
  "os"
)

func cleanInput(text string) []string {
  lowered_sliced_text := strings.Fields(strings.ToLower(text))

  return lowered_sliced_text
}

func main() {
  scanner := bufio.NewScanner(os.Stdin)

  for {
      fmt.Print("Pokedex > ")

      scanner.Scan()
      scannedText := scanner.Text() 

      displaySlice := cleanInput(scannedText) 
      displayWord := displaySlice[0]

      fmt.Println("Your command was: " + displayWord)

      if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "error in reading standard input:", err)
      }

    }
  }
