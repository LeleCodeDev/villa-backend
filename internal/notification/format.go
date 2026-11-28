// Package notification
package notification

import (
	"fmt"
	"strconv"
	"strings"
)

func formatRupiah(amount int64) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	amountString := strconv.FormatInt(amount, 10)

	var builder strings.Builder
	for i, char := range amountString {
		if i != 0 && (len(amountString)-i)%3 == 0 {
			builder.WriteByte('.')
		}

		builder.WriteRune(char)
	}

	return sign + "Rp " + builder.String()
}

func line(name, qty string, price, subtotal int64) string {
	return fmt.Sprintf("%s\n%s x %s = %s", name, qty, formatRupiah(price), formatRupiah(subtotal))
}
