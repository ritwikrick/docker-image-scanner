package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

// SendEmailNotification sends an email with the scan report attached as PDF.
func SendEmailNotification(to, image, filePath string) {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	if from == "" || password == "" {
		fmt.Println("❌ Missing SMTP credentials")
		return
	}

	// Generate PDF path from the JSON path
	pdfPath := strings.TrimSuffix(filePath, ".json") + ".pdf"

	// Generate PDF from scan result
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("❌ Failed to read JSON file:", err)
		return
	}
	if err := GeneratePDFReport(jsonData, pdfPath); err != nil {
		fmt.Println("❌ Failed to generate PDF:", err)
		return
	}

	// Read PDF file to attach
	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		fmt.Println("❌ Failed to read PDF file:", err)
		return
	}

	// Build MIME email body with attachment
	var body bytes.Buffer
	boundary := "my-boundary-unique"

	body.WriteString("From: " + from + "\r\n")
	body.WriteString("To: " + to + "\r\n")
	body.WriteString("Subject: Docker Scan Completed\r\n")
	body.WriteString("MIME-Version: 1.0\r\n")
	body.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
	body.WriteString("\r\n--" + boundary + "\r\n")
	body.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	body.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	body.WriteString(fmt.Sprintf("Scan completed for image: %s\nReport attached as PDF.\n", image))
	body.WriteString("\r\n--" + boundary + "\r\n")
	body.WriteString("Content-Type: application/pdf\r\n")
	body.WriteString("Content-Disposition: attachment; filename=\"" + filepath.Base(pdfPath) + "\"\r\n")
	body.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")

	// Base64 encode PDF
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(pdfBytes)))
	base64.StdEncoding.Encode(encoded, pdfBytes)
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		body.Write(encoded[i:end])
		body.WriteString("\r\n")
	}
	body.WriteString("--" + boundary + "--")

	// Send the email
	err = smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from,
		[]string{to},
		body.Bytes(),
	)

	if err != nil {
		fmt.Println("❌ Email failed:", err)
		return
	}

	fmt.Println("✅ Email sent to", to, "with PDF attachment")
}
