package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Wesley-Nunes/KnowledgeVault-Backend/internal/durable"
)

func TestMain(m *testing.M) {
	// Run tests
	exitCode := m.Run()

	// Run cleanup after all tests
	createBookCleanUp()
	createDetailsCleanUp()

	// Exit with the result of the test run
	os.Exit(exitCode)
}

func createBookCleanUp() {
	sqlDelete := "DELETE FROM books WHERE author = 'Author fake to test' AND title = 'Title fake to test' AND pages = 999999999;"
	pool := durable.GetConnPool()
	defer pool.Close()

	_, err := pool.Exec(context.Background(), sqlDelete)
	if err != nil {
		fmt.Println("Error cleaning up:", err)
	}
}

func createDetailsCleanUp() {
	sqlDelete := "DELETE FROM details WHERE status = 'Wishlist' AND pages = 999999999;"
	pool := durable.GetConnPool()
	defer pool.Close()

	_, err := pool.Exec(context.Background(), sqlDelete)
	if err != nil {
		fmt.Println("Error cleaning up:", err)
	}
}
