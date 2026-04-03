import * as v from "valibot";

export const DetailMasjidSchema = v.object({
  sdm_masjid: v.object({
    jumlah_pengurus: v.pipe(v.number(), v.minValue(1, "Wajib diisi")),
    jumlah_imam: v.pipe(v.number(), v.minValue(1, "Wajib diisi")),
    jumlah_khotib: v.pipe(v.number(), v.minValue(1, "Wajib diisi")),
    jumlah_muadzin: v.pipe(v.number(), v.minValue(1, "Wajib diisi")),
    jumlah_remaja: v.pipe(v.number(), v.minValue(1, "Wajib diisi")),
    jumlah_jemaah: v.pipe(v.number(), v.minValue(1, "Wajib diisi")),
  }),
  pengurus_dkm: v.object({
    ketua: v.pipe(v.string(), v.minLength(1, "Ketua DKM wajib diisi")),
    sekretaris: v.pipe(
      v.array(v.string()),
      v.minLength(1, "Minimal 1 sekretaris"),
    ),
    bendahara: v.pipe(
      v.array(v.string()),
      v.minLength(1, "Minimal 1 bendahara"),
    ),
    nama_imam: v.pipe(v.array(v.string()), v.minLength(1, "Minimal 1 imam")),
    nama_muazdin: v.pipe(
      v.array(v.string()),
      v.minLength(1, "Minimal 1 muadzin"),
    ),
  }),
});

export const DEFAULT_DETAIL_MASJID = {
  perpustakaan: {
    luas_m2: 0,
    jumlah_pengurus: 0,
    kondisi: "",
    jenis_buku: [],
  },
  kalibrasi_arah_kiblat: {
    azimut: "",
    tanggal_kalibrasi: "",
  },
  sdm_masjid: {
    jumlah_pengurus: 0,
    jumlah_imam: 0,
    jumlah_khotib: 0,
    jumlah_muadzin: 0,
    jumlah_remaja: 0,
    jumlah_jemaah: 0,
  },
  pengurus_dkm: {
    ketua: "",
    sekretaris: [],
    bendahara: [],
    nama_imam: [],
    nama_muazdin: [],
  },
  fasilitas_umum: [],
  fasilitas_raman_anak: [],
  fasilitas_disabilitas: [],
  kegiatan: [],
};

export const DetailMusholaSchema = v.object({
  sdm_masjid: v.object({
    jumlah_imam: v.pipe(v.number(), v.minValue(1, "Wajib diisi")),
    jumlah_muadzin: v.pipe(v.number(), v.minValue(1, "Wajib diisi")),
    jumlah_jemaah: v.pipe(v.number(), v.minValue(1, "Wajib diisi")),
  }),
  nama_imam: v.pipe(v.array(v.string()), v.minLength(1, "Minimal 1 imam")),
  nama_muadzin: v.pipe(
    v.array(v.string()),
    v.minLength(1, "Minimal 1 muadzin"),
  ),
});

export const DEFAULT_DETAIL_MUSHOLA = {
  kalibrasi_arah_kiblat: {
    azimut: "",
    tanggal_kalibrasi: "",
  },
  sdm_masjid: {
    jumlah_imam: 0,
    jumlah_muadzin: 0,
    jumlah_jemaah: 0,
  },
  fasilitas_umum: [],
  nama_imam: [],
  nama_muadzin: [],
  kegiatan: [],
};

export const DetailMTSchema = v.object({
  pengurus: v.object({
    ketua: v.pipe(v.string(), v.minLength(1, "Ketua wajib diisi")),
    sekretaris: v.pipe(
      v.array(v.string()),
      v.minLength(1, "Minimal 1 sekretaris"),
    ),
    bendahara: v.pipe(
      v.array(v.string()),
      v.minLength(1, "Minimal 1 bendahara"),
    ),
  }),
  penceramah: v.pipe(
    v.array(
      v.object({
        nama: v.pipe(v.string(), v.minLength(1)),
        materi: v.pipe(v.string(), v.minLength(1)),
      }),
    ),
    v.minLength(1, "Minimal 1 penceramah wajib"),
  ),
  kehadiran_pria: v.optional(
    v.object({
      jumlah_dewasa: v.pipe(v.number(), v.minValue(0)),
      jumlah_remaja: v.pipe(v.number(), v.minValue(0)),
      waktu_pengajian: v.string(),
    }),
  ),
  kehadiran_wanita: v.optional(
    v.object({
      jumlah_dewasa: v.pipe(v.number(), v.minValue(0)),
      jumlah_remaja: v.pipe(v.number(), v.minValue(0)),
      waktu_pengajian: v.string(),
    }),
  ),
});

export const DEFAULT_DETAIL_MT = {
  pengurus: {
    ketua: "",
    sekretaris: [],
    bendahara: [],
  },
  kehadiran_pria: {
    jumlah_dewasa: 0,
    jumlah_remaja: 0,
    waktu_pengajian: "",
  },
  kehadiran_wanita: {
    jumlah_dewasa: 0,
    jumlah_remaja: 0,
    waktu_pengajian: "",
  },
  penceramah: [],
};

export const DetailPesantrenSchema = v.object({
  nama_yayasan: v.optional(v.string()),
  pimpinan_pesantren: v.pipe(
    v.string(),
    v.minLength(1, "Nama Pimpinan wajib diisi"),
  ),
  kepengurusan: v.object({
    ketua: v.pipe(v.string(), v.minLength(1, "Ketua wajib diisi")),
    sekretaris: v.pipe(
      v.array(v.string()),
      v.minLength(1, "Minimal 1 sekretaris"),
    ),
    bendahara: v.pipe(
      v.array(v.string()),
      v.minLength(1, "Minimal 1 bendahara"),
    ),
  }),
  santri: v.object({
    mondok: v.object({
      pria: v.pipe(v.number(), v.minValue(0)),
      wanita: v.pipe(v.number(), v.minValue(0)),
    }),
    tidak_mondok: v.object({
      pria: v.pipe(v.number(), v.minValue(0)),
      wanita: v.pipe(v.number(), v.minValue(0)),
    }),
    total: v.object({
      pria: v.pipe(v.number(), v.minValue(0)),
      wanita: v.pipe(v.number(), v.minValue(0)),
    }),
  }),
});

export const DEFAULT_DETAIL_PESANTREN = {
  nama_yayasan: "",
  pimpinan_pesantren: "",
  kepengurusan: {
    ketua: "",
    sekretaris: [],
    bendahara: [],
  },
  santri: {
    mondok: { pria: 0, wanita: 0 },
    tidak_mondok: { pria: 0, wanita: 0 },
    total: { pria: 0, wanita: 0 },
  },
  fasilitas_pondok: [],
};

export const SitusInputSchema = v.object({
  jenis_situs_id: v.pipe(
    v.string(),
    v.minLength(1, "Jenis situs wajib dipilih"),
  ),
  nama: v.pipe(v.string(), v.minLength(1, "Nama situs wajib diisi")),
  nomor_telepon: v.optional(v.string()),
  email: v.union([
    v.literal(""),
    v.pipe(v.string(), v.email("Format email tidak valid")),
  ]),
  website: v.union([
    v.literal(""),
    v.pipe(
      v.string(),
      v.url(
        "Format URL/Website tidak valid (harus pakai http:// atau https://)",
      ),
    ),
  ]),
  nomor_badan_hukum: v.optional(v.string()),
  jenis_tipologi: v.pipe(v.string(), v.minLength(1, "Tipologi wajib dipilih")),
  tahun_berdiri: v.pipe(v.number(), v.minValue(1000, "Tahun tidak valid")),
  provinsi: v.pipe(v.string(), v.minLength(1)),
  kabupaten_kota: v.pipe(v.string(), v.minLength(1)),
  kecamatan: v.pipe(v.string(), v.minLength(1)),
  desa: v.pipe(v.string(), v.minLength(1, "Desa wajib diisi")),
  alamat_lengkap: v.pipe(v.string(), v.minLength(5, "Alamat kurang lengkap")),
  latitude: v.pipe(v.number("Tentukan lokasi di peta")),
  longitude: v.pipe(v.number("Tentukan lokasi di peta")),
  luas_tanah: v.pipe(v.number(), v.minValue(1, "Luas tanah wajib diisi")),
  luas_bangunan: v.pipe(v.number(), v.minValue(1, "Luas bangunan wajib diisi")),
  nomor_telpon_pengurus: v.pipe(
    v.array(v.string()),
    v.minLength(1, "Minimal 1 nomor pengurus"),
  ),
  status_tanah: v.pipe(
    v.string(),
    v.minLength(1, "Status tanah wajib dipilih"),
  ),
  daya_tampung_max: v.pipe(
    v.number(),
    v.minValue(1, "Daya tampung wajib diisi"),
  ),
  nomor_aiw: v.optional(v.string()),
  nomor_sertifikat_wakaf: v.optional(v.string()),
  detail: v.optional(v.any()),
});

export const DEFAULT_SITUS_FORM_DATA = {
  jenis_situs_id: "",
  nama: "",
  jenis_tipologi: "",
  nomor_telepon: "",
  nomor_telpon_pengurus: [],
  email: "",
  website: "",
  nomor_badan_hukum: "",
  tahun_berdiri: 2024,
  provinsi: "Jawa Barat",
  kabupaten_kota: "Kab. Sukabumi",
  kecamatan: "Ciemas",
  desa: "",
  alamat_lengkap: "",
  latitude: null as number | null,
  longitude: null as number | null,
  luas_tanah: 0,
  luas_bangunan: 0,
  status_tanah: "",
  nomor_aiw: "",
  nomor_sertifikat_wakaf: "",
  daya_tampung_max: 0,
  detail: {},
};

export const SitusEditSchema = v.object({
  nama: v.pipe(v.string(), v.minLength(1, "Nama situs wajib diisi")),
  nomor_telepon: v.optional(v.string()),
  email: v.union([
    v.literal(""),
    v.pipe(v.string(), v.email("Format email tidak valid")),
  ]),
  website: v.union([
    v.literal(""),
    v.pipe(
      v.string(),
      v.url(
        "Format URL/Website tidak valid (harus pakai http:// atau https://)",
      ),
    ),
  ]),
  nomor_badan_hukum: v.optional(v.string()),
  jenis_tipologi: v.pipe(v.string(), v.minLength(1, "Tipologi wajib dipilih")),
  tahun_berdiri: v.pipe(v.number(), v.minValue(1000, "Tahun tidak valid")),
  provinsi: v.pipe(v.string(), v.minLength(1)),
  kabupaten_kota: v.pipe(v.string(), v.minLength(1)),
  kecamatan: v.pipe(v.string(), v.minLength(1)),
  desa: v.pipe(v.string(), v.minLength(1, "Desa wajib diisi")),
  alamat_lengkap: v.pipe(v.string(), v.minLength(5, "Alamat kurang lengkap")),
  latitude: v.pipe(v.number("Tentukan lokasi di peta")),
  longitude: v.pipe(v.number("Tentukan lokasi di peta")),
  luas_tanah: v.pipe(v.number(), v.minValue(1, "Luas tanah wajib diisi")),
  luas_bangunan: v.pipe(v.number(), v.minValue(1, "Luas bangunan wajib diisi")),
  nomor_telpon_pengurus: v.pipe(
    v.array(v.string()),
    v.minLength(1, "Minimal 1 nomor pengurus"),
  ),
  status_tanah: v.pipe(
    v.string(),
    v.minLength(1, "Status tanah wajib dipilih"),
  ),
  daya_tampung_max: v.pipe(
    v.number(),
    v.minValue(1, "Daya tampung wajib diisi"),
  ),
  nomor_aiw: v.optional(v.string()),
  nomor_sertifikat_wakaf: v.optional(v.string()),
  detail: v.optional(v.any()),
});
