<script setup lang="ts">
import LoginForm from "@/components/login/LoginForm.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { useProfileStore } from "@/stores/profile";
import { useHead } from "@unhead/vue";
import { useI18n } from "vue-i18n";

const profile = useProfileStore();
const { t } = useI18n();

useHead({
  title: t("meta.login"),
});
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
