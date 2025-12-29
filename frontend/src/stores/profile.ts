import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import type { Profile } from "@/types";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useProfileStore = defineStore("profile", () => {
  const profile = ref<Profile>();
  const loading = ref(false);
  const error = ref();
  const { authFetch } = useAuthFetch();

  const isAdmin = computed(() => profile.value && profile.value.roles.includes("admin"));
  const isFetching = computed(() => loading.value);
  const hasProfile = computed(() => profile.value != undefined);
  const hasFetched = computed(() => error.value != undefined || profile.value != undefined);

  function setProfile(prof: Profile | undefined) {
    profile.value = prof;
  }

  const fetchProfile = async () => {
    loading.value = true;

    try {
      const res = await authFetch(endpoints.self);
      if (res.ok) {
        setProfile(await res.json());
      }
    } finally {
      loading.value = false;
    }
  };

  return { profile, isAdmin, isFetching, setProfile, hasProfile, hasFetched, fetchProfile };
});
