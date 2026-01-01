<script setup lang="ts">
import { useProfileStore } from "@/stores/profile";
import {
  LucideAtSign,
  LucideHandCoins,
  LucideHouse,
  LucideLayoutDashboard,
  LucideLogOut,
  LucidePackageSearch,
  LucideRocket,
  LucideStar,
  LucideUser,
} from "lucide-vue-next";
import { computed, ref } from "vue";

const profile = useProfileStore();
const menuOpen = ref(false);

const urlEncodedName = computed(() => {
  return `https://ui-avatars.com/api/?name=${encodeURIComponent(profile.profile?.name || "")}`;
});

const links = [
  {
    name: "navigation.home",
    href: "/",
    icon: LucideHouse,
  },
  {
    name: "navigation.all_products",
    href: "/search",
    icon: LucidePackageSearch,
  },
  {
    name: "navigation.acknowledgements",
    href: "/acknowledgements",
    icon: LucideAtSign,
  },
  {
    name: "separator",
    href: "separator",
  },
  {
    name: "navigation.my_bids",
    href: "/bids",
    icon: LucideHandCoins,
  },
  {
    name: "navigation.my_reviews",
    href: "/reviews",
    icon: LucideStar,
  },
  {
    name: "navigation.my_subscriptions",
    href: "/subscriptions",
    icon: LucideRocket,
  },
  {
    name: "separator",
    href: "separator",
  },
  {
    href: "/admin",
    name: "navigation.admin",
    icon: LucideLayoutDashboard,
    admin: true,
  },
  {
    name: "navigation.profile",
    href: "/profile",
    icon: LucideUser,
  },
  {
    name: "navigation.logout",
    href: "/logout",
    icon: LucideLogOut,
  },
];
</script>

<template>
  <div
    class="border-claret-100 flex w-full flex-row items-center justify-between border-b pb-4 lg:px-6"
  >
    <a
      class="via-watermelon-600 to-claret-600 flex flex-row items-center gap-2 bg-linear-to-r from-pink-600 bg-clip-text text-xl font-black duration-200 hover:text-transparent md:text-2xl"
      href="/"
    >
      <img src="/icon.png" alt="CherryAuctions" class="size-8" width="32" height="32" />
      CherryAuctions
    </a>

    <div @click="menuOpen = !menuOpen" v-if="profile.hasProfile" class="relative">
      <img
        :src="urlEncodedName"
        class="hover:ring-claret-600 aspect-square h-10 w-auto cursor-pointer rounded-full hover:ring-2"
      />

      <div
        v-if="menuOpen"
        class="absolute right-0 -bottom-1 z-20 flex w-fit translate-y-full flex-col rounded-xl border border-zinc-500 bg-white shadow-md"
      >
        <template v-for="link in links" :key="link.href">
          <div
            v-if="link.href == 'separator'"
            class="h-2 w-full border-b border-zinc-300 bg-zinc-100"
          ></div>
          <a
            v-else-if="!link.admin || profile.isAdmin"
            :href="link.href"
            class="flex flex-row items-center gap-2 border-b border-zinc-300 bg-white px-4 py-2 whitespace-nowrap duration-200 first-of-type:rounded-t-xl last-of-type:rounded-b-xl last-of-type:border-0 hover:bg-zinc-200"
          >
            <component :is="link.icon" class="size-4 min-w-fit translate-y-0.5" />
            <span class="min-w-fit">
              {{ $t(link.name!) }}
            </span>
          </a>
        </template>
      </div>
    </div>
    <a
      href="/login"
      class="bg-claret-600 hover:bg-claret-700 flex h-full w-fit min-w-fit items-center justify-center rounded-lg px-4 py-2 font-semibold text-white duration-200"
      v-else
      >{{ $t("general.login") }}</a
    >
  </div>
</template>
