<script setup lang="ts">
import NavigationBar from "@/components/shared/NavigationBar.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { useTokenStore } from "@/stores/token";
import type { ChatMessage, ChatSession } from "@/types";
import { computed, nextTick, onMounted, onUnmounted, ref, useTemplateRef, watch } from "vue";
import { LucideSend, LucideImage } from "lucide-vue-next";
import { useProfileStore } from "@/stores/profile";
import OverlayScreen from "@/components/shared/OverlayScreen.vue";
import ErrorDialog from "@/components/shared/ErrorDialog.vue";

const { authFetch, tryRefresh } = useAuthFetch({ json: false });
const authToken = useTokenStore();
const profile = useProfileStore();

const chatSessions = ref<ChatSession[]>([]);
const currentChatSession = ref<number>(0);
const currentChatMessages = ref<ChatMessage[]>([]);
const message = ref("");
const es = ref<EventSource>();
const image = ref<Blob>();
const error = ref("");

const hiddenInput = useTemplateRef<HTMLInputElement>("hiddenInput");
const messageCon = useTemplateRef<HTMLDivElement>("messageContainer");
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
watch(
  currentChatMessages,
  async () => {
    await nextTick();
    if (messageCon.value) {
      messageCon.value.scrollTo({
        top: messageCon.value.scrollHeight,
        behavior: "smooth",
      });
    }
  },
  { deep: true },
);

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

  es.value.addEventListener("transaction", (event) => {
    const json = JSON.parse(event.data);
    chatSessions.value = chatSessions.value.map((sess) => {
      if (sess.id != json.chat_session_id) {
        return sess;
      }
      if (!sess.product.transaction) {
        return sess;
      }

      return {
        ...sess,
        product: {
          ...sess.product,
          transaction: {
            ...sess.product.transaction,
            transaction_status: json.transaction_status,
          },
        },
      };
    });
  });

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

// Proceed the transaction with a status
async function proceed(status: string) {
  if (status == "start") {
    await authFetch(endpoints.transactions.index, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ product_id: currentSessionData.value?.product.id }),
    });
  } else {
    await authFetch(endpoints.transactions.id(currentSessionData.value?.product.transaction?.id), {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ status }),
    });
  }
}
</script>

<template>
  <OverlayScreen shown v-if="error">
    <ErrorDialog :title="$t('messages.error')" :description="error" @close="error = ''" />
  </OverlayScreen>

  <WhiteContainer class="min-h-0 justify-start gap-6">
    <NavigationBar />
    <input ref="hiddenInput" class="hidden" type="file" accept="image/*" @change="onInputChange" />

    <div
      class="grid h-[calc(100vh-200px)] min-h-0 w-full grid-cols-1 gap-4 lg:grid-cols-3 xl:grid-cols-4"
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

            <div class="flex flex-col items-start justify-center">
              <span class="text-base font-semibold"
                >{{ session.seller.name }}, {{ session.buyer.name }}</span
              >
              <span class="text-sm">{{ session.product.name }}</span>
              <span class="text-sm italic" v-if="session.product.transaction">{{
                $t(`messages.transaction_status_${session.product.transaction.transaction_status}`)
              }}</span>
            </div>
          </button>
        </template>
      </aside>

      <!-- Main chat session -->
      <div
        class="flex h-full min-h-0 w-full flex-col gap-4 lg:col-span-2 xl:col-span-3"
        v-if="currentSessionData"
      >
        <!-- Transaction card -->
        <div class="flex w-full flex-row items-center justify-between">
          <span>
            <span class="font-semibold">
              {{ $t("messages.transaction_status") }}
            </span>
            :
            {{
              $t(
                `messages.transaction_status_${currentSessionData.product.transaction?.transaction_status || "none"}`,
              )
            }}
          </span>

          <div class="flex flex-row items-center gap-2">
            <button
              v-if="
                currentSessionData.product.transaction?.transaction_status == 'pending' &&
                profile.profile?.id == currentSessionData.product.transaction.seller_id
              "
              class="flex rounded-full bg-zinc-200 px-2 py-1 font-semibold text-black hover:bg-zinc-300"
              @click="proceed('cancel')"
            >
              {{ $t("messages.transaction_cancel") }}
            </button>

            <button
              class="bg-claret-600 hover:bg-claret-700 flex rounded-full px-2 py-1 font-semibold text-white"
              v-if="
                currentSessionData.seller.id == profile.profile?.id &&
                currentSessionData.product.transaction == null
              "
              @click="proceed('start')"
            >
              {{ $t("messages.transaction_start") }}
            </button>

            <button
              class="bg-claret-600 hover:bg-claret-700 flex rounded-full px-2 py-1 font-semibold text-white"
              v-if="
                currentSessionData.seller.id == profile.profile?.id &&
                currentSessionData.product.transaction?.transaction_status == 'pending'
              "
              @click="proceed('paid')"
            >
              {{ $t("messages.transaction_paid") }}
            </button>

            <button
              class="bg-claret-600 hover:bg-claret-700 flex rounded-full px-2 py-1 font-semibold text-white"
              v-if="
                currentSessionData.seller.id == profile.profile?.id &&
                currentSessionData.product.transaction?.transaction_status == 'paid'
              "
              @click="proceed('delivered')"
            >
              {{ $t("messages.transaction_deliver") }}
            </button>

            <button
              class="bg-claret-600 hover:bg-claret-700 flex rounded-full px-2 py-1 font-semibold text-white"
              v-if="
                currentSessionData.buyer.id == profile.profile?.id &&
                currentSessionData.product.transaction?.transaction_status == 'delivered'
              "
              @click="proceed('completed')"
            >
              {{ $t("messages.transaction_complete") }}
            </button>
          </div>
        </div>

        <hr class="w-full border-zinc-300" />

        <div class="flex min-h-0 flex-col gap-2 overflow-y-auto" ref="messageContainer">
          <!-- Spacer to allow inner scrolling -->
          <div class="flex-1"></div>

          <div class="flex w-fit flex-col self-center rounded-xl">
            <img :src="currentSessionData.product.thumbnail_url" class="aspect-video w-64" />
          </div>

          <template v-for="chatMsg in currentChatMessages" :key="chatMsg.id">
            <div
              :class="[
                'flex min-h-10 w-fit max-w-80 shrink-0 flex-col gap-2 rounded-xl',
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
