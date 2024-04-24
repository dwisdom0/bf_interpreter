package main

import (
	"fmt"
	"testing"
)

func TestParenSearch(t *testing.T) {
	s := "[12]12]"
	got, _ := matching_paren(s)
	want := 3
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestParenSearchNested(t *testing.T) {
	s := "[1234[12]12]"
	got, _ := matching_paren(s)
	want := len(s) - 1
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestParenSearchDeeplyNested(t *testing.T) {
	s := "[1234[12[[]1]]12]1234"
	got, _ := matching_paren(s)
	want := len(s) - 5
	fmt.Println(got, string(s[got]))
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestParenSearchBackward(t *testing.T) {
	s := "[1234[12]"
	got, _ := matching_paren_backward(s)
	want := 5
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestParenSearchBackwardNested(t *testing.T) {
	s := "[1234[12]12]"
	got, _ := matching_paren_backward(s)
	want := 0
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestBreakingLoop(t *testing.T) {
	s := `-,+[                         Read first character and start outer character reading loop
    -[                       Skip forward if character is 0
        >>++++[>++++++++<-]  Set up divisor (32) for division loop
                               (MEMORY LAYOUT: dividend copy remainder divisor quotient zero zero)
        <+<-[                Set up dividend (x minus 1) and enter division loop
            >+>+>-[>>>]`
	got, _ := matching_paren_backward(s)
	want := len(s) - 5
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
