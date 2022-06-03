Extracts payment information from an email mbox archive.

1. File all of your credit card payment notices into a single label in Gmail
2. Collect this data over a (long) time
3. Download the data using [Google Takeout](https://takeout.google.com)
4. `go get github.com/jbaikge/mail-payments`
5. `mail-payments path/to/mbox/file`
6. Output is a date and amount, separated by a tab. This allows copy and paste into any spreadsheet for further analysis
