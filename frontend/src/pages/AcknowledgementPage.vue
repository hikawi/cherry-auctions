<script setup lang="ts">
import GithubProfile from "@/components/acknowledgements/GithubProfile.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";

const contributors = [
  {
    username: "hikawi",
    name: "Nguyệt Ánh",
  },
  {
    username: "mrgeraffe-kun",
    name: "MrGeraffe",
  },
];
</script>

<template>
  <WhiteContainer class="justify-start gap-8">
    <h1 class="text-center text-3xl font-bold">{{ $t("acknowledgements.title") }}</h1>

    <section class="flex w-full max-w-3xl flex-col gap-4">
      <h2 class="text-xl font-semibold">{{ $t("acknowledgements.fonts") }}</h2>
      <p>{{ $t("acknowledgements.fonts_paragraph") }}</p>
    </section>

    <section class="flex w-full max-w-3xl flex-col gap-4">
      <h2 class="text-xl font-semibold">{{ $t("acknowledgements.contributors") }}</h2>
      <div class="flex w-full flex-row flex-wrap gap-4">
        <GithubProfile
          v-for="contributor in contributors"
          :key="contributor.username"
          :name="contributor.name"
          :username="contributor.username"
        />
      </div>
    </section>
  </WhiteContainer>
</template>
