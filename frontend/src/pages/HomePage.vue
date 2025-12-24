<script setup lang="ts">
import { endpoints } from "@/consts";
import { useFetch } from "@/hooks/use-fetch";
import { onMounted } from "vue";

const { data: profile, error: profileError, loading: profileLoading, doFetch } = useFetch();

onMounted(async () => {
  await doFetch(endpoints.self);
  if (profileError) {
    console.log(profileError);
  }
});
</script>

<template>
  <div
    class="flex w-full flex-col items-center justify-center self-stretch rounded-2xl bg-white p-6"
  >
    <h1 class="text-2xl font-bold">Hello, world</h1>
    <p v-if="profileLoading">Loading...</p>
    <p v-else-if="profileError">
      Uh oh, an error happened. Probably unauthorized? Try
      <a class="underline" href="/login">logging in</a>.
    </p>
    <p v-else>{{ profile }}</p>
  </div>
</template>
