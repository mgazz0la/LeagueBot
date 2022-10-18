package sleeper

import (
	"encoding/json"
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

	players map[playerID]player
}

const (
	apiBaseURL = "https://api.sleeper.app/v1/"
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

func (s *sleeperClient) maybeLoadPlayers() (bool, error) {
	path, err := filepath.Abs("internal/platform/sleeper/players.json")
	if err != nil {
		return false, err
	}

	playerJsonFile, err := os.Open(path)
	if err != nil {
		return false, err
	}

	fileInfo, err := playerJsonFile.Stat()
	if err != nil {
		return false, err
	}

	if time.Now().Sub(fileInfo.ModTime()) > 24*time.Hour {
		fmt.Print("LOAD AGAIN PLEASE")
		// return false, err
	}

	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(fileData, &s.players)
	if err != nil {
		return false, err
	}

	for _, v := range s.players {
		if v.Position == "" {
			fmt.Println(v)
		}
	}

	return true, nil
}
