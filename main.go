package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"math"
	"strings"
	"time"
)

type Mahasiswa struct {
	NIM           string
	Nama          string
	Jurusan       string
	Angkatan      int
	TotalSKS      int
	IPK           float64
	Prestasi      int
	NilaiSemester map[int]SemesterData
}

type SemesterData struct {
	MataKuliah []MataKuliahData
	IPS        float64
	TotalSKS   int
}

type MataKuliahData struct {
	Kode  string
	Nama  string
	SKS   int
	Nilai NilaiKomponen
}

type NilaiKomponen struct {
	UTS   float64
	UAS   float64
	Quiz  float64
	Total float64
	Grade string
	IP    float64
}

type LogAktivitas struct {
	Tanggal   time.Time
	JenisAksi string
	Deskripsi string
	Pengguna  string
}

type RegressionResult struct {
	Slope      float64
	Intercept  float64
	R2         float64
	Prediction float64
}

const (
	MAX_MAHASISWA = 100
	MAX_LOG       = 500
)

var (
	mahasiswa       [MAX_MAHASISWA]Mahasiswa
	logAktivitas    [MAX_LOG]LogAktivitas
	jumlahMahasiswa int
	jumlahLog       int
)

func main() {
	for {
		fmt.Println("\n=== Sistem Pengelolaan Nilai Mahasiswa ===")
		fmt.Println("1. Kelola Data Mahasiswa")
		fmt.Println("2. Kelola Nilai Per Semester")
		fmt.Println("3. Tampilkan Data")
		fmt.Println("4. Log Aktivitas")
		fmt.Println("5. Ekspor Transkrip")
		fmt.Println("6. Prediksi Kinerja Akademik")
		fmt.Println("0. Keluar")
		fmt.Print("Pilihan: ")
		var pilihan int
		fmt.Scan(&pilihan)

		switch pilihan {
		case 1:
			menuKelolaMahasiswa()
		case 2:
			menuKelolaNilaiSemester()
		case 3:
			menuTampilkanData()
		case 4:
			tampilkanLog()
		case 5:
			eksporTranskrip()
		case 6:
			prediksiKinerjaAkademik()
		case 0:
			fmt.Println("Keluar dari program. Terima kasih!")
			return
		default:
			fmt.Println("Pilihan tidak valid. Silakan coba lagi.")
		}
	}
}

func menuKelolaMahasiswa() {
	for {
		fmt.Println("\n=== Kelola Data Mahasiswa ===")
		fmt.Println("1. Tambah Mahasiswa")
		fmt.Println("2. Ubah Mahasiswa")
		fmt.Println("3. Hapus Mahasiswa")
		fmt.Println("4. Kembali")
		fmt.Print("Pilihan: ")
		var pilihan int
		fmt.Scan(&pilihan)

		switch pilihan {
		case 1:
			tambahMahasiswa()
		case 2:
			ubahMahasiswa()
		case 3:
			hapusMahasiswa()
		case 4:
			return
		default:
			fmt.Println("Pilihan tidak valid. Silakan coba lagi.")
		}
	}
}

func tambahMahasiswa() {
	if jumlahMahasiswa >= MAX_MAHASISWA {
		fmt.Println("Data mahasiswa penuh.")
		return
	}

	var m Mahasiswa
	var buffer string

	fmt.Print("NIM: ")
	fmt.Scan(&m.NIM)
	fmt.Scanln(&buffer)

	fmt.Print("Nama: ")
	fmt.Scanln(&m.Nama)

	fmt.Print("Jurusan: ")
	fmt.Scanln(&m.Jurusan)

	fmt.Print("Angkatan: ")
	fmt.Scan(&m.Angkatan)

	m.NilaiSemester = make(map[int]SemesterData)

	mahasiswa[jumlahMahasiswa] = m
	jumlahMahasiswa++
	tambahLog("Tambah", fmt.Sprintf("Mahasiswa %s - %s", m.NIM, m.Nama), "Admin")
	fmt.Println("Mahasiswa berhasil ditambahkan.")
}

func ubahMahasiswa() {
	fmt.Print("Masukkan NIM mahasiswa yang akan diubah: ")
	var nim string
	fmt.Scan(&nim)
	idx := cariMahasiswaSequential(nim)
	if idx == -1 {
		fmt.Println("Mahasiswa tidak ditemukan.")
		return
	}

	var buffer string
	fmt.Scanln(&buffer)

	fmt.Print("Nama baru: ")
	fmt.Scanln(&mahasiswa[idx].Nama)

	fmt.Print("Jurusan baru: ")
	fmt.Scanln(&mahasiswa[idx].Jurusan)

	fmt.Print("Angkatan baru: ")
	fmt.Scan(&mahasiswa[idx].Angkatan)

	tambahLog("Ubah", fmt.Sprintf("Mahasiswa %s - %s", nim, mahasiswa[idx].Nama), "Admin")
	fmt.Println("Data mahasiswa berhasil diubah.")
}

func hapusMahasiswa() {
	fmt.Print("Masukkan NIM mahasiswa yang akan dihapus: ")
	var nim string
	fmt.Scan(&nim)
	idx := cariMahasiswaSequential(nim)
	if idx == -1 {
		fmt.Println("Mahasiswa tidak ditemukan.")
		return
	}

	for i := idx; i < jumlahMahasiswa-1; i++ {
		mahasiswa[i] = mahasiswa[i+1]
	}
	jumlahMahasiswa--
	tambahLog("Hapus", fmt.Sprintf("Mahasiswa %s", nim), "Admin")
	fmt.Println("Data mahasiswa berhasil dihapus.")
}

func cariMahasiswaSequential(nim string) int {
	for i := 0; i < jumlahMahasiswa; i++ {
		if mahasiswa[i].NIM == nim {
			return i
		}
	}
	return -1
}

func tambahLog(jenisAksi, deskripsi, pengguna string) {
	if jumlahLog >= MAX_LOG {
		fmt.Println("Log penuh, tidak dapat mencatat aktivitas baru.")
		return
	}

	logAktivitas[jumlahLog] = LogAktivitas{
		Tanggal:   time.Now(),
		JenisAksi: jenisAksi,
		Deskripsi: deskripsi,
		Pengguna:  pengguna,
	}
	jumlahLog++
}

func tampilkanLog() {
	fmt.Println("\n=== Log Aktivitas ===")
	for i := 0; i < jumlahLog; i++ {
		log := logAktivitas[i]
		fmt.Printf("[%s] %s - %s oleh %s\n",
			log.Tanggal.Format("2006-01-02 15:04:05"),
			log.JenisAksi,
			log.Deskripsi,
			log.Pengguna)
	}
}

func menuKelolaNilaiSemester() {
	fmt.Print("Masukkan NIM mahasiswa: ")
	var nim string
	fmt.Scan(&nim)
	idx := cariMahasiswaSequential(nim)
	if idx == -1 {
		fmt.Println("Mahasiswa tidak ditemukan.")
		return
	}

	fmt.Print("Masukkan semester (1-8): ")
	var semester int
	fmt.Scan(&semester)
	if semester < 1 || semester > 8 {
		fmt.Println("Semester tidak valid.")
		return
	}

	for {
		fmt.Printf("\n=== Kelola Nilai Semester %d ===\n", semester)
		fmt.Println("1. Tambah Mata Kuliah dan Nilai")
		fmt.Println("2. Edit Nilai Mata Kuliah")
		fmt.Println("3. Hapus Mata Kuliah")
		fmt.Println("4. Tampilkan Nilai Semester")
		fmt.Println("5. Kembali")
		fmt.Print("Pilihan: ")
		var pilihan int
		fmt.Scan(&pilihan)

		switch pilihan {
		case 1:
			tambahMataKuliahDanNilai(idx, semester)
		case 2:
			editNilaiMataKuliah(idx, semester)
		case 3:
			hapusMataKuliah(idx, semester)
		case 4:
			tampilkanNilaiSemester(idx, semester)
		case 5:
			return
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

func tambahMataKuliahDanNilai(mahasiswaIdx, semester int) {
	var mk MataKuliahData
	var buffer string

	fmt.Scanln(&buffer)

	fmt.Print("Kode Mata Kuliah: ")
	fmt.Scanln(&mk.Kode)

	fmt.Print("Nama Mata Kuliah: ")
	fmt.Scanln(&mk.Nama)

	fmt.Print("SKS: ")
	fmt.Scan(&mk.SKS)

	fmt.Print("Nilai UTS (0-100): ")
	fmt.Scan(&mk.Nilai.UTS)
	fmt.Print("Nilai UAS (0-100): ")
	fmt.Scan(&mk.Nilai.UAS)
	fmt.Print("Nilai Quiz (0-100): ")
	fmt.Scan(&mk.Nilai.Quiz)

	if mk.Nilai.UTS < 0 || mk.Nilai.UTS > 100 ||
		mk.Nilai.UAS < 0 || mk.Nilai.UAS > 100 ||
		mk.Nilai.Quiz < 0 || mk.Nilai.Quiz > 100 {
		fmt.Println("Nilai tidak valid. Harus antara 0 dan 100.")
		return
	}

	mk.Nilai.Total = hitungTotal(mk.Nilai)
	mk.Nilai.Grade = hitungGrade(mk.Nilai.Total)
	mk.Nilai.IP = hitungIP(mk.Nilai.Grade)

	if _, ok := mahasiswa[mahasiswaIdx].NilaiSemester[semester]; !ok {
		mahasiswa[mahasiswaIdx].NilaiSemester[semester] = SemesterData{
			MataKuliah: make([]MataKuliahData, 0),
		}
	}

	semesterData := mahasiswa[mahasiswaIdx].NilaiSemester[semester]
	semesterData.MataKuliah = append(semesterData.MataKuliah, mk)

	updateNilaiSemester(&semesterData)
	mahasiswa[mahasiswaIdx].NilaiSemester[semester] = semesterData

	updateIPKTotal(mahasiswaIdx)

	fmt.Println("Mata kuliah dan nilai berhasil ditambahkan.")
}

func editNilaiMataKuliah(mahasiswaIdx, semester int) {
	fmt.Print("Masukkan kode mata kuliah: ")
	var kodeMK string
	fmt.Scan(&kodeMK)

	semesterData, ok := mahasiswa[mahasiswaIdx].NilaiSemester[semester]
	if !ok {
		fmt.Println("Belum ada data untuk semester ini.")
		return
	}

	for i, mk := range semesterData.MataKuliah {
		if mk.Kode == kodeMK {
			fmt.Print("Nilai UTS baru (0-100): ")
			fmt.Scan(&semesterData.MataKuliah[i].Nilai.UTS)
			fmt.Print("Nilai UAS baru (0-100): ")
			fmt.Scan(&semesterData.MataKuliah[i].Nilai.UAS)
			fmt.Print("Nilai Quiz baru (0-100): ")
			fmt.Scan(&semesterData.MataKuliah[i].Nilai.Quiz)

			if semesterData.MataKuliah[i].Nilai.UTS < 0 || semesterData.MataKuliah[i].Nilai.UTS > 100 ||
				semesterData.MataKuliah[i].Nilai.UAS < 0 || semesterData.MataKuliah[i].Nilai.UAS > 100 ||
				semesterData.MataKuliah[i].Nilai.Quiz < 0 || semesterData.MataKuliah[i].Nilai.Quiz > 100 {
				fmt.Println("Nilai tidak valid. Harus antara 0 dan 100.")
				return
			}

			semesterData.MataKuliah[i].Nilai.Total = hitungTotal(semesterData.MataKuliah[i].Nilai)
			semesterData.MataKuliah[i].Nilai.Grade = hitungGrade(semesterData.MataKuliah[i].Nilai.Total)
			semesterData.MataKuliah[i].Nilai.IP = hitungIP(semesterData.MataKuliah[i].Nilai.Grade)

			updateNilaiSemester(&semesterData)
			mahasiswa[mahasiswaIdx].NilaiSemester[semester] = semesterData
			updateIPKTotal(mahasiswaIdx)

			fmt.Println("Nilai berhasil diubah.")
			return
		}
	}
	fmt.Println("Mata kuliah tidak ditemukan.")
}

func hapusMataKuliah(mahasiswaIdx, semester int) {
	fmt.Print("Masukkan kode mata kuliah: ")
	var kodeMK string
	fmt.Scan(&kodeMK)

	semesterData, ok := mahasiswa[mahasiswaIdx].NilaiSemester[semester]
	if !ok {
		fmt.Println("Belum ada data untuk semester ini.")
		return
	}

	for i, mk := range semesterData.MataKuliah {
		if mk.Kode == kodeMK {
			semesterData.MataKuliah = append(semesterData.MataKuliah[:i], semesterData.MataKuliah[i+1:]...)
			updateNilaiSemester(&semesterData)
			mahasiswa[mahasiswaIdx].NilaiSemester[semester] = semesterData
			updateIPKTotal(mahasiswaIdx)
			fmt.Println("Mata kuliah berhasil dihapus.")
			return
		}
	}
	fmt.Println("Mata kuliah tidak ditemukan.")
}

func updateNilaiSemester(semester *SemesterData) {
	var totalIP float64
	semester.TotalSKS = 0

	for _, mk := range semester.MataKuliah {
		totalIP += mk.Nilai.IP * float64(mk.SKS)
		semester.TotalSKS += mk.SKS
	}

	if semester.TotalSKS > 0 {
		semester.IPS = totalIP / float64(semester.TotalSKS)
	}
}

func updateIPKTotal(mahasiswaIdx int) {
	var totalIP float64
	var totalSKS int

	for _, semesterData := range mahasiswa[mahasiswaIdx].NilaiSemester {
		if semesterData.TotalSKS > 0 {
			totalIP += semesterData.IPS * float64(semesterData.TotalSKS)
			totalSKS += semesterData.TotalSKS
		}
	}

	if totalSKS > 0 {
		mahasiswa[mahasiswaIdx].IPK = totalIP / float64(totalSKS)
		mahasiswa[mahasiswaIdx].TotalSKS = totalSKS
	}
}

func tampilkanNilaiSemester(mahasiswaIdx, semester int) {
	mhs := mahasiswa[mahasiswaIdx]
	semesterData, ok := mhs.NilaiSemester[semester]
	if !ok {
		fmt.Printf("Belum ada data nilai untuk semester %d\n", semester)
		return
	}

	fmt.Printf("\n=== Nilai Semester %d ===\n", semester)
	fmt.Printf("Nama: %s\n", mhs.Nama)
	fmt.Printf("NIM: %s\n", mhs.NIM)
	fmt.Printf("IPS: %.2f\n", semesterData.IPS)
	fmt.Printf("Total SKS: %d\n\n", semesterData.TotalSKS)

	fmt.Printf("%-6s %-30s %-3s %-7s %-7s %-7s %-7s %-7s %-5s\n",
		"Kode", "Mata Kuliah", "SKS", "UTS", "UAS", "Quiz", "Total", "Grade", "IP")
	fmt.Println(strings.Repeat("-", 85))

	for _, mk := range semesterData.MataKuliah {
		fmt.Printf("%-6s %-30s %-3d %-7.2f %-7.2f %-7.2f %-7.2f %-7s %-5.2f\n",
			mk.Kode, mk.Nama, mk.SKS,
			mk.Nilai.UTS, mk.Nilai.UAS, mk.Nilai.Quiz,
			mk.Nilai.Total, mk.Nilai.Grade, mk.Nilai.IP)
	}
}

func eksporTranskrip() {
	fmt.Print("Masukkan NIM mahasiswa: ")
	var nim string
	fmt.Scan(&nim)
	idx := cariMahasiswaSequential(nim)
	if idx == -1 {
		fmt.Println("Mahasiswa tidak ditemukan.")
		return
	}

	mhs := mahasiswa[idx]

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetName := "Transkrip Nilai"
	f.SetSheetName("Sheet1", sheetName)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#A9A9A9"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	borderStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	f.SetCellValue(sheetName, "A1", "Nama:")
	f.SetCellValue(sheetName, "B1", mhs.Nama)
	f.SetCellValue(sheetName, "A2", "NIM:")
	f.SetCellValue(sheetName, "B2", mhs.NIM)
	f.SetCellValue(sheetName, "A3", "IPK:")
	f.SetCellValue(sheetName, "B3", fmt.Sprintf("%.2f", mhs.IPK))
	f.SetCellValue(sheetName, "A4", "Total SKS:")
	f.SetCellValue(sheetName, "B4", mhs.TotalSKS)

	f.SetColWidth(sheetName, "A", "A", 10)
	f.SetColWidth(sheetName, "B", "B", 30)
	f.SetColWidth(sheetName, "C", "C", 8)
	f.SetColWidth(sheetName, "D", "D", 10)
	f.SetColWidth(sheetName, "E", "E", 10)
	f.SetColWidth(sheetName, "F", "F", 10)
	f.SetColWidth(sheetName, "G", "G", 10)

	currentRow := 6

	numberStyle, _ := f.NewStyle(&excelize.Style{
		NumFmt:    2,
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	for semester := 1; semester <= 8; semester++ {
		if semesterData, ok := mhs.NilaiSemester[semester]; ok && len(semesterData.MataKuliah) > 0 {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow),
				fmt.Sprintf("Semester %d (IPS: %.2f)", semester, semesterData.IPS))
			f.MergeCell(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("G%d", currentRow))
			currentRow += 2

			headers := []string{"Kode", "Mata Kuliah", "SKS", "UTS", "UAS", "Quiz", "Grade"}
			for i, header := range headers {
				cell := fmt.Sprintf("%c%d", 'A'+i, currentRow)
				f.SetCellValue(sheetName, cell, header)
				f.SetCellStyle(sheetName, cell, cell, headerStyle)
			}
			currentRow++

			for _, mk := range semesterData.MataKuliah {
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), mk.Kode)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), mk.Nama)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", currentRow), mk.SKS)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), mk.Nilai.UTS)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", currentRow), mk.Nilai.UAS)
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", currentRow), mk.Nilai.Quiz)
				f.SetCellValue(sheetName, fmt.Sprintf("G%d", currentRow), mk.Nilai.Grade)

				for col := 'D'; col <= 'F'; col++ {
					cell := fmt.Sprintf("%c%d", col, currentRow)
					f.SetCellStyle(sheetName, cell, cell, numberStyle)
				}

				for col := 'A'; col <= 'G'; col++ {
					cell := fmt.Sprintf("%c%d", col, currentRow)
					if col == 'A' || col == 'B' {
						style, _ := f.NewStyle(&excelize.Style{
							Border: []excelize.Border{
								{Type: "left", Color: "000000", Style: 1},
								{Type: "top", Color: "000000", Style: 1},
								{Type: "bottom", Color: "000000", Style: 1},
								{Type: "right", Color: "000000", Style: 1},
							},
							Alignment: &excelize.Alignment{Horizontal: "left"},
						})
						f.SetCellStyle(sheetName, cell, cell, style)
					} else {
						f.SetCellStyle(sheetName, cell, cell, borderStyle)
					}
				}
				currentRow++
			}

			currentRow += 2
		}
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("transkrip_%s_%s.xlsx", mhs.NIM, timestamp)

	if err := f.SaveAs(filename); err != nil {
		fmt.Printf("Error saving excel file: %v\n", err)
		return
	}

	fmt.Printf("Transkrip berhasil diekspor ke file: %s\n", filename)
}

func menuTampilkanData() {
	for {
		fmt.Println("\n=== Tampilkan Data ===")
		fmt.Println("1. Daftar Mahasiswa")
		fmt.Println("2. Detail Mahasiswa")
		fmt.Println("3. Mahasiswa Berprestasi")
		fmt.Println("4. Kembali")
		fmt.Print("Pilihan: ")
		var pilihan int
		fmt.Scan(&pilihan)

		switch pilihan {
		case 1:
			tampilkanDaftarMahasiswa()
		case 2:
			tampilkanDetailMahasiswa()
		case 3:
			tampilkanMahasiswaBerprestasi()
		case 4:
			return
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

func tampilkanDaftarMahasiswa() {
	fmt.Println("\n=== Daftar Mahasiswa ===")
	if jumlahMahasiswa == 0 {
		fmt.Println("Belum ada data mahasiswa.")
		return
	}

	fmt.Printf("%-10s %-30s %-20s %-8s %-6s %-6s\n",
		"NIM", "Nama", "Jurusan", "Angkatan", "SKS", "IPK")
	fmt.Println(strings.Repeat("-", 85))

	for i := 0; i < jumlahMahasiswa; i++ {
		m := mahasiswa[i]
		fmt.Printf("%-10s %-30s %-20s %-8d %-6d %-6.2f\n",
			m.NIM, m.Nama, m.Jurusan, m.Angkatan, m.TotalSKS, m.IPK)
	}
}

func tampilkanDetailMahasiswa() {
	fmt.Print("Masukkan NIM mahasiswa: ")
	var nim string
	fmt.Scan(&nim)
	idx := cariMahasiswaSequential(nim)
	if idx == -1 {
		fmt.Println("Mahasiswa tidak ditemukan.")
		return
	}

	m := mahasiswa[idx]
	fmt.Printf("\n=== Detail Mahasiswa ===\n")
	fmt.Printf("NIM: %s\n", m.NIM)
	fmt.Printf("Nama: %s\n", m.Nama)
	fmt.Printf("Jurusan: %s\n", m.Jurusan)
	fmt.Printf("Angkatan: %d\n", m.Angkatan)
	fmt.Printf("Total SKS: %d\n", m.TotalSKS)
	fmt.Printf("IPK: %.2f\n", m.IPK)

	fmt.Println("\nRiwayat Nilai per Semester:")
	for semester := 1; semester <= 8; semester++ {
		if semesterData, ok := m.NilaiSemester[semester]; ok {
			fmt.Printf("\nSemester %d - IPS: %.2f\n", semester, semesterData.IPS)
		}
	}
}

func tampilkanMahasiswaBerprestasi() {
	fmt.Println("\n=== Mahasiswa Berprestasi ===")
	if jumlahMahasiswa == 0 {
		fmt.Println("Belum ada data mahasiswa.")
		return
	}

	type mahasiswaIPK struct {
		idx int
		ipk float64
	}

	temp := make([]mahasiswaIPK, jumlahMahasiswa)
	for i := 0; i < jumlahMahasiswa; i++ {
		temp[i] = mahasiswaIPK{i, mahasiswa[i].IPK}
	}

	for i := 0; i < len(temp)-1; i++ {
		for j := 0; j < len(temp)-i-1; j++ {
			if temp[j].ipk < temp[j+1].ipk {
				temp[j], temp[j+1] = temp[j+1], temp[j]
			}
		}
	}

	fmt.Printf("%-10s %-30s %-20s %-6s\n",
		"NIM", "Nama", "Jurusan", "IPK")
	fmt.Println(strings.Repeat("-", 70))

	limit := 10
	if jumlahMahasiswa < limit {
		limit = jumlahMahasiswa
	}

	for i := 0; i < limit; i++ {
		m := mahasiswa[temp[i].idx]
		fmt.Printf("%-10s %-30s %-20s %-6.2f\n",
			m.NIM, m.Nama, m.Jurusan, m.IPK)
	}
}

func hitungTotal(nilai NilaiKomponen) float64 {
	return (nilai.UTS * 0.35) + (nilai.UAS * 0.35) + (nilai.Quiz * 0.30)
}

func hitungGrade(nilai float64) string {
	switch {
	case nilai >= 85:
		return "A"
	case nilai >= 80:
		return "AB"
	case nilai >= 75:
		return "B"
	case nilai >= 70:
		return "BC"
	case nilai >= 65:
		return "C"
	case nilai >= 55:
		return "D"
	default:
		return "E"
	}
}

func hitungIP(grade string) float64 {
	switch grade {
	case "A":
		return 4.0
	case "AB":
		return 3.5
	case "B":
		return 3.0
	case "BC":
		return 2.5
	case "C":
		return 2.0
	case "D":
		return 1.0
	default:
		return 0.0
	}
}

func prediksiKinerjaAkademik() {
	fmt.Print("Masukkan NIM mahasiswa: ")
	var nim string
	fmt.Scan(&nim)
	idx := cariMahasiswaSequential(nim)
	if idx == -1 {
		fmt.Println("Mahasiswa tidak ditemukan.")
		return
	}

	var (
		x []float64
		y []float64
	)

	mhs := mahasiswa[idx]
	for semester := 1; semester <= 8; semester++ {
		if data, ok := mhs.NilaiSemester[semester]; ok && data.TotalSKS > 0 {
			x = append(x, float64(semester))
			y = append(y, data.IPS)
		}
	}

	if len(x) < 2 {
		fmt.Println("Data tidak cukup untuk melakukan prediksi (minimal 2 semester).")
		return
	}

	result := hitungRegresiLinear(x, y)
	nextSemester := float64(len(x) + 1)

	rawPrediction := result.Slope*nextSemester + result.Intercept

	result.Prediction = math.Max(0.0, math.Min(4.0, rawPrediction))

	tampilkanHasilPrediksi(mhs, x, y, result, rawPrediction)
}

func hitungRegresiLinear(x, y []float64) RegressionResult {
	n := float64(len(x))

	var sumX, sumY, sumXY, sumX2 float64
	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n

	var ssTot, ssRes float64
	yMean := sumY / n

	for i := 0; i < len(y); i++ {
		predicted := slope*x[i] + intercept
		ssTot += math.Pow(y[i]-yMean, 2)
		ssRes += math.Pow(y[i]-predicted, 2)
	}

	r2 := 1 - (ssRes / ssTot)

	return RegressionResult{
		Slope:     slope,
		Intercept: intercept,
		R2:        r2,
	}
}

func tampilkanHasilPrediksi(mhs Mahasiswa, x, y []float64, result RegressionResult, rawPrediction float64) {
	fmt.Printf("\n=== Hasil Prediksi Kinerja Akademik ===\n")
	fmt.Printf("Nama: %s\n", mhs.Nama)
	fmt.Printf("NIM: %s\n", mhs.NIM)
	fmt.Printf("Jurusan: %s\n", mhs.Jurusan)
	fmt.Printf("\nData Historis IPS:\n")
	fmt.Printf("%-10s %-10s\n", "Semester", "IPS")
	fmt.Println(strings.Repeat("-", 20))

	for i := 0; i < len(x); i++ {
		fmt.Printf("%-10.0f %-10.2f\n", x[i], y[i])
	}

	nextSemester := len(x) + 1

	fmt.Printf("\nPrediksi IPS untuk Semester %d: %.2f", nextSemester, result.Prediction)

}
