package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

// RunTrivyScan executes a Trivy scan on the given Docker image
// and returns the JSON output as bytes.
func RunTrivyScan(image string) ([]byte, error) {
	cmd := exec.Command("trivy", "image", "--quiet", "--format", "json", image)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Trivy scan failed:", err)
		fmt.Println("🔴 STDERR:", stderr.String())
		return nil, fmt.Errorf("trivy error: %s", stderr.String())
	}

	return stdout.Bytes(), nil
}
