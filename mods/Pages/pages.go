package pages

//Abilities represents a character's ability scores
type Abilities struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

//Coin represents the money a character has
type Coin struct {
	CP int `json:"cp"`
	SP int `json:"sp"`
	EP int `json:"ep"`
	GP int `json:"gp"`
	PP int `json:"pp"`
}

//Item represents an item in a character's inventory
type Item struct {
	Name        string `json:"name"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

//Feat represents a feat a character might have
type Feat struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

//Ally represents an ally a character has
type Ally struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

//HitDice represents what the character's hit dice are
type HitDice struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

//Spell represents a spell a character might have
type Spell struct {
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
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
	Scores            Abilities `json:"scores"`
	Saves             []string  `json:"saves"`
	ProficientSkills  []string  `json:"proficientSkills"`
	ExpertSkills      []string  `json:"expertSkills"`
	Languages         []string  `json:"languages"`
	Tools             []string  `json:"tools"`
	Vehicles          []string  `json:"vehicles"`
	Weapons           []string  `json:"weapons"`
	Armor             []string  `json:"armor"`
	Inventory         []Item    `json:"inventory"`
	AC                int       `json:"ac"`
	Initiative        int       `json:"initiative"`
	Speed             int       `json:"speed"`
	Ideals            string    `json:"ideals"`
	Bonds             string    `json:"bonds"`
	Flaw              string    `json:"flaw"`
	Feats             []Feat    `json:"feats"`
	Money             Coin      `json:"money"`
	PassivePerception int       `json:"passivePerception"`
	Backstory         string    `json:"backstory"`
	Allies            []Ally    `json:"allies"`
	HitDie            HitDice   `json:"hitDie"`
	Health            int       `json:"health"`
	Spells            []Spell   `json:"spells"`
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

//DeletePage holds the data that fills the delete page
type DeletePage struct {
	SheetName string
}

//FailPage holds the data that fills the error page
type FailPage struct {
	Error string
}
