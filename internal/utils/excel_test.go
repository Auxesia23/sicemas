package utils

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
	"sicemas/internal/dto"
)

func TestExportSitusToExcel(t *testing.T) {
	detailRaw := json.RawMessage(`{"fasilitas": {"wifi": true, "parkir": "luas"}, "kapasitas": 500, "listrik": "4500W"}`)

	situsID := "ST.01.01.001"

	t.Run("success with single row", func(t *testing.T) {
		data := []dto.SitusKeagamaanDetailResponse{
			{
				ID:                  uuid.New(),
				StatusVerifikasi:    "Terverifikasi",
				JenisSitus:          "Gereja",
				SitusID:             &situsID,
				Nama:                "Gereja Test",
				JenisTipologi:       "A",
				NomorTelepon:        "021-12345678",
				NomorTelponPengurus: pq.StringArray{"0811111111", "0822222222"},
				Email:               "test@example.com",
				Website:             "https://example.com",
				NomorBadanHukum:     "AHU-001",
				TahunBerdiri:        2020,
				Provinsi:            "Jawa Barat",
				KabupatenKota:       "Sukabumi",
				Kecamatan:           "Ciemas",
				Desa:                "Mekarjaya",
				AlamatLengkap:       "Jl. Raya Ciemas No. 1",
				Latitude:            -6.1234,
				Longitude:           106.5678,
				LuasTanah:           500,
				LuasBangunan:        200,
				StatusTanah:         "Milik Sendiri",
				NomorAIW:            "AIW-001",
				NomorSertifikatWakaf: "SW-001",
				DayaTampungMax:      100,
				Detail:              &detailRaw,
			},
		}

		f, err := ExportSitusToExcel(data)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		sheetName := "Data Situs"
		val, err := f.GetCellValue(sheetName, "A1")
		assert.NoError(t, err)
		assert.Equal(t, "ID", val)

		val, err = f.GetCellValue(sheetName, "A2")
		assert.NoError(t, err)
		assert.Equal(t, situsID, val)

		val, err = f.GetCellValue(sheetName, "B2")
		assert.NoError(t, err)
		assert.Equal(t, "Terverifikasi", val)

		// Check dynamic headers exist
		rows, err := f.GetRows(sheetName)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(rows[0]), 24)
	})

	t.Run("success with multiple rows", func(t *testing.T) {
		data := []dto.SitusKeagamaanDetailResponse{
			{
				ID:               uuid.New(),
				StatusVerifikasi: "Terverifikasi",
				JenisSitus:       "Masjid",
				Nama:             "Masjid Al-Hidayah",
				JenisTipologi:    "B",
				Provinsi:         "Jawa Barat",
				KabupatenKota:    "Sukabumi",
				Kecamatan:        "Ciemas",
				Desa:             "Mekarjaya",
				AlamatLengkap:    "Jl. Masjid No. 1",
				Latitude:         -6.5,
				Longitude:        106.8,
			},
			{
				ID:               uuid.New(),
				StatusVerifikasi: "Belum Terverifikasi",
				JenisSitus:       "Pura",
				Nama:             "Pura Hindu",
				JenisTipologi:    "C",
				Provinsi:         "Bali",
				KabupatenKota:    "Badung",
				Kecamatan:        "Kuta",
				Desa:             "Kuta Selatan",
				AlamatLengkap:    "Jl. Pantai No. 1",
				Latitude:         -8.7,
				Longitude:        115.2,
			},
		}

		f, err := ExportSitusToExcel(data)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		rows, err := f.GetRows("Data Situs")
		assert.NoError(t, err)
		assert.Len(t, rows, 3) // header + 2 data rows
	})

	t.Run("success with empty data list", func(t *testing.T) {
		f, err := ExportSitusToExcel([]dto.SitusKeagamaanDetailResponse{})
		assert.NoError(t, err)
		assert.NotNil(t, f)

		rows, err := f.GetRows("Data Situs")
		assert.NoError(t, err)
		assert.Len(t, rows, 1) // only header
	})

	t.Run("success with nil SitusID", func(t *testing.T) {
		data := []dto.SitusKeagamaanDetailResponse{
			{
				ID:               uuid.New(),
				StatusVerifikasi: "Terverifikasi",
				JenisSitus:       "Vihara",
				Nama:             "Vihara Dharma",
				JenisTipologi:    "A",
				Provinsi:         "DKI Jakarta",
				KabupatenKota:    "Jakarta Pusat",
				Kecamatan:        "Gambir",
				Desa:             "Petojo",
				AlamatLengkap:    "Jl. Vihara No. 1",
				Latitude:         -6.2,
				Longitude:        106.8,
				SitusID:          nil,
			},
		}

		f, err := ExportSitusToExcel(data)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		val, err := f.GetCellValue("Data Situs", "A2")
		assert.NoError(t, err)
		assert.Equal(t, "", val)
	})

	t.Run("success with complex nested detail JSON", func(t *testing.T) {
		complexDetail := json.RawMessage(`{
			"fasilitas": {
				"ruangan": {
					"utama": "Aula Besar",
					"kelas": ["Kelas A", "Kelas B", "Kelas C"]
				},
				"parkir": "luas",
				"toilet": 5
			},
			"kegiatan": ["Ibadah", "Pengajian", "Sosial"]
		}`)

		data := []dto.SitusKeagamaanDetailResponse{
			{
				ID:               uuid.New(),
				StatusVerifikasi: "Terverifikasi",
				JenisSitus:       "Masjid",
				Nama:             "Masjid Jami",
				JenisTipologi:    "A",
				Provinsi:         "Jawa Barat",
				KabupatenKota:    "Sukabumi",
				Kecamatan:        "Ciemas",
				Desa:             "Mekarjaya",
				AlamatLengkap:    "Jl. Masjid Jami No. 1",
				Latitude:         -6.5,
				Longitude:        106.8,
				Detail:           &complexDetail,
			},
		}

		f, err := ExportSitusToExcel(data)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		rows, err := f.GetRows("Data Situs")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(rows[0]), 24) // base headers + dynamic headers
	})

	t.Run("file has bold header style", func(t *testing.T) {
		data := []dto.SitusKeagamaanDetailResponse{
			{
				ID:               uuid.New(),
				StatusVerifikasi: "Terverifikasi",
				JenisSitus:       "Klenteng",
				Nama:             "Klenteng Test",
				JenisTipologi:    "A",
				Provinsi:         "Jawa Timur",
				KabupatenKota:    "Surabaya",
				Kecamatan:        "Gubeng",
				Desa:             "Gubeng",
				AlamatLengkap:    "Jl. Klenteng No. 1",
				Latitude:         -7.2,
				Longitude:        112.7,
			},
		}

		f, err := ExportSitusToExcel(data)
		assert.NoError(t, err)

		styleID, err := f.GetCellStyle("Data Situs", "A1")
		assert.NoError(t, err)
		style, err := f.GetStyle(styleID)
		assert.NoError(t, err)
		assert.NotNil(t, style.Font)
		assert.True(t, style.Font.Bold)
	})

	t.Run("exported file can be saved and reopened", func(t *testing.T) {
		data := []dto.SitusKeagamaanDetailResponse{
			{
				ID:               uuid.New(),
				StatusVerifikasi: "Terverifikasi",
				JenisSitus:       "Masjid",
				Nama:             "Masjid Testing",
				JenisTipologi:    "B",
				Provinsi:         "Jawa Barat",
				KabupatenKota:    "Sukabumi",
				Kecamatan:        "Ciemas",
				Desa:             "Mekarjaya",
				AlamatLengkap:    "Jl. Testing No. 1",
				Latitude:         -6.5,
				Longitude:        106.8,
				NomorTelepon:     "021-99999999",
				LuasTanah:        1000,
				LuasBangunan:     300,
			},
		}

		f, err := ExportSitusToExcel(data)
		assert.NoError(t, err)

		buf, err := f.WriteToBuffer()
		assert.NoError(t, err)
		assert.NotEmpty(t, buf.Bytes())

		reopened, err := excelize.OpenReader(buf)
		assert.NoError(t, err)
		assert.NotNil(t, reopened)
		reopened.Close()
	})

	t.Run("handles nil detail gracefully", func(t *testing.T) {
		data := []dto.SitusKeagamaanDetailResponse{
			{
				ID:               uuid.New(),
				StatusVerifikasi: "Terverifikasi",
				JenisSitus:       "Masjid",
				Nama:             "Masjid No Detail",
				JenisTipologi:    "A",
				Provinsi:         "Jawa Barat",
				KabupatenKota:    "Sukabumi",
				Kecamatan:        "Ciemas",
				Desa:             "Mekarjaya",
				AlamatLengkap:    "Jl. No Detail No. 1",
				Latitude:         -6.5,
				Longitude:        106.8,
				Detail:           nil,
			},
		}

		f, err := ExportSitusToExcel(data)
		assert.NoError(t, err)
		assert.NotNil(t, f)
	})
}
