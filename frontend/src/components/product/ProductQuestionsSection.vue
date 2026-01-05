<script setup lang="ts">
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { useProfileStore } from "@/stores/profile";
import type { Product } from "@/types";
import dayjs from "dayjs";
import { LucideReply } from "lucide-vue-next";
import { computed, ref } from "vue";
import AvatarCircle from "../shared/AvatarCircle.vue";
import ProductAnswerBlock from "./ProductAnswerBlock.vue";

const profile = useProfileStore();
const { authFetch } = useAuthFetch();

const props = defineProps<{
  data: Product & { similar_products?: Product[]; categories: { id: number; name: string }[] };
}>();

const emit = defineEmits<{
  refresh: [];
}>();

const question = ref<string>();
const loading = ref(false);
const isExpired = computed(() => dayjs(props.data.expired_at).isBefore(dayjs()));

async function ask() {
  loading.value = true;

  try {
    const res = await authFetch(endpoints.questions.index, {
      method: "POST",
      body: JSON.stringify({ product_id: props.data.id, content: question.value }),
    });
    if (res.ok) {
      emit("refresh");
    }
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <section class="flex w-full flex-col gap-4">
    <h2 class="text-xl font-bold">{{ $t("products.questions") }}</h2>

    <div class="flex w-full flex-col gap-2" v-if="data.questions && data.questions.length > 0">
      <div
        v-for="(question, idx) in data.questions"
        :key="idx"
        class="flex w-full flex-col gap-2 rounded-2xl border border-zinc-500 p-4"
      >
        <span class="flex flex-row items-center gap-2 text-sm font-semibold">
          <AvatarCircle
            :name="question.user.name"
            :avatar_url="question.user.avatar_url"
            class="size-8"
          />
          {{ $t("products.asked_by", { name: question.user.name }) }}
        </span>

        <p>{{ question.content }}</p>

        <div v-if="question.answer" class="flex flex-col gap-1 px-4">
          <span class="flex flex-row items-center gap-2 text-sm text-black/50">
            <LucideReply class="size-4" />

            {{ $t("products.replied_by_seller") }}
          </span>

          <p>{{ question.answer }}</p>
        </div>

        <ProductAnswerBlock
          v-if="!question.answer && profile.profile?.id == data.seller.id && !isExpired"
          :question
          @refresh="$emit('refresh')"
        />
      </div>
    </div>
    <p v-else class="w-full px-4 py-6 text-center">
      {{ $t("products.no_questions") }}
    </p>

    <!-- Part to ask a question -->
    <div
      class="flex w-full flex-col items-center justify-center gap-2 sm:flex-row"
      v-if="profile.profile && data.seller.id != profile.profile.id"
    >
      <div class="flex w-full flex-row items-center gap-2">
        <AvatarCircle :name="profile.profile.name" :avatar_url="profile.profile.avatar_url" />

        <input
          type="text"
          required
          v-model="question"
          :placeholder="$t('products.ask_a_question')"
          class="focus:ring-claret-600 h-full w-full rounded-xl border border-zinc-300 px-4 py-2 duration-200 outline-none hover:border-zinc-500 focus:ring-2"
        />
      </div>

      <button
        class="bg-claret-600 enabled:hover:text-claret-600 border-claret-600 h-full w-full min-w-fit cursor-pointer rounded-xl border-2 px-4 py-2 font-semibold text-white duration-200 enabled:hover:bg-white disabled:cursor-not-allowed disabled:opacity-50 sm:w-fit"
        :disabled="loading"
        @click="ask"
      >
        {{ loading ? $t("products.asking") : $t("products.ask") }}
      </button>
    </div>
  </section>
</template>
