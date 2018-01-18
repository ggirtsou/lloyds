# Lloyds spend report

Thanks for stopping by! This is a quick tool I wrote using Go, that gives you more
visibility to your spending with Lloyds. This program does not have a GUI (Graphical
User Interface), and runs on the terminal.

## How it works

It parses a CSV file, or data from `stdin` (if a file was not provided) and produces a
spending report.

**WARNING:** For security reasons, do not install a binary file that's supposed to do
analysis on your bank statement! Please, do not post the output of this program online,
as it's really bad practice revealing your financial history!

## How to get transactions export

* download monthly csv exports from internet banking
* create a new directory: mkdir `~/Downloads/lloyds` and move the CSV exports there `mv ~/Downloads/*.csv ~/Downloads/lloyds`
* concat them: `cat ~/Downloads/lloyds/*.csv > transactions.csv`
* delete the multiple headers: `sed '1!{/^Transaction Date/d;}' transactions.csv > clean.csv`

If you're using a ref/description for your transactions, you can easily find them.

* see rent transaction payments: `cat clean.csv | grep -i "rent"`
* how many times you paid rent: `cat clean.csv | grep -i "rent" | wc -l`

Pretty cool, right?

## Detailed Report

In the detailed report, you will see the following information:

* total transactions per day with daily spend
* total spend
* total earnings
* average spend & earnings
* total number of days you didn't spend money at all
* transaction types by count
* monthly spend

## Setup

* You have to have [Go](https://golang.org/dl/) installed on your system.
* Use `go get` to get the package: `go get github.com/ggirtsou/lloyds` and then `cd` to your project directory: `cd $GOPATH/src/github.com/ggirtsou/lloyds`
* Alternatively, clone the repository: `git clone git@github.com:ggirtsou/lloyds.git`

## Build

`go build` or `go run main.go`

## How to use

Either pass a filename in CLI argument: `./lloyds --file="/tmp/clean.csv"`, or pipe CSV contents
to stdin: `cat ~/Downloads/lloyds/clean.csv | ./lloyds`.

## Example output

```text
date		total transactions	spend
16/05/2016		1		£0.00
[...]
-----------
Total transactions:	1 in 1 days
-----------
Total spend:		£1
-----------
Total earnings:		£1
-----------
Average Spend / day:	£1
Average Earnings / day:	£1
-----------
You didn't spend money *at all* for 1 days!
-----------
Transactions types by count:
* DEB:	1	Debit Card
* CPT:	1	Cashpoint
* DD:	1	Direct Debit
* TFR:	1	Transfer
* PAY:	1	Payment
* DEP:	1	Deposit
* SO:	1	Standing Order
* OTH:	1	Other
-----------
Total Spend By Month:
* 05/2016:	£1
...
```

## Does this support other banks?

It all boils down to the CSV contents, having the right information in the right columns.
If you get an export with your transactions from another bank and columns are in different
order, you can reorder them and use this tool for that (althought it might not report
the transaction types correctly as these are probably specific to Lloyds).

Feel free to make any changes, and if you feel it would be useful for other cases, please
submit a Pull Request to make the tool more generic.

## Feature requests / Bugs

If you want a feature to be added, or a different calculation on the output, or found a bug
please create an issue here on Github, and depending on my time constraints I will add it.

## Can I make changes to this project?

Sure! Clone the repository with `git`, or download the source code and hack away! You're
free to do whatever you want with it, as it's under MIT license. Please note, I'm not
responsible for any issues with this software - read LICENSE.
