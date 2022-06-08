all: sign dilithium verify

sign: sign.go connect.go util.go
	go build sign.go connect.go util.go

verify: verify.go connect.go util.go
	go build verify.go connect.go util.go

dilithium: dilithium.go connect.go util.go
	go build dilithium.go connect.go util.go
