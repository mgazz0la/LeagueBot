package domain

type (
	PlayerStats struct {
		Player

		// Passing stats
		PassingAttempts    uint
		PassingCompletions uint
		PassingYards       uint
		PassingTouchdowns  uint

		// Receiving stats
		Targets             uint
		Receptions          uint
		ReceivingYards      uint
		ReceivingTouchdowns uint

		// Rushing stats
		Carries           uint
		RushingYards      uint
		RushingTouchdowns uint

		// Turnover stats
		InterceptionsThrown uint
		FumblesLost         uint

		// Defensive stats
		Sacks               uint
		InterceptionsCaught uint
		FumblesRecovered    uint
		TacklesForLoss      uint

		// Kicking stats
		ExtraPointsMade    uint
		ExtraPointsMissed  uint
		FieldGoalsMade     uint
		FieldGoalMadeYards uint
		FieldGoalsMissed   uint
		KickoffReturnYards uint
		PuntReturnYards    uint
		PuntYards          uint
	}
)

/*
type (
	Down     int
	Quarter  int
	PlayType int

	PlayTimestamp struct {
		Quarter       Quarter
		TimeRemaining time.Duration
	}

	Play struct {
		PlayClock     PlayTimestamp
		Down          Down
		DistanceYards uint
		Type          PlayType
	}
)

const (
	FirstQuarter Quarter = iota
	SecondQuarter
	ThirdQuarter
	FourthQuarter
	Overtime

	FirstDown Down = iota
	SecondDown
	ThirdDown
	FourthDown

	Rush PlayType = iota
	Pass
	Kickoff
	Punt
	FieldGoal
	ExtraPoint
	TwoPointConversion
)*/
