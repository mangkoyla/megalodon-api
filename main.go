package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/FoolVPN-ID/megalodon-api/api"
	"github.com/FoolVPN-ID/megalodon-api/modules/db/users"
	_ "github.com/joho/godotenv"
)

func main() {
	// Initialization
	users.MakeUsersTableClient().CreateTableSafe()

	go api.StartApi()

	// Channel untuk menangkap sinyal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// Tangkap sinyal interrupt (Ctrl+C)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println("Received signal:", sig)

		done <- true
	}()

	fmt.Println("Running... Press Ctrl+C to stop.")
	<-done
	fmt.Println("Exiting.")
}
