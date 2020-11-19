module DB

go 1.15

replace Pages => ../Pages

require (
	Pages v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.4.3
	golang.org/x/crypto v0.0.0-20190530122614-20be4c3c3ed5
)
