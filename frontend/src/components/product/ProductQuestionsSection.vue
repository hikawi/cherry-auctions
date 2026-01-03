<script setup lang="ts">
import { useProfileStore } from "@/stores/profile";
import type { Product } from "@/types";
import { LucideReply } from "lucide-vue-next";
import AvatarCircle from "../shared/AvatarCircle.vue";
import ProductAnswerBlock from "./ProductAnswerBlock.vue";

const profile = useProfileStore();

defineProps<{
  data: Product & { similar_products?: Product[]; categories: { id: number; name: string }[] };
}>();

defineEmits<{
  refresh: [];
}>();
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
          v-if="!question.answer && profile.profile?.id == data.seller.id"
          :question
          @refresh="$emit('refresh')"
        />
      </div>
    </div>
    <p v-else class="w-full px-4 py-6 text-center">
      {{ $t("products.no_questions") }}
    </p>

    <!-- Part to ask a question -->
    <div class="flex w-full" v-if="data.seller.id != profile.profile?.id"></div>
  </section>
</template>
