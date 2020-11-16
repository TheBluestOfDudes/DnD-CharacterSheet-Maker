module Handlers

go 1.15

replace Pages => ../Pages/

replace DB => ../DB/

require (
	DB v0.0.0-00010101000000-000000000000
	Pages v0.0.0-00010101000000-000000000000
	github.com/gorilla/sessions v1.2.1
)
