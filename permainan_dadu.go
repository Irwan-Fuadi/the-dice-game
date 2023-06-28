package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Player struct {
	Pemain        int
	JumlahDadu    int
	Poin          int
	PlayingStatus bool
}

func initiateGame(pemain, dadu int) map[int]*Player {
	game := make(map[int]*Player)
	if pemain <= 1 {
		fmt.Println("Error: Permainan hanya bisa dilakukan jika pemain lebih dari 1!")
		return game
	}

	for i := 1; i <= pemain; i++ {
		game[i] = &Player{
			Pemain:        i,
			JumlahDadu:    dadu,
			Poin:          0,
			PlayingStatus: true,
		}
	}

	return game
}

func fetchPlayer(game map[int]*Player, initialfetcher, fetcher, points int) {
	nextFetch := fetcher + 1
	if nextFetch > len(game) {
		nextFetch = 1
	}

	if game[nextFetch].PlayingStatus {
		game[nextFetch].JumlahDadu += points
		game[initialfetcher].JumlahDadu -= points
	} else {
		fetchPlayer(game, initialfetcher, nextFetch, points)
	}
}

func beginRoll(game map[int]*Player) map[int][]int {
	rollObj := make(map[int][]int)
	fmt.Println("PUTARAN BARU")
	for i := 1; i <= len(game); i++ {
		rollResult := make([]int, 0)
		for j := 1; j <= game[i].JumlahDadu; j++ {
			time.Sleep(time.Millisecond * 50) // Jeda untuk mensimulasikan lemparan dadu
			rand.Seed(time.Now().UnixNano())
			roll := rand.Intn(6) + 1
			rollResult = append(rollResult, roll)
		}
		rollObj[i] = rollResult
		fmt.Printf("Pemain %d mengocok dadu: %v\n", i, rollResult)
	}
	return rollObj
}

func evaluatePlayer(game map[int]*Player, p Player, rollRound map[int][]int) {
	fetch := 0
	for _, val := range rollRound[p.Pemain] {
		if val == 1 {
			fetch++
		}
	}

	poin := 0
	for _, val := range rollRound[p.Pemain] {
		if val == 6 {
			poin++
		}
	}

	if p.JumlahDadu != 0 && p.PlayingStatus {
		game[p.Pemain].Poin += poin
		game[p.Pemain].JumlahDadu -= poin
		if fetch > 0 {
			fetchPlayer(game, p.Pemain, p.Pemain, fetch)
		}
	}
}

func findWinner(game map[int]*Player) []int {
	pointsVal := make([]int, 0)
	for _, player := range game {
		pointsVal = append(pointsVal, player.Poin)
	}

	winnerVal := pointsVal[0]
	for _, val := range pointsVal {
		if val > winnerVal {
			winnerVal = val
		}
	}

	winnerArr := make([]int, 0)
	for i, player := range game {
		if player.Poin == winnerVal {
			winnerArr = append(winnerArr, i)
		}
	}

	return winnerArr
}

func play(game map[int]*Player, currentState int) {
	if currentState == 1 {
		winner := findWinner(game)
		fmt.Println("Permainan Selesai! Pemenangnya adalah:", winner)
		return
	}

	roundResult := beginRoll(game)
	for _, player := range game {
		evaluatePlayer(game, *player, roundResult)
	}

	for i, player := range game {
		if player.JumlahDadu == 0 && game[i].PlayingStatus {
			game[i].PlayingStatus = false
			currentState = currentState - 1
		}
	}

	fmt.Println("EVALUASI:")
	for _, player := range game {
		gameLog := fmt.Sprintf("Pemain %d memiliki poin %d dengan sisa dadu %d", player.Pemain, player.Poin, player.JumlahDadu)
		if player.PlayingStatus {
			fmt.Println(gameLog)
		} else {
			fmt.Println(gameLog + " (SELESAI BERMAIN)")
		}
	}

	play(game, currentState)
}

func main() {
	pemain := 3
	dadu := 4
	game := initiateGame(pemain, dadu)
	play(game, pemain)
}
