package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)

	player := "Jerome"
	wantedCount := 1000

	wg := sync.WaitGroup{}
	wg.Add(wantedCount)
	for i := 0; i < wantedCount; i++ {
		go func(w *sync.WaitGroup) {
			server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
			w.Done()
		}(&wg)
	}
	wg.Wait()

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), strconv.Itoa(wantedCount))
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetLeagueRequest())

		wantedLeague := []Player{{Name: player, Wins: wantedCount}}
		got := getLeagueFromResponse(t, response.Body)

		assertLeague(t, got, wantedLeague)
	})
}