import * as v from "valibot";

export const nipSchema = v.pipe(
  v.string(),
  v.regex(/^(\d{16}|\d{18})$/, "NIP/NIK harus berupa 16 atau 18 digit angka"),
);

export const otpSchema = v.pipe(
  v.string(),
  v.regex(/^\d{6}$/, "OTP harus berupa 6 digit angka"),
);
