package pages

type abilities struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

type money struct {
	CP int `json:"cp"`
	SP int `json:"sp"`
	EP int `json:"ep"`
	GP int `json:"gp"`
	PP int `json:"pp"`
}

type item struct {
	Name        string `json:"name"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type feat struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ally struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type hitDie struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

//Sheet holds all the character sheet data
type Sheet struct {
	Owner             string    `json:"owner"`
	Name              string    `json:"name"`
	CharacterName     string    `json:"characterName"`
	Age               int       `json:"age"`
	Weight            string    `json:"weight"`
	Height            string    `json:"height"`
	Size              string    `json:"size"`
	Gender            string    `json:"gender"`
	EyeColor          string    `json:"eyeColor"`
	Skin              string    `json:"skin"`
	Class             string    `json:"class"`
	Race              string    `json:"race"`
	Level             int       `json:"level"`
	Allignment        string    `json:"allignment"`
	Background        string    `json:"background"`
	CurrentExpirience int       `json:"currentExpirience"`
	NextExpirience    int       `json:"nextExpirience"`
	Proficiency       int       `json:"proficiency"`
	Scores            abilities `json:"scores"`
	Saves             []string  `json:"saves"`
	ProficientSkills  []string  `json:"proficientSkills"`
	ExpertSkills      []string  `json:"expertSkills"`
	Languages         []string  `json:"languages"`
	Tools             []string  `json:"tools"`
	Vehicles          []string  `json:"vehicles"`
	Weapons           []string  `json:"weapons"`
	Armor             []string  `json:"armor"`
	Inventory         []item    `json:"inventory"`
	AC                int       `json:"ac"`
	Initiative        int       `json:"initiative"`
	Speed             int       `json:"speed"`
	Ideals            string    `json:"ideals"`
	Bonds             string    `json:"bonds"`
	Flaw              string    `json:"flaw"`
	Feats             []feat    `json:"feats"`
	Money             money     `json:"money"`
	PassivePerception int       `json:"passivePerception"`
	Backstory         string    `json:"backstory"`
	Allies            []ally    `json:"allies"`
	HitDie            hitDie    `json:"hitDie"`
}

//Index holds the data that fills our index page.
type Index struct {
	Sheets   []string
	Title    string
	LoggedIn bool
}

//SheetPage holds the data that fills the sheet page
type SheetPage struct {
	CharacterSheet Sheet
	LoggedIn       bool
}
