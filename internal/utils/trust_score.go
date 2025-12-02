package utils

import (
	"net"
	"situs-keagamaan/internal/dto"
	"strings"
)

const (
	WeightUserAgent   float64 = 30.0
	WeightIPAddress   float64 = 40.0
	WeightGeoLocation float64 = 20.0
	WeightDeviceID    float64 = 10.0
)

func getIPDifferenceLevel(currentIP, trustedIP string) string {
	// Jika IP sama, tidak ada pengurangan.
	if currentIP == trustedIP {
		return "None"
	}

	// Menguraikan IP
	ip1 := net.ParseIP(currentIP)
	ip2 := net.ParseIP(trustedIP)

	// Jika salah satu tidak valid, atau formatnya berbeda (misalnya v4 vs v6), anggap total change
	if ip1 == nil || ip2 == nil || len(ip1) != len(ip2) {
		return "Total"
	}

	// Membandingkan oktet/byte (kita hanya akan bandingkan 3 oktet pertama untuk menentukan Minor/Major)
	// Untuk IPv4: [oktet1, oktet2, oktet3, oktet4]
	if len(ip1.To4()) == 4 {
		// Bandingkan 3 oktet pertama (misal: 192.168.1)
		if ip1.To4()[0] == ip2.To4()[0] && ip1.To4()[1] == ip2.To4()[1] && ip1.To4()[2] == ip2.To4()[2] {
			// Hanya oktet terakhir yang berbeda (Minor Change)
			return "Minor"
		}

		// Bandingkan 2 oktet pertama (misal: 192.168)
		if ip1.To4()[0] == ip2.To4()[0] && ip1.To4()[1] == ip2.To4()[1] {
			// Oktet ke-3/ke-4 berbeda (Major Change)
			return "Major"
		}
	}

	// Selain kondisi di atas, dianggap perubahan total (Total Change)
	return "Total"
}

func CalculateTrustScore(current *dto.SessionRequest, trusted *dto.SessionValue) float64 {
	score := 100.0
	ipDiff := getIPDifferenceLevel(current.IPAddress, trusted.IPAddress)

	ipDeduction := 0.0
	switch ipDiff {
	case "Minor":
		ipDeduction = WeightIPAddress * 0.25
	case "Major":
		ipDeduction = WeightIPAddress * 0.60
	case "Total":
		ipDeduction = WeightIPAddress * 1.00
	}
	score -= ipDeduction

	if !strings.EqualFold(current.UserAgent, trusted.UserAgent) {
		score -= 30.0
	}

	if current.DeviceID != trusted.DeviceID {
		score -= 10.0
	}

	if current.GeoLocation != trusted.GeoLocation {
		score -= 20.0
	}

	if score < 0.0 {
		return 0.0
	}

	return score
}
