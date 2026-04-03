import * as v from "valibot";

export const UserSchema = v.object({
  nama_lengkap: v.string(),
  permissions: v.array(v.string()),
});

export type User = v.InferOutput<typeof UserSchema> | null;

export const UserInputSchema = v.object({
  nip: v.pipe(
    v.string("NIP/NIK harus berupa teks"),
    v.regex(
      /^(\d{16}|\d{18})$/,
      "Harus berupa NIK (16 digit) atau NIP (18 digit angka)",
    ),
  ),
  nama_lengkap: v.pipe(v.string(), v.minLength(1, "Nama lengkap wajib diisi")),
  email: v.optional(
    v.union([
      v.pipe(v.string(), v.email("Format email tidak valid")),
      v.literal(""),
    ]),
  ),
  nomor_telepon: v.pipe(
    v.string(),
    v.minLength(1, "Nomor telepon wajib diisi"),
    v.regex(
      /^[0-9+\-\s]+$/,
      "Format nomor telepon tidak valid (hanya angka, +, atau -)",
    ),
  ),
  jabatan: v.pipe(v.string(), v.minLength(1, "Jabatan wajib diisi")),
  unit_kerja: v.optional(v.string()),
  peran: v.pipe(v.string(), v.minLength(1, "Peran wajib dipilih")),
});

export type UserInput = v.InferOutput<typeof UserInputSchema>;

export const RoleResponseSchema = v.object({
  id: v.string(),
  name: v.string(),
  created_at: v.string(),
});

export type RoleResponse = v.InferOutput<typeof RoleResponseSchema>;

export const UserResponseSchema = v.object({
  id: v.string(),
  nip: v.string(),
  nama_lengkap: v.string(),
  email: v.nullable(v.string()),
  nomor_telepon: v.string(),
  jabatan: v.string(),
  unit_kerja: v.nullable(v.string()),
  peran: v.string(),
});

export type UserResponse = v.InferOutput<typeof UserResponseSchema>;

export const ProfileResponseSchema = v.object({
  id: v.string(),
  nip: v.string(),
  nama_lengkap: v.string(),
  email: v.nullable(v.string()),
  nomor_telepon: v.string(),
  jabatan: v.string(),
  unit_kerja: v.nullable(v.string()),
  peran: v.string(),
  created_at: v.string(),
});

export type ProfileResponse = v.InferOutput<typeof ProfileResponseSchema>;
