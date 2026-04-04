package utils

import (
	"encoding/json"
	"fmt"
	"sicemas/internal/dto"
	"sort"
	"strings"

	"github.com/xuri/excelize/v2"
)

func flattenJSON(prefix string, v any, result map[string]string) {
	switch val := v.(type) {
	case map[string]any:
		for k, child := range val {
			newKey := k
			if prefix != "" {
				newKey = prefix + " - " + k
			}
			flattenJSON(newKey, child, result)
		}
	case []any:
		var strVals []string
		for _, item := range val {
			if obj, isMap := item.(map[string]any); isMap {
				var objStr []string
				for k, childVal := range obj {
					objStr = append(objStr, fmt.Sprintf("%s: %v", formatColumnName(k), childVal))
				}
				strVals = append(strVals, "["+strings.Join(objStr, " | ")+"]")
			} else {
				strVals = append(strVals, fmt.Sprintf("%v", item))
			}
		}
		result[prefix] = strings.Join(strVals, ", ")

	default:
		if val != nil {
			result[prefix] = fmt.Sprintf("%v", val)
		} else {
			result[prefix] = ""
		}
	}
}
func formatColumnName(s string) string {
	words := strings.Split(strings.ReplaceAll(s, "_", " "), " ")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func ExportSitusToExcel(dataList []dto.SitusKeagamaanDetailResponse) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Data Situs"
	f.SetSheetName("Sheet1", sheetName)

	baseHeaders := []string{
		"ID", "Status Verifikasi", "Jenis Situs", "Nama", "Jenis Tipologi",
		"Nomor Telepon", "No Telp Pengurus", "Email", "Website", "No Badan Hukum",
		"Tahun Berdiri", "Provinsi", "Kabupaten/Kota", "Kecamatan", "Desa",
		"Alamat Lengkap", "Latitude", "Longitude", "Luas Tanah (m2)",
		"Luas Bangunan (m2)", "Status Tanah", "Nomor AIW", "Sertifikat Wakaf",
		"Daya Tampung",
	}

	dynamicHeadersSet := make(map[string]bool)
	var dynamicHeaders []string

	allRowDetails := make([]map[string]string, len(dataList))

	for i, data := range dataList {
		flatDetail := make(map[string]string)
		if data.Detail != nil {
			var m map[string]any
			if err := json.Unmarshal(*data.Detail, &m); err == nil {
				flattenJSON("", m, flatDetail)
			}
		}
		allRowDetails[i] = flatDetail

		for k := range flatDetail {
			if !dynamicHeadersSet[k] {
				dynamicHeadersSet[k] = true
				dynamicHeaders = append(dynamicHeaders, k)
			}
		}
	}

	sort.Strings(dynamicHeaders)

	allHeaders := append(baseHeaders, dynamicHeaders...)

	for i, header := range allHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)

		if i >= len(baseHeaders) {
			header = formatColumnName(header)
		}
		f.SetCellValue(sheetName, cell, header)
	}

	style, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	f.SetRowStyle(sheetName, 1, 1, style)

	for i, data := range dataList {
		row := i + 2

		noPengurusStr := strings.Join(data.NomorTelponPengurus, ", ")
		var situsIDStr string
		if data.SitusID != nil {
			situsIDStr = *data.SitusID
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), situsIDStr)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), data.StatusVerifikasi)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), data.JenisSitus)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), data.Nama)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), data.JenisTipologi)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), data.NomorTelepon)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), noPengurusStr)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), data.Email)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), data.Website)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), data.NomorBadanHukum)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), data.TahunBerdiri)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), data.Provinsi)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), data.KabupatenKota)
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), data.Kecamatan)
		f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), data.Desa)
		f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), data.AlamatLengkap)
		f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), data.Latitude)
		f.SetCellValue(sheetName, fmt.Sprintf("R%d", row), data.Longitude)
		f.SetCellValue(sheetName, fmt.Sprintf("S%d", row), data.LuasTanah)
		f.SetCellValue(sheetName, fmt.Sprintf("T%d", row), data.LuasBangunan)
		f.SetCellValue(sheetName, fmt.Sprintf("U%d", row), data.StatusTanah)
		f.SetCellValue(sheetName, fmt.Sprintf("V%d", row), data.NomorAIW)
		f.SetCellValue(sheetName, fmt.Sprintf("W%d", row), data.NomorSertifikatWakaf)
		f.SetCellValue(sheetName, fmt.Sprintf("X%d", row), data.DayaTampungMax)

		flatDetail := allRowDetails[i]
		for j, dh := range dynamicHeaders {
			colIndex := len(baseHeaders) + j + 1
			cell, _ := excelize.CoordinatesToCellName(colIndex, row)

			if val, ok := flatDetail[dh]; ok {
				f.SetCellValue(sheetName, cell, val)
			}
		}
	}

	return f, nil
}
