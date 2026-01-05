<script setup lang="ts">
import { useProfileStore } from "@/stores/profile";
import type { Product } from "@/types";
import { LucidePencil, LucideX } from "lucide-vue-next";
import OverlayScreen from "../shared/OverlayScreen.vue";
import { computed, ref } from "vue";
import WYSIWYGInput from "../shared/inputs/WYSIWYGInput.vue";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { endpoints } from "@/consts";
import dayjs from "dayjs";
import { useI18n } from "vue-i18n";

const profile = useProfileStore();
const { locale } = useI18n();
const { authFetch } = useAuthFetch();

const props = defineProps<{
  data: Product & { similar_products?: Product[]; categories: { id: number; name: string }[] };
}>();

const emits = defineEmits<{
  reload: [];
}>();

const editDialogShown = ref(false);
const editLoading = ref(false);
const editError = ref<string>();
const descriptionChange = ref<string>();

const isExpired = computed(() => dayjs(props.data.expired_at).isBefore(dayjs()));

function makeLocalTime(timeString: string) {
  return dayjs(timeString).locale(locale.value).format("lll");
}

async function createDescriptionChange() {
  editLoading.value = true;
  editError.value = undefined;

  try {
    const res = await authFetch(endpoints.products.description(props.data.id), {
      method: "POST",
      body: JSON.stringify({
        description: descriptionChange.value,
      }),
    });

    if (res.ok) {
      emits("reload");
    }
  } catch {
    editError.value = "products.cant_update_description";
  } finally {
    editLoading.value = false;
  }
}
</script>

<template>
  <OverlayScreen :shown="editDialogShown" class="p-6">
    <div class="flex w-full max-w-xl flex-col gap-4 rounded-2xl bg-white p-6 shadow-md">
      <div class="flex w-full flex-row items-center justify-between">
        <h2 class="text-xl font-bold">{{ $t("products.add_remark") }}</h2>
        <button
          class="cursor-pointer rounded-full p-2 hover:bg-zinc-200"
          @click="editDialogShown = false"
        >
          <LucideX class="size-4 text-black" />
        </button>
      </div>

      <hr class="h-px w-full border-zinc-300" />

      <div class="flex w-full flex-col gap-4">
        <WYSIWYGInput :label="$t('products.description')" required v-model="descriptionChange" />
      </div>

      <div class="flex w-full flex-row items-center justify-end gap-2 font-semibold">
        <button
          type="submit"
          class="bg-claret-600 hover:bg-claret-700 cursor-pointer rounded-full px-4 py-1 text-white disabled:cursor-progress disabled:opacity-50"
          :disabled="editLoading && !isExpired"
          v-if="!isExpired"
          @click.prevent="createDescriptionChange"
        >
          {{ editLoading ? $t("products.loading") : $t("products.change_description") }}
        </button>
      </div>
    </div>
  </OverlayScreen>

  <section class="flex w-full flex-col gap-4">
    <div class="flex flex-row items-center gap-2">
      <h2 class="text-xl font-bold">{{ $t("products.description") }}</h2>

      <button
        class="cursor-pointer text-black/50 duration-200 hover:text-black"
        @click="editDialogShown = true"
      >
        <LucidePencil v-if="profile.profile?.id == data.seller.id" class="size-4" />
      </button>
    </div>

    <div class="prose w-full text-justify" v-html="data.description"></div>

    <div v-if="data.description_changes" class="flex w-full flex-col gap-4">
      <div
        v-for="change in data.description_changes"
        :key="change.id"
        class="flex w-full flex-col gap-2"
      >
        <span class="text-sm font-semibold">{{ makeLocalTime(change.created_at) }}</span>
        <div v-html="change.changes" class="prose"></div>
      </div>
    </div>
  </section>
</template>
