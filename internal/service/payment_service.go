package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/client"
	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/payment"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"gorm.io/gorm"
)

type PaymentService struct {
	txManager        *repository.TxManager
	paymentRepo      *repository.PaymentRepository
	bookingRepo      *repository.BookingRepository
	eventBookingRepo *repository.EventBookingRepository
	applicationRepo  *repository.ApplicationRepository
	voucherService   *VoucherService
	gateway          *payment.MidtransGateway
	client           *client.WhatsappClient
}

func NewPaymentService(
	txManager *repository.TxManager,
	paymentRepo *repository.PaymentRepository,
	bookingRepo *repository.BookingRepository,
	eventBookingRepo *repository.EventBookingRepository,
	applicationRepo *repository.ApplicationRepository,
	voucherService *VoucherService,
	gateway *payment.MidtransGateway,
	client *client.WhatsappClient,
) *PaymentService {
	return &PaymentService{
		txManager:        txManager,
		paymentRepo:      paymentRepo,
		bookingRepo:      bookingRepo,
		eventBookingRepo: eventBookingRepo,
		applicationRepo:  applicationRepo,
		voucherService:   voucherService,
		gateway:          gateway,
		client:           client,
	}
}

func (s *PaymentService) HandleNotificationCallback(ctx context.Context, n dto.MidtransCallbackRequest) error {
	p, err := s.paymentRepo.GetByOrderID(ctx, n.OrderID)
	if err != nil {
		return err
	}
	if p == nil {
		return appError.NotFound(fmt.Sprintf("Payment not found with order ID : %s", n.OrderID))
	}

	status, err := s.gateway.CheckTransaction(ctx, n.OrderID)
	if err != nil {
		return err
	}

	grossAmount, err := strconv.ParseFloat(n.GrossAmount, 64)
	if err != nil || int64(math.Round(grossAmount)) != p.Amount {
		return appError.Conflict("Amount mismatch")
	}

	target, ok := mapPaymentStatus(status.TransactionStatus, status.FraudStatus)
	if !ok || p.Status == target {
		return nil
	}

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txPaymentRepo := s.paymentRepo.WithTx(tx)
		txBookingRepo := s.bookingRepo.WithTx(tx)
		txEventBookingRepo := s.eventBookingRepo.WithTx(tx)

		updated, err := txPaymentRepo.TransitionStatus(ctx, p.ID, model.PaymentStatusPending, target)
		if err != nil {
			return err
		}
		if !updated {
			if target == model.PaymentStatusPaid {
				log.Printf("REVIEW NEEDED: payment %s paid but was already %s", p.OrderID, p.Status)
			}

			return nil
		}

		if p.BookingID != nil {
			from, to, ok := bookingTransition(p.Kind, target)
			if !ok {
				return nil
			}

			moved, err := txBookingRepo.TransitionStatus(ctx, *p.BookingID, from, to)
			if err != nil {
				return err
			}

			if !moved && target == model.PaymentStatusPaid {
				log.Printf("REVIEW NEEDED: payment %s paid but booking not in status %s", p.OrderID, from)
			}

			if p.Booking != nil && (to == model.BookingStatusPaid || to == model.BookingStatusPartiallyPaid) {
				go s.notify(p, p.Booking, nil)
			}

		}

		if p.EventBookingID != nil {
			from, to, ok := eventBookingTransition(p.Kind, target)
			if !ok {
				return nil
			}

			moved, err := txEventBookingRepo.TransitionStatus(ctx, *p.EventBookingID, from, to)
			if err != nil {
				return err
			}

			if !moved && target == model.PaymentStatusPaid {
				log.Printf("REVIEW NEEDED: payment %s paid but event booking not in status %s", p.OrderID, from)
			}

			if p.EventBooking != nil && to == model.EventBookingStatusPaid {
				go s.notify(p, nil, p.EventBooking)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

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

func (s *PaymentService) notify(payment *model.Payment, booking *model.Booking, eventBooking *model.EventBooking) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("notify panic for %s: %v", payment.OrderID, r)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	app, err := s.applicationRepo.Get(ctx)
	if err != nil {
		log.Printf("notify %s: get application: %v", payment.OrderID, err)
	}

	var reward *model.Voucher
	if booking != nil && payment.Kind == model.PaymentKindFull {
		voucher, err := s.voucherService.CreateRewardVoucher(ctx, booking)
		if err != nil {
			log.Printf("notify %s: failed to create voucher: %v", payment.OrderID, err)
		}

		reward = voucher
	}

	userMsg, adminMsg, userPhone := buildPaidMessages(app, payment, booking, eventBooking, reward)
	if userMsg == "" {
		return
	}

	if err := s.client.SendText(ctx, toChatID(userPhone), userMsg); err != nil {
		log.Printf("notify %s: send to user %s: %v", payment.OrderID, userPhone, err)
	}

	if app != nil {
		if err := s.client.SendText(ctx, toChatID(app.PhoneNumber), adminMsg); err != nil {
			log.Printf("notify %s: send to admin %s: %v", payment.OrderID, app.PhoneNumber, err)
		}
	}
}

func buildPaidMessages(app *model.Application, p *model.Payment, b *model.Booking, e *model.EventBooking, reward *model.Voucher) (userMsg, adminMsg, phone string) {
	loc := config.Env.Location

	var appTitle, appPhone string
	if app != nil {
		appTitle, appPhone = app.Title, app.PhoneNumber
	}

	var (
		title, customer  string
		schedule, items  []string
		total, remaining int64
		dueAt            *time.Time
	)

	var paid string

	const divider = "=============================="

	switch {
	case b != nil:
		customer, phone, total = b.User.Username, b.User.Phone, b.TotalPrice

		nights := max(int(b.CheckOutDate.Sub(b.CheckInDate).Hours()/24), 1)
		unit := "-"
		if b.RentalOption != nil {
			unit = b.RentalOption.Title
		}

		schedule = []string{
			"Check-in: " + b.CheckInDate.Format("02/01/2006"),
			fmt.Sprintf("Check-out: %s (%d malam)\n", b.CheckOutDate.Format("02/01/2006"), nights),
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
			paid = "\nDP dibayar: " + formatRupiah(p.Amount)
		case b.PaymentType == model.PaymentTypeDeposit:
			title = "Bukti Pelunasan"
			paid = "\nDP dibayar: " + formatRupiah(total-p.Amount) + "\nPelunasan: " + formatRupiah(p.Amount)
		default:
			title = "Bukti Pembayaran"
			paid = "\nDibayar: " + formatRupiah(p.Amount)
		}

	case e != nil:
		customer, phone, total, title = e.User.Username, e.User.Phone, e.TotalPrice, "Bukti Pembayaran Acara"

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
			"Mulai: " + e.StartAt.In(loc).Format("02/01/2006 15:04"),
			"Selesai: " + e.EndAt.In(loc).Format("02/01/2006 15:04"),
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
		return "", "", ""
	}

	section := func(sb *strings.Builder, heading string, body ...string) {
		fmt.Fprintf(sb, "%s\n%s\n%s\n\n", heading, divider, strings.Join(body, "\n"))
	}

	summary := []string{
		"Total tagihan: " + formatRupiah(total),
		paid,
		"Dibayar: " + formatRupiah(p.Amount),
	}
	if remaining > 0 {
		summary = append(summary, "", "Sisa tagihan: "+formatRupiah(remaining))
		if dueAt != nil {
			summary = append(summary, "Jatuh tempo: "+dueAt.In(loc).Format("02/01/2006"))
		}
	}

	statusHeading, statusBody := "✅ *Status*", "*LUNAS*"
	if remaining > 0 {
		statusHeading, statusBody = "⏳ *Status*", "*BELUM LUNAS (DP)*"
	}

	var user strings.Builder

	fmt.Fprintf(&user, "*%s*\n%s\n\n", strings.ToUpper(appTitle), title)

	section(&user, "📋 *Detail Pesanan*",
		"No. Pesanan: "+p.OrderID+"\n",
		"Pelanggan: "+customer,
		"Tanggal: "+time.Now().In(loc).Format("02/01/2006, 15:04"),
	)
	section(&user, "📅 *Jadwal*", schedule...)
	section(&user, "🛒 *Rincian*", strings.Join(items, "\n\n"))
	section(&user, "💰 *Pembayaran*", summary...)
	section(&user, statusHeading, statusBody)

	if reward != nil {
		voucherLines := []string{"Kode: *" + reward.Code + "*"}
		voucherLines = append(voucherLines, fmt.Sprintf("Diskon %d%%", reward.DiscountValue))

		if reward.ExpiresAt != nil {
			voucherLines = append(voucherLines, "Berlaku sampai: "+reward.ExpiresAt.In(loc).Format("02/01/2006"))
		}

		voucherLines = append(voucherLines, "Bisa dipakai Anda atau keluarga.")
		section(&user, "🎁 *Voucher Untuk Anda*", voucherLines...)
	}

	if appPhone != "" {
		fmt.Fprintf(&user, "📞 Pertanyaan? Hubungi %s \n", appPhone)
	}
	fmt.Fprintf(&user, "🙏 Terima kasih telah memilih %s.", appTitle)

	var admin strings.Builder

	fmt.Fprintf(&admin, "🔔 *Pembayaran Masuk*\n%s\n\n", title)

	section(&admin, "👤 *Pelanggan*", customer, "WhatsApp: "+phone)
	section(&admin, "📋 *Pesanan*",
		"No. Pesanan: "+p.OrderID,
		"Tanggal: "+time.Now().In(loc).Format("02/01/2006, 15:04"),
	)
	section(&admin, "📅 *Jadwal*", schedule...)
	section(&admin, "🛒 *Rincian*", strings.Join(items, "\n\n"))
	section(&admin, "💰 *Pembayaran*", summary...)
	section(&admin, statusHeading, statusBody)

	return user.String(), strings.TrimRight(admin.String(), "\n"), phone
}

func line(name, qty string, price, subtotal int64) string {
	return fmt.Sprintf("%s\n%s x %s = %s", name, qty, formatRupiah(price), formatRupiah(subtotal))
}

func toChatID(phone string) string {
	return "62" + strings.TrimPrefix(phone, "0") + "@c.us"
}

func mapPaymentStatus(transactionStatus, fraudStatus string) (model.PaymentStatus, bool) {
	switch transactionStatus {
	case "settlement":
		return model.PaymentStatusPaid, true
	case "capture":
		if fraudStatus == "accept" {
			return model.PaymentStatusPaid, true
		}
	case "expire":
		return model.PaymentStatusExpired, true
	case "cancel":
		return model.PaymentStatusCanceled, true
	case "deny":
		return model.PaymentStatusFailed, true
	}
	return "", false
}

func bookingTransition(kind model.PaymentKind, target model.PaymentStatus) (from, to model.BookingStatus, ok bool) {
	switch target {
	case model.PaymentStatusPaid:
		switch kind {
		case model.PaymentKindDeposit:
			return model.BookingStatusPending, model.BookingStatusPartiallyPaid, true
		case model.PaymentKindFull:
			return model.BookingStatusPending, model.BookingStatusPaid, true
		case model.PaymentKindBalance:
			return model.BookingStatusPartiallyPaid, model.BookingStatusPaid, true
		}

	case model.PaymentStatusExpired, model.PaymentStatusFailed, model.PaymentStatusCanceled:
		if kind == model.PaymentKindBalance {
			return "", "", false
		}

		switch target {
		case model.PaymentStatusExpired:
			return model.BookingStatusPending, model.BookingStatusExpired, true
		default:
			return model.BookingStatusPending, model.BookingStatusCanceled, true
		}
	}

	return "", "", false
}

func eventBookingTransition(kind model.PaymentKind, target model.PaymentStatus) (from, to model.EventBookingStatus, ok bool) {
	if kind == model.PaymentKindFull && target == model.PaymentStatusPaid {
		return model.EventBookingStatusPending, model.EventBookingStatusPaid, true
	}
	return "", "", false
}
