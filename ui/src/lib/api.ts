import { auth } from "$lib/states/auth.svelte";

interface FetchOptions extends RequestInit {
  customFetch?: typeof fetch;
}

function getCsrfFromCookie(): string | null {
  const match = document.cookie.match(/(?:^|; )csrf_token=([^;]*)/);
  return match ? decodeURIComponent(match[1]) : null;
}

type RefreshSubscriber = (success: boolean) => void;

class TokenRefreshManager {
  private isRefreshing = false;
  private subscribers: RefreshSubscriber[] = [];

  async refresh(customFetch: typeof fetch): Promise<boolean> {
    if (this.isRefreshing) {
      return this.waitForRefresh();
    }

    this.isRefreshing = true;

    try {
      const response = await customFetch("/api/auth/refresh", {
        method: "POST",
        headers: { "X-Device-Id": auth.device_id || "" },
      });

      const success = response.ok;
      this.notifySubscribers(success);

      if (!success) {
        auth.user = null;
      }

      return success;
    } finally {
      this.isRefreshing = false;
      this.subscribers = [];
    }
  }

  private waitForRefresh(): Promise<boolean> {
    return new Promise((resolve) => {
      this.subscribers.push(resolve);
    });
  }

  private notifySubscribers(success: boolean): void {
    this.subscribers.forEach((cb) => cb(success));
  }
}

const refreshManager = new TokenRefreshManager();

function buildRequestConfig(
  options: FetchOptions,
  deviceId: string | null,
): RequestInit {
  const { customFetch: _, ...init } = options;
  const headers = new Headers(init.headers);

  if (deviceId) {
    headers.set("X-Device-Id", deviceId);
  }

  const csrfToken = getCsrfFromCookie();
  if (csrfToken) {
    headers.set("X-CSRF-Token", csrfToken);
  }

  return { ...init, headers };
}

export async function apiFetch(
  endpoint: string,
  options: FetchOptions = {},
): Promise<Response> {
  const { customFetch = fetch } = options;
  const config = buildRequestConfig(options, auth.device_id);

  let response = await customFetch(endpoint, config);

  if (response.status !== 401) {
    return response;
  }

  if (endpoint === "/api/auth/refresh") {
    return response;
  }

  const refreshed = await refreshManager.refresh(customFetch);

  if (!refreshed) {
    return response;
  }

  return customFetch(endpoint, config);
}

export function createPageLoadFetch(eventFetch: typeof fetch): typeof fetch {
  return (input: RequestInfo | URL, init?: RequestInit) => {
    return apiFetch(input.toString(), {
      ...init,
      customFetch: eventFetch,
    });
  };
}
