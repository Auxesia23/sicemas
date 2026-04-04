import * as v from "valibot";
import { SitusKeagamaanSchema } from "./sites";

export const DashboardStatsSchema = v.object({
  total_situs: v.number(),
  total_jenis: v.number(),
  petugas_aktif: v.number(),
  menunggu_verifikasi: v.number(),
});

export const StatistikJenisSchema = v.object({
  nama: v.string(),
  jumlah_situs: v.number(),
});

export const RecentActivitySchema = v.object({
  id: v.pipe(v.string(), v.uuid()),
  user: v.string(),
  action: v.string(),
  target: v.string(),
  created_at: v.string(),
});

export const DashboardResponseSchema = v.object({
  stats: DashboardStatsSchema,
  statistik_jenis: v.optional(v.array(StatistikJenisSchema), []),
  recent_activities: v.optional(v.array(RecentActivitySchema), []),
  recent_sites: v.optional(v.array(SitusKeagamaanSchema), []),
});

export type DashboardResponse = v.InferOutput<typeof DashboardResponseSchema>;
