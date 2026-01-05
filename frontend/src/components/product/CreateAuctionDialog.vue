<script setup lang="ts">
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import type { Category } from "@/types";
import { LucideX } from "lucide-vue-next";
import { computed, ref } from "vue";
import z from "zod";
import CheckboxInput from "../shared/inputs/CheckboxInput.vue";
import MoneyInput from "../shared/inputs/MoneyInput.vue";
import MultiImageInput from "../shared/inputs/MultiImageInput.vue";
import SelectInput from "../shared/inputs/SelectInput.vue";
import TextInput from "../shared/inputs/TextInput.vue";
import WYSIWYGInput from "../shared/inputs/WYSIWYGInput.vue";

const { authFetch } = useAuthFetch({ json: false });

const props = defineProps<{
  loading?: boolean;
  categories?: Category[];
}>();

const emits = defineEmits<{
  close: [];
  status: [code: number];
}>();

const allowedFileTypes = ["image/png", "image/jpeg", "image/webp"];

const name = ref<string>();
const description = ref<string>();
const productImages = ref<File[]>([]);
const categoriesSelected = ref<string[]>([]);
const startingBid = ref<number>();
const stepBidValue = ref<number>();
const binPrice = ref<number>();
const allowUnrated = ref<boolean>(true);
const autoExtends = ref<boolean>(true);
const expiredAt = ref<string>();

const error = ref<string>();
const categoryOptions = computed(() => {
  const options: { value: string; label: string }[] = [];

  function traverse(categories: Category[]) {
    for (const cat of categories) {
      options.push({ value: cat.id.toString(), label: cat.name });
      traverse(cat.subcategories);
    }
  }

  if (props.categories) {
    traverse(props.categories);
  }
  return options;
});

async function confirm() {
  error.value = "";

  const schema = z.object({
    name: z.string(),
    description: z.string().min(50),
    product_images: z.array(z.file().refine((file) => allowedFileTypes.includes(file.type))).min(3),
    categories: z.array(z.coerce.number()).min(1),
    starting_bid: z.coerce.number().min(0),
    step_bid_value: z.coerce.number().min(0),
    bin_price: z.coerce.number().min(0),
    allows_unrated: z.coerce.boolean().default(true),
    auto_extends: z.coerce.boolean().default(true),
    expired_at: z.coerce.date().min(new Date()),
  });

  const data = schema.safeParse({
    name: name.value,
    description: description.value,
    product_images: productImages.value,
    categories: categoriesSelected.value,
    starting_bid: startingBid.value,
    step_bid_value: stepBidValue.value,
    bin_price: binPrice.value,
    allows_unrated: allowUnrated.value,
    auto_extends: autoExtends.value,
    expired_at: expiredAt.value,
  });

  if (data.error) {
    error.value = `auctions.error_${data.error.issues[0].path}`;
    return;
  }

  const formData = new FormData();
  formData.append("name", data.data.name);
  formData.append("description", data.data.description);
  data.data.product_images.forEach((file) => {
    formData.append("product_images", file);
  });
  data.data.categories.forEach((cat) => formData.append("categories", cat.toString()));
  formData.append("starting_bid", data.data.starting_bid.toString());
  formData.append("step_bid_value", data.data.step_bid_value.toString());
  formData.append("bin_price", data.data.bin_price.toString());
  formData.append("allows_unrated", data.data.allows_unrated.toString());
  formData.append("auto_extends", data.data.auto_extends.toString());
  formData.append("expired_at", data.data.expired_at.toISOString());

  const res = await authFetch(endpoints.products.get, {
    method: "POST",
    body: formData,
  });
  emits("status", res.status);
}
</script>

<template>
  <div
    class="flex max-h-4/5 w-full max-w-xl flex-col gap-4 overflow-x-visible overflow-y-scroll rounded-2xl bg-white p-6 shadow-md"
  >
    <div class="flex w-full flex-row items-center justify-between">
      <h2 class="text-xl font-bold">{{ $t("auctions.new") }}</h2>
      <button class="cursor-pointer rounded-full p-2 hover:bg-zinc-200" @click="$emit('close')">
        <LucideX class="size-4 text-black" />
      </button>
    </div>

    <hr class="h-px w-full border-zinc-300" />

    <form class="flex w-full flex-col gap-4">
      <div class="flex w-full flex-col gap-4">
        <TextInput :label="$t('auctions.product_name')" type="text" required v-model="name" />

        <WYSIWYGInput :label="$t('auctions.product_description')" required v-model="description" />

        <MultiImageInput :label="$t('auctions.product_images')" required v-model="productImages" />

        <SelectInput
          :options="categoryOptions"
          :label="$t('auctions.categories')"
          @add="(id) => (categoriesSelected = [...categoriesSelected, id])"
          @remove="(id) => (categoriesSelected = categoriesSelected.filter((val) => val != id))"
        />

        <MoneyInput
          :label="$t('auctions.starting_bid')"
          placeholder="9.99"
          required
          @change="(n) => (startingBid = n)"
        />

        <MoneyInput
          :label="$t('auctions.step_bid_value')"
          placeholder="2.0"
          required
          @change="(n) => (stepBidValue = n)"
        />

        <MoneyInput
          :label="$t('auctions.bin_price')"
          placeholder="100.0"
          required
          @change="(n) => (binPrice = n)"
        />

        <CheckboxInput :label="$t('auctions.allows_unrated_buyers')" v-model="allowUnrated" />
        <CheckboxInput :label="$t('auctions.auto_extends')" v-model="autoExtends" />

        <TextInput
          :label="$t('auctions.expired_at')"
          v-model="expiredAt"
          type="datetime-local"
          required
        />
      </div>

      <div
        v-if="error"
        class="border-watermelon-600 bg-watermelon-100 text-watermelon-600 w-full rounded-lg border-2 px-4 py-2 font-semibold"
      >
        {{ $t(error) }}
      </div>

      <div class="flex w-full flex-row items-center justify-end gap-2 font-semibold">
        <button
          type="reset"
          class="cursor-pointer rounded-full bg-zinc-200 px-4 py-1 hover:bg-zinc-300"
          @click="$emit('close')"
        >
          {{ $t("auctions.cancel") }}
        </button>

        <button
          type="submit"
          class="bg-claret-600 hover:bg-claret-700 cursor-pointer rounded-full px-4 py-1 text-white disabled:cursor-progress disabled:opacity-50"
          :disabled="loading"
          @click.prevent="confirm"
        >
          {{ loading ? $t("auctions.loading") : $t("auctions.create") }}
        </button>
      </div>
    </form>
  </div>
</template>
