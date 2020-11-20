package main

import (
	db "DB"
	pages "Pages"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
)

type featOrAlly struct {
	Name        string
	Description string
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func main() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/index/", indexHandler)
	http.HandleFunc("/loginpage/", loginPageHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/logout/", logoutHandler)
	http.HandleFunc("/sheet/", sheetHandler)
	http.HandleFunc("/register/", registerHandler)
	http.HandleFunc("/registerpage/", registerPageHandler)
	http.HandleFunc("/newsheet/", newSheetHandler)
	http.HandleFunc("/newsheetpage/", newSheetPageHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/deletepage/", deletePageHandler)
	log.Printf("Listening on %s...\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

func actionFailed(w http.ResponseWriter, message string) {
	fail := pages.FailPage{}
	t, _ := template.ParseFiles("./templates/fail.html")
	fail.Error = message
	w.WriteHeader(http.StatusInternalServerError)
	t.Execute(w, fail)
}

func parseFeatsAndAllies(list []string) []featOrAlly {
	featsOrAllies := []featOrAlly{}
	for i := 0; i < len(list); i++ {
		feOAl := featOrAlly{}
		info := strings.Split(list[i], ":")
		if len(info) == 2 {
			for j := 0; j < len(info); j++ {
				switch j {
				case 0:
					{
						feOAl.Name = info[j]
						break
					}
				case 1:
					{
						feOAl.Description = info[j]
						break
					}
				}
			}
			featsOrAllies = append(featsOrAllies, feOAl)
		}
	}
	return featsOrAllies
}

func parseItems(itemList []string) []pages.Item {
	inventory := []pages.Item{}
	for i := 0; i < len(itemList); i++ {
		item := pages.Item{}
		itemInfo := strings.Split(itemList[i], ":")
		if len(itemInfo) == 3 {
			for j := 0; j < len(itemInfo); j++ {
				switch j {
				case 0:
					{
						num, err := strconv.Atoi(itemInfo[j])
						if err != nil {
							item.Amount = -1
						} else {
							item.Amount = num
						}
						break
					}
				case 1:
					{
						item.Name = itemInfo[j]
						break
					}
				case 2:
					{
						item.Description = itemInfo[j]
						break
					}
				}
			}
			inventory = append(inventory, item)
		}
	}
	return inventory
}

func parseSpells(spellList []string) []pages.Spell {
	spells := []pages.Spell{}
	for i := 0; i < len(spellList); i++ {
		spell := pages.Spell{}
		spellInfo := strings.Split(spellList[i], ":")
		if len(spellInfo) == 3 {
			for j := 0; j < len(spellInfo); j++ {
				switch j {
				case 0:
					{
						spell.Name = spellInfo[j]
						break
					}
				case 1:
					{
						num, _ := strconv.Atoi(spellInfo[j])
						spell.Level = num
						break
					}
				case 2:
					{
						spell.Description = spellInfo[j]
						break
					}
				}
			}
			spells = append(spells, spell)
		}
	}
	return spells
}

func makeHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func setSession(userName string, w http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:   "session",
			Value:  encoded,
			Path:   "/",
			MaxAge: 3600,
		}
		http.SetCookie(w, cookie)
	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func getUserName(r *http.Request) (userName string) {
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := pages.Index{}
	username := getUserName(r)
	if username == "" {
		data.Title = "Please log in or register"
		data.LoggedIn = false
	} else {
		data.Title = "Welcome " + username
		data.LoggedIn = true
		sheets, err := db.GetSheets(username)
		if err != nil {
			data.Sheets = []string{"Could not load sheets from database"}
		} else if len(sheets) == 0 {
			data.Sheets = []string{"You have no saved sheets"}
		} else {
			data.Sheets = sheets
		}
	}
	pageData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	t, _ := template.ParseFiles("./templates/index.html")
	t.Execute(w, string(pageData))
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/login.html")
	t.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username != "" && password != "" {
		ok, err := db.CheckUser(username, password)
		if err != nil || !ok {
			actionFailed(w, `{"message":"Could not match password or username"}`)
		} else if ok {
			setSession(username, w)
			http.Redirect(w, r, "/index/", 303)
		}
	} else {
		http.Redirect(w, r, "/index/", 303)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/index/", 303)
}

func sheetHandler(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	page := pages.SheetPage{}
	if username == "" {
		http.Redirect(w, r, "/index/", 303)
	} else {
		page.LoggedIn = true
		sheet, err := db.GetSheet(username, r.FormValue("sheet"))
		if err != nil {
			actionFailed(w, `{"message":"`+err.Error()+`"}`)
		} else {
			page.CharacterSheet = sheet
			pageData, err := json.Marshal(page)
			if err != nil {
				panic(err)
			}
			t, _ := template.ParseFiles("./templates/sheet.html")
			t.Execute(w, string(pageData))
		}
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user := db.User{}
	if username == "" || password == "" {
		http.Redirect(w, r, "/index/", 303)
	} else if username != "" && password != "" {
		exist, err := db.CheckUserName(r.FormValue("username"))
		if err != nil {
			actionFailed(w, `{"message":"`+err.Error()+`"}`)
		} else if !exist {
			user.Username = r.FormValue("username")
			user.Password = makeHash([]byte(r.FormValue("password")))
			user.Sheets = []string{}
			err := db.RegisterUser(user)
			if err != nil {
				actionFailed(w, `{"message":"`+err.Error()+`"}`)
			} else {
				http.Redirect(w, r, "/loginpage/", 303)
			}
		} else {
			actionFailed(w, `{"message":"Username is already taken"}`)
		}
	}
}

func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/register.html")
	t.Execute(w, nil)
}

func newSheetHandler(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	if username != "" {
		sheet := pages.Sheet{}
		r.ParseForm()
		age, _ := strconv.Atoi(r.Form["age"][0])
		level, _ := strconv.Atoi(r.Form["level"][0])
		curEx, _ := strconv.Atoi(r.Form["currentExpirience"][0])
		nexEx, _ := strconv.Atoi(r.Form["nextExpirience"][0])
		prof, _ := strconv.Atoi(r.Form["proficiency"][0])
		str, _ := strconv.Atoi(r.Form["strength"][0])
		dex, _ := strconv.Atoi(r.Form["dexterity"][0])
		con, _ := strconv.Atoi(r.Form["constitution"][0])
		intel, _ := strconv.Atoi(r.Form["intelligence"][0])
		wis, _ := strconv.Atoi(r.Form["wisdom"][0])
		cha, _ := strconv.Atoi(r.Form["charisma"][0])
		ac, _ := strconv.Atoi(r.Form["ac"][0])
		spd, _ := strconv.Atoi(r.Form["speed"][0])
		init, _ := strconv.Atoi(r.Form["initiative"][0])
		cp, _ := strconv.Atoi(r.Form["cp"][0])
		sp, _ := strconv.Atoi(r.Form["sp"][0])
		ep, _ := strconv.Atoi(r.Form["ep"][0])
		gp, _ := strconv.Atoi(r.Form["gp"][0])
		pp, _ := strconv.Atoi(r.Form["pp"][0])
		passivePer, _ := strconv.Atoi(r.Form["passivePerception"][0])
		health, _ := strconv.Atoi(r.Form["health"][0])
		scores := pages.Abilities{
			Strength:     str,
			Dexterity:    dex,
			Constitution: con,
			Intelligence: intel,
			Wisdom:       wis,
			Charisma:     cha,
		}
		money := pages.Coin{
			CP: cp,
			SP: sp,
			EP: ep,
			GP: gp,
			PP: pp,
		}
		hitDie := pages.HitDice{
			Name:   r.Form["hitDie"][0],
			Amount: level,
		}
		var inventory = []pages.Item{}
		if len(r.Form["inventory"]) > 0 {
			inventory = parseItems(strings.Split(r.Form["inventory"][0], ","))
		}
		var feats = []pages.Feat{}
		if len(r.Form["feats"]) > 0 {
			featsOrAllies := parseFeatsAndAllies(strings.Split(r.Form["feats"][0], ","))
			for _, feat := range featsOrAllies {
				feats = append(feats, pages.Feat{Name: feat.Name, Description: feat.Description})
			}
		}
		var allies = []pages.Ally{}
		if len(r.Form["allies"]) > 0 {
			featsOrAllies := parseFeatsAndAllies(strings.Split(r.Form["allies"][0], ","))
			for _, ally := range featsOrAllies {
				allies = append(allies, pages.Ally{Name: ally.Name, Description: ally.Description})
			}
		}
		var spells = []pages.Spell{}
		if len(r.Form["spells"]) > 0 {
			spells = parseSpells(strings.Split(r.Form["spells"][0], ","))
		}
		sheet.Owner = username
		sheet.Name = r.Form["name"][0]
		sheet.CharacterName = r.Form["characterName"][0]
		sheet.Age = age
		sheet.Weight = r.Form["weight"][0]
		sheet.Height = r.Form["height"][0]
		sheet.Size = r.Form["size"][0]
		sheet.Gender = r.Form["gender"][0]
		sheet.EyeColor = r.Form["eyeColor"][0]
		sheet.Skin = r.Form["skin"][0]
		sheet.Class = r.Form["class"][0]
		sheet.Race = r.Form["race"][0]
		sheet.Level = level
		sheet.Allignment = r.Form["allignment"][0]
		sheet.Race = r.Form["race"][0]
		sheet.CurrentExpirience = curEx
		sheet.NextExpirience = nexEx
		sheet.Proficiency = prof
		sheet.Scores = scores
		sheet.Saves = r.Form["saves"]
		sheet.ProficientSkills = r.Form["proficientSkills"]
		sheet.ExpertSkills = r.Form["expertSkills"]
		sheet.Languages = r.Form["languages"]
		sheet.Tools = r.Form["tools"]
		sheet.Vehicles = r.Form["vehicles"]
		sheet.Weapons = r.Form["weapons"]
		sheet.Armor = r.Form["armor"]
		sheet.Inventory = inventory
		sheet.AC = ac
		sheet.Initiative = init
		sheet.Speed = spd
		sheet.Ideals = r.Form["ideals"][0]
		sheet.Bonds = r.Form["bonds"][0]
		sheet.Flaw = r.Form["flaw"][0]
		sheet.Feats = feats
		sheet.Money = money
		sheet.PassivePerception = passivePer
		sheet.Backstory = r.Form["backstory"][0]
		sheet.Allies = allies
		sheet.HitDie = hitDie
		sheet.Health = health
		sheet.Spells = spells
		err := db.RegisterSheet(username, sheet)
		if err != nil {
			actionFailed(w, `{"message":"`+err.Error()+`"}`)
		} else {
			http.Redirect(w, r, "/index/", 303)
		}
	}
	http.Redirect(w, r, "/index/", 303)
}

func newSheetPageHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/newSheet.html")
	t.Execute(w, nil)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	sheet := r.FormValue("yes")
	err := db.DeleteSheet(username, sheet)
	if err != nil {
		actionFailed(w, `{"message":"`+err.Error()+`"}`)
	} else {
		http.Redirect(w, r, "/index/", 303)
	}
}

func deletePageHandler(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	page := pages.DeletePage{}
	if username != "" {
		page.SheetName = r.FormValue("delete")
		t, _ := template.ParseFiles("./templates/delete.html")
		t.Execute(w, page)
	} else {
		http.Redirect(w, r, "/index/", 303)
	}
}
