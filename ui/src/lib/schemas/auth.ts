import * as v from "valibot";

export const loginSchema = v.object({
  nip: v.pipe(
    v.string(),
    v.regex(/^(\d{16}|\d{18})$/, "NIP/NIK harus berupa 16 atau 18 digit angka"),
  ),
  password: v.pipe(v.string(), v.minLength(8, "Password minimal 8 karakter")),
});

export type UserLogin = v.InferOutput<typeof loginSchema>;

export const otpSchema = v.pipe(
  v.string(),
  v.regex(/^\d{6}$/, "OTP harus berupa 6 digit angka"),
);
