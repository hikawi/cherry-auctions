import { endpoints } from "@/consts";
import { useTokenStore } from "@/stores/token";

export function useAuthFetch({ json = true }: { json?: boolean } = {}) {
  const tokenStore = useTokenStore();

  async function authFetch(url: RequestInfo | URL, options: RequestInit = {}): Promise<Response> {
    // 1. Prepare Headers (Merge defaults with user overrides)
    const headers = new Headers(options.headers);
    if (tokenStore.token && !headers.has("Authorization")) {
      headers.set("Authorization", `Bearer ${tokenStore.token}`);
    }
    if (json && !headers.has("Content-Type")) {
      headers.set("Content-Type", "application/json");
    }

    const finalOptions = { ...options, headers };

    // 2. Initial Request
    let response = await fetch(url, finalOptions);

    // 3. Handle 401 and Token Refresh
    if (response.status === 401) {
      const refreshed = await tryRefresh();

      if (refreshed) {
        // Update headers with the NEW token
        const retryHeaders = new Headers(finalOptions.headers);
        retryHeaders.set("Authorization", `Bearer ${tokenStore.token}`);

        // Retry the request one time
        response = await fetch(url, { ...finalOptions, headers: retryHeaders });
      }
    }

    // 4. Return the whole response regardless of success/failure
    return response;
  }

  async function tryRefresh(): Promise<boolean> {
    try {
      const resp = await fetch(endpoints.auth.refresh, { method: "POST" });
      if (!resp.ok) throw new Error("Refresh failed");
      const data = await resp.json();

      // Assuming your store has an action to update the token
      tokenStore.setToken(data.access_token);
      return true;
    } catch (err) {
      console.error("Auth refresh failed:", err);
      // Optional: tokenStore.logout() if refresh fails
      return false;
    }
  }

  return { authFetch };
}
