package main

import (
	"flag"
	"fmt"
	"io"
	"net/mail"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/emersion/go-mbox"
)

type Payment struct {
	Date    time.Time
	Subject string
	Amount  string
}

type ByDate []Payment

func (s ByDate) Len() int           { return len(s) }
func (s ByDate) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByDate) Less(i, j int) bool { return s[i].Date.Before(s[j].Date) }

var moneyRegex = regexp.MustCompile(`(&#36;|\$)[0-9,.]+`)

func isPayment(msg *mail.Message) bool {
	subject := strings.ToLower(msg.Header.Get("Subject"))
	return strings.Contains(subject, "payment") && strings.Contains(subject, "scheduling")
}

func extractAmount(msg *mail.Message) (amount string, err error) {
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return
	}

	match := moneyRegex.Find(body)
	if match == nil {
		// Dump body, float error to top and kill execution
		fmt.Println(string(body))
		err = fmt.Errorf("could not find money in body")
		return
	}

	amount = string(match)
	amount = strings.Replace(amount, "$", "", 1)
	amount = strings.Replace(amount, "&#36;", "", 1)
	amount = strings.Replace(amount, ",", "", 1)

	return
}

func getPayments(path string) (payments []Payment, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}

	payments = make([]Payment, 0, 1024)
	r := mbox.NewReader(f)
	for {
		mboxMsg, err := r.NextMessage()
		if err == io.EOF {
			break
		}
		if err != nil {
			return payments, err
		}

		msg, err := mail.ReadMessage(mboxMsg)
		if err != nil {
			return payments, err
		}

		if !isPayment(msg) {
			continue
		}

		amount, err := extractAmount(msg)
		if err != nil {
			// continue
			return payments, err
		}

		date, err := msg.Header.Date()
		if err != nil {
			return payments, err
		}

		p := Payment{
			Date:    date,
			Subject: msg.Header.Get("Subject"),
			Amount:  amount,
		}
		payments = append(payments, p)
	}
	return
}

func main() {
	flag.Parse()

	path := flag.Arg(0)
	if path == "" {
		fmt.Printf("Usage: %s path/to/mbox\n", flag.Args()[0])
		os.Exit(1)
	}

	payments, err := getPayments(path)
	if err != nil {
		fmt.Printf("Error during processing: %v\n", err)
		os.Exit(1)
	}

	sort.Sort(ByDate(payments))

	for _, payment := range payments {
		fmt.Printf("%s\t%s\n", payment.Date.Format("01/02/2006"), payment.Amount)
	}
}
