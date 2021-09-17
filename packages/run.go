package packages

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	FILESIZELIMIT = 2 * GIGABYTES
	DIVIDEND      = "dividend"
	INTEREST      = "interest"
	CONTRIBUTION  = "contribution"
	OPTION        = "optn"
	EQUITY        = "eq"
	YEARMARKER    = "'"
	PUT           = "put"
	CALL          = "call"
)

type EtradeTransactions struct {
	accountNumber      string
	totalDividends     float64
	totalInterest      float64
	totalCommissions   float64
	totalContributions float64
	headers            []string
	lines              []ETradeTransaction
	options            []option
}

type ETradeTransaction struct {
	rowIdentifier int
	transDate     string
	transType     string
	securityType  string
	symbol        string
	quantity      int
	amount        float64
	price         float64
	commission    float64
	description   string
}

type option struct {
	shortSymbolDate       string
	shortSymbolAmount     string
	shortSymbolOptionType string
	rowIdentifier         int
}

var (
	report EtradeTransactions
)

func init() {
}

func Run(fileName string) (err error) {

	var (
		preProcessFileData string
	)

	fileSize := getFileSize(fileName)

	if result := sizeWithinLimit(fileSize, FILESIZELIMIT); result {
		fmt.Printf("The file size (%v MB) is within the size limit (%v MB).\n", ConvertBytesTo(fileSize, MB), ConvertBytesTo(getSystemMemory(), MB))
		preProcessFileData, err = LoadFileAsString(fileName)
		buildEtradeTransactions(&preProcessFileData)
		preProcessFileData = ""
		analyzeTransactions()
		generateReport()
	} else {
		fmt.Printf("The file size (%v MB) exceeds the size limit (%v MB).\n", ConvertBytesTo(fileSize, MB), ConvertBytesTo(getSystemMemory(), MB))
	}

	return
}

func buildEtradeTransactions(csvData *string) {

	var (
		eTT = ETradeTransaction{
			rowIdentifier: -1,
			transDate:     "",
			transType:     "",
			securityType:  "",
			symbol:        "",
			quantity:      0,
			amount:        0,
			price:         0,
			commission:    0,
			description:   "",
		}
	)

	r := csv.NewReader(strings.NewReader(*csvData))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if !strings.Contains(err.Error(), "wrong number of fields") {
				log.Fatal(err)
			}
		}

		if len(record) == 2 {
			report.accountNumber = record[1]
		}
		if len(record) == 9 {
			if strings.Contains(record[0], "TransactionDate") {
				report.headers = record
			} else {
				eTT.rowIdentifier++
				eTT.transDate = record[0]
				eTT.transType = record[1]
				eTT.securityType = record[2]
				eTT.symbol = record[3]
				eTT.quantity, _ = strconv.Atoi(record[4])
				eTT.amount, _ = strconv.ParseFloat(record[5], 64)
				eTT.price, _ = strconv.ParseFloat(record[6], 64)
				eTT.commission, _ = strconv.ParseFloat(record[7], 64)
				eTT.description = record[8]
				report.lines = append(report.lines, eTT)
			}
		}
	}

}

func analyzeTransactions() {

	wg := new(sync.WaitGroup)
	wg.Add(4)
	go collectContribution(wg)
	go collectDividends(wg)
	go collectInterest(wg)
	go collectCommission(wg)
	wg.Wait()

	wg.Add(1)
	go collectOptions(wg)
	wg.Wait()

}

func collectDividends(wg *sync.WaitGroup) {

	for _, transaction := range report.lines {
		if strings.ToLower(transaction.transType) == DIVIDEND {
			report.totalDividends += transaction.amount
		}
	}

	defer wg.Done()

}

func collectInterest(wg *sync.WaitGroup) {

	for _, transaction := range report.lines {
		if strings.ToLower(transaction.transType) == INTEREST {
			report.totalInterest += transaction.amount
		}
	}

	defer wg.Done()
}

func collectCommission(wg *sync.WaitGroup) {

	for _, transaction := range report.lines {
		if transaction.commission > 0 {
			report.totalCommissions += transaction.commission
		}
	}

	defer wg.Done()
}

func collectContribution(wg *sync.WaitGroup) {

	for _, transaction := range report.lines {
		if strings.ToLower(transaction.transType) == CONTRIBUTION {
			report.totalContributions += transaction.amount
		}
	}

	defer wg.Done()
}

func collectOptions(wg *sync.WaitGroup) {

	var (
		w       = 0
		x, y, z string
	)

	for _, line := range report.lines {
		if strings.ToLower(line.securityType) == OPTION {
			if strings.Contains(line.symbol, YEARMARKER) {
				w = strings.Index(line.symbol, YEARMARKER) + 3
				x = line.symbol[:w]
				y = line.symbol[w:]
				if strings.Contains(strings.ToLower(y), PUT) {
					z = PUT
				} else {
					z = CALL
				}
				report.options = append(report.options, option{shortSymbolDate: x, shortSymbolAmount: strings.TrimSpace(y), shortSymbolOptionType: z, rowIdentifier: line.rowIdentifier})
			}
		}
	}

	sort.Slice(report.options, func(i, j int) bool {
		return report.options[i].shortSymbolDate < report.options[j].shortSymbolDate
	})

	defer wg.Done()
}

func generateReport() {

	var (
		lastShortSymbolDate       = ""
		optionTotal, runningTotal float64
	)

	fmt.Printf("\nReport for Account: %v\n\n", report.accountNumber)

	fmt.Printf("Total Contributions for the reporting period:\t %10.2f\n", report.totalContributions)
	fmt.Printf("Total Dividends for the reporting period:\t %10.2f\n", report.totalDividends)
	fmt.Printf("Total Interest for the reporting period:\t %10.2f\n", report.totalInterest)
	fmt.Printf("Total Commissions for the reporting period:\t %10.2f\n", report.totalCommissions)

	fmt.Println("\nOptions:")
	for _, option := range report.options {
		if option.shortSymbolDate == lastShortSymbolDate {
			optionTotal += report.lines[option.rowIdentifier].amount
			printOptionLine(option.rowIdentifier)
		} else {
			if len(lastShortSymbolDate) > 0 {
				printTotalLines(optionTotal)
			}
			lastShortSymbolDate = option.shortSymbolDate
			runningTotal += optionTotal
			optionTotal = 0
			optionTotal += report.lines[option.rowIdentifier].amount
			fmt.Printf("%v\n", report.lines[option.rowIdentifier].symbol)
			printOptionLine(option.rowIdentifier)
		}
	}
	fmt.Printf("%52v \n", "____________")
	fmt.Printf("%38v %12.2f\n", " ", optionTotal)

	fmt.Printf("\n %52v \n", "============")
	fmt.Printf("%38v %12.2f\n", " ", runningTotal)
}

func printOptionLine(idx int) {
	fmt.Printf("\t %v %20v %12.2f \n", report.lines[idx].transDate, report.lines[idx].transType, report.lines[idx].amount)
}

func printTotalLines(optionTotal float64) {
	fmt.Printf("%52v \n", "____________")
	fmt.Printf("%38v %12.2f\n", " ", optionTotal)
}
