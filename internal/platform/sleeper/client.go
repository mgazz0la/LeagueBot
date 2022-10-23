package sleeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type sleeperClient struct {
	c *http.Client
}

const (
	apiBaseURL  = "https://api.sleeper.app/v1/"
	playersJSON = "internal/platform/sleeper/players.json"
)

func newSleeperClient() *sleeperClient {
	s := new(sleeperClient)
	s.c = new(http.Client)

	return s
}

func (s *sleeperClient) GetUsers(id LeagueID) ([]user, error) {
	b, err := s.get(fmt.Sprintf("league/%s/users", id))
	if err != nil {
		return nil, err
	}

	var us []user
	err = json.Unmarshal(b, &us)
	if err != nil {
		return nil, err
	}

	// sort by User ID
	sort.Slice(us, func(i, j int) bool {
		return us[i].UserID < us[j].UserID
	})

	return us, nil
}

func (s *sleeperClient) GetRosters(id LeagueID) ([]roster, error) {
	b, err := s.get(fmt.Sprintf("league/%s/rosters", id))
	if err != nil {
		return nil, err
	}

	var rs []roster
	err = json.Unmarshal(b, &rs)
	if err != nil {
		return nil, err
	}

	// sort by User ID
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].UserID < rs[j].UserID
	})

	return rs, nil
}

func (s *sleeperClient) GetTransactions(id LeagueID, week uint) ([]transaction, error) {
	b, err := s.get(fmt.Sprintf("league/%s/transactions/%d", id, week))
	if err != nil {
		return nil, err
	}

	var ts []transaction
	err = json.Unmarshal(b, &ts)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (s *sleeperClient) get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", apiBaseURL, path), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (s *sleeperClient) GetPlayers() (map[playerID]*player, error) {
	path, err := filepath.Abs(playersJSON)
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) || time.Now().Sub(stat.ModTime()) > 24*time.Hour {
		if err = s.fetchPlayerList(); err != nil {
			return nil, err
		}
	}

	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var players map[playerID]*player
	err = json.Unmarshal(fileData, &players)
	if err != nil {
		return nil, err
	}

	return players, nil
}

// This endpoint should be called at most once per day, therefore we save the results
// in a known file location
func (s *sleeperClient) fetchPlayerList() error {
	fmt.Println("ALERT: fetching new player list!")
	resp, err := s.get("players/nfl")
	if err != nil {
		return err
	}

	f, err := os.Create(playersJSON)
	if err != nil {
		return err
	}

	_, err = f.Write(resp)
	if err != nil {
		return err
	}

	return nil
}
