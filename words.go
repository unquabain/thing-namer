package main

import (
  "math/rand"
  "fmt"
)

type WeightedWord struct {
  Word string
  Weight int
}


type WordList struct {
  WeightedWords []WeightedWord
}

func (wl *WordList) Choose() string {
  var winner string
  var sum int
  for _, ww := range wl.WeightedWords {
    sum += ww.Weight
    w := rand.Intn(sum)
    if w <= ww.Weight {
      winner = ww.Word
    }
  }
  return winner
}

func (wl *WordList) Add(other *WordList) *WordList {
  sum := new(WordList)
  sum.WeightedWords = make([]WeightedWord, len(wl.WeightedWords), len(wl.WeightedWords) + len(other.WeightedWords))
  copy(sum.WeightedWords, wl.WeightedWords)
  sum.WeightedWords = append(sum.WeightedWords, other.WeightedWords...)
  return sum
}

func (wl *WordList) UnmarshalYAML(unmarshal func(interface{}) error) error {
  m := make(map[string]int)
  if err := unmarshal(&m); err != nil {
    return fmt.Errorf(`unable to read a wordlist: %w`, err)
  }
  for word, val := range m {
    wl.WeightedWords = append(
      wl.WeightedWords,
      WeightedWord{
        Word: word,
        Weight: val,
      },
    )
  }
  return nil
}

type WordFile map[string]WordList

func (wf WordFile) Choose(lists ...string) (string, error) {
  wl := new(WordList)
  for _, listName := range lists {
    newList, ok := wf[listName]
    if !ok {
      return ``, fmt.Errorf(`no such word list %q`, listName)
    }
    wl = wl.Add(&newList)
  }
  return wl.Choose(), nil
}
