<script setup lang="ts">
import CreateAuctionDialog from "@/components/auctions/CreateAuctionDialog.vue";
import NavigationBar from "@/components/shared/NavigationBar.vue";
import OverlayScreen from "@/components/shared/OverlayScreen.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { LucidePackage } from "lucide-vue-next";
import { ref } from "vue";

const createDialogShown = ref(false);

function onCreate(status: number) {
  createDialogShown.value = false;
  if (status != 201) {
  }
}
</script>

<template>
  <WhiteContainer class="justify-start gap-8">
    <NavigationBar />

    <OverlayScreen :shown="createDialogShown">
      <CreateAuctionDialog @close="createDialogShown = false" @status="onCreate" />
    </OverlayScreen>

    <section class="flex w-full max-w-4xl flex-col gap-4">
      <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
        <h2 class="text-2xl font-bold">{{ $t("auctions.title") }}</h2>

        <button
          class="bg-claret-600 hover:bg-claret-700 mt-2 flex w-fit cursor-pointer flex-row items-center justify-center gap-2 self-end rounded-full px-4 py-1 font-semibold text-white duration-200"
          @click="createDialogShown = true"
        >
          <LucidePackage class="size-4 text-white" />
          {{ $t("auctions.new") }}
        </button>
      </div>
    </section>
  </WhiteContainer>
</template>
