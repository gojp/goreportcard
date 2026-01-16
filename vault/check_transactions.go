// Package vault provides functionality for processing PayPal transaction data.
// It reads CSV files containing transaction records, categorizes them by type,
// and generates formatted ledger reports.
package vault

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// TransactionType represents the category of a PayPal transaction.
type TransactionType string

const (
	// PaymentTransaction represents incoming payments from customers.
	PaymentTransaction TransactionType = "Payments"
	// TransferTransaction represents money transfers to/from accounts.
	TransferTransaction TransactionType = "Transfers"
	// FeeTransaction represents PayPal processing and service fees.
	FeeTransaction TransactionType = "Fees"
)

// Transaction represents a single PayPal transaction record with all relevant details.
type Transaction struct {
	Date          string          // Date of the transaction
	Type          TransactionType // Category: Payments, Transfers, or Fees
	Amount        string          // Transaction amount (can be negative)
	Description   string          // Human-readable description
	TransactionID string          // Unique PayPal transaction identifier
}

// TransactionProcessor handles reading, categorizing, and reporting on PayPal transactions.
type TransactionProcessor struct {
	vaultDir  string      // Directory containing CSV transaction files
	ledgerDir string      // Directory for generated ledger reports
	logger    *log.Logger // Logger for operational messages
}

// NewTransactionProcessor creates a new processor with the specified directories.
// It initializes logging and validates that the vault directory exists.
func NewTransactionProcessor(vaultDir, ledgerDir string) (*TransactionProcessor, error) {
	logger := log.New(os.Stdout, "[TransactionProcessor] ", log.LstdFlags)

	// Validate vault directory exists
	if _, err := os.Stat(vaultDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("vault directory does not exist: %s", vaultDir)
	}

	// Create ledger directory if it doesn't exist
	if err := os.MkdirAll(ledgerDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create ledger directory: %w", err)
	}

	return &TransactionProcessor{
		vaultDir:  vaultDir,
		ledgerDir: ledgerDir,
		logger:    logger,
	}, nil
}

// ReadCSVFiles reads all CSV files from the vault directory and returns parsed transactions.
// It handles file reading errors gracefully and logs any issues encountered.
func (tp *TransactionProcessor) ReadCSVFiles() ([]Transaction, error) {
	var allTransactions []Transaction

	// Find all CSV files in vault directory
	files, err := filepath.Glob(filepath.Join(tp.vaultDir, "*.csv"))
	if err != nil {
		return nil, fmt.Errorf("failed to search for CSV files: %w", err)
	}

	if len(files) == 0 {
		tp.logger.Printf("Warning: No CSV files found in %s", tp.vaultDir)
		return allTransactions, nil
	}

	tp.logger.Printf("Found %d CSV file(s) to process", len(files))

	// Process each CSV file
	for _, filename := range files {
		transactions, err := tp.readSingleCSV(filename)
		if err != nil {
			// Log error but continue processing other files
			tp.logger.Printf("Error reading %s: %v", filepath.Base(filename), err)
			continue
		}
		allTransactions = append(allTransactions, transactions...)
		tp.logger.Printf("Successfully processed %s: %d transactions", filepath.Base(filename), len(transactions))
	}

	return allTransactions, nil
}

// readSingleCSV reads and parses a single CSV file.
// It expects a header row with: Date, Type, Amount, Description, Transaction ID
func (tp *TransactionProcessor) readSingleCSV(filename string) ([]Transaction, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read header row
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate header structure
	if len(headers) < 5 {
		return nil, fmt.Errorf("invalid CSV format: expected at least 5 columns, got %d", len(headers))
	}

	var transactions []Transaction
	lineNum := 1 // Header is line 1, data starts at line 2

	// Read data rows
	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			tp.logger.Printf("Warning: Error reading line %d in %s: %v", lineNum, filepath.Base(filename), err)
			continue
		}

		// Validate record has enough fields
		if len(record) < 5 {
			tp.logger.Printf("Warning: Line %d in %s has insufficient fields (%d), skipping", lineNum, filepath.Base(filename), len(record))
			continue
		}

		// Parse transaction type
		transactionType := tp.categorizeTransaction(record[1], record[2], record[3])

		transaction := Transaction{
			Date:          strings.TrimSpace(record[0]),
			Type:          transactionType,
			Amount:        strings.TrimSpace(record[2]),
			Description:   strings.TrimSpace(record[3]),
			TransactionID: strings.TrimSpace(record[4]),
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// categorizeTransaction determines the transaction category based on type, amount, and description.
// It uses heuristics to classify transactions as Payments, Transfers, or Fees.
func (tp *TransactionProcessor) categorizeTransaction(rawType, amount, description string) TransactionType {
	typeStr := strings.ToLower(strings.TrimSpace(rawType))
	descStr := strings.ToLower(strings.TrimSpace(description))

	// Check for fee indicators
	if typeStr == "fee" || strings.Contains(descStr, "fee") || strings.Contains(descStr, "charge") {
		return FeeTransaction
	}

	// Check for transfer indicators
	if typeStr == "transfer" || strings.Contains(descStr, "transfer") ||
		strings.Contains(descStr, "withdrawal") || strings.Contains(descStr, "bank") {
		return TransferTransaction
	}

	// Default to payment for anything else
	return PaymentTransaction
}

// CategorizeTransactions groups transactions by their type.
// Returns a map with transaction types as keys and transaction slices as values.
func (tp *TransactionProcessor) CategorizeTransactions(transactions []Transaction) map[TransactionType][]Transaction {
	categorized := make(map[TransactionType][]Transaction)

	for _, txn := range transactions {
		categorized[txn.Type] = append(categorized[txn.Type], txn)
	}

	tp.logger.Printf("Categorization complete: %d Payments, %d Transfers, %d Fees",
		len(categorized[PaymentTransaction]),
		len(categorized[TransferTransaction]),
		len(categorized[FeeTransaction]))

	return categorized
}

// GenerateLedger creates a markdown-formatted ledger report and writes it to the specified file.
// The report includes a summary table with all transactions organized by category.
func (tp *TransactionProcessor) GenerateLedger(transactions []Transaction, outputFilename string) error {
	if len(transactions) == 0 {
		return fmt.Errorf("no transactions to write to ledger")
	}

	outputPath := filepath.Join(tp.ledgerDir, outputFilename)

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create ledger file: %w", err)
	}
	defer file.Close()

	// Write header
	if _, err := file.WriteString("# FK Master Ledger\n\n"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	generatedAt := time.Now().Format("2006-01-02 15:04:05")
	if _, err := file.WriteString(fmt.Sprintf("**Generated:** %s\n\n", generatedAt)); err != nil {
		return fmt.Errorf("failed to write timestamp: %w", err)
	}

	if _, err := file.WriteString(fmt.Sprintf("**Total Transactions:** %d\n\n", len(transactions))); err != nil {
		return fmt.Errorf("failed to write transaction count: %w", err)
	}

	// Categorize transactions
	categorized := tp.CategorizeTransactions(transactions)

	// Write each category
	categories := []TransactionType{PaymentTransaction, TransferTransaction, FeeTransaction}
	for _, category := range categories {
		txns := categorized[category]
		if len(txns) == 0 {
			continue
		}

		if err := tp.writeCategory(file, category, txns); err != nil {
			return fmt.Errorf("failed to write category %s: %w", category, err)
		}
	}

	tp.logger.Printf("Successfully generated ledger: %s", outputPath)
	return nil
}

// writeCategory writes a single category section to the ledger file.
// It includes a header and a formatted table of transactions.
func (tp *TransactionProcessor) writeCategory(w io.Writer, category TransactionType, transactions []Transaction) error {
	// Write category header
	if _, err := fmt.Fprintf(w, "## %s\n\n", category); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "**Count:** %d\n\n", len(transactions)); err != nil {
		return err
	}

	// Write table header
	if _, err := w.Write([]byte("| Dagsetning | Tegund | Upphæð | Lýsing | PayPal Transaction ID |\n")); err != nil {
		return err
	}
	if _, err := w.Write([]byte("|------------|--------|---------|--------|-----------------------|\n")); err != nil {
		return err
	}

	// Sort transactions by date for better readability
	sortedTxns := make([]Transaction, len(transactions))
	copy(sortedTxns, transactions)
	sort.Slice(sortedTxns, func(i, j int) bool {
		return sortedTxns[i].Date < sortedTxns[j].Date
	})

	// Write transaction rows
	for _, txn := range sortedTxns {
		if _, err := fmt.Fprintf(w, "| %s | %s | %s | %s | %s |\n",
			txn.Date,
			txn.Type,
			txn.Amount,
			txn.Description,
			txn.TransactionID); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

// Process is the main entry point that orchestrates the entire transaction processing workflow.
// It reads CSV files, categorizes transactions, and generates the ledger report.
func (tp *TransactionProcessor) Process() error {
	tp.logger.Println("Starting transaction processing...")

	// Read all CSV files
	transactions, err := tp.ReadCSVFiles()
	if err != nil {
		return fmt.Errorf("failed to read CSV files: %w", err)
	}

	if len(transactions) == 0 {
		tp.logger.Println("No transactions found to process")
		return nil
	}

	tp.logger.Printf("Total transactions read: %d", len(transactions))

	// Generate ledger
	if err := tp.GenerateLedger(transactions, "FK_MASTER_LEDGER.md"); err != nil {
		return fmt.Errorf("failed to generate ledger: %w", err)
	}

	tp.logger.Println("Transaction processing completed successfully")
	return nil
}

// Run is a convenience function that creates a processor and runs the complete workflow.
// It's the primary entry point for using this package.
func Run(vaultDir, ledgerDir string) error {
	processor, err := NewTransactionProcessor(vaultDir, ledgerDir)
	if err != nil {
		return fmt.Errorf("failed to initialize processor: %w", err)
	}

	return processor.Process()
}
