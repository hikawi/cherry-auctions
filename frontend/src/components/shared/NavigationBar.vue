<script setup lang="ts">
import { useProfileStore } from "@/stores/profile";
import {
  LucideHandCoins,
  LucideLayoutDashboard,
  LucideLogOut,
  LucideMenu,
  LucideRocket,
  LucideStar,
  LucideStore,
  LucideUser,
  LucideX,
} from "lucide-vue-next";
import { ref, type Component } from "vue";
import { useRoute } from "vue-router";
import AvatarCircle from "./AvatarCircle.vue";
import OverlayScreen from "./OverlayScreen.vue";

type Link =
  | {
      name: string;
      href: string;
      icon: Component;
      admin?: boolean;
    }
  | "separator";

const profile = useProfileStore();
const route = useRoute();

const navMenuOpen = ref(false);
const menuOpen = ref(false);

const guestLinks = [
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

const profileLinks: Link[] = [
  {
    name: "navigation.my_auctions",
    href: "/auctions",
    icon: LucideStore,
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
  "separator",
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
  <OverlayScreen shown v-if="navMenuOpen" class="items-end md:hidden">
    <div class="relative flex w-full flex-col gap-4 rounded-t-2xl bg-white p-6">
      <h1 class="text-center text-xl font-semibold">{{ $t("navigation.bar_title") }}</h1>

      <button
        @click="navMenuOpen = false"
        class="absolute top-5 right-6 cursor-pointer rounded-full p-2 duration-200 hover:bg-zinc-200"
      >
        <LucideX class="size-5 text-black" />
      </button>

      <hr class="rounded-full border-zinc-300" />

      <div class="flex w-full flex-col">
        <router-link
          v-for="(link, idx) in guestLinks"
          :key="idx"
          :to="{ path: link.href }"
          class="hover:text-claret-600 py-2 duration-200"
          :class="{
            'text-claret-600 font-semibold': route.path == link.href,
          }"
        >
          {{ $t(link.name) }}
        </router-link>
      </div>
    </div>
  </OverlayScreen>

  <div
    class="border-claret-100 flex w-full flex-row items-center justify-between border-b pb-4 lg:px-6"
  >
    <router-link
      class="via-watermelon-600 to-claret-600 flex flex-row items-center gap-2 bg-linear-to-r from-pink-600 bg-clip-text text-xl font-black duration-200 hover:text-transparent md:text-2xl"
      to="/"
    >
      <img src="/icon.png" alt="CherryAuctions" class="size-8" width="32" height="32" />
      <span class="hidden sm:block">CherryAuctions</span>
    </router-link>

    <div class="flex flex-row items-center justify-center gap-2 md:gap-4">
      <nav class="hidden flex-row items-center gap-4 md:flex">
        <router-link
          v-for="(link, idx) in guestLinks"
          :key="idx"
          :to="{ path: link.href }"
          class="underline-offset-8 hover:underline"
          :class="{ underline: route.path == link.href }"
        >
          {{ $t(link.name) }}
        </router-link>
      </nav>
      <button
        class="block cursor-pointer rounded-full p-2 duration-200 hover:bg-zinc-200 md:hidden"
        @click="navMenuOpen = true"
      >
        <LucideMenu class="size-5 text-black" />
      </button>

      <div @click="menuOpen = !menuOpen" v-if="profile.hasProfile" class="relative shrink-0">
        <AvatarCircle
          :name="profile.profile?.name"
          :avatar_url="profile.profile?.avatar_url"
          hover
        />

        <div
          v-if="menuOpen"
          class="absolute right-0 -bottom-1 z-20 flex w-fit translate-y-full flex-col rounded-xl border border-zinc-500 bg-white shadow-md"
        >
          <template v-for="(link, idx) in profileLinks" :key="idx">
            <div
              v-if="link == 'separator'"
              class="h-2 w-full border-b border-zinc-300 bg-zinc-100"
            ></div>
            <router-link
              v-else-if="!link.admin || profile.isAdmin"
              :to="{ path: link.href }"
              class="flex flex-row items-center gap-2 border-b border-zinc-300 bg-white px-4 py-2 whitespace-nowrap duration-200 first-of-type:rounded-t-xl last-of-type:rounded-b-xl last-of-type:border-0 hover:bg-zinc-200"
              :class="{ 'bg-zinc-200': route.path == link.href }"
            >
              <component :is="link.icon" class="size-4 min-w-fit translate-y-0.5" />
              <span class="min-w-fit">
                {{ $t(link.name!) }}
              </span>
            </router-link>
          </template>
        </div>
      </div>
      <router-link
        :to="{ name: 'login' }"
        class="bg-claret-600 hover:bg-claret-700 flex h-full w-fit min-w-fit items-center justify-center rounded-lg px-4 py-2 font-semibold text-white duration-200"
        v-else
        >{{ $t("general.login") }}</router-link
      >
    </div>
  </div>
</template>
