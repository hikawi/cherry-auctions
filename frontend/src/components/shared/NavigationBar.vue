<script setup lang="ts">
import { useProfileStore } from "@/stores/profile";
import { LucideMenu, LucideX } from "lucide-vue-next";
import { computed, ref } from "vue";
import { useRoute } from "vue-router";
import OverlayScreen from "./OverlayScreen.vue";

const route = useRoute();
const profile = useProfileStore();
const menuOpen = ref(false);

const urlEncodedName = computed(() => {
  return `https://ui-avatars.com/api/?name=${encodeURIComponent(profile.profile?.name || "")}`;
});

const links = [
  {
    name: "navigation.home",
    href: "/",
  },
  {
    name: "navigation.all_products",
    href: "/search",
  },
  {
    name: "navigation.acknowledgements",
    href: "/acknowledgements",
  },
];
</script>

<template>
  <div
    class="border-claret-100 flex w-full flex-row items-center justify-between border-b pb-4 lg:px-6"
  >
    <a
      class="via-watermelon-600 to-claret-600 flex flex-row items-center gap-2 bg-linear-to-r from-pink-600 bg-clip-text text-2xl font-black duration-200 hover:text-transparent"
      href="/"
    >
      <img src="/icon.png" alt="CherryAuctions" class="size-8" width="32" height="32" />
      CherryAuctions
    </a>

    <button
      class="flex cursor-pointer items-center justify-center duration-200 hover:rotate-90 md:hidden"
      @click="menuOpen = true"
    >
      <LucideMenu class="size-6 text-black" />
    </button>

    <!-- The sheet to move up for a mobile nav bar -->
    <OverlayScreen :shown="menuOpen" class="flex-col items-center justify-end p-0">
      <div
        class="flex h-4/5 w-full flex-col gap-4 rounded-t-2xl bg-white px-6 py-4 duration-200"
        :class="{
          '-translate-y-full': !menuOpen,
          'translate-y-0': menuOpen,
        }"
      >
        <div class="relative flex items-center justify-end">
          <h2 class="absolute inset-y-0 left-1/2 -translate-x-1/2 py-2 text-xl font-bold">
            {{ $t("navigation.bar_title") }}
          </h2>

          <button class="rounded-full p-2 duration-200 hover:bg-black/20" @click="menuOpen = false">
            <LucideX class="size-6 text-black" />
          </button>
        </div>

        <div class="flex w-full flex-col gap-2">
          <template v-for="link in links" :key="link.href">
            <a
              v-if="route.path == link.href"
              :href="link.href"
              class="w-full cursor-pointer rounded-xl bg-black/10 p-2 px-4"
            >
              {{ $t(link.name) }}
            </a>
            <a
              v-else
              :href="link.href"
              class="w-full cursor-pointer rounded-xl p-2 px-4 hover:bg-black/10"
            >
              {{ $t(link.name) }}
            </a>
          </template>
        </div>
      </div>
    </OverlayScreen>

    <!-- Classic nav bar on desktop -->
    <div class="hidden flex-row items-center gap-4 md:flex">
      <nav class="flex flex-row items-center gap-4">
        <template v-for="link in links" :key="link.href">
          <a
            v-if="route.path == link.href"
            :href="link.href"
            class="cursor-pointer underline underline-offset-8"
          >
            {{ $t(link.name) }}
          </a>
          <a
            v-else
            :href="link.href"
            class="cursor-pointer text-black/50 duration-200 hover:text-black/75 hover:underline hover:underline-offset-8"
          >
            {{ $t(link.name) }}
          </a>
        </template>
      </nav>

      <img
        v-if="profile.hasProfile"
        :src="urlEncodedName"
        class="aspect-square h-10 w-auto rounded-full"
      />
      <a
        href="/login"
        class="bg-claret-600 hover:bg-claret-700 flex h-full w-fit min-w-fit items-center justify-center rounded-lg px-4 py-2 font-semibold text-white duration-200"
        v-else
        >{{ $t("general.login") }}</a
      >
    </div>
  </div>
</template>
