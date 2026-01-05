<script setup lang="ts">
import GithubProfile from "@/components/acknowledgements/GithubProfile.vue";
import NavigationBar from "@/components/shared/NavigationBar.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";

type Contributor = {
  username: string;
  name: string;
  tasks: string[];
};

const contributors: Contributor[] = [
  {
    username: "hikawi",
    name: "Nguyệt Ánh",
    tasks: ["frontend", "backend", "design", "infra"],
  },
  {
    username: "mrgeraffe-kun",
    name: "MrGeraffe",
    tasks: ["forgot_password_frontend"],
  },
];
</script>

<template>
  <WhiteContainer class="justify-start gap-8">
    <NavigationBar />

    <h1 class="text-center text-3xl font-bold">{{ $t("acknowledgements.title") }}</h1>

    <section class="flex w-full max-w-3xl flex-col gap-4">
      <h2 class="text-xl font-semibold">{{ $t("acknowledgements.fonts") }}</h2>
      <p>{{ $t("acknowledgements.fonts_paragraph") }}</p>
    </section>

    <section class="flex w-full max-w-3xl flex-col gap-4">
      <h2 class="text-xl font-semibold">{{ $t("acknowledgements.contributors") }}</h2>
      <div class="flex w-full flex-col items-center justify-center gap-4 md:justify-start">
        <div
          class="flex w-full flex-col items-center gap-4 rounded-xl border border-zinc-300 p-4 shadow-md sm:flex-row sm:gap-8"
          v-for="contributor in contributors"
          :key="contributor.username"
        >
          <GithubProfile :name="contributor.name" :username="contributor.username" />

          <ul class="list-inside list-disc">
            <li v-for="task in contributor.tasks" :key="task">
              {{ $t(`acknowledgements.tasks.${task}`) }}
            </li>
          </ul>
        </div>
      </div>
    </section>
  </WhiteContainer>
</template>
