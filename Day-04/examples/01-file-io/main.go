package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	// ================================
	// 1. WRITING A FILE
	// ================================

	fmt.Println("=== Writing a File ===")

	content := []byte("Hello from Go!\nThis is line 2\nThis is line 3\n")
	err := os.WriteFile("output.txt", content, 0644)
	if err != nil {
		log.Fatal("Cannot write file:", err)
	}
	fmt.Println("File written: output.txt")

	// ================================
	// 2. READING ENTIRE FILE
	// ================================

	fmt.Println("\n=== Reading Entire File ===")

	data, err := os.ReadFile("output.txt")
	if err != nil {
		log.Fatal("Cannot read file:", err)
	}
	fmt.Println("Content:")
	fmt.Println(string(data))

	// ================================
	// 3. APPENDING TO A FILE
	// ================================

	fmt.Println("=== Appending to a File ===")

	// O_APPEND = add to end, O_CREATE = create if missing, O_WRONLY = write only
	file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Cannot open file:", err)
	}

	_, err = file.WriteString("This line was appended!\n")
	if err != nil {
		file.Close()
		log.Fatal("Cannot write:", err)
	}
	file.Close()

	// Read again to verify
	data, _ = os.ReadFile("output.txt")
	fmt.Println("After append:")
	fmt.Println(string(data))

	// ================================
	// 4. READING LINE BY LINE (bufio)
	// ================================

	fmt.Println("=== Reading Line by Line ===")

	file, err = os.Open("output.txt")
	if err != nil {
		log.Fatal("Cannot open:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		fmt.Printf("  Line %d: %s\n", lineNum, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading:", err)
	}

	fmt.Printf("Total: %d lines\n", lineNum)

	// ================================
	// 5. CHECKING IF FILE EXISTS
	// ================================

	fmt.Println("\n=== Checking File Existence ===")

	files := []string{"output.txt", "missing.txt", ".."}

	for _, name := range files {
		info, err := os.Stat(name)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("  ❌ %s - does not exist\n", name)
			} else {
				fmt.Printf("  ⚠️  %s - error: %s\n", name, err)
			}
		} else {
			if info.IsDir() {
				fmt.Printf("  📁 %s - directory\n", name)
			} else {
				fmt.Printf("  ✅ %s - %d bytes\n", name, info.Size())
			}
		}
	}

	// ================================
	// 6. CREATING A FILE WITH os.Create
	// ================================

	fmt.Println("\n=== Creating File with os.Create ===")

	newFile, err := os.Create("server.log")
	if err != nil {
		log.Fatal("Cannot create:", err)
	}
	defer newFile.Close()

	// Write some log lines
	logLines := []string{
		"2024-01-15 10:00:00 INFO Server started",
		"2024-01-15 10:00:01 INFO Listening on port 8080",
		"2024-01-15 10:00:05 WARN High memory usage: 78%",
		"2024-01-15 10:00:10 ERROR Connection to db-02 failed",
		"2024-01-15 10:00:15 INFO Failover to db-01 complete",
	}

	writer := bufio.NewWriter(newFile)
	for _, line := range logLines {
		writer.WriteString(line + "\n")
	}
	writer.Flush() // Don't forget to flush buffered writer!

	fmt.Println("Created server.log with", len(logLines), "lines")

	// Cleanup - remove the files we created
	os.Remove("output.txt")
	os.Remove("server.log")
	fmt.Println("\nCleaned up temporary files")
}
