module main

go 1.15

replace DB => ./mods/DB/

replace Handlers => ./mods/Handlers/

replace Pages => ./mods/Pages/

require (
	DB v0.0.0-00010101000000-000000000000
	Handlers v0.0.0-00010101000000-000000000000
	Pages v0.0.0-00010101000000-000000000000
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/lib/pq v1.8.0 // indirect
	github.com/pkg/errors v0.9.1
	go.mongodb.org/mongo-driver v1.4.3
)
