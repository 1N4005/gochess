package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/notnil/chess"
)

var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)

func Input() string {
	if scanner.Scan() {
		line := scanner.Text()
		return line
	}
	return ""
}

func pass() {

}

func main() {
	rand.Seed(time.Now().Unix())
	game := chess.NewGame(chess.UseNotation(chess.UCINotation{}))

	input := Input()
	if !(input == "uci" || input == "gui") {
		fmt.Println("no.")
		os.Exit(1)
	}

	if input == "gui" {
		game.MoveStr("e2e4")
		game.MoveStr("e7e5")
		fmt.Println(game.Position().Board().Draw())

		os.Exit(0)
	}

	fmt.Println("uciok")
	for {
		input = Input()

		tokens := strings.Split(input, " ")
		command := tokens[0]

		switch command {
		case "quit":
			os.Exit(0)
		case "stop":
			moves := game.ValidMoves()
			move := moves[rand.Intn(len(moves))]
			fmt.Print("bestmove ")
			fmt.Println(move)
			game.Move(move)
		case "position":
			if tokens[1] == "startpos" {
				game = chess.NewGame(chess.UseNotation(chess.UCINotation{}))

				if len(tokens) > 3 && tokens[2] == "moves" {
					moves := tokens[3:]
					for i := 0; i < len(moves); i++ {
						game.MoveStr(moves[i])
					}
				}
			} else if tokens[1] == "fen" {
				fen, _ := chess.FEN(strings.Join(tokens[2:], " "))
				game = chess.NewGame(fen, chess.UseNotation(chess.UCINotation{}))
				if len(tokens) > 4 && tokens[3] == "moves" {
					moves := tokens[4:]
					for i := 0; i < len(moves); i++ {
						game.MoveStr(moves[i])
					}
				}
			}
		case "go":
		case "isready":
			fmt.Println("readyok")
		default:
			fmt.Println("awiohgweiohioghweioheg")
			os.Exit(1)
		}
		// fmt.Println(game.Position().Board().Draw())
	}
}
