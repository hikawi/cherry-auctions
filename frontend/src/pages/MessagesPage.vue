<script setup lang="ts">
import NavigationBar from "@/components/shared/NavigationBar.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { onMounted } from "vue";

const { authFetch } = useAuthFetch();

onMounted(() => fetchChatSessions());

async function fetchChatSessions() {
  const res = await authFetch(endpoints.chat.index);
  console.log(await res.json());
}
</script>

<template>
  <WhiteContainer class="justify-start gap-6">
    <NavigationBar />
  </WhiteContainer>
</template>
