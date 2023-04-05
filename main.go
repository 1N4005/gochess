package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dylhunn/dragontoothmg"
)

var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)

type eval_result struct {
	best_move dragontoothmg.Move
}

var result eval_result = eval_result{}

func input() string {
	if scanner.Scan() {
		line := scanner.Text()
		return line
	}
	return ""
}

func draw_board(fen string) {
	board := strings.Split(strings.Split(fen, " ")[0], "/")
	for _, s := range board {
		for _, x := range strings.Split(s, "") {
			switch x {
			case "1":
				fmt.Print(".")
			case "2":
				fmt.Print("..")
			case "3":
				fmt.Print("...")
			case "4":
				fmt.Print("....")
			case "5":
				fmt.Print(".....")
			case "6":
				fmt.Print("......")
			case "7":
				fmt.Print(".......")
			case "8":
				fmt.Print("........")
			default:
				fmt.Print(x)
			}
		}
		fmt.Println()
	}
}

func evaluate(board dragontoothmg.Board) int {
	eval := 0
	eval_board := strings.Split(strings.Split(board.ToFen(), " ")[0], "/")
	for _, s := range eval_board {
		for _, x := range strings.Split(s, "") {
			switch x {
			case "R":
				eval += 5
			case "N":
				eval += 3
			case "B":
				eval += 3
			case "Q":
				eval += 9
			case "P":
				eval += 1
			case "r":
				eval -= 5
			case "n":
				eval -= 3
			case "b":
				eval -= 3
			case "q":
				eval -= 9
			case "p":
				eval -= 1
			}
		}
	}
	if board.Wtomove {
		return eval
	} else {
		return -eval
	}
}

func search(depth int, board dragontoothmg.Board, depth_from_root int, alpha int, beta int) int {
	if depth == 0 {
		return evaluate(board)
	}
	legal_moves := board.GenerateLegalMoves()
	if len(legal_moves) == 0 {
		if board.OurKingInCheck() {
			return -1000 - depth
		}
		return 0
	}

	for _, move := range legal_moves {
		undo := board.Apply(move)
		eval := -search(depth-1, board, depth_from_root+1, -beta, -alpha)
		undo()
		if eval > alpha {
			alpha = eval
			if depth_from_root == 0 {
				result.best_move = move
			}
		}

		if eval >= beta {
			return beta
		}

	}
	return alpha
}

func main() {
	board := dragontoothmg.ParseFen("8/2rnk1BP/3p2B1/8/5p2/5P2/ppPN1P2/1K6 w - - 0 37")
	for {
		move, _ := dragontoothmg.ParseMove(input())
		board.Apply(move)
		draw_board(board.ToFen())
		search(5, board, 0, -1000000, 1000000)
		board.Apply(result.best_move)
		draw_board(board.ToFen())
		fmt.Println(result.best_move.String())
	}
}
