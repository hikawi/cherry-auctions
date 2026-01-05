import HomePage from "@/pages/HomePage.vue";
import { useProfileStore } from "@/stores/profile";
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
      name: "products-home",
      path: "/products",
      redirect: {
        name: "search",
      },
    },
    {
      name: "product-details",
      path: "/products/:id",
      component: () => import("../pages/ProductDetailsPage.vue"),
    },
    {
      name: "auctions",
      path: "/auctions",
      component: () => import("../pages/AuctionsPage.vue"),
      meta: {
        requiresAuth: true,
      },
    },
    {
      name: "bids",
      path: "/bids",
      component: () => import("../pages/BidsPage.vue"),
      meta: {
        requiresAuth: true,
      },
    },
    {
      name: "profile",
      path: "/profile",
      component: () => import("../pages/ProfilePage.vue"),
      meta: {
        requiresAuth: true,
      },
    },
    {
      name: "subscriptions",
      path: "/subscriptions",
      component: () => import("../pages/SubscriptionPage.vue"),
      meta: {
        requiresAuth: true,
      },
    },
    {
      name: "search",
      path: "/search",
      component: () => import("../pages/SearchPage.vue"),
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
      name: "logout",
      path: "/logout",
      component: () => import("../pages/LogoutPage.vue"),
    },
    {
      name: "forgot",
      path: "/forgot",
      component: () => import("../pages/ForgotPasswordPage.vue"),
    },
    {
      name: "acknowledgements",
      path: "/acknowledgements",
      component: () => import("../pages/AcknowledgementPage.vue"),
    },
    {
      name: "admin",
      path: "/admin",
      meta: {
        requiresAuth: true,
      },
      component: () => import("../pages/admin/AdminIndexPage.vue"),
      children: [
        {
          meta: {
            requiresAuth: true,
          },
          name: "admin-categories",
          path: "/admin/categories",
          component: () => import("../pages/admin/AdminCategoriesPage.vue"),
        },
        {
          meta: {
            requiresAuth: true,
          },
          name: "admin-users",
          path: "/admin/users",
          component: () => import("../pages/admin/AdminUsersPage.vue"),
        },
      ],
    },
    {
      name: "403",
      path: "/403",
      component: () => import("../pages/403Page.vue"),
    },
    {
      name: "500",
      path: "/500",
      component: () => import("../pages/500Page.vue"),
    },
    {
      path: "/:pathMatch(.*)*",
      component: () => import("../pages/404Page.vue"),
    },
  ],
});

router.beforeEach(async (to) => {
  const profile = useProfileStore();

  if (!profile.hasFetched) {
    await profile.fetchProfile();
  }

  // 未ログイン
  if (!profile.profile && to.meta.requiresAuth) {
    return { name: "login" };
  }

  // 管理者専用
  if (to.path.startsWith("/admin") && !profile.isAdmin) {
    return { name: "403" };
  }
});

export default router;
