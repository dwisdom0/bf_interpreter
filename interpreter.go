package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type StateMachine struct {
	prog string
	tape []byte
	ip   int
	dp   int
}

func (sm StateMachine) repr() string {
	return "ip: " + strconv.Itoa(sm.ip) + ", dp: " + strconv.Itoa(sm.dp) + ", tape: " + fmt.Sprint(sm.tape)
}

func main() {
	// first commandline argument is the path to the brainfuck program we want to run
	bf_filepath := os.Args[1]
	program, error := ioutil.ReadFile(bf_filepath)
	if error != nil {
		log.Fatal(error)
	}

	prog := string(program)
	tape := []byte{0}
	sm := StateMachine{prog, tape, 0, 0}

	fmt.Println()

	for sm.ip != len(sm.prog)-1 {
		process_next_token(&sm)
	}

	fmt.Println()
	fmt.Println()
}

func process_next_token(sm *StateMachine) {
	valid_tokens := [9]rune{'>', '<', '+', '-', '.', ',', '[', ']'}
	// search for the next valid instruction
	stop_searching := false
	for i, candidate := range sm.prog[sm.ip:] {
		if stop_searching {
			break
		}
		for _, token := range valid_tokens {
			if candidate == token {
				sm.ip = sm.ip + i
				stop_searching = true
				break
			}
		}
	}

	// if there aren't any valid instructions left in the file,
	// set the instruction pointer to the end of the program so the main loop exits
	if stop_searching == false {
		sm.ip = len(sm.prog) - 1
		return
	}
	next_token := sm.prog[sm.ip]

	// now we have the next token, so we can apply it to the state machine
	switch next_token {
	case '>':
		// move the data pointer to the right
		sm.dp++
		if sm.dp == len(sm.tape) {
			sm.tape = append(sm.tape, 0)
		}
	case '<':
		// move the data pointer to the left
		sm.dp--
		if sm.dp < 0 {
			errors.New("data pointer cannot be negative")
		}
	case '+':
		// increment the current cell
		sm.tape[sm.dp]++
	case '-':
		// decrement the current cell
		sm.tape[sm.dp]--
	case '.':
		// output the current cell as a UTF-8 character
		fmt.Print(string(rune(sm.tape[sm.dp])))
	case ',':
		// input one byte and store it in the current cell
		var b byte
		fmt.Scanf("%c", &b)
		sm.tape[sm.dp] = b
	case '[':
		// jump the instruction pointer forward to the matching ']' if the current cell is 0
		if sm.tape[sm.dp] == 0 {
			idx, err := matching_paren(sm.prog[sm.ip:])
			if err != nil {
				log.Fatal("failed to find end of loop")
			}
			// We have to add the index to the instruction pointer
			// becuase we only sent the program from the instruction pointer onwards
			// so the index is based on where the instruction pointer was
			// when we called matching_paren()
			sm.ip = sm.ip + idx
		}
	case ']':
		// jump the instruction pointer backward to the matching '[' if the current cell is nonzero
		if sm.tape[sm.dp] != 0 {
			s := sm.prog[:sm.ip+1]
			idx, err := matching_paren_backward(s)
			if err != nil {
				log.Fatal("failed to find start of loop")
			}
			sm.ip = idx
		}

	}
	// increment the instruction pointer after we've processed the instruction
	sm.ip++
}

func matching_paren(s string) (int, error) {
	if s[0] != '[' {
		return -1, errors.New("matching paren string doesn't start with an open paren")
	}
	stack := []rune{}
	for i, char := range s {
		if char == '[' {
			stack = append(stack, char)
		} else if char == ']' {
			if len(stack) == 0 {
				return -1, errors.New("extra closing paren")
			}
			// pop the matching paren off the top
			stack = stack[:len(stack)-1]
			// if the stack is empty after popping the close paren,
			// we've found the index of the matching paren we want
			if len(stack) == 0 {
				return i, nil
			}
		}
	}
	// we shouldn't ever get here
	return -1, errors.New("failed to find matching paren")
}

func matching_paren_backward(s string) (int, error) {
	if s[len(s)-1] != ']' {
		return -1, errors.New("matching paren string doesn't end with a close paren")
	}
	stack := []rune{}
	
	// search backwards
	i := len(s) - 1
	for i >= 0 {
		char := rune(s[i])
		if char == ']' {
			stack = append(stack, char)
		} else if char == '[' {
			if len(stack) == 0 {
				return -1, errors.New("extra opening paren")
			}
			// pop the matching paren off the top
			stack = stack[:len(stack)-1]
			// if the stack is empty after popping the open paren,
			// we've found the index of the matching paren we want
			if len(stack) == 0 {
				return i, nil
			}
		}
	i--
	}
	// we shouldn't ever get here
	return -1, errors.New("failed to find matching paren")
}
