package dto

type Stats struct {
	TotalSitus         int `json:"total_situs" db:"total_situs"`
	TotalJenis         int `json:"total_jenis" db:"total_jenis"`
	PetugasAktif       int `json:"petugas_aktif" db:"petugas_aktif"`
	MenungguVerifikasi int `json:"menunggu_verifikasi" db:"menunggu_verifikasi"`
}

type StatistikJenis struct {
	Nama        string `json:"nama" db:"nama"`
	JumlahSitus int    `json:"jumlah_situs" db:"jumlah_situs"`
}

type DashboardResponse struct {
	Stats *Stats `json:"stats"`

	// Ini untuk grafik (Pakai DTO jenis situs yang udah lu punya)
	StatistikJenis []StatistikJenis `json:"statistik_jenis"`

	// Ini untuk list Aktivitas (Maksimal 5)
	RecentActivities []ActivityResponse `json:"recent_activities"`

	// Ini untuk list tabel Situs Terbaru (Maksimal 5)
	RecentSites []SitusKeagamaanResponse `json:"recent_sites"`
}

type LandingStatsResponse struct {
	TotalSitus     int `db:"total_situs" json:"total_situs"`
	TotalDesa      int `db:"total_desa" json:"total_desa"`
	TotalKategori  int `db:"total_kategori" json:"total_kategori"`
	TotalKapasitas int `db:"total_kapasitas" json:"total_kapasitas"`
}
