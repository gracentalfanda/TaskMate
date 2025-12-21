package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	dataFile = "tasks.json"
)

var (
	taskManager *TaskManager
	scanner     *bufio.Scanner
)

func main() {
	// Inisialisasi
	taskManager = NewTaskManager(dataFile)
	scanner = bufio.NewScanner(os.Stdin)

	// Load data dari file
	err := taskManager.LoadTasks()
	if err != nil {
		fmt.Printf("⚠️  Error loading tasks: %v\n", err)
	}

	// Tampilkan header
	displayHeader()

	// Loop menu utama
	for {
		displayMenu()
		choice := getInput("\n🔹 Pilih menu (1-9): ")

		switch choice {
		case "1":
			addTaskMenu()
		case "2":
			taskManager.ViewTasks()
		case "3":
			editTaskMenu()
		case "4":
			deleteTaskMenu()
		case "5":
			markCompleteMenu()
		case "6":
			filterBySemesterMenu()
		case "7":
			filterByDosenMenu()
		case "8":
			filterByStatusMenu()
		case "9":
			exitProgram()
			return
		default:
			fmt.Println("\n❌ Pilihan tidak valid! Silakan pilih 1-9.")
		}

		// Tunggu user sebelum lanjut
		fmt.Print("\n🔄 Tekan ENTER untuk melanjutkan...")
		scanner.Scan()
	}
}

// displayHeader menampilkan header aplikasi
func displayHeader() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 TASKMATE - Sistem Pengelola Tugas Perkuliahan Mahasiswa")
	fmt.Println(strings.Repeat("=", 80))
}

// displayMenu menampilkan menu utama
func displayMenu() {
	fmt.Println("\n┌─────────────────────────────────────────┐")
	fmt.Println("│          MENU UTAMA TASKMATE            │")
	fmt.Println("├─────────────────────────────────────────┤")
	fmt.Println("│ 1. ➕ Tambah Tugas                      │")
	fmt.Println("│ 2. 📋 Lihat Semua Tugas                 │")
	fmt.Println("│ 3. ✏️  Edit Tugas                       │")
	fmt.Println("│ 4. 🗑️  Hapus Tugas                      │")
	fmt.Println("│ 5. ✅ Tandai Tugas Selesai              │")
	fmt.Println("│ 6. 🔍 Filter Berdasarkan Semester       │")
	fmt.Println("│ 7. 🔍 Filter Berdasarkan Dosen          │")
	fmt.Println("│ 8. 🔍 Filter Berdasarkan Status         │")
	fmt.Println("│ 9. 🚪 Keluar                            │")
	fmt.Println("└─────────────────────────────────────────┘")
}

// getInput membaca input dari user
func getInput(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// addTaskMenu menu untuk menambah tugas baru
func addTaskMenu() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("➕ TAMBAH TUGAS BARU")
	fmt.Println(strings.Repeat("=", 80))

	var task Task

	// Input judul
	task.Title = getInput("📝 Judul Tugas: ")
	if task.Title == "" {
		fmt.Println("❌ Judul tidak boleh kosong!")
		return
	}

	// Input deskripsi
	task.Description = getInput("📄 Deskripsi: ")

	// Input deadline
	deadlineStr := getInput("📅 Deadline (format: DD/MM/YYYY HH:MM, contoh: 25/12/2024 23:59): ")
	deadline, err := time.Parse("02/01/2006 15:04", deadlineStr)
	if err != nil {
		fmt.Printf("❌ Format deadline salah! Error: %v\n", err)
		return
	}
	task.Deadline = deadline

	// Input semester
	task.Semester = getInput("📚 Semester (contoh: Semester 5 atau Semester 3 Ulang): ")
	if task.Semester == "" {
		fmt.Println("❌ Semester tidak boleh kosong!")
		return
	}

	// Input dosen
	task.DosenName = getInput("👨‍🏫 Nama Dosen: ")
	task.DosenCode = getInput("📖 Kode Mata Kuliah: ")

	// Simpan tugas
	taskManager.AddTask(task)
	fmt.Println("\n✅ Tugas berhasil ditambahkan!")
}

// editTaskMenu menu untuk mengedit tugas
func editTaskMenu() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✏️  EDIT TUGAS")
	fmt.Println(strings.Repeat("=", 80))

	// Tampilkan daftar tugas
	taskManager.ViewTasks()

	// Input ID tugas yang akan diedit
	idStr := getInput("\n🔢 Masukkan ID tugas yang akan diedit (0 untuk batal): ")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		fmt.Println("❌ Batal edit tugas.")
		return
	}

	// Cari tugas
	oldTask, err := taskManager.GetTaskByID(id)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	fmt.Printf("\n📝 Edit Tugas: %s\n", oldTask.Title)
	fmt.Println("💡 Tekan ENTER untuk mempertahankan nilai lama\n")

	var newTask Task

	// Input judul baru
	title := getInput(fmt.Sprintf("📝 Judul [%s]: ", oldTask.Title))
	if title == "" {
		newTask.Title = oldTask.Title
	} else {
		newTask.Title = title
	}

	// Input deskripsi baru
	desc := getInput(fmt.Sprintf("📄 Deskripsi [%s]: ", oldTask.Description))
	if desc == "" {
		newTask.Description = oldTask.Description
	} else {
		newTask.Description = desc
	}

	// Input deadline baru
	deadlineStr := getInput(fmt.Sprintf("📅 Deadline [%s]: ", oldTask.Deadline.Format("02/01/2006 15:04")))
	if deadlineStr == "" {
		newTask.Deadline = oldTask.Deadline
	} else {
		deadline, err := time.Parse("02/01/2006 15:04", deadlineStr)
		if err != nil {
			fmt.Printf("❌ Format deadline salah! Menggunakan deadline lama.\n")
			newTask.Deadline = oldTask.Deadline
		} else {
			newTask.Deadline = deadline
		}
	}

	// Input semester baru
	semester := getInput(fmt.Sprintf("📚 Semester [%s]: ", oldTask.Semester))
	if semester == "" {
		newTask.Semester = oldTask.Semester
	} else {
		newTask.Semester = semester
	}

	// Input dosen baru
	dosenName := getInput(fmt.Sprintf("👨‍🏫 Nama Dosen [%s]: ", oldTask.DosenName))
	if dosenName == "" {
		newTask.DosenName = oldTask.DosenName
	} else {
		newTask.DosenName = dosenName
	}

	dosenCode := getInput(fmt.Sprintf("📖 Kode Mata Kuliah [%s]: ", oldTask.DosenCode))
	if dosenCode == "" {
		newTask.DosenCode = oldTask.DosenCode
	} else {
		newTask.DosenCode = dosenCode
	}

	// Update tugas
	err = taskManager.EditTask(id, newTask)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Println("\n✅ Tugas berhasil diupdate!")
}

// deleteTaskMenu menu untuk menghapus tugas
func deleteTaskMenu() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🗑️  HAPUS TUGAS")
	fmt.Println(strings.Repeat("=", 80))

	// Tampilkan daftar tugas
	taskManager.ViewTasks()

	// Input ID tugas yang akan dihapus
	idStr := getInput("\n🔢 Masukkan ID tugas yang akan dihapus (0 untuk batal): ")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		fmt.Println("❌ Batal hapus tugas.")
		return
	}

	// Konfirmasi
	confirm := getInput(fmt.Sprintf("⚠️  Yakin ingin menghapus tugas ID %d? (y/n): ", id))
	if strings.ToLower(confirm) != "y" {
		fmt.Println("❌ Batal hapus tugas.")
		return
	}

	// Hapus tugas
	err = taskManager.DeleteTask(id)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	fmt.Println("\n✅ Tugas berhasil dihapus!")
}

// markCompleteMenu menu untuk menandai tugas selesai
func markCompleteMenu() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✅ TANDAI TUGAS SELESAI")
	fmt.Println(strings.Repeat("=", 80))

	// Tampilkan daftar tugas
	taskManager.ViewTasks()

	// Input ID tugas yang akan ditandai selesai
	idStr := getInput("\n🔢 Masukkan ID tugas yang sudah selesai (0 untuk batal): ")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		fmt.Println("❌ Batal tandai tugas.")
		return
	}

	// Tandai selesai
	err = taskManager.MarkTaskComplete(id)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	fmt.Println("\n✅ Tugas berhasil ditandai selesai!")
}

// filterBySemesterMenu menu untuk filter berdasarkan semester
func filterBySemesterMenu() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 FILTER BERDASARKAN SEMESTER")
	fmt.Println(strings.Repeat("=", 80))

	semester := getInput("📚 Masukkan semester (contoh: 5): ")
	if semester == "" {
		fmt.Println("❌ Semester tidak boleh kosong!")
		return
	}

	taskManager.FilterBySemester(semester)
}

// filterByDosenMenu menu untuk filter berdasarkan dosen
func filterByDosenMenu() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 FILTER BERDASARKAN DOSEN")
	fmt.Println(strings.Repeat("=", 80))

	dosen := getInput("👨‍🏫 Masukkan nama dosen: ")
	if dosen == "" {
		fmt.Println("❌ Nama dosen tidak boleh kosong!")
		return
	}

	taskManager.FilterByDosen(dosen)
}

// filterByStatusMenu menu untuk filter berdasarkan status
func filterByStatusMenu() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 FILTER BERDASARKAN STATUS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\nStatus yang tersedia:")
	fmt.Println("1. belum")
	fmt.Println("2. sudah")
	fmt.Println("3. terlambat")

	status := getInput("\n📊 Masukkan status: ")
	if status == "" {
		fmt.Println("❌ Status tidak boleh kosong!")
		return
	}

	taskManager.FilterByStatus(status)
}

// exitProgram keluar dari program
func exitProgram() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("👋 Terima kasih telah menggunakan TaskMate!")
	fmt.Println("📚 Semoga tugas-tugasmu terkelola dengan baik!")
	fmt.Println(strings.Repeat("=", 80))
}
