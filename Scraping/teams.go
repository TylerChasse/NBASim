package Scraping

type Team struct {
	Abbreviation            string
	Division                string
	Conference              string
	FullName                string
	OREBPerc                float64
	DREBPerc                float64
	OppOREBPerc             float64
	OppDREBPerc             float64
	TurnoverPerc            float64 // % of plays (how likely they are to turn the ball over each play)
	OppTurnoverPerc         float64
	DefFoulChance           float64
	OppDefFoulChance        float64
	TwoPerc                 float64
	OppTwoPerc              float64
	ThreePerc               float64
	OppThreePerc            float64
	FreeThrowPerc           float64
	ShootingFoulChance      float64
	ShootingFoulDrawnChance float64
	TwoFrequency            float64
	ThreeFrequency          float64
	OppTwoFrequency         float64
	OppThreeFrequency       float64
	TwoPointAnd1Chance      float64
	ThreePointAnd1Chance    float64
	ShotFrequency           float64
	Wins                    int
	Losses                  int
	GamesPlayed             int
	TotalPoints             int

	Score         int
	Fouls         int
	Last2MinFouls int

	IsDivisionWinner bool
	IsPlayoffTeam    bool
}

func CreateTeamsMap() map[string]*Team {
	var teamsInfo = [][]string{
		{"ATL", "East", "Southeast", "atlanta-hawks"},
		{"BOS", "East", "Atlantic", "boston-celtics"},
		{"CHA", "East", "Southeast", "charlotte-hornets"},
		{"CHI", "East", "Central", "chicago-bulls"},
		{"CLE", "East", "Central", "cleveland-cavaliers"},
		{"DAL", "West", "Southwest", "dallas-mavericks"},
		{"DEN", "West", "Northwest", "denver-nuggets"},
		{"DET", "East", "Central", "detroit-pistons"},
		{"GSW", "West", "Pacific", "golden-state-warriors"},
		{"HOU", "West", "Southwest", "houston-rockets"},
		{"IND", "East", "Central", "indiana-pacers"},
		{"LAC", "West", "Pacific", "los-angeles-clippers"},
		{"LAL", "West", "Pacific", "los-angeles-lakers"},
		{"MEM", "West", "Southwest", "memphis-grizzlies"},
		{"MIA", "East", "Southeast", "miami-heat"},
		{"MIL", "East", "Central", "milwaukee-bucks"},
		{"MIN", "West", "Northwest", "minnesota-timberwolves"},
		{"NOP", "West", "Southwest", "new-orleans-pelicans"},
		{"NYK", "East", "Atlantic", "new-york-knicks"},
		{"BKN", "East", "Atlantic", "brooklyn-nets"},
		{"OKC", "West", "Northwest", "oklahoma-city-thunder"},
		{"ORL", "East", "Southeast", "orlando-magic"},
		{"PHI", "East", "Atlantic", "philadelphia-76ers"},
		{"PHX", "West", "Pacific", "phoenix-suns"},
		{"POR", "West", "Northwest", "portland-trail-blazers"},
		{"SAC", "West", "Pacific", "sacramento-kings"},
		{"SAS", "West", "Southwest", "san-antonio-spurs"},
		{"TOR", "East", "Atlantic", "toronto-raptors"},
		{"UTA", "West", "Northwest", "utah-jazz"},
		{"WAS", "East", "Southeast", "washington-wizards"}}

	var teams = make(map[string]*Team)
	for _, teamInfo := range teamsInfo {
		abbreviation, conference, division, fullName := teamInfo[0], teamInfo[1], teamInfo[2], teamInfo[3]
		teams[abbreviation] = &Team{Abbreviation: abbreviation, Conference: conference, Division: division, FullName: fullName}
	}

	return teams
}
