# indobank-go

Daftar bank di Indonesia (nama dan sandi bank) untuk Go. Data tertanam di binary lewat `go:embed`: tanpa database, tanpa jaringan, tanpa dependensi.

Versi PHP/Laravel: [andes2912/indobank](https://github.com/andes2912/indobank). Versi Flutter/Dart: [andes2912/indobank-dart](https://github.com/andes2912/indobank-dart).

Data diambil dari **Tabel Sandi Bank resmi BCA per 31 Maret 2026** (123 entri bank umum dan Unit Usaha Syariah).

## Instalasi

```
go get github.com/andes2912/indobank-go
```

Membutuhkan Go 1.21 atau lebih baru.

## Pemakaian

```go
package main

import (
	"fmt"

	indobank "github.com/andes2912/indobank-go"
)

func main() {
	// Semua bank
	all := indobank.Banks()

	// Pencarian sebagian nama (abaikan huruf besar/kecil, tanda baca, spasi ganda)
	hasil := indobank.Search("bri")

	// Kecocokan persis berdasarkan kolom
	bca, ok := indobank.Find("sandi_bank", "014")

	// Semua bank dengan sandi tertentu (sandi tidak unik)
	cimb := indobank.FindBySandi("022")

	// Bank syariah (heuristik berdasarkan nama)
	syariah := indobank.Syariah()

	// Ekspor JSON
	data, _ := indobank.JSON()

	fmt.Println(len(all), len(hasil), bca, ok, len(cimb), len(syariah), len(data))
}
```

## API

| Fungsi | Keterangan |
|---|---|
| `Banks() []Bank` | Semua bank (salinan). |
| `Search(name string) []Bank` | Cari berdasarkan sebagian nama. Query kosong mengembalikan semua. |
| `Find(key, value string) (Bank, bool)` | Cocok persis pada `nama_bank` atau `sandi_bank`; mengembalikan yang pertama. |
| `FindBySandi(code string) []Bank` | Semua bank dengan sandi tersebut. |
| `Syariah() []Bank` | Bank yang namanya mengandung "syariah". |
| `(Bank).IsSyariah() bool` | Heuristik nama, bukan data resmi. |
| `JSON() ([]byte, error)` | Semua bank sebagai array JSON. |

`Bank` memiliki field `Nama` (`nama_bank`) dan `Sandi` (`sandi_bank`). `Sandi` adalah string agar nol di depan tetap ada (mis. `001`).

> Sandi bank tidak unik. Contoh: `022` dipakai Bank CIMB Niaga dan Unit Usaha Syariah-nya. Gunakan `FindBySandi` bila perlu semua yang cocok.

## Memperbarui data

1. Ganti `data/bank.csv` (format sama dengan repo PHP: `"nama_bank",sandi_bank`).
2. Sesuaikan jumlah di `TestBanksCount` bila jumlah bank berubah.
3. `go test ./...`

## Lisensi

MIT
