package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewTransactionProcessor tests the initialization of the transaction processor.
func TestNewTransactionProcessor(t *testing.T) {
	// Create temporary directories for testing
	tmpDir := t.TempDir()
	vaultDir := filepath.Join(tmpDir, "vault")
	ledgerDir := filepath.Join(tmpDir, "ledger")

	// Create vault directory
	if err := os.MkdirAll(vaultDir, 0755); err != nil {
		t.Fatalf("Failed to create test vault directory: %v", err)
	}

	processor, err := NewTransactionProcessor(vaultDir, ledgerDir)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}

	if processor == nil {
		t.Fatal("Processor is nil")
	}

	if processor.vaultDir != vaultDir {
		t.Errorf("Expected vault dir %s, got %s", vaultDir, processor.vaultDir)
	}

	// Check that ledger directory was created
	if _, err := os.Stat(ledgerDir); os.IsNotExist(err) {
		t.Error("Ledger directory was not created")
	}
}

// TestNewTransactionProcessorInvalidVault tests error handling for non-existent vault directory.
func TestNewTransactionProcessorInvalidVault(t *testing.T) {
	tmpDir := t.TempDir()
	vaultDir := filepath.Join(tmpDir, "nonexistent")
	ledgerDir := filepath.Join(tmpDir, "ledger")

	_, err := NewTransactionProcessor(vaultDir, ledgerDir)
	if err == nil {
		t.Error("Expected error for non-existent vault directory, got nil")
	}
}

// TestCategorizeTransaction tests the transaction categorization logic.
func TestCategorizeTransaction(t *testing.T) {
	tmpDir := t.TempDir()
	vaultDir := filepath.Join(tmpDir, "vault")
	ledgerDir := filepath.Join(tmpDir, "ledger")
	os.MkdirAll(vaultDir, 0755)

	processor, err := NewTransactionProcessor(vaultDir, ledgerDir)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}

	tests := []struct {
		name        string
		rawType     string
		amount      string
		description string
		expected    TransactionType
	}{
		{
			name:        "Fee transaction",
			rawType:     "fee",
			amount:      "-2.99",
			description: "PayPal processing fee",
			expected:    FeeTransaction,
		},
		{
			name:        "Fee by description",
			rawType:     "other",
			amount:      "-5.00",
			description: "Service charge fee",
			expected:    FeeTransaction,
		},
		{
			name:        "Transfer transaction",
			rawType:     "transfer",
			amount:      "-100.00",
			description: "Bank transfer",
			expected:    TransferTransaction,
		},
		{
			name:        "Transfer by description",
			rawType:     "other",
			amount:      "-50.00",
			description: "Withdrawal to bank",
			expected:    TransferTransaction,
		},
		{
			name:        "Payment transaction",
			rawType:     "payment",
			amount:      "100.50",
			description: "Product sale",
			expected:    PaymentTransaction,
		},
		{
			name:        "Default to payment",
			rawType:     "other",
			amount:      "75.00",
			description: "Some income",
			expected:    PaymentTransaction,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.categorizeTransaction(tt.rawType, tt.amount, tt.description)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestReadCSVFiles tests reading and parsing CSV files.
func TestReadCSVFiles(t *testing.T) {
	tmpDir := t.TempDir()
	vaultDir := filepath.Join(tmpDir, "vault")
	ledgerDir := filepath.Join(tmpDir, "ledger")

	if err := os.MkdirAll(vaultDir, 0755); err != nil {
		t.Fatalf("Failed to create vault directory: %v", err)
	}

	// Create a test CSV file
	csvContent := `Date,Type,Amount,Description,Transaction ID
2024-01-15,Payment,100.50,Product sale,TXN001
2024-01-16,Transfer,-50.00,Bank transfer,TXN002
2024-01-17,Fee,-2.99,Processing fee,TXN003
`
	csvPath := filepath.Join(vaultDir, "test.csv")
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	processor, err := NewTransactionProcessor(vaultDir, ledgerDir)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}

	transactions, err := processor.ReadCSVFiles()
	if err != nil {
		t.Fatalf("Failed to read CSV files: %v", err)
	}

	if len(transactions) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(transactions))
	}

	// Verify first transaction
	if transactions[0].Date != "2024-01-15" {
		t.Errorf("Expected date 2024-01-15, got %s", transactions[0].Date)
	}
	if transactions[0].TransactionID != "TXN001" {
		t.Errorf("Expected ID TXN001, got %s", transactions[0].TransactionID)
	}
}

// TestReadCSVFilesNoFiles tests handling of empty vault directory.
func TestReadCSVFilesNoFiles(t *testing.T) {
	tmpDir := t.TempDir()
	vaultDir := filepath.Join(tmpDir, "vault")
	ledgerDir := filepath.Join(tmpDir, "ledger")

	if err := os.MkdirAll(vaultDir, 0755); err != nil {
		t.Fatalf("Failed to create vault directory: %v", err)
	}

	processor, err := NewTransactionProcessor(vaultDir, ledgerDir)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}

	transactions, err := processor.ReadCSVFiles()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(transactions) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(transactions))
	}
}

// TestGenerateLedger tests the ledger generation functionality.
func TestGenerateLedger(t *testing.T) {
	tmpDir := t.TempDir()
	vaultDir := filepath.Join(tmpDir, "vault")
	ledgerDir := filepath.Join(tmpDir, "ledger")

	if err := os.MkdirAll(vaultDir, 0755); err != nil {
		t.Fatalf("Failed to create vault directory: %v", err)
	}

	processor, err := NewTransactionProcessor(vaultDir, ledgerDir)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}

	// Create test transactions
	transactions := []Transaction{
		{
			Date:          "2024-01-15",
			Type:          PaymentTransaction,
			Amount:        "100.50",
			Description:   "Product sale",
			TransactionID: "TXN001",
		},
		{
			Date:          "2024-01-16",
			Type:          TransferTransaction,
			Amount:        "-50.00",
			Description:   "Bank transfer",
			TransactionID: "TXN002",
		},
		{
			Date:          "2024-01-17",
			Type:          FeeTransaction,
			Amount:        "-2.99",
			Description:   "Processing fee",
			TransactionID: "TXN003",
		},
	}

	err = processor.GenerateLedger(transactions, "test_ledger.md")
	if err != nil {
		t.Fatalf("Failed to generate ledger: %v", err)
	}

	// Verify ledger file was created
	ledgerPath := filepath.Join(ledgerDir, "test_ledger.md")
	content, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatalf("Failed to read ledger file: %v", err)
	}

	contentStr := string(content)

	// Check that the file contains expected sections
	expectedStrings := []string{
		"# FK Master Ledger",
		"**Generated:**",
		"**Total Transactions:** 3",
		"## Payments",
		"## Transfers",
		"## Fees",
		"TXN001",
		"TXN002",
		"TXN003",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Ledger content missing expected string: %s", expected)
		}
	}
}

// TestGenerateLedgerNoTransactions tests error handling for empty transaction list.
func TestGenerateLedgerNoTransactions(t *testing.T) {
	tmpDir := t.TempDir()
	vaultDir := filepath.Join(tmpDir, "vault")
	ledgerDir := filepath.Join(tmpDir, "ledger")

	if err := os.MkdirAll(vaultDir, 0755); err != nil {
		t.Fatalf("Failed to create vault directory: %v", err)
	}

	processor, err := NewTransactionProcessor(vaultDir, ledgerDir)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}

	err = processor.GenerateLedger([]Transaction{}, "test_ledger.md")
	if err == nil {
		t.Error("Expected error for empty transaction list, got nil")
	}
}

// TestCategorizeTransactions tests the grouping of transactions by type.
func TestCategorizeTransactions(t *testing.T) {
	tmpDir := t.TempDir()
	vaultDir := filepath.Join(tmpDir, "vault")
	ledgerDir := filepath.Join(tmpDir, "ledger")

	if err := os.MkdirAll(vaultDir, 0755); err != nil {
		t.Fatalf("Failed to create vault directory: %v", err)
	}

	processor, err := NewTransactionProcessor(vaultDir, ledgerDir)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}

	transactions := []Transaction{
		{Type: PaymentTransaction, TransactionID: "TXN001"},
		{Type: PaymentTransaction, TransactionID: "TXN002"},
		{Type: TransferTransaction, TransactionID: "TXN003"},
		{Type: FeeTransaction, TransactionID: "TXN004"},
		{Type: FeeTransaction, TransactionID: "TXN005"},
	}

	categorized := processor.CategorizeTransactions(transactions)

	if len(categorized[PaymentTransaction]) != 2 {
		t.Errorf("Expected 2 payments, got %d", len(categorized[PaymentTransaction]))
	}
	if len(categorized[TransferTransaction]) != 1 {
		t.Errorf("Expected 1 transfer, got %d", len(categorized[TransferTransaction]))
	}
	if len(categorized[FeeTransaction]) != 2 {
		t.Errorf("Expected 2 fees, got %d", len(categorized[FeeTransaction]))
	}
}
