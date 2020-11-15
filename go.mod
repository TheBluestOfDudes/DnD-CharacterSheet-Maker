module main

go 1.15

replace DB => ./mods/DB/

replace Pages => ./mods/Pages/

require (
	DB v0.0.0-00010101000000-000000000000
	Pages v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.8.0 // indirect
	go.mongodb.org/mongo-driver v1.4.3
)
