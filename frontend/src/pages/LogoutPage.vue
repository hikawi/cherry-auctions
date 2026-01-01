<script setup lang="ts">
import LoadingSpinner from "@/components/shared/LoadingSpinner.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { useProfileStore } from "@/stores/profile";
import { useTokenStore } from "@/stores/token";
import { onMounted } from "vue";
import { useRouter } from "vue-router";

const { authFetch } = useAuthFetch();
const router = useRouter();
const authToken = useTokenStore();
const profile = useProfileStore();

onMounted(async () => {
  await authFetch(endpoints.auth.logout);
  authToken.setToken(undefined);
  profile.setProfile(undefined);
  router.push({ name: "login" });
});
</script>

<template>
  <WhiteContainer>
    <LoadingSpinner />
  </WhiteContainer>
</template>
