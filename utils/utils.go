package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/phpdave11/gofpdf"
)

// SaveScanResultToFile saves the raw Trivy JSON output to a timestamped file.
func SaveScanResultToFile(image string, data []byte) (string, error) {
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		fmt.Println("❌ Failed to create uploads dir:", err)
		return "", err
	}

	// Replace unsupported filename characters
	safeName := strings.ReplaceAll(image, ":", "_")
	safeName = strings.ReplaceAll(safeName, "/", "_")
	timestamp := time.Now().Format("20060102_150405")
	file := filepath.Join("uploads", fmt.Sprintf("%s_%s.json", safeName, timestamp))

	if len(data) == 0 {
		fmt.Println("❌ Refused to write: empty data buffer")
		return "", fmt.Errorf("empty scan data")
	}

	fmt.Printf("💾 Writing %d bytes to %s\n", len(data), file)
	err = os.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Println("❌ Write failed:", err)
		return "", err
	}

	fmt.Println("✅ File written successfully.")
	return file, nil
}

// GeneratePDFReport takes Trivy's JSON output and saves it as a pretty-printed PDF.
func GeneratePDFReport(jsonData []byte, outputPath string) error {
	var prettyJSON map[string]interface{}
	err := json.Unmarshal(jsonData, &prettyJSON)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Courier", "", 10) // Courier works better for structured data

	jsonBytes, _ := json.MarshalIndent(prettyJSON, "", "  ")
	lines := strings.Split(string(jsonBytes), "\n")

	for _, line := range lines {
		pdf.MultiCell(0, 5, line, "", "L", false)
	}

	err = pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return fmt.Errorf("failed to save PDF: %w", err)
	}

	fmt.Println("📄 PDF report saved to:", outputPath)
	return nil
}
