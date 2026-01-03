<script setup lang="ts">
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import type { Question } from "@/types";
import { LucideLoaderCircle, LucideReply } from "lucide-vue-next";
import { ref } from "vue";

const { authFetch } = useAuthFetch();

const props = defineProps<{
  question: Question;
}>();

const emits = defineEmits<{
  refresh: [];
}>();

const answer = ref<string>();
const loading = ref(false);

async function replyToQuestion() {
  loading.value = true;

  try {
    const res = await authFetch(endpoints.questions.id(props.question.id), {
      method: "PUT",
      body: JSON.stringify({ answer: answer.value }),
    });
    if (res.ok) {
      emits("refresh");
    }
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div
    class="focus-within:border-claret-600 flex w-full flex-row items-center gap-2 rounded-lg border border-zinc-300 px-2 py-2 duration-200 outline-none hover:border-zinc-400"
  >
    <button
      :disabled="!answer || loading"
      class="rounded-full p-2 enabled:cursor-pointer enabled:hover:bg-zinc-300 disabled:cursor-not-allowed disabled:text-black/50"
      @click="replyToQuestion"
    >
      <LucideLoaderCircle v-if="loading" class="size-4 animate-spin" />
      <LucideReply v-else class="size-4" />
    </button>

    <input
      type="text"
      :placeholder="$t('products.answer_question')"
      class="w-full outline-none"
      v-model="answer"
    />
  </div>
</template>
