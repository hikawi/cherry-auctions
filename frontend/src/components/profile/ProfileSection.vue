<script setup lang="ts">
import { useProfileStore } from "@/stores/profile";
import AvatarCircle from "../shared/AvatarCircle.vue";
import { onUnmounted, ref } from "vue";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { endpoints } from "@/consts";

const profile = useProfileStore();
const { authFetch } = useAuthFetch({ json: false });

const fileInput = ref<HTMLInputElement | null>(null);
const filePicked = ref<File | null>(null);
const previewUrl = ref<string>();
const loading = ref(false);

const name = ref<string | undefined>(profile.profile?.name);
const address = ref<string | undefined>(profile.profile?.address);

// Trigger the hidden file input when the styled button is clicked
function triggerPicker() {
  fileInput.value?.click();
}

const onFileChange = (e: Event) => {
  const files = (e.target as HTMLInputElement).files;
  if (files && files[0]) {
    if (previewUrl.value) {
      URL.revokeObjectURL(previewUrl.value);
    }

    filePicked.value = files[0];
    previewUrl.value = URL.createObjectURL(filePicked.value);
  }
};

async function updateAvatar() {
  if (!filePicked.value) {
    return;
  }

  loading.value = true;

  try {
    const formData = new FormData();
    formData.append("avatar", filePicked.value);

    const res = await authFetch(endpoints.users.avatar, {
      method: "POST",
      body: formData,
    });
    if (res.ok) {
      profile.fetchProfile();
    }
  } finally {
    loading.value = false;
  }
}

async function updateProfile() {
  loading.value = true;

  try {
    const res = await authFetch(endpoints.users.profile, {
      method: "PUT",
      headers: {
        "content-type": "application/json",
      },
      body: JSON.stringify({ name: name.value, address: address.value }),
    });
    if (res.ok) {
      await profile.fetchProfile();
    }
  } finally {
    loading.value = false;
  }
}

onUnmounted(() => {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value);
  }
});
</script>

<template>
  <h2 class="text-2xl font-semibold">{{ $t("profile.profile") }}</h2>

  <div class="flex w-full flex-col items-center gap-4 sm:flex-row sm:items-start sm:gap-8">
    <!-- Image thing -->
    <div class="flex w-fit flex-col gap-2">
      <AvatarCircle
        class="h-auto max-w-80 min-w-64"
        hover
        :name="profile.profile?.name"
        :avatar_url="profile.profile?.avatar_url"
        v-if="!filePicked"
      />
      <img
        :src="previewUrl"
        class="aspect-square h-auto w-full max-w-64 rounded-full object-cover"
        v-else
      />

      <input ref="fileInput" type="file" class="hidden" accept="image/*" @change="onFileChange" />

      <button
        class="w-full min-w-fit cursor-pointer rounded-xl border-2 border-zinc-500 px-4 py-2 text-black duration-200 hover:bg-zinc-200"
        type="submit"
        @click.prevent="triggerPicker"
      >
        {{ $t("profile.upload_avatar") }}
      </button>

      <button
        class="border-claret-500 hover:bg-claret-200 w-full cursor-pointer rounded-xl border-2 px-4 py-2 text-black duration-200"
        type="submit"
        @click.prevent="updateAvatar"
        v-if="filePicked"
      >
        {{ $t("general.confirm") }}
      </button>
    </div>

    <div class="flex w-full flex-col gap-2">
      <label class="flex w-full flex-col gap-2">
        {{ $t("profile.email") }}

        <input
          disabled
          readonly
          type="email"
          :value="profile.profile?.email"
          class="hover:ring-claret-200 focus:ring-claret-600 w-full rounded-lg border border-zinc-300 px-4 py-2 duration-200 outline-none placeholder:text-black/50 hover:ring-2 focus:ring-2 disabled:cursor-not-allowed disabled:bg-zinc-200"
        />
      </label>

      <label class="flex w-full flex-col gap-2">
        {{ $t("profile.name") }}

        <input
          type="text"
          name="name"
          v-model="name"
          class="hover:ring-claret-200 focus:ring-claret-600 w-full rounded-lg border border-zinc-300 px-4 py-2 duration-200 outline-none placeholder:text-black/50 hover:ring-2 focus:ring-2 disabled:cursor-not-allowed disabled:bg-zinc-200"
        />
      </label>

      <label class="flex w-full flex-col gap-2">
        {{ $t("profile.address") }}

        <input
          type="text"
          name="address"
          v-model="address"
          class="hover:ring-claret-200 focus:ring-claret-600 w-full rounded-lg border border-zinc-300 px-4 py-2 duration-200 outline-none placeholder:text-black/50 hover:ring-2 focus:ring-2 disabled:cursor-not-allowed disabled:bg-zinc-200"
        />
      </label>

      <button
        class="bg-claret-600 hover:bg-claret-700 mt-2 flex w-fit cursor-pointer flex-row items-center justify-center gap-2 self-end rounded-full px-4 py-1 font-semibold text-white duration-200"
        @click="updateProfile"
      >
        {{ loading ? $t("general.loading") : $t("general.confirm") }}
      </button>
    </div>
  </div>
</template>
