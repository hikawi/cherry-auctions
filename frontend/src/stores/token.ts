import { defineStore } from "pinia";
import { ref } from "vue";

export const useTokenStore = defineStore(
  "accessToken",
  () => {
    const token = ref<string>();
    function setToken(t: string | undefined) {
      token.value = t;
    }
    return { token, setToken };
  },
  {
    persist: true,
  },
);
