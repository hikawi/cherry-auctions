import { endpoints } from "@/consts";
import { useTokenStore } from "@/stores/token";
import { ref } from "vue";

export function useFetch<T>() {
  const data = ref<T | undefined>();
  const error = ref<unknown>();
  const loading = ref(false);
  const status = ref(0);
  const token = useTokenStore();

  const doFetch = async (url: string, options: RequestInit = {}) => {
    loading.value = true;
    error.value = undefined;

    // Attach token if we have one
    if (token.token) {
      options.headers = {
        ...options.headers,
        Authorization: `Bearer ${token.token}`,
      };
    }

    try {
      let response = await fetch(url, options);

      // If Unauthorized, try to refresh
      if (response.status === 401) {
        const success = await refreshAccessToken();
        if (success) {
          // Retry original call with NEW token
          options.headers = {
            ...options.headers,
            Authorization: `Bearer ${token.token}`,
          };
          response = await fetch(url, options);
          status.value = response.status;
        } else {
          status.value = response.status;
          throw new Error("Session expired. Please login again.");
        }
      }

      status.value = response.status;
      if (!response.ok) {
        throw new Error("Request failed");
      }

      data.value = await response.json();
    } catch (err: unknown) {
      error.value = err;
    } finally {
      loading.value = false;
    }
  };

  // Dedicated function to call the /refresh endpoint
  const refreshAccessToken = async () => {
    try {
      // Browser sends HttpOnly cookies automatically
      const res = await fetch(endpoints.auth.refresh, {
        method: "POST",
        credentials: "include",
        mode: "cors",
      });
      if (!res.ok) {
        return false;
      }

      const body = await res.json();
      token.setToken(body.access_token); // Update in-memory token
      return true;
    } catch {
      return false;
    }
  };

  return { data, error, loading, doFetch, status };
}
