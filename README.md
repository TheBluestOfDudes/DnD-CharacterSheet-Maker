# DnD-CharacterSheet-Maker
 A web application that lets you fill out a dnd 5E character sheet, and save it in a database to be retrieved

# Usage
The app is deployed on heroku on [this link](https://shrouded-hamlet-99487.herokuapp.com/index/). You will be asked to log in/register. Once you make an account, you have the option to save new sheets, view sheets youve made or delete the sheets.

## Filling out a sheet
The app assumes the user knows the rules of DnD and doesn't do much to vet the values the user inputs. You fill out the sheet registration form, and the app tries to save it, assuming all those values were valid. You will then be able to see it appear in your login page.

### Editing:
It was intended to be possilbe to later edit the sheet you've made, which is why some of the fields in the character sheet like the inventory and feat list are not mandatory. The editing feature was never implemented though.