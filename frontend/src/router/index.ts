import HomePage from "@/pages/HomePage.vue";
import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      name: "home",
      path: "/",
      component: HomePage,
    },
    {
      name: "login",
      path: "/login",
      component: () => import("../pages/LoginPage.vue"),
    },
    {
      name: "register",
      path: "/register",
      component: () => import("../pages/RegisterPage.vue"),
    },
  ],
});

export default router;
