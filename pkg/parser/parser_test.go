package parser_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ggirtsou/lloyds/pkg/model"
	"github.com/ggirtsou/lloyds/pkg/parser"
)

func TestExtractTransactionsTooFewColumnsFail(t *testing.T) {
	r := strings.NewReader("col1,col2,col3\n1,2,3")
	_, err := parser.ExtractTransactions(r)
	if err == nil {
		t.Fatal("expected too few columns error, got: <nil>")
	}

	if !strings.Contains(err.Error(), "expected 8 columns, got: 3") {
		t.Fatalf("expected too few columns error, got: %v", err)
	}
}

var recordTests = []struct {
	in         string
	out        []model.Transaction
	errorMatch string // will use strings.Contains to see if there's a match
}{
	{
		// empty transaction type defaults to OTH
		in: "13/12/2018,,12-34-56,12345678,rent,100,0,50",
		out: []model.Transaction{

			{
				Date:          time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC),
				Type:          "OTH",
				SortCode:      "12-34-56",
				AccountNumber: 12345678,
				Description:   "rent",
				DebitAmount:   model.Amount(100),
				CreditAmount:  model.Amount(0),
				Balance:       model.Amount(50),
			},
		},
	},
	{
		// non-empty transaction type
		in: "13/12/2018,DEB,12-34-56,12345678,rent,100,0,50",
		out: []model.Transaction{

			{
				Date:          time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC),
				Type:          "DEB",
				SortCode:      "12-34-56",
				AccountNumber: 12345678,
				Description:   "rent",
				DebitAmount:   model.Amount(100),
				CreditAmount:  model.Amount(0),
				Balance:       model.Amount(50),
			},
		},
	},
	{
		// invalid date format
		in:         "12/13/2018,,12-34-56,12345678,rent,100,0,50",
		errorMatch: "could not parse date",
	},
}

// Because function returns a map, max 1 result is returned, so we can assert on struct fields.
func TestExtractTransactions(t *testing.T) {
	const header = "date,type,sortCode,accNo,desc,debitAmount,creditAmount,balance"
	for _, tt := range recordTests {
		r := strings.NewReader(fmt.Sprintf("%s\n%s", header, tt.in))
		transactions, err := parser.ExtractTransactions(r)

		if tt.errorMatch != "" {
			if err == nil {
				t.Fatalf("expected err: [%v], got: <nil>", tt.errorMatch)
			}

			if !strings.Contains(err.Error(), tt.errorMatch) {
				t.Fatalf("expected err: [%v], got: [%v]", tt.errorMatch, err.Error())
			}
		}

		if err != nil && tt.errorMatch == "" {
			t.Fatalf("unexpected err: %v", err)
		}

		if len(transactions) != len(tt.out) {
			t.Fatalf("expected %d transactions, got: %d", len(tt.out), len(transactions))
		}

		for _, fixture := range tt.out {
			got := transactions[0]
			if fixture.AccountNumber != got.AccountNumber {
				t.Fatalf("expected [%d] account number, got: [%d]", fixture.AccountNumber, got.AccountNumber)
			}

			if fixture.Balance != got.Balance {
				t.Fatalf("expected [%v] balance, got: [%v]", fixture.Balance, got.Balance)
			}

			if fixture.CreditAmount != got.CreditAmount {
				t.Fatalf("expected [%v] credit amount, got: [%v]", fixture.CreditAmount, got.CreditAmount)
			}

			if fixture.DebitAmount != got.DebitAmount {
				t.Fatalf("expected [%v] debit amount, got: [%v]", fixture.DebitAmount, got.DebitAmount)
			}

			if !fixture.Date.Equal(got.Date) {
				t.Fatalf("expected [%v] balance, got: [%v]", fixture.Date, got.Date)
			}

			if fixture.Type != got.Type {
				t.Fatalf("expected [%v] balance, got: [%v]", fixture.Type, got.Type)
			}
		}
	}
}
