package indobank

import (
	"encoding/json"
	"testing"
)

func TestBanksCount(t *testing.T) {
	if n := len(Banks()); n != 123 {
		t.Fatalf("jumlah bank = %d, want 123", n)
	}
}

func TestBanksReturnsCopy(t *testing.T) {
	a := Banks()
	a[0].Nama = "DIUBAH"
	if Banks()[0].Nama == "DIUBAH" {
		t.Fatal("Banks() harus mengembalikan salinan")
	}
}

func TestSandiKeepsLeadingZero(t *testing.T) {
	b, ok := Find("nama_bank", "BANK INDONESIA")
	if !ok || b.Sandi != "001" {
		t.Fatalf("got %+v ok=%v, want sandi 001", b, ok)
	}
}

func TestSearchCaseInsensitivePartial(t *testing.T) {
	if len(Search("mandiri")) == 0 {
		t.Fatal("mandiri harus ditemukan")
	}
	if len(Search("MANDIRI")) != len(Search("mandiri")) {
		t.Fatal("pencarian harus case-insensitive")
	}
}

func TestSearchIgnoresPunctuationAndSpaces(t *testing.T) {
	if len(Search("bank   bri")) == 0 {
		t.Fatal("spasi ganda harus diabaikan")
	}
	if len(Search("(bank rakyat indonesia)")) == 0 {
		t.Fatal("tanda baca harus diabaikan")
	}
}

func TestSearchEmptyReturnsAll(t *testing.T) {
	if len(Search("")) != len(Banks()) {
		t.Fatal("query kosong harus mengembalikan semua bank")
	}
}

func TestFindBCA(t *testing.T) {
	b, ok := Find("sandi_bank", "014")
	if !ok || b.Nama == "" {
		t.Fatalf("sandi 014 tidak ditemukan: %+v", b)
	}
}

func TestFindMissing(t *testing.T) {
	if _, ok := Find("sandi_bank", "xyz"); ok {
		t.Fatal("tidak boleh ditemukan")
	}
	if _, ok := Find("kolom_tidak_ada", "014"); ok {
		t.Fatal("kolom tidak dikenal tidak boleh cocok")
	}
}

func TestFindBySandiNonUnique(t *testing.T) {
	if n := len(FindBySandi("022")); n < 2 {
		t.Fatalf("sandi 022 dipakai >= 2 bank, got %d", n)
	}
}

func TestSyariah(t *testing.T) {
	s := Syariah()
	if len(s) == 0 {
		t.Fatal("harus ada bank syariah")
	}
	for _, b := range s {
		if !b.IsSyariah() {
			t.Fatalf("%q bukan syariah", b.Nama)
		}
	}
}

func TestJSON(t *testing.T) {
	data, err := JSON()
	if err != nil {
		t.Fatal(err)
	}
	var out []map[string]string
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if len(out) != 123 || out[0]["sandi_bank"] != "001" {
		t.Fatalf("JSON tidak sesuai: len=%d first=%v", len(out), out[0])
	}
}
