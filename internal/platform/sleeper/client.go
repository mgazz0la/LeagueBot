package sleeper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type sleeperClient struct {
	c *http.Client
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

	return rs, nil
}

func (s *sleeperClient) GetTransactions(id LeagueID, week uint) ([]*transaction, error) {
	b, err := s.get(fmt.Sprintf("league/%s/transactions/%d", id, week))
	if err != nil {
		return nil, err
	}

	var ts []*transaction
	err = json.Unmarshal(b, &ts)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

// THIS METHOD SHOULD ONLY BE CALLED AT MOST ONCE PER DAY
// The response is like 10 MB, and I don't want Sleeper to get mad at me.
func (s *sleeperClient) GetPlayers() (map[playerID]*player, error) {
	fmt.Println("ALERT: fetching new player list!")
	b, err := s.get("players/nfl")
	if err != nil {
		return nil, err
	}

	var pmap map[playerID]*player
	err = json.Unmarshal(b, &pmap)
	if err != nil {
		return nil, err
	}

	return pmap, nil
}

func (s *sleeperClient) get(path string) ([]byte, error) {
	url := apiBaseURL + path
	req, err := http.NewRequest("GET", url, nil)
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
