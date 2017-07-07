package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type transactions []transaction
type dates []time.Time

type pair struct {
	Key   string
	Value int
}

type pairList []pair

type transaction struct {
	transactionDate        time.Time
	transactionType        string
	sortCode               string
	accountNumber          int
	transactionDescription string
	debitAmount            float64
	creditAmount           float64
	balance                float64
}

var transactionCodes = map[string]string{
	"BGC": "Bank Giro Credit",
	"BNS": "Bonus",
	"BP":  "Bill Payment",
	"CHG": "Charge",
	"CHQ": "Cheque",
	"COM": "Commission",
	"COR": "Correction",
	"CPT": "Cashpoint",
	"CSH": "Cash",
	"CSQ": "Cash/Cheque",
	"DD":  "Direct Debit",
	"DEB": "Debit Card",
	"DEP": "Deposit",
	"EFT": "EFTPOS (electronic funds transfer at point of sale)",
	"EUR": "Euro Cheque",
	"FE":  "Foreign Exchange",
	"FEE": "Fixed Service Charge",
	"FPC": "Faster Payment charge",
	"FPI": "Faster Payment incoming",
	"FPO": "Faster Payment outgoing",
	"IB":  "Internet Banking",
	"INT": "Interest",
	"MPI": "Mobile Payment incoming",
	"MPO": "Mobile Payment outgoing",
	"MTG": "Mortgage",
	"NS":  "National Savings Dividend",
	"NSC": "National Savings Certificates",
	"OTH": "Other",
	"PAY": "Payment",
	"PSB": "Premium Savings Bonds",
	"PSV": "Paysave",
	"SAL": "Salary",
	"SPB": "Cashpoint",
	"SO":  "Standing Order",
	"STK": "Stocks/Shares",
	"TD":  "Dep Term Dec",
	"TDG": "Term Deposit Gross Interest",
	"TDI": "Dep Term Inc",
	"TDN": "Term Deposit Net Interest",
	"TFR": "Transfer",
	"UT":  "Unit Trust",
	"SUR": "Excess Reject",
}

// credit: https://groups.google.com/d/msg/golang-nuts/FT7cjmcL7gw/Gj4_aEsE_IsJ
func rankByWordCount(wordFrequencies map[string]int) pairList {
	pl := make(pairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (d dates) Len() int              { return len(d) }
func (d dates) Less(i, j int) bool    { return d[i].Before(d[j]) }
func (d dates) Swap(i, j int)         { d[i], d[j] = d[j], d[i] }

func main() {
	file := flag.String("file", "", "Path to transactions CSV. If empty, will use stdin")
	flag.Parse()

	r := os.Stdin
	if *file != "" {
		r, err := os.Open(*file)
		if err != nil {
			log.Fatalf("could not open file: %v", err)
		}
		defer r.Close()
	}

	dates, records, err := getRecords(r)
	if err != nil {
		log.Fatalf("could not get records: %v", err)
	}

	fmt.Printf("date\t\ttotal transactions\tspend\n")
	var totalTransactions int
	var totalSpend, totalEarnings float64
	sort.Sort(dates)
	transactionTypes := make(map[string]int)
	totalDaysNotSpend := 0
	monthlySpend := make(map[string]float64)

	for _, date := range dates {
		var daySpend float64
		for _, record := range records[date] {
			totalSpend += record.debitAmount
			totalEarnings += record.creditAmount
			daySpend += record.debitAmount
			transactionTypes[record.transactionType]++
			monthlySpend[date.Format("01/2006")] += record.debitAmount
			if record.debitAmount == 0 {
				totalDaysNotSpend++
			}
		}

		dayTransactions := len(records[date])
		fmt.Printf("%v\t\t%d\t\t£%.2f\n", date.Format("02/01/2006"), dayTransactions, daySpend)
		totalTransactions += dayTransactions
		fmt.Println("-----------")

	}
	fmt.Printf("Total transactions:\t%d in %d days\n", totalTransactions, len(dates))
	fmt.Println("-----------")
	fmt.Printf("Total spend:\t\t£%.2f\n", totalSpend)
	fmt.Println("-----------")
	fmt.Printf("Total earnings:\t\t£%.2f\n", totalEarnings)
	fmt.Println("-----------")
	spendInt, err := strconv.Atoi(fmt.Sprintf("%.0f", totalSpend))
	if err != nil {
		log.Fatalln("could not convert float to int: %v", err)
	}
	earnInt, err := strconv.Atoi(fmt.Sprintf("%.0f", totalEarnings))
	if err != nil {
		log.Fatalln("could not convert float to int: %v", err)
	}
	fmt.Printf("Average Spend / day:\t£%d\n", spendInt/len(dates))
	fmt.Printf("Average Earnings / day:\t£%d\n", earnInt/len(dates))
	fmt.Println("-----------")
	fmt.Printf("You didn't spend money *at all* for %d days!\n", totalDaysNotSpend)
	fmt.Println("-----------")
	fmt.Println("Transactions types by count:")
	mostFrequentTypes := rankByWordCount(transactionTypes)
	for _, transactionType := range mostFrequentTypes {
		fmt.Printf("* %v:\t%v\t%v\n", transactionType.Key, transactionType.Value, transactionCodes[transactionType.Key])
	}

	fmt.Println("-----------")
	fmt.Println("Total Spend By Month:")
	for month, totalMonthlySpend := range monthlySpend {
		fmt.Printf("* %v:\t£%.2f\n", month, totalMonthlySpend)
	}
}

func getRecords(r io.Reader) (dates, map[time.Time]transactions, error) {
	s := bufio.NewScanner(r)
	var dates dates
	collection := make(map[time.Time]transactions)
	lineNumber := 0
	for s.Scan() {
		// skip header (or we could just get CSVs without header and not have to do that)
		if lineNumber == 0 {
			lineNumber++
			continue
		}
		line := strings.Split(s.Text(), ",")
		date, err := time.Parse("02/01/2006", line[0])
		if err != nil {
			return nil, nil, fmt.Errorf("could not parse date: %v, line: %d", err, lineNumber)
		}

		aNumber, err := strconv.Atoi(line[3])
		if err != nil {
			return nil, nil, fmt.Errorf("could not parse account number: %v, line: %d", err, lineNumber)
		}

		if line[5] == "" {
			line[5] = "0"
		}
		debitAmount, err := strconv.ParseFloat(line[5], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("could not parse debit amount: %v, line: %d", err, lineNumber)
		}

		if line[6] == "" {
			line[6] = "0"
		}
		creditAmount, err := strconv.ParseFloat(line[6], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("could not parse credit amount: %v, line: %d", err, lineNumber)
		}

		balance, err := strconv.ParseFloat(line[7], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("could not parse balance: %v, line: %d", err, lineNumber)
		}

		transactionType := "OTH"
		if line[1] != "" {
			transactionType = line[1]
		}

		record := transaction{
			transactionDate:        date,
			transactionType:        transactionType,
			sortCode:               line[2],
			accountNumber:          aNumber,
			transactionDescription: line[4],
			debitAmount:            debitAmount,
			creditAmount:           creditAmount,
			balance:                balance,
		}
		if _, ok := collection[date]; !ok {
			dates = append(dates, date)
		}
		collection[date] = append(collection[date], record)
		lineNumber++
	}
	if err := s.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading: %v", err)
	}

	return dates, collection, nil
}
