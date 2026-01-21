// Package main provides a command-line utility for processing PayPal transactions.
// It reads CSV files from the vault directory and generates a ledger report.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gojp/goreportcard/vault"
)

func main() {
	// Define command-line flags
	vaultDir := flag.String("vault", "vault", "Directory containing PayPal CSV transaction files")
	ledgerDir := flag.String("ledger", "ledger", "Directory for generated ledger reports")
	help := flag.Bool("help", false, "Show usage information")

	flag.Parse()

	if *help {
		fmt.Println("PayPal Transaction Processor")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Println("  go run vault/cmd/main.go -vault=./vault -ledger=./ledger")
		os.Exit(0)
	}

	// Convert to absolute paths
	absVaultDir, err := filepath.Abs(*vaultDir)
	if err != nil {
		log.Fatalf("Invalid vault directory path: %v", err)
	}

	absLedgerDir, err := filepath.Abs(*ledgerDir)
	if err != nil {
		log.Fatalf("Invalid ledger directory path: %v", err)
	}

	fmt.Printf("Processing transactions...\n")
	fmt.Printf("Vault Directory:  %s\n", absVaultDir)
	fmt.Printf("Ledger Directory: %s\n\n", absLedgerDir)

	// Run the transaction processor
	if err := vault.Run(absVaultDir, absLedgerDir); err != nil {
		log.Fatalf("Error processing transactions: %v", err)
	}

	fmt.Println("\n✓ Transaction processing completed successfully!")
}
