package main

import (
  _ "embed"
  yaml "gopkg.in/yaml.v2"
  "fmt"
  "math/rand"
  "time"
  "strings"
  "flag"
)

//go:embed data/words.yaml
var words []byte

func projectName(wf WordFile) (string, error) {
  var adjective,
      substantive string
  adjective, err := wf.Choose(`common`, `adjective`)
  if err != nil {
    return ``, fmt.Errorf("could not pick an adjective: %w", err)
  }
  for {
    substantive, err = wf.Choose(`common`, `substantive`)
    if err != nil {
      return ``, fmt.Errorf("could not pick an substantive: %w", err)
    }
    if substantive != adjective {
      break
    }
  }
  return fmt.Sprintf(`%s %s`,strings.Title(adjective), strings.Title(substantive)), nil
}

func main() {
  rand.Seed(time.Now().UnixNano())
  wf := make(WordFile)

  yaml.Unmarshal(words, wf)

  var num int
  flag.IntVar(&num, `n`, 1, `Number of project titles to output`)
  flag.Parse()

  if num == 1 {
    pName, err := projectName(wf)
    if err != nil {
      fmt.Printf("Your project cannot be named: %v\n", err)
      return
    }
    fmt.Printf("Your project is now called \"%s\"\n", pName)
    return
  }

  for ; num > 0; num-- {
    pName, err := projectName(wf)
    if err != nil {
      fmt.Printf("Your project cannot be named: %v\n", err)
      return
    }
    fmt.Println(pName) 
  }
}
