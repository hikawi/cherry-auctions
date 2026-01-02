<script setup lang="ts">
import type { Product } from "@/types";

defineProps<{
  data: Product & { similar_products?: Product[]; categories: { id: number; name: string }[] };
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
      </div>
    </div>
    <p v-else class="w-full px-4 py-6 text-center">
      {{ $t("products.no_questions") }}
    </p>
  </section>
</template>
