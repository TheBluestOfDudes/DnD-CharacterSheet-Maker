package pages

type abilities struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

type money struct {
	CP int
	SP int
	EP int
	GP int
	PP int
}

type item struct {
	Name        string
	Amount      int
	Description string
}

type feat struct {
	Name        string
	Description string
}

type ally struct {
	Name        string
	Description string
}

//Index holds the data that fills our index page.
type Index struct {
	Sheets   []string
	Title    string
	LoggedIn bool
}

//Sheet holds the data that fills the sheet page
type Sheet struct {
	Owner             string
	Name              string
	CharacterName     string
	Age               int
	Weight            string
	Height            string
	Size              string
	Gender            string
	EyeColor          string
	Skin              string
	Class             string
	Race              string
	Level             int
	Allignment        string
	Background        string
	CurrentExpirience int
	NextExpirience    int
	Proficiency       int
	Scores            abilities
	ProficientSkills  []string
	ExpertSkills      []string
	Languages         []string
	Tools             []string
	Vehicles          []string
	Weapons           []string
	Armor             []string
	Inventory         []item
	AC                int
	Initiative        int
	Speed             int
	Ideals            string
	Bonds             string
	Flaw              string
	Feats             []feat
	Money             money
	PassivePerception int
	Backstory         string
	Allies            []ally
}
