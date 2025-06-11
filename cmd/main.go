package main

import (
	"SportHub-Forum/db"
	"SportHub-Forum/internal/server"
	"log"
)

func main() {
	// Initialise and start the server
	dsn := "root:root@tcp(127.0.0.1:3306)/forumdb"
	if err := db.InitDB(dsn); err != nil {
		log.Fatalf("❌ Connexion DB échouée : %v", err)
	}
	srv := server.New()
	log.Fatal(srv.Start())
	log.Println("✅ Base connectée avec succès.")
}
