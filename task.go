package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Task merepresentasikan struktur data tugas mahasiswa
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Semester    string    `json:"semester"`
	DosenName   string    `json:"dosen_name"`
	DosenCode   string    `json:"dosen_code"`
	Status      string    `json:"status"` // "belum", "sudah", "terlambat"
	CreatedAt   time.Time `json:"created_at"`
}

// TaskManager mengelola daftar tugas
type TaskManager struct {
	Tasks    []Task
	FilePath string
}

// NewTaskManager membuat instance baru TaskManager
func NewTaskManager(filePath string) *TaskManager {
	return &TaskManager{
		Tasks:    []Task{},
		FilePath: filePath,
	}
}

// LoadTasks memuat data tugas dari file JSON
func (tm *TaskManager) LoadTasks() error {
	file, err := os.ReadFile(tm.FilePath)
	if err != nil {
		// Jika file tidak ada, buat file baru
		if os.IsNotExist(err) {
			return tm.SaveTasks()
		}
		return err
	}

	err = json.Unmarshal(file, &tm.Tasks)
	if err != nil {
		return err
	}

	// Update status terlambat secara otomatis
	tm.UpdateLateStatus()
	return nil
}

// SaveTasks menyimpan data tugas ke file JSON
func (tm *TaskManager) SaveTasks() error {
	data, err := json.MarshalIndent(tm.Tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tm.FilePath, data, 0644)
}

// AddTask menambahkan tugas baru
func (tm *TaskManager) AddTask(task Task) {
	// Generate ID baru
	task.ID = tm.getNextID()
	task.CreatedAt = time.Now()
	task.Status = "belum"
	
	tm.Tasks = append(tm.Tasks, task)
	tm.SaveTasks()
}

// getNextID menghasilkan ID unik untuk tugas baru
func (tm *TaskManager) getNextID() int {
	maxID := 0
	for _, task := range tm.Tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

// ViewTasks menampilkan semua tugas
func (tm *TaskManager) ViewTasks() {
	if len(tm.Tasks) == 0 {
		fmt.Println("\n❌ Belum ada tugas yang tercatat.")
		return
	}

	fmt.Println("\n📚 DAFTAR TUGAS")
	fmt.Println(strings.Repeat("=", 80))

	for _, task := range tm.Tasks {
		tm.displayTask(task)
	}
}

// displayTask menampilkan detail satu tugas
func (tm *TaskManager) displayTask(task Task) {
	statusIcon := tm.getStatusIcon(task.Status)
	timeLeft := tm.calculateTimeRemaining(task.Deadline)

	fmt.Printf("\n[ID: %d] %s %s\n", task.ID, statusIcon, task.Title)
	fmt.Printf("📝 Deskripsi: %s\n", task.Description)
	fmt.Printf("📅 Deadline: %s\n", task.Deadline.Format("02 Jan 2006 15:04"))
	fmt.Printf("⏰ Sisa Waktu: %s\n", timeLeft)
	fmt.Printf("📚 Semester: %s\n", task.Semester)
	fmt.Printf("👨‍🏫 Dosen: %s (%s)\n", task.DosenName, task.DosenCode)
	fmt.Printf("📊 Status: %s\n", task.Status)
	fmt.Println(strings.Repeat("-", 80))
}

// getStatusIcon mengembalikan icon berdasarkan status
func (tm *TaskManager) getStatusIcon(status string) string {
	switch status {
	case "sudah":
		return "✅"
	case "terlambat":
		return "🚨"
	default:
		return "⏳"
	}
}

// calculateTimeRemaining menghitung sisa waktu hingga deadline
func (tm *TaskManager) calculateTimeRemaining(deadline time.Time) string {
	now := time.Now()
	duration := deadline.Sub(now)

	if duration < 0 {
		return "🚨 TERLAMBAT"
	}

	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%d hari, %d jam, %d menit", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%d jam, %d menit", hours, minutes)
	} else {
		return fmt.Sprintf("%d menit", minutes)
	}
}

// UpdateLateStatus memperbarui status tugas yang sudah melewati deadline
func (tm *TaskManager) UpdateLateStatus() {
	now := time.Now()
	updated := false

	for i := range tm.Tasks {
		if tm.Tasks[i].Status == "belum" && tm.Tasks[i].Deadline.Before(now) {
			tm.Tasks[i].Status = "terlambat"
			updated = true
		}
	}

	if updated {
		tm.SaveTasks()
	}
}

// EditTask mengubah detail tugas berdasarkan ID
func (tm *TaskManager) EditTask(id int, updatedTask Task) error {
	for i := range tm.Tasks {
		if tm.Tasks[i].ID == id {
			// Pertahankan ID, CreatedAt, dan Status asli
			updatedTask.ID = tm.Tasks[i].ID
			updatedTask.CreatedAt = tm.Tasks[i].CreatedAt
			updatedTask.Status = tm.Tasks[i].Status
			
			tm.Tasks[i] = updatedTask
			return tm.SaveTasks()
		}
	}
	return fmt.Errorf("tugas dengan ID %d tidak ditemukan", id)
}

// DeleteTask menghapus tugas berdasarkan ID
func (tm *TaskManager) DeleteTask(id int) error {
	for i, task := range tm.Tasks {
		if task.ID == id {
			tm.Tasks = append(tm.Tasks[:i], tm.Tasks[i+1:]...)
			return tm.SaveTasks()
		}
	}
	return fmt.Errorf("tugas dengan ID %d tidak ditemukan", id)
}

// MarkTaskComplete menandai tugas sebagai selesai
func (tm *TaskManager) MarkTaskComplete(id int) error {
	for i := range tm.Tasks {
		if tm.Tasks[i].ID == id {
			if tm.Tasks[i].Status == "terlambat" {
				return fmt.Errorf("tugas yang sudah terlambat tidak bisa diubah statusnya")
			}
			tm.Tasks[i].Status = "sudah"
			return tm.SaveTasks()
		}
	}
	return fmt.Errorf("tugas dengan ID %d tidak ditemukan", id)
}

// FilterBySemester menampilkan tugas berdasarkan semester
func (tm *TaskManager) FilterBySemester(semester string) {
	fmt.Printf("\n📚 TUGAS SEMESTER: %s\n", semester)
	fmt.Println(strings.Repeat("=", 80))

	found := false
	for _, task := range tm.Tasks {
		if strings.Contains(strings.ToLower(task.Semester), strings.ToLower(semester)) {
			tm.displayTask(task)
			found = true
		}
	}

	if !found {
		fmt.Printf("\n❌ Tidak ada tugas untuk semester %s\n", semester)
	}
}

// FilterByDosen menampilkan tugas berdasarkan dosen
func (tm *TaskManager) FilterByDosen(dosenName string) {
	fmt.Printf("\n👨‍🏫 TUGAS DARI DOSEN: %s\n", dosenName)
	fmt.Println(strings.Repeat("=", 80))

	found := false
	for _, task := range tm.Tasks {
		if strings.Contains(strings.ToLower(task.DosenName), strings.ToLower(dosenName)) {
			tm.displayTask(task)
			found = true
		}
	}

	if !found {
		fmt.Printf("\n❌ Tidak ada tugas dari dosen %s\n", dosenName)
	}
}

// FilterByStatus menampilkan tugas berdasarkan status
func (tm *TaskManager) FilterByStatus(status string) {
	fmt.Printf("\n📊 TUGAS DENGAN STATUS: %s\n", status)
	fmt.Println(strings.Repeat("=", 80))

	found := false
	for _, task := range tm.Tasks {
		if strings.ToLower(task.Status) == strings.ToLower(status) {
			tm.displayTask(task)
			found = true
		}
	}

	if !found {
		fmt.Printf("\n❌ Tidak ada tugas dengan status %s\n", status)
	}
}

// GetTaskByID mencari tugas berdasarkan ID
func (tm *TaskManager) GetTaskByID(id int) (*Task, error) {
	for i := range tm.Tasks {
		if tm.Tasks[i].ID == id {
			return &tm.Tasks[i], nil
		}
	}
	return nil, fmt.Errorf("tugas dengan ID %d tidak ditemukan", id)
}
