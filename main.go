package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/dylhunn/dragontoothmg"
)

var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
var table map[uint64]transposition_entry = make(map[uint64]transposition_entry)

type transposition_entry struct {
	eval  float32
	depth int
}

type eval_result struct {
	best_move dragontoothmg.Move
}

var piece_values []int = []int{0, 1, 3, 3, 5, 9, 100}

var result eval_result = eval_result{}

func determinePieceType(ourBitboardPtr *dragontoothmg.Bitboards, squareMask uint64) (dragontoothmg.Piece, *uint64) {
	var pieceType dragontoothmg.Piece = dragontoothmg.Nothing
	pieceTypeBitboard := &(ourBitboardPtr.All)
	if squareMask&ourBitboardPtr.Pawns != 0 {
		pieceType = dragontoothmg.Pawn
		pieceTypeBitboard = &(ourBitboardPtr.Pawns)
	} else if squareMask&ourBitboardPtr.Knights != 0 {
		pieceType = dragontoothmg.Knight
		pieceTypeBitboard = &(ourBitboardPtr.Knights)
	} else if squareMask&ourBitboardPtr.Bishops != 0 {
		pieceType = dragontoothmg.Bishop
		pieceTypeBitboard = &(ourBitboardPtr.Bishops)
	} else if squareMask&ourBitboardPtr.Rooks != 0 {
		pieceType = dragontoothmg.Rook
		pieceTypeBitboard = &(ourBitboardPtr.Rooks)
	} else if squareMask&ourBitboardPtr.Queens != 0 {
		pieceType = dragontoothmg.Queen
		pieceTypeBitboard = &(ourBitboardPtr.Queens)
	} else if squareMask&ourBitboardPtr.Kings != 0 {
		pieceType = dragontoothmg.King
		pieceTypeBitboard = &(ourBitboardPtr.Kings)
	}
	return pieceType, pieceTypeBitboard
}

func piece_at(board dragontoothmg.Board, square uint8) (dragontoothmg.Piece, bool) {
	var white_bitboard_ptr, black_bitboard_ptr *dragontoothmg.Bitboards
	white_bitboard_ptr = &(board.White)
	black_bitboard_ptr = &(board.Black)
	white_piece_type, _ := determinePieceType(white_bitboard_ptr, (uint64(1) << square))
	black_piece_type, _ := determinePieceType(black_bitboard_ptr, (uint64(1) << square))
	if white_piece_type != 0 {
		return white_piece_type, true
	}
	if black_piece_type != 0 {
		return black_piece_type, false
	}
	return 0, true
}

func order_moves(moves []dragontoothmg.Move, board dragontoothmg.Board) []dragontoothmg.Move {
	type move_guess struct {
		move  dragontoothmg.Move
		guess int
	}
	move_guesses := []move_guess{}
	for _, move := range moves {
		move_quality_guess := move_guess{}
		move_quality_guess.move = move
		if dragontoothmg.IsCapture(move, &board) {
			to_piece, _ := piece_at(board, move.To())
			from_piece, _ := piece_at(board, move.From())
			move_quality_guess.guess = piece_values[to_piece] - piece_values[from_piece]
		} else {
			move_quality_guess.guess = 0
		}
		move_guesses = append(move_guesses, move_quality_guess)
	}
	sort.Slice(move_guesses, func(i, j int) bool {
		return move_guesses[i].guess < move_guesses[j].guess
	})
	new_moves := []dragontoothmg.Move{}
	for _, x := range move_guesses {
		new_moves = append(new_moves, x.move)
	}
	return new_moves
}

func input() string {
	if scanner.Scan() {
		line := scanner.Text()
		return line
	}
	return ""
}

func draw_board(fen string) {
	fmt.Println(fen)
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

func evaluate(board dragontoothmg.Board) float32 {
	var eval float32 = 0
	eval_board := strings.Split(strings.Split(board.ToFen(), " ")[0], "/")
	for rank, s := range eval_board {
		for _, x := range strings.Split(s, "") {
			switch x {
			case "R":
				eval += 5
			case "N":
				switch rank {
				case 0:
					eval += 3
				case 1:
					eval += 3.05
				case 2:
					eval += 3.1
				case 3:
					eval += 3.15
				case 4:
					eval += 3.15
				case 5:
					eval += 3.1
				case 6:
					eval += 3.05
				case 7:
					eval += 3
				}
			case "B":
				eval += 3.3
			case "Q":
				eval += 9
			case "P":
				eval += 1 + 0.1*float32(7-rank)
			case "r":
				eval -= 5
			case "n":
				switch rank {
				case 0:
					eval -= 3
				case 1:
					eval -= 3.05
				case 2:
					eval -= 3.1
				case 3:
					eval -= 3.15
				case 4:
					eval -= 3.15
				case 5:
					eval -= 3.1
				case 6:
					eval -= 3.05
				case 7:
					eval -= 3
				}
			case "b":
				eval -= 3.3
			case "q":
				eval -= 9
			case "p":
				eval -= 1 + 0.1*float32(rank)
			}
		}
	}
	if board.Wtomove {
		return eval
	} else {
		return -eval
	}
}

func search(depth int, board dragontoothmg.Board, depth_from_root int, alpha float32, beta float32) float32 {
	e, ok := table[board.Hash()]
	if ok && e.depth >= depth {
		return e.eval
	}
	if depth == 0 {
		eval := evaluate(board)
		table[board.Hash()] = transposition_entry{eval, 0}
		return eval
	}
	legal_moves := order_moves(board.GenerateLegalMoves(), board)
	// legal_moves := board.GenerateLegalMoves()
	if len(legal_moves) == 0 {
		if board.OurKingInCheck() {
			return float32(-1000 - depth)
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
			table[board.Hash()] = transposition_entry{beta, depth}
			return beta
		}

	}
	table[board.Hash()] = transposition_entry{alpha, depth}
	return alpha
}

func main() {
	board := dragontoothmg.ParseFen(strings.Join(os.Args[1:], " "))
	search(5, board, 0, -1000000, 1000000)
	fmt.Println(result.best_move.String())
	// board := dragontoothmg.ParseFen(dragontoothmg.Startpos)
	// playing_white := true
	// for {
	// 	if board.Wtomove && playing_white {
	// 		move, _ := dragontoothmg.ParseMove(input())
	// 		board.Apply(move)
	// 	} else {
	// 		table = make(map[uint64]transposition_entry)
	// 		start := time.Now()
	// 		fmt.Println(search(5, board, 0, -1000000, 1000000))
	// 		elapsed := time.Since(start)
	// 		board.Apply(result.best_move)
	// 		fmt.Println(elapsed)
	// 		fmt.Println(result.best_move.String())
	// 	}
	// 	draw_board(board.ToFen())
	// 	fmt.Println()
	// }
}
