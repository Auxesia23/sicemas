import type { User } from "$lib/schemas/user";
import { UserSchema } from "$lib/schemas/user";
import * as v from "valibot";
import { apiFetch } from "$lib/api";
import fpjs from "@fingerprintjs/fingerprintjs";
import { goto } from "$app/navigation";

class AuthState {
  user = $state<User | null>(null);
  device_id = $state<string | null>(null);
  loading = $state(true);

  async initDevice() {
    if (this.device_id) return;

    try {
      const fp = await fpjs.load();
      const result = await fp.get();
      this.device_id = result.visitorId;
    } catch (e) {
      console.error(e);
    }
  }

  isAuthenticated() {
    return this.user !== null;
  }

  hasAnyPermission(requiredPermissions: string[]): boolean {
    if (!this.user || !Array.isArray(this.user.permissions)) return false;

    return requiredPermissions.some((permission) =>
      this.user!.permissions.includes(permission),
    );
  }

  hasAllPermissions(requiredPermissions: string[]) {
    if (!this.user || !Array.isArray(this.user.permissions)) return false;

    return requiredPermissions.every((permission) =>
      this.user!.permissions.includes(permission),
    );
  }

  async login(nip: string): Promise<string | null> {
    try {
      const res = await apiFetch("/api/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ nip }),
      });

      if (!res.ok) return await res.text();
      return null;
    } catch {
      return "Network error";
    }
  }

  async verifyOtp(otp: string, nip: string): Promise<string | null> {
    try {
      const res = await apiFetch("/api/auth/verify-otp", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          nip,
          otp,
        }),
      });

      if (!res.ok) return await res.text();
      await this.fetchUser();
      return null;
    } catch {
      return "Network error";
    }
  }

  async fetchUser(customFetch?: typeof fetch): Promise<boolean> {
    this.loading = true;
    try {
      const res = await apiFetch("/api/users/me", { customFetch });

      if (!res.ok) {
        this.user = null;
        return false;
      }

      const json = await res.json();
      const parsed = v.safeParse(UserSchema, json);

      if (!parsed.success) {
        this.user = null;
        return false;
      }

      this.user = parsed.output;
      return true;
    } catch {
      this.user = null;
      return false;
    } finally {
      this.loading = false;
    }
  }

  async logout(): Promise<void> {
    try {
      await apiFetch("/api/auth/logout", {
        method: "POST",
      });
    } finally {
      this.user = null;
      goto("/login");
    }
  }
}

export const auth = new AuthState();
