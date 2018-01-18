package parser

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"

	"github.com/ggirtsou/lloyds/pkg/model"
	"github.com/pkg/errors"
)

func ExtractTransactions(reader io.Reader) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	r := csv.NewReader(reader)
	i := 0

	for {
		record, err := r.Read()
		// we don't care about the header
		if i == 0 {
			i++
			continue
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, errors.Wrap(err, "failed to read csv")
		}

		transaction, err := transformToModel(record)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse line %d", i)
		}

		transactions = append(transactions, transaction)
		i++
	}

	return transactions, nil
}

func transformToModel(columns []string) (*model.Transaction, error) {
	const expectedCols = 8
	if len(columns) != expectedCols {
		return nil, errors.Errorf("expected %d columns, got: %d", expectedCols, len(columns))
	}

	accountNumber, err := strconv.Atoi(columns[3])
	if err != nil {
		return nil, errors.Errorf("could not convert account number to int: %v, err: %v", accountNumber, err)
	}

	balance, err := strconv.ParseFloat(columns[7], 64)
	if err != nil {
		return nil, errors.Errorf("could not convert balance to float: %v, err: %v", balance, err)
	}

	creditAmount, err := strconv.ParseFloat(columns[6], 64)
	if err != nil {
		return nil, errors.Errorf("could not convert credit amount to float: %v, err: %v", creditAmount, err)
	}

	debitAmount, err := strconv.ParseFloat(columns[5], 64)
	if err != nil {
		return nil, errors.Errorf("could not convert debit amount to float: %v, err: %v", debitAmount, err)
	}

	date, err := time.Parse("02/01/2006", columns[0])
	if err != nil {
		return nil, errors.Errorf("could not parse date: %v, err: %v", date, err)
	}

	const defaultTransactionType = "OTH" // @todo use it from a map somewhere
	transactionType := defaultTransactionType
	if columns[1] != "" {
		transactionType = columns[1]
	}

	return &model.Transaction{
		AccountNumber: accountNumber,
		Balance:       model.Amount(balance),
		CreditAmount:  model.Amount(creditAmount),
		DebitAmount:   model.Amount(debitAmount),
		SortCode:      columns[2],
		Date:          date,
		Description:   columns[4],
		Type:          model.TransactionType(transactionType),
	}, nil
}
