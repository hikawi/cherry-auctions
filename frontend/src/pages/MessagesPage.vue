<script setup lang="ts">
import NavigationBar from "@/components/shared/NavigationBar.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { useTokenStore } from "@/stores/token";
import type { ChatMessage, ChatSession } from "@/types";
import { computed, onMounted, onUnmounted, ref, useTemplateRef, watch } from "vue";
import { LucideSend, LucideImage } from "lucide-vue-next";
import { useProfileStore } from "@/stores/profile";

const { authFetch, tryRefresh } = useAuthFetch({ json: false });
const authToken = useTokenStore();
const profile = useProfileStore();

const chatSessions = ref<ChatSession[]>([]);
const currentChatSession = ref<number>(0);
const currentChatMessages = ref<ChatMessage[]>([]);
const message = ref("");
const es = ref<EventSource>();
const image = ref<Blob>();

const hiddenInput = useTemplateRef<HTMLInputElement>("hiddenInput");
const currentSessionData = computed(() =>
  chatSessions.value.length > 0
    ? chatSessions.value.filter((sess) => currentChatSession.value == sess.id)[0]
    : undefined,
);

onMounted(async () => {
  await fetchChatSessions();
  await createSSEStream();
});

watch(currentChatSession, (cb) => fetchChatMessages(cb));

// Courtesy of Gemini. Thank you!
async function createSSEStream() {
  // 1. Clean up any existing connection before starting a new one
  if (es.value) {
    es.value.close();
  }

  const url = `${endpoints.chat.stream}?token=${authToken.token}`;
  es.value = new EventSource(url);

  es.value.onopen = () => {
    console.log("SSE stream connected");
  };

  // 2. Automated Reconnection with fresh token
  es.value.onerror = async () => {
    console.warn("SSE Connection lost. Attempting to refresh token and reconnect...");
    es.value?.close();

    // Optional: Add a small delay so you don't spam the server if it's down
    setTimeout(async () => {
      // Assuming you have a method to refresh the token if it expired
      // await authStore.refreshToken();
      const refresh = await tryRefresh();
      if (refresh) {
        createSSEStream();
      } else {
        console.error("logged out");
      }
    }, 3000);
  };

  // 3. Proper Message Handling
  es.value.onmessage = (event) => {
    try {
      const json = JSON.parse(event.data);

      if (json.chat_session_id === currentChatSession.value) {
        currentChatMessages.value.push(json);
      }
    } catch (err) {
      console.error("Failed to parse SSE message:", err);
    }
  };
}

onUnmounted(() => {
  if (es.value) {
    console.log("Cleaning up SSE connection...");
    es.value.close();
  }
});

function onInputChange(e: Event) {
  const el = e.target as HTMLInputElement;
  if (el.files && el.files[0]) {
    image.value = el.files[0];
  }
  sendMessage();
}

async function fetchChatSessions() {
  const res = await authFetch(endpoints.chat.index);
  if (res.ok) {
    const json = await res.json();
    chatSessions.value = json.data || [];
  }
}

async function fetchChatMessages(sessionID: number) {
  if (sessionID <= 0) {
    return;
  }

  const res = await authFetch(endpoints.chat.id(sessionID));
  if (res.ok) {
    const json = await res.json();
    currentChatMessages.value = json.data.reverse() || [];
  }
}

async function sendMessage() {
  const formData = new FormData();
  formData.append("content", message.value);
  if (image.value) {
    formData.append("image", image.value);
  }

  message.value = "";
  image.value = undefined;

  await authFetch(endpoints.chat.id(currentChatSession.value), {
    method: "POST",
    body: formData,
  });
}

function openImage() {
  hiddenInput.value?.click();
}
</script>

<template>
  <WhiteContainer class="justify-start gap-6 overflow-y-auto">
    <NavigationBar />
    <input ref="hiddenInput" class="hidden" type="file" accept="image/*" @change="onInputChange" />

    <div
      class="grid h-full max-h-full min-h-0 w-full grid-cols-1 gap-4 lg:grid-cols-3 xl:grid-cols-4"
    >
      <!-- Chat Sessions List -->
      <aside class="flex w-full flex-col gap-4">
        <template v-for="session in chatSessions" :key="session.id">
          <button
            class="flex flex-row gap-2 rounded-xl border border-zinc-300 hover:cursor-pointer hover:border-zinc-500"
            @click="currentChatSession = session.id"
          >
            <img
              :src="session.product.thumbnail_url"
              class="aspect-square h-20 rounded-l-xl object-cover object-center"
            />

            <div class="flex flex-col items-start py-4">
              <span class="text-base font-semibold"
                >{{ session.seller.name }}, {{ session.buyer.name }}</span
              >
              <span class="text-sm">{{ session.product.name }}</span>
            </div>
          </button>
        </template>
      </aside>

      <!-- Main chat session -->
      <div
        class="flex h-full min-h-0 w-full flex-col gap-4 lg:col-span-2 xl:col-span-3"
        v-if="currentSessionData"
      >
        <div class="w-full">Lol</div>

        <div class="flex h-full min-h-0 flex-1 flex-col justify-end gap-2 overflow-y-auto">
          <div class="flex w-fit flex-col self-center rounded-xl">
            <img :src="currentSessionData.product.thumbnail_url" class="aspect-video w-64" />
          </div>

          <template v-for="chatMsg in currentChatMessages" :key="chatMsg.id">
            <div
              :class="[
                'flex h-fit min-h-10 w-fit max-w-80 flex-col gap-2 rounded-xl',
                chatMsg.sender.id == profile.profile?.id
                  ? 'bg-claret-600 self-end px-4 py-2 text-white'
                  : 'self-start bg-zinc-200 px-4 py-2 text-black',
              ]"
            >
              <img
                v-if="chatMsg.image_url"
                :src="chatMsg.image_url"
                class="block h-auto w-full rounded-lg object-contain"
              />
              <span v-if="chatMsg.content">{{ chatMsg.content }}</span>
            </div>
          </template>
        </div>

        <div
          class="flex h-fit w-full flex-row items-center justify-center gap-2 rounded-xl border border-zinc-300 px-4 py-2 focus-within:border-zinc-500 hover:border-zinc-500"
        >
          <button class="cursor-pointer" @click="openImage">
            <LucideImage class="size-5 text-black/50 hover:text-black" />
          </button>
          <input type="text" class="w-full outline-none" v-model="message" />
          <button class="cursor-pointer" @click="sendMessage">
            <LucideSend class="size-5 text-black/50 hover:text-black" />
          </button>
        </div>
      </div>
      <div class="flex size-full items-center justify-center lg:col-span-2 xl:col-span-3" v-else>
        {{ $t("messages.choose_session") }}
      </div>
    </div>
  </WhiteContainer>
</template>
