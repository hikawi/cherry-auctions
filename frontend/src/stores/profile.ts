import { endpoints } from "@/consts";
import { useFetch } from "@/hooks/use-fetch";
import type { Profile } from "@/types";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useProfileStore = defineStore("profile", () => {
  const profile = ref<Profile>();
  const { data, error, loading, doFetch } = useFetch<Profile>();

  const isAdmin = computed(() => profile.value && profile.value.roles.includes("admin"));
  const isFetching = computed(() => loading.value);
  const hasProfile = computed(() => profile.value != undefined);
  const hasFetched = computed(() => error.value != undefined || data.value != undefined);

  function setProfile(prof: Profile | undefined) {
    profile.value = prof;
  }

  const fetchProfile = async () => {
    await doFetch(endpoints.self);
    if (!error.value && data.value) {
      setProfile(data.value);
    }
  };

  return { profile, isAdmin, isFetching, setProfile, hasProfile, hasFetched, fetchProfile };
});
