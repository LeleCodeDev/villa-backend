package notification

import (
	"fmt"
	"strings"

	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func requestDetails(r *model.EventBookingRequest) []string {
	loc := config.Env.Location

	pkg := "Acara Kustom"
	if r.VillaPackage != nil {
		pkg = r.VillaPackage.Title
	}

	lines := []string{
		"Acara: " + r.Title,
		"Paket: " + pkg,
		"",
		"Mulai: " + r.StartAt.In(loc).Format("02/01/2006 15:04"),
		"Selesai: " + r.EndAt.In(loc).Format("02/01/2006 15:04"),
		"",
		fmt.Sprintf("Tamu: %d orang", r.GuestCount),
	}
	if r.Budget != nil {
		lines = append(lines, "Budget: "+formatRupiah(*r.Budget))
	}
	if r.Description != "" {
		lines = append(lines, "", "Catatan: "+r.Description)
	}

	return lines
}

func RequestApproved(app *model.Application, r *model.EventBookingRequest) []Message {
	var sb strings.Builder

	fmt.Fprintf(&sb, "🧾 *%s*\nPermintaan Acara Disetujui\n%s\n\n", strings.ToUpper(app.Title), divider)
	section(&sb, "✅ *Status*", "*DISETUJUI*")
	section(&sb, "🎟️ *Detail Permintaan*", requestDetails(r)...)
	sb.WriteString("📌 Silakan lanjutkan ke pembayaran agar jadwal Anda terkonfirmasi.\n\n")
	fmt.Fprintf(&sb, "📞 Pertanyaan? Hubungi %s\n", app.PhoneNumber)
	fmt.Fprintf(&sb, "🙏 Terima kasih telah memilih %s.", app.Title)

	return []Message{{Phone: r.User.Phone, Text: sb.String()}}
}

func RequestRejected(app *model.Application, r *model.EventBookingRequest) []Message {
	status := []string{"*DITOLAK*"}
	if r.RejectReason != nil && *r.RejectReason != "" {
		status = append(status, "", "Alasan: "+*r.RejectReason)
	}

	var sb strings.Builder

	fmt.Fprintf(&sb, "🧾 *%s*\nPermintaan Acara Ditolak\n%s\n\n", strings.ToUpper(app.Title), divider)
	section(&sb, "❌ *Status*", status...)
	section(&sb, "🎟️ *Detail Permintaan*", requestDetails(r)...)
	fmt.Fprintf(&sb, "📞 Pertanyaan? Hubungi %s\n", app.PhoneNumber)
	fmt.Fprintf(&sb, "🙏 Terima kasih telah memilih %s.", app.Title)

	return []Message{{Phone: r.User.Phone, Text: sb.String()}}
}

func RequestCreated(app *model.Application, r *model.EventBookingRequest) []Message {
	var sb strings.Builder

	fmt.Fprintf(&sb, "🔔 *%s*\nPermintaan Acara Baru\n%s\n\n", strings.ToUpper(app.Title), divider)
	section(&sb, "👤 *Pemohon*", "Nama: "+r.User.Username, "No. HP: "+r.User.Phone)
	section(&sb, "🎟️ *Detail Permintaan*", requestDetails(r)...)
	sb.WriteString("⏳ Menunggu persetujuan.")

	return []Message{{Phone: app.PhoneNumber, Text: sb.String()}}
}

func RequestCanceled(app *model.Application, r *model.EventBookingRequest) []Message {
	var sb strings.Builder

	fmt.Fprintf(&sb, "🔔 *%s*\nPermintaan Acara Dibatalkan\n%s\n\n", strings.ToUpper(app.Title), divider)
	section(&sb, "👤 *Pemohon*", "Nama: "+r.User.Username, "No. HP: "+r.User.Phone)
	section(&sb, "🎟️ *Detail Permintaan*", requestDetails(r)...)
	sb.WriteString("🚫 Dibatalkan oleh pemohon.")

	return []Message{{Phone: app.PhoneNumber, Text: sb.String()}}
}
