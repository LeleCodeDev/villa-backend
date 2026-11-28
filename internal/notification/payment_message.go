package notification

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/lelecodedev/villa-backend/internal/model"
)

const (
	dateLayout     = "02/01/2006"
	dateTimeLayout = "02/01/2006, 15:04"
)

func PaidMessages(app *model.Application, p *model.Payment, b *model.Booking, e *model.EventBooking, reward *model.Voucher) []Message {
	loc := config.Env.Location

	var appTitle, appPhone string
	if app != nil {
		appTitle, appPhone = app.Title, app.PhoneNumber
	}

	var (
		title, customer, phone string
		schedule, items, paid  []string
		total, remaining       int64
		dueAt                  *time.Time
	)

	switch {
	case b != nil:
		customer, phone, total = b.User.Username, b.User.Phone, b.TotalPrice

		nights := max(int(b.CheckOutDate.Sub(b.CheckInDate).Hours()/24), 1)
		unit := "-"
		if b.RentalOption != nil {
			unit = b.RentalOption.Title
		}

		schedule = []string{
			"Check-in: " + b.CheckInDate.Format(dateLayout),
			fmt.Sprintf("Check-out: %s (%d malam)", b.CheckOutDate.Format(dateLayout), nights),
			"",
			fmt.Sprintf("Tamu menginap: %d orang", b.GuestCount),
			fmt.Sprintf("Tamu: %d orang", b.VisitorCount),
		}
		items = []string{line(unit, fmt.Sprintf("%d malam", nights), b.BookingPrice/int64(nights), b.BookingPrice)}
		for _, a := range b.BookingAddons {
			items = append(items, line(a.Addon.Title, strconv.Itoa(a.Quantity), a.Addon.Price, a.Addon.Price*int64(a.Quantity)))
		}

		switch {
		case p.Kind == model.PaymentKindDeposit:
			title, remaining, dueAt = "Bukti Pembayaran DP", total-p.Amount, b.BalanceDueAt
			paid = []string{"DP dibayar: " + formatRupiah(p.Amount)}
		case b.PaymentType == model.PaymentTypeDeposit:
			title = "Bukti Pelunasan"
			paid = []string{"DP dibayar: " + formatRupiah(total-p.Amount), "Pelunasan: " + formatRupiah(p.Amount)}
		default:
			title = "Bukti Pembayaran"
			paid = []string{"Dibayar: " + formatRupiah(p.Amount)}
		}

	case e != nil:
		customer, phone, total, title = e.User.Username, e.User.Phone, e.TotalPrice, "Bukti Pembayaran Acara"
		paid = []string{"Dibayar: " + formatRupiah(p.Amount)}

		unit, name := "", e.Title
		if e.VillaPackage == nil {
			name += " (Acara Kustom)"
		} else {
			switch e.VillaPackage.Type {
			case model.ChargePerHour:
				unit = "jam"
			case model.ChargePerDay:
				unit = "hari"
			}
		}

		schedule = []string{
			"Mulai: " + e.StartAt.In(loc).Format(dateTimeLayout),
			"Selesai: " + e.EndAt.In(loc).Format(dateTimeLayout),
		}
		if unit != "" {
			schedule = append(schedule, fmt.Sprintf("Durasi: %d %s", e.Quantity, unit))
		}
		schedule = append(schedule, fmt.Sprintf("Tamu: %d orang", e.GuestCount))

		switch {
		case unit != "":
			items = []string{line(name, fmt.Sprintf("%d %s", e.Quantity, unit), e.Price, e.TotalPrice)}
		case e.Quantity > 1:
			items = []string{line(name, fmt.Sprintf("%d paket", e.Quantity), e.Price, e.TotalPrice)}
		default:
			items = []string{name + "\n" + formatRupiah(e.TotalPrice)}
		}

	default:
		return nil
	}

	summary := append([]string{"Total tagihan: " + formatRupiah(total), ""}, paid...)
	statusHeading, statusBody := "✅ *Status*", "*LUNAS*"
	if remaining > 0 {
		statusHeading, statusBody = "⏳ *Status*", "*BELUM LUNAS (DP)*"
		summary = append(summary, "", "Sisa tagihan: "+formatRupiah(remaining))
		if dueAt != nil {
			summary = append(summary, "Jatuh tempo: "+dueAt.In(loc).Format(dateLayout))
		}
	}

	now := time.Now().In(loc).Format(dateTimeLayout)

	// Sections shared by the customer and admin messages.
	details := func(sb *strings.Builder) {
		section(sb, "📅 *Jadwal*", schedule...)
		section(sb, "🛒 *Rincian*", strings.Join(items, "\n\n"))
		section(sb, "💰 *Pembayaran*", summary...)
		section(sb, statusHeading, statusBody)
	}

	var user strings.Builder
	fmt.Fprintf(&user, "*%s*\n%s\n\n", strings.ToUpper(appTitle), title)
	section(&user, "📋 *Detail Pesanan*", "No. Pesanan: "+p.OrderID, "", "Pelanggan: "+customer, "Tanggal: "+now)
	details(&user)

	if reward != nil {
		lines := []string{"Kode: *" + reward.Code + "*", fmt.Sprintf("Diskon %d%%", reward.DiscountValue)}
		if reward.ExpiresAt != nil {
			lines = append(lines, "Berlaku sampai: "+reward.ExpiresAt.In(loc).Format(dateLayout))
		}
		lines = append(lines, "Bisa dipakai Anda atau keluarga.")
		section(&user, "🎁 *Voucher Untuk Anda*", lines...)
	}

	if appPhone != "" {
		fmt.Fprintf(&user, "📞 Pertanyaan? Hubungi %s\n", appPhone)
	}
	fmt.Fprintf(&user, "🙏 Terima kasih telah memilih %s.", appTitle)

	var admin strings.Builder
	fmt.Fprintf(&admin, "🔔 *Pembayaran Masuk*\n%s\n\n", title)
	section(&admin, "👤 *Pelanggan*", customer, "WhatsApp: "+phone)
	section(&admin, "📋 *Pesanan*", "No. Pesanan: "+p.OrderID, "Tanggal: "+now)
	details(&admin)

	return []Message{
		{Phone: phone, Text: user.String()},
		{Phone: appPhone, Text: admin.String()},
	}
}
