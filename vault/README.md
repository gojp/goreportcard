# Vault Transaction Processor

A robust Go package for processing PayPal CSV transaction files and generating ledger reports.

## Features

- **CSV Parsing**: Reads PayPal transaction CSV files from the vault directory
- **Transaction Categorization**: Automatically categorizes transactions into:
  - **Payments**: Incoming payments from customers
  - **Transfers**: Money transfers to/from accounts
  - **Fees**: PayPal processing and service fees
- **Error Handling**: Robust error handling for file and data issues with detailed logging
- **Ledger Generation**: Generates formatted markdown ledger reports with transaction tables
- **High Code Quality**: Follows Go best practices with comprehensive documentation

## Installation

```bash
go get github.com/gojp/goreportcard/vault
```

## Usage

### As a Library

```go
package main

import (
    "log"
    "github.com/gojp/goreportcard/vault"
)

func main() {
    // Process transactions from vault/ and generate ledger in ledger/
    if err := vault.Run("vault", "ledger"); err != nil {
        log.Fatalf("Error: %v", err)
    }
}
```

### As a Command-Line Tool

```bash
# Use default directories (vault/ and ledger/)
go run vault/cmd/main.go

# Specify custom directories
go run vault/cmd/main.go -vault=./my-transactions -ledger=./my-reports

# Show help
go run vault/cmd/main.go -help
```

## CSV Format

The processor expects CSV files with the following header:

```
Date,Type,Amount,Description,Transaction ID
```

Example:

```csv
Date,Type,Amount,Description,Transaction ID
2024-01-15,Payment,100.50,Product sale payment,TXN001
2024-01-16,Transfer,-50.00,Bank transfer,TXN002
2024-01-17,Fee,-2.99,PayPal processing fee,TXN003
```

## Output

The processor generates a markdown ledger file (`FK_MASTER_LEDGER.md`) with:

- Transaction summary statistics
- Categorized transaction tables
- Icelandic column headers: Dagsetning, Tegund, Upphæð, Lýsing, PayPal Transaction ID

Example output:

```markdown
# FK Master Ledger

**Generated:** 2026-01-16 01:35:50
**Total Transactions:** 7

## Payments

**Count:** 3

| Dagsetning | Tegund | Upphæð | Lýsing | PayPal Transaction ID |
|------------|--------|---------|--------|-----------------------|
| 2024-01-15 | Payments | 100.50 | Product sale payment | TXN001 |
```

## Testing

```bash
# Run tests
go test ./vault/...

# Run tests with coverage
go test ./vault/... -cover

# Run tests verbosely
go test -v ./vault/...
```

## Code Quality

This package follows Go best practices and passes all standard quality checks:

- ✓ `go fmt` - Code formatting
- ✓ `go vet` - Static analysis
- ✓ `staticcheck` - Advanced static analysis
- ✓ `golint` - Code style checking
- ✓ `gocyclo` - Cyclomatic complexity (all functions < 15)
- ✓ `misspell` - Spelling checks
- ✓ Test coverage: 66.7%

## Package Structure

```
vault/
├── check_transactions.go      # Main transaction processor implementation
├── check_transactions_test.go # Comprehensive test suite
├── cmd/
│   └── main.go               # Command-line interface
├── sample_transactions.csv   # Example CSV file
└── README.md                 # This file
```

## API Documentation

### Types

- `TransactionType`: Enum for transaction categories (Payments, Transfers, Fees)
- `Transaction`: Represents a single transaction record
- `TransactionProcessor`: Main processor for handling transactions

### Functions

- `NewTransactionProcessor(vaultDir, ledgerDir string)`: Create a new processor
- `Run(vaultDir, ledgerDir string)`: Convenience function to run the full workflow

### Methods

- `ReadCSVFiles()`: Read all CSV files from vault directory
- `CategorizeTransactions(transactions)`: Group transactions by type
- `GenerateLedger(transactions, outputFilename)`: Generate markdown ledger
- `Process()`: Run the complete processing workflow

## Error Handling

The processor includes comprehensive error handling:

- Validates vault directory exists before processing
- Creates ledger directory if it doesn't exist
- Logs warnings for malformed CSV rows (continues processing)
- Returns errors for critical issues (file access, write failures)
- Provides detailed error messages with context

## Logging

The processor logs all operations to stdout with timestamps:

```
[TransactionProcessor] 2026/01/16 01:35:50 Starting transaction processing...
[TransactionProcessor] 2026/01/16 01:35:50 Found 1 CSV file(s) to process
[TransactionProcessor] 2026/01/16 01:35:50 Successfully processed sample_transactions.csv: 7 transactions
```

## License

This package is part of the goreportcard project and follows the same Apache 2.0 license.
