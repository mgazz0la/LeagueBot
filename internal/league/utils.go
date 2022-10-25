package league

import (
	"sort"
	"time"

	"github.com/mgazz0la/leaguebot/internal/domain"
	"github.com/mgazz0la/leaguebot/internal/utils"
)

type ByRank []*domain.Squad

func (sqs ByRank) Len() int { return len(sqs) }
func (sqs ByRank) Less(i, j int) bool {
	return sqs[i].Wins > sqs[j].Wins ||
		(sqs[i].Wins == sqs[j].Wins && sqs[i].PointsFor > sqs[j].PointsFor)
}
func (sqs ByRank) Swap(i, j int) { sqs[i], sqs[j] = sqs[j], sqs[i] }

type ByPF []*domain.Squad

func (sqs ByPF) Len() int           { return len(sqs) }
func (sqs ByPF) Less(i, j int) bool { return sqs[i].PointsFor > sqs[j].PointsFor }
func (sqs ByPF) Swap(i, j int)      { sqs[i], sqs[j] = sqs[j], sqs[i] }

func ApplySeeds(sqs []*domain.Squad) {
	sort.Sort(ByRank(sqs))
	for i := uint(1); i <= 4; i++ {
		sqs[i-1].Seed = i
	}

	wcs := sqs[4:]
	sort.Sort(ByPF(wcs))
	wcs[0].Seed = 5
	wcs[1].Seed = 6

	sackos := wcs[2:]
	sort.Sort(ByRank(sackos))
	for i := range sackos {
		sackos[i].Seed = uint(i + 7)
	}
}

func GetCurrentWeek() uint {
	return uint(utils.Second(time.Now().AddDate(0, 0, -2).ISOWeek()) - 35)
}
