package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"docker-image-scanner/middleware"
	"docker-image-scanner/utils"
)

type ScanRequest struct {
	Image string `json:"image"`
}

type ScanResponse struct {
	Image   string      `json:"image"`
	Message string      `json:"message"`
	File    string      `json:"file"`
	Report  interface{} `json:"report"`
}

func ScanDockerImageHandler(w http.ResponseWriter, r *http.Request) {
	// ✅ Extract JWT claims from context
	claims, ok := r.Context().Value(middleware.UserKey).(*middleware.Claims)
	if !ok || claims == nil {
		http.Error(w, "Unauthorized: no valid token found", http.StatusUnauthorized)
		return
	}
	email := claims.Email

	// ✅ Parse scan request
	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// ✅ Run Trivy scan
	cmd := exec.Command("trivy", "image", "-q", "-f", "json", req.Image)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("Trivy scan failed: %s", stderr.String()), http.StatusInternalServerError)
		return
	}

	// ✅ Save to file
	filePath, err := utils.SaveScanResultToFile(req.Image, stdout.Bytes())
	if err != nil {
		http.Error(w, "Failed to save scan result", http.StatusInternalServerError)
		return
	}

	// ✅ Send Email in Background
	go utils.SendEmailNotification(email, req.Image, filePath)

	// ✅ Respond to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ScanResponse{
		Image:   req.Image,
		Message: "Scan successful",
		File:    filePath,
		Report:  json.RawMessage(stdout.Bytes()),
	})
}
