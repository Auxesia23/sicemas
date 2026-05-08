import * as v from "valibot";

export const ActivityResponseSchema = v.object({
  id: v.string(),
  user: v.string(),
  action: v.string(),
  target: v.string(),
  entity: v.string(),
  created_at: v.string(),
});

export type ActivityResponse = v.InferOutput<typeof ActivityResponseSchema>;
