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
        {
      name: "forgot-password",
      path: "/forgot-password",
      component: () => import("../pages/ForgotPasswordPage.vue"),
    },
  ],
});

export default router;
