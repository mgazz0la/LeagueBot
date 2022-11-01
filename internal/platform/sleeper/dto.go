package sleeper

type (
	LeagueID          string
	userID            string
	rosterID          int
	playerID          string
	transactionType   string
	transactionStatus string
	transactionID     string

	player struct {
		PlayerID   playerID `json:"player_id"`
		FirstName  string   `json:"first_name"`
		LastName   string   `json:"last_name"`
		Position   string   `json:"position"`
		Team       string   `json:"team"`
		SearchRank uint     `json:"search_rank"`
	}

	user struct {
		UserID      userID `json:"user_id"`
		DisplayName string `json:"display_name"`
		Metadata    struct {
			TeamName string `json:"team_name"`
		} `json:"metadata"`
		AvatarID string `json:"avatar"`
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
		Type            transactionType   `json:"type"`
		Status          transactionStatus `json:"status"`
		TransactionID   transactionID     `json:"transaction_id"`
		TimestampMillis uint              `json:"status_updated"`

		InvolvedRosters []rosterID            `json:"roster_ids"`
		Adds            map[playerID]rosterID `json:"adds"`
		Drops           map[playerID]rosterID `json:"drops"`
		WaiverBudget    []struct {
			Sender   rosterID `json:"sender"`
			Receiver rosterID `json:"receiver"`
			Amount   uint     `json:"amount"`
		} `json:"waiver_budget"`

		Settings struct {
			WaiverBid uint `json:"waiver_bid"`
		} `json:"settings"`
	}

	nflState struct {
		Week uint `json:"display_week"`
	}
)

const (
	transactionTypeFreeAgent transactionType = "free_agent"
	transactionTypeWaiver    transactionType = "waiver"
	transactionTypeTrade     transactionType = "trade"

	transactionStatusSuccess transactionStatus = "complete"
	transactionStatusFailed  transactionStatus = "failed"
)

func (r roster) PointsFor() float32 {
	return float32(r.Settings.PF) + (float32(r.Settings.PFDecimal) / 100.0)
}

func (r roster) PointsAgainst() float32 {
	return float32(r.Settings.PA) + (float32(r.Settings.PADecimal) / 100.0)
}
