<script setup lang="ts">
import LoginForm from "@/components/login/LoginForm.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { useProfileStore } from "@/stores/profile";

const profile = useProfileStore();
</script>

<template>
  <WhiteContainer>
    <div
      v-if="profile.profile && profile.hasProfile"
      class="bg-watermelon-100 border-watermelon-500 text-watermelon-500 max-w-lg rounded-xl border-2 px-4 py-2"
    >
      {{
        $t("login.already_logged_in", { name: profile.profile.name, email: profile.profile.email })
      }}
    </div>

    <LoginForm />
  </WhiteContainer>
</template>
