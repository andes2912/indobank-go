// Package indobank menyediakan daftar bank di Indonesia (nama dan sandi bank)
// dengan data tertanam, tanpa database dan tanpa akses jaringan.
//
// Sumber data: Tabel Sandi Bank resmi BCA per 31 Maret 2026.
package indobank

import (
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"strings"
	"sync"
	"unicode"
)

//go:embed data/bank.csv
var rawCSV string

// Bank adalah satu entri bank.
type Bank struct {
	Nama  string `json:"nama_bank"`
	Sandi string `json:"sandi_bank"`
}

var (
	once  sync.Once
	banks []Bank
)

func load() {
	once.Do(func() {
		r := csv.NewReader(strings.NewReader(rawCSV))
		r.FieldsPerRecord = 2
		records, err := r.ReadAll()
		if err != nil {
			panic("indobank: data/bank.csv tidak valid: " + err.Error())
		}
		for i, rec := range records {
			if i == 0 { // header
				continue
			}
			banks = append(banks, Bank{Nama: rec[0], Sandi: rec[1]})
		}
	})
}

// Banks mengembalikan semua bank (salinan, aman diubah pemanggil).
func Banks() []Bank {
	load()
	out := make([]Bank, len(banks))
	copy(out, banks)
	return out
}

// normalize mengabaikan huruf besar/kecil, tanda baca, dan spasi ganda.
func normalize(s string) string {
	var b strings.Builder
	space := false
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if space && b.Len() > 0 {
				b.WriteRune(' ')
			}
			space = false
			b.WriteRune(r)
		} else {
			space = true
		}
	}
	return b.String()
}

// Search mencari bank berdasarkan sebagian nama. Huruf besar/kecil, tanda baca,
// dan spasi ganda diabaikan. Query kosong mengembalikan semua bank.
func Search(name string) []Bank {
	load()
	needle := normalize(name)
	if needle == "" {
		return Banks()
	}
	var out []Bank
	for _, b := range banks {
		if strings.Contains(normalize(b.Nama), needle) {
			out = append(out, b)
		}
	}
	return out
}

// Find mencari satu bank dengan kecocokan persis pada kolom key
// ("nama_bank" atau "sandi_bank"). Mengembalikan bank pertama yang cocok.
func Find(key, value string) (Bank, bool) {
	load()
	for _, b := range banks {
		switch key {
		case "nama_bank":
			if b.Nama == value {
				return b, true
			}
		case "sandi_bank":
			if b.Sandi == value {
				return b, true
			}
		}
	}
	return Bank{}, false
}

// FindBySandi mengembalikan semua bank dengan sandi tertentu. Sandi tidak unik
// (mis. 022 dipakai CIMB Niaga dan Unit Usaha Syariah-nya).
func FindBySandi(code string) []Bank {
	load()
	var out []Bank
	for _, b := range banks {
		if b.Sandi == code {
			out = append(out, b)
		}
	}
	return out
}

// IsSyariah memperkirakan apakah bank syariah berdasarkan namanya
// (mengandung kata "syariah"). Ini heuristik, bukan data resmi.
func (b Bank) IsSyariah() bool {
	return strings.Contains(strings.ToLower(b.Nama), "syariah")
}

// Syariah mengembalikan bank yang namanya mengandung "syariah".
func Syariah() []Bank {
	load()
	var out []Bank
	for _, b := range banks {
		if b.IsSyariah() {
			out = append(out, b)
		}
	}
	return out
}

// JSON mengekspor semua bank sebagai array JSON.
func JSON() ([]byte, error) {
	load()
	return json.Marshal(banks)
}
