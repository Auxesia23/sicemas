import * as v from "valibot";

export const DetailMasjidSchema = v.object({
  sdm_masjid: v.object({
    jumlah_pengurus: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
    jumlah_imam: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
    jumlah_khotib: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
    jumlah_muadzin: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
    jumlah_remaja: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
    jumlah_jemaah: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
  }),
  pengurus_dkm: v.object({
    ketua: v.pipe(
      v.string("Ketua DKM wajib diisi"),
      v.minLength(1, "Ketua DKM wajib diisi"),
    ),
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
    luas_m2: null,
    jumlah_pengurus: null,
    kondisi: null,
    jenis_buku: [],
  },
  kalibrasi_arah_kiblat: {
    azimut: null,
    tanggal_kalibrasi: null,
  },
  sdm_masjid: {
    jumlah_pengurus: null,
    jumlah_imam: null,
    jumlah_khotib: null,
    jumlah_muadzin: null,
    jumlah_remaja: null,
    jumlah_jemaah: null,
  },
  pengurus_dkm: {
    ketua: null,
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
    jumlah_imam: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
    jumlah_muadzin: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
    jumlah_jemaah: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
  }),
  nama_imam: v.pipe(v.array(v.string()), v.minLength(1, "Minimal 1 imam")),
  nama_muadzin: v.pipe(
    v.array(v.string()),
    v.minLength(1, "Minimal 1 muadzin"),
  ),
});

export const DEFAULT_DETAIL_MUSHOLA = {
  kalibrasi_arah_kiblat: {
    azimut: null,
    tanggal_kalibrasi: null,
  },
  sdm_masjid: {
    jumlah_imam: null,
    jumlah_muadzin: null,
    jumlah_jemaah: null,
  },
  fasilitas_umum: [],
  nama_imam: [],
  nama_muadzin: [],
  kegiatan: [],
};

export const DetailMTSchema = v.object({
  pengurus: v.object({
    ketua: v.pipe(
      v.string("Ketua wajib diisi"),
      v.minLength(1, "Ketua wajib diisi"),
    ),
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
      jumlah_dewasa: v.nullable(v.number()),
      jumlah_remaja: v.nullable(v.number()),
      waktu_pengajian: v.nullable(v.string()),
    }),
  ),
  kehadiran_wanita: v.optional(
    v.object({
      jumlah_dewasa: v.nullable(v.number()),
      jumlah_remaja: v.nullable(v.number()),
      waktu_pengajian: v.nullable(v.string()),
    }),
  ),
});

export const DEFAULT_DETAIL_MT = {
  pengurus: {
    ketua: null,
    sekretaris: [],
    bendahara: [],
  },
  kehadiran_pria: {
    jumlah_dewasa: null,
    jumlah_remaja: null,
    waktu_pengajian: null,
  },
  kehadiran_wanita: {
    jumlah_dewasa: null,
    jumlah_remaja: null,
    waktu_pengajian: "",
  },
  penceramah: [],
};

export const DetailPesantrenSchema = v.object({
  nama_yayasan: v.nullable(v.string()),
  pimpinan_pesantren: v.pipe(
    v.string("Nama Pimpinan wajib diisi"),
    v.minLength(1, "Nama Pimpinan wajib diisi"),
  ),
  kepengurusan: v.object({
    ketua: v.pipe(
      v.string("Ketua wajib diisi"),
      v.minLength(1, "Ketua wajib diisi"),
    ),
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
      pria: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
      wanita: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
    }),
    tidak_mondok: v.object({
      pria: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
      wanita: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
    }),
    total: v.object({
      pria: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
      wanita: v.pipe(v.number("Wajib diisi"), v.minValue(0)),
    }),
  }),
});

export const DEFAULT_DETAIL_PESANTREN = {
  nama_yayasan: null,
  pimpinan_pesantren: null,
  kepengurusan: {
    ketua: null,
    sekretaris: [],
    bendahara: [],
  },
  santri: {
    mondok: { pria: null, wanita: null },
    tidak_mondok: { pria: null, wanita: null },
    total: { pria: null, wanita: null },
  },
  fasilitas_pondok: [],
};

export const SitusInputSchema = v.object({
  jenis_situs_id: v.pipe(
    v.string("Jenis situs wajib dipilih"),
    v.minLength(1, "Jenis situs wajib dipilih"),
  ),
  nama: v.pipe(
    v.string("Nama situs wajib diisi"),
    v.minLength(1, "Nama situs wajib diisi"),
  ),
  nomor_telepon: v.pipe(v.string(), v.minLength(0)),
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
  nomor_badan_hukum: v.pipe(v.string(), v.minLength(0)),
  jenis_tipologi: v.pipe(v.string(), v.minLength(1, "Tipologi wajib dipilih")),
  tahun_berdiri: v.pipe(
    v.number("Wajib diisi"),
    v.minValue(1000, "Tahun tidak valid"),
  ),

  provinsi: v.pipe(v.string(), v.minLength(1)),
  kabupaten_kota: v.pipe(v.string(), v.minLength(1)),
  kecamatan: v.pipe(v.string(), v.minLength(1)),
  desa: v.pipe(v.string(), v.minLength(1, "Desa wajib diisi")),
  alamat_lengkap: v.pipe(v.string(), v.minLength(5, "Alamat kurang lengkap")),
  latitude: v.pipe(v.number("Tentukan lokasi di peta")),
  longitude: v.pipe(v.number("Tentukan lokasi di peta")),
  luas_tanah: v.pipe(
    v.number("Luas tanah wajib diisi"),
    v.minValue(1, "Luas tanah wajib diisi"),
  ),
  luas_bangunan: v.pipe(
    v.number("Luas bangunan wajib diisi"),
    v.minValue(1, "Luas bangunan wajib diisi"),
  ),
  nomor_telpon_pengurus: v.pipe(
    v.array(v.string()),
    v.minLength(1, "Minimal 1 nomor pengurus"),
  ),
  status_tanah: v.pipe(
    v.string(),
    v.minLength(1, "Status tanah wajib dipilih"),
  ),
  daya_tampung_max: v.pipe(
    v.number("Daya tampung wajib diisi"),
    v.minValue(1, "Daya tampung wajib diisi"),
  ),
  nomor_aiw: v.pipe(v.string(), v.minLength(0)),
  nomor_sertifikat_wakaf: v.pipe(v.string(), v.minLength(0)),
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
  tahun_berdiri: null,
  provinsi: "Jawa Barat",
  kabupaten_kota: "Kab. Sukabumi",
  kecamatan: "Ciemas",
  desa: "",
  alamat_lengkap: "",
  latitude: null as number | null,
  longitude: null as number | null,
  luas_tanah: null as number | null,
  luas_bangunan: null as number | null,
  status_tanah: "",
  nomor_aiw: "",
  nomor_sertifikat_wakaf: "",
  daya_tampung_max: null as number | null,
  detail: {},
};

export const SitusEditSchema = v.object({
  nama: v.pipe(v.string(), v.minLength(1, "Nama situs wajib diisi")),
  nomor_telepon: v.pipe(v.string(), v.minLength(0)),
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
  nomor_badan_hukum: v.pipe(v.string(), v.minLength(0)),
  jenis_tipologi: v.pipe(v.string(), v.minLength(1, "Tipologi wajib dipilih")),
  tahun_berdiri: v.pipe(
    v.number("Wajib diisi"),
    v.minValue(1000, "Tahun tidak valid"),
  ),
  provinsi: v.pipe(v.string(), v.minLength(1)),
  kabupaten_kota: v.pipe(v.string(), v.minLength(1)),
  kecamatan: v.pipe(v.string(), v.minLength(1)),
  desa: v.pipe(v.string(), v.minLength(1, "Desa wajib diisi")),
  alamat_lengkap: v.pipe(v.string(), v.minLength(5, "Alamat kurang lengkap")),
  latitude: v.pipe(v.number("Tentukan lokasi di peta")),
  longitude: v.pipe(v.number("Tentukan lokasi di peta")),
  luas_tanah: v.pipe(
    v.number("Luas tanah wajib diisi"),
    v.minValue(1, "Luas tanah wajib diisi"),
  ),
  luas_bangunan: v.pipe(
    v.number("Luas tanah wajib diisi"),
    v.minValue(1, "Luas bangunan wajib diisi"),
  ),
  nomor_telpon_pengurus: v.pipe(
    v.array(v.string()),
    v.minLength(1, "Minimal 1 nomor pengurus"),
  ),
  status_tanah: v.pipe(
    v.string(),
    v.minLength(1, "Status tanah wajib dipilih"),
  ),
  daya_tampung_max: v.pipe(
    v.number("Daya tampung wajib diisi"),
    v.minValue(1, "Daya tampung wajib diisi"),
  ),
  nomor_aiw: v.pipe(v.string(), v.minLength(0)),
  nomor_sertifikat_wakaf: v.pipe(v.string(), v.minLength(0)),
  detail: v.optional(v.any()),
});
