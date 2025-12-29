<script setup lang="ts">
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const links = [
  {
    to: "/admin/categories",
    name: "admin-categories",
    label: "general.categories",
  },
  {
    to: "/admin/products",
    name: "admin-products",
    label: "general.products",
  },
  {
    to: "/admin/users",
    name: "admin-users",
    label: "general.users",
  },
];
</script>

<template>
  <WhiteContainer class="justify-start gap-4">
    <div class="max-w-2xl rounded-2xl border-2 border-pink-600 bg-pink-200 p-2 px-4 text-pink-600">
      <span class="font-semibold">Caution!</span> You are currently working on the administrative
      panel of this application. Any changes here may be destructive.
    </div>

    <router-view></router-view>
  </WhiteContainer>

  <div class="fixed inset-x-0 bottom-8 z-50 flex items-center justify-center">
    <footer
      class="bg-claret-50/20 flex-row items-center justify-center rounded-full p-1 backdrop-blur-md"
    >
      <template v-for="link in links" :key="link.to">
        <router-link
          class="rounded-full px-4 py-2 text-xs font-semibold sm:text-sm"
          :class="{
            'bg-pink-200': link.name == $route.name,
            'hover:bg-pink-100': link.name != $route.name,
          }"
          :to="link.to"
        >
          {{ t(link.label) }}
        </router-link>
      </template>
    </footer>
  </div>
</template>
