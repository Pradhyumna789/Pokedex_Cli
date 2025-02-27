package main

import (
  "testing"
)

func TestCleanInput(t *testing.T) {
  cases := []struct {
    input string
    expected []string
  }{
    {
      input: "  hello world!  ",
      expected: []string{"hello", "world!"},
    },

    {
      input: "who are you?",
      expected: []string{"who", "are", "you?"},
    },

    {
      input: "existence is a paradox",
      expected: []string{"existence", "is", "a", "paradox"},
    },

  }

  for _, c := range cases {
    actual := cleanInput(c.input)
    if len(actual) != len(c.expected) {
      t.Errorf("Make sure the lengths of the actual slice and the expected slice match")
    } 

    for i := range actual {
      actualWord := actual[i]
      expectedWord := c.expected[i]
      if actualWord != expectedWord {
        t.Errorf("Make sure the actual word is the same as the expected word")
      }
    }
  }
}


