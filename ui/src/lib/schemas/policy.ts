import * as v from "valibot";

export const PolicySchema = v.object({
  subject: v.string(),
  object: v.string(),
  action: v.string(),
});

export type PolicyResponse = v.InferOutput<typeof PolicySchema>;
