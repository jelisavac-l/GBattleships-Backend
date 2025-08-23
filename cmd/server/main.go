package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jelisavac-l/GBattleships/internal/routes"
)

func main() {
	fmt.Println("Battleships server started!")
	routes.RegisterServerRoutes()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
