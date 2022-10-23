package sleeper

type (
	LeagueID string
	userID   string
	rosterID int
	playerID string

	player struct {
		PlayerID  playerID `json:"player_id"`
		FirstName string   `json:"first_name"`
		LastName  string   `json:"last_name"`
		Position  string   `json:"position"`
		Team      string   `json:"team"`
	}

	user struct {
		UserID      userID `json:"user_id"`
		DisplayName string `json:"display_name"`
		Metadata    struct {
			TeamName string `json:"team_name"`
		} `json:"metadata"`
	}

	roster struct {
		UserID   userID     `json:"owner_id"`
		RosterID rosterID   `json:"roster_id"`
		Starters []playerID `json:"starters"`
		Settings struct {
			Wins               uint `json:"wins"`
			Losses             uint `json:"losses"`
			WaiverBudgetUsed   uint `json:"waiver_budget_used"`
			PF                 uint `json:"fpts"`
			PFDecimal          uint `json:"fpts_decimal"`
			PA                 uint `json:"fpts_against"`
			PADecimal          uint `json:"fpts_against_decimal"`
			ProjectedPF        uint `json:"ppts"`
			ProjectedPFDecimal uint `json:"ppts_decimal"`
		} `json:"settings"`
		IR       []playerID `json:"reserve"`
		Players  []playerID `json:"players"`
		Metadata struct {
			Record string `json:"record"` // e.g. "LLWWW"
		} `json:"metadata"`
	}

	transaction struct {
		Type string `json:"type"`
	}
)

func (r roster) PointsFor() float32 {
	return float32(r.Settings.PF) + (float32(r.Settings.PFDecimal) / 100.0)
}

func (r roster) PointsAgainst() float32 {
	return float32(r.Settings.PA) + (float32(r.Settings.PADecimal) / 100.0)
}
