-- +goose Up
-- +goose StatementBegin
CREATE TABLE situs_keagamaan (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    status_verifikasi VARCHAR(100) DEFAULT 'menunggu',

    -- Relasi
    jenis_situs_id UUID REFERENCES jenis_situs(id) ON DELETE RESTRICT,
    pendata_id UUID REFERENCES users(id) ON DELETE RESTRICT,

    situs_id VARCHAR(50),
    nama VARCHAR(255) NOT NULL,
    jenis_tipologi VARCHAR(100),


    -- [DATA UMUM]
    nomor_telepon TEXT,
    email TEXT,
    website VARCHAR(255),
    nomor_badan_hukum TEXT,
    tahun_berdiri INT,

    -- [LOKASI]
    provinsi VARCHAR(100) DEFAULT 'Jawa Barat',
    kabupaten_kota VARCHAR(100),
    kecamatan VARCHAR(100),
    desa VARCHAR(100),
    alamat_lengkap TEXT,
    koordinat geometry(Point, 4326),

    -- [INFORMASI KAPASITAS & TANAH]
    luas_tanah DECIMAL(10, 2),
    luas_bangunan DECIMAL(10, 2),
    status_tanah VARCHAR(100),
    nomor_aiw TEXT,
    nomor_sertifikat_wakaf TEXT,
    daya_tampung_max INT,

    -- =========================================================
    -- [SEMUA SISA DATA FORM MENGUMPUL DI SINI]
    -- =========================================================
    detail JSONB DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ  DEFAULT NOW(),
    updated_at TIMESTAMPTZ  DEFAULT NOW()
);
-- Indexing
CREATE INDEX idx_situs_koordinat ON situs_keagamaan USING GIST (koordinat);
-- Cukup 1 Index GIN untuk seluruh struktur data dinamis!
CREATE INDEX idx_situs_detail ON situs_keagamaan USING GIN (detail);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS situs_keagamaan;
-- +goose StatementEnd
