package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// SendWhatsAppOTP mengirim pesan teks via API gowhatsapp
func SendWhatsAppOTP(phone string, otp string) error {
	formattedPhone := formatPhoneNumber(phone)
	apiURL := os.Getenv("GOWA_URL")

	payload := map[string]any{
		"phone":        fmt.Sprintf("%s@s.whatsapp.net", formattedPhone),
		"message":      fmt.Sprintf("Halo! Kode OTP Anda untuk KUA Kecamatan Ciemas adalah: *%s*.\n\nMohon jangan berikan kode ini kepada siapapun.", otp),
		"is_forwarded": false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gagal parse payload: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("gagal buat request HTTP: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Devic-Id", os.Getenv("GOWA_DEVICE_ID"))

	req.SetBasicAuth("kuaci_emas", "kuaci_emas")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("gagal nge-hit server wa: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server wa merespon dengan status error: %d", resp.StatusCode)
	}

	return nil
}

func formatPhoneNumber(phone string) string {
	phone = strings.TrimSpace(phone)

	if after, found := strings.CutPrefix(phone, "62"); found {
		return "+62" + after
	}

	if after, found := strings.CutPrefix(phone, "0"); found {
		return "+62" + after
	}

	return phone
}
