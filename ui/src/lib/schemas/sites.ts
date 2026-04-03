import * as v from "valibot";

export const SitusKeagamaanSchema = v.object({
  id: v.pipe(v.string(), v.uuid()),
  situs_id: v.nullable(v.string()),
  nama: v.string(),
  jenis: v.string(),
  lokasi: v.string(),
  status: v.string(),
  pendata: v.string(),
  updated_at: v.string(),
});

export type SitusKeagamaan = v.InferOutput<typeof SitusKeagamaanSchema>;

export const JenisSitusSchema = v.object({
  id: v.pipe(v.string(), v.uuid()),
  nama: v.string(),
  deskripsi: v.string(),
  jumlah_situs: v.pipe(v.number(), v.integer()),
});

export type JenisSitus = v.InferOutput<typeof JenisSitusSchema>;

export const JenisSitusInputSchema = v.object({
  nama: v.pipe(
    v.string("Nama harus berupa teks"),
    v.minLength(1, "Nama jenis situs wajib diisi"),
  ),
  deskripsi: v.pipe(
    v.string("Deskripsi harus berupa teks"),
    v.minLength(1, "Deskripsi jenis situs wajib diisi"),
  ),
});

export type JenisSitusInput = v.InferOutput<typeof JenisSitusInputSchema>;

export const FotoSchema = v.object({
  id: v.string(),
  image_url: v.string(),
});

export type Foto = v.InferOutput<typeof FotoSchema>;

export const SitusDetailResponseSchema = v.object({
  id: v.string(),
  situs_id: v.nullable(v.string()),
  nama: v.string(),
  jenis_tipologi: v.string(),
  jenis_situs: v.string(),
  status_verifikasi: v.string(),
  updated_at: v.nullable(v.string()),
  nomor_telpon_pengurus: v.array(v.string()),
  nomor_telepon: v.nullable(v.string()),
  email: v.nullable(v.string()),
  website: v.nullable(v.string()),
  tahun_berdiri: v.nullable(v.number()),
  nomor_badan_hukum: v.nullable(v.string()),
  alamat_lengkap: v.string(),
  desa: v.string(),
  kecamatan: v.string(),
  kabupaten_kota: v.string(),
  provinsi: v.string(),
  latitude: v.nullable(v.number()),
  longitude: v.nullable(v.number()),
  luas_tanah: v.nullable(v.number()),
  luas_bangunan: v.nullable(v.number()),
  status_tanah: v.nullable(v.string()),
  daya_tampung_max: v.nullable(v.number()),
  nomor_aiw: v.nullable(v.string()),
  nomor_sertifikat_wakaf: v.nullable(v.string()),
  foto: v.nullable(v.array(FotoSchema)),
  detail: v.any(),
});

export type SitusDetailResponse = v.InferOutput<
  typeof SitusDetailResponseSchema
>;

export const SitusPublikSchema = v.object({
  id: v.string(),
  nama: v.string(),
  jenis: v.string(),
  desa: v.string(),
  latitude: v.nullable(v.number()),
  longitude: v.nullable(v.number()),
  foto_utama: v.nullable(v.string()),
  jarak_meter: v.optional(v.nullable(v.number())),
});

export const SitusPublikResponseSchema = v.array(SitusPublikSchema);

export type SitusPublik = v.InferOutput<typeof SitusPublikSchema>;

export const SitusPublikDetailResponseSchema = v.object({
  id: v.string(),
  nama: v.string(),
  jenis: v.string(),
  alamat_lengkap: v.string(),
  desa: v.string(),
  latitude: v.number(),
  longitude: v.number(),
  luas_tanah: v.number(),
  daya_tampung: v.number(),
  galeri: v.array(v.string()),
  fasilitas: v.record(v.string(), v.any()),
});

export type SitusPublikDetail = v.InferOutput<
  typeof SitusPublikDetailResponseSchema
>;
