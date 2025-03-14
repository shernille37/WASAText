import { createRouter, createWebHashHistory } from "vue-router";
import Login from "../views/Login.vue";
import Home from "../views/Home.vue";
import MainDashboard from "../views/MainDashboard.vue";
import Profile from "../views/Profile.vue";

import { authStore } from "../stores/authStore";
import EditGroupConversation from "../views/EditGroupConversation.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
    children: [
      {
        path: "",
        name: "main-dashboard",
        component: MainDashboard,
      },
      {
        name: "profile",
        path: "profile",
        component: Profile,
      },
      {
        name: "EditGroupConversation",
        path: "/group-conversations/:conversationID/edit",
        component: EditGroupConversation,
      },
    ],
    meta: { private: true },
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
  },
];

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes,
});

router.beforeEach((to, from, next) => {
  // Check if the user is authenticated
  const isAuthenticated = !!authStore.user.data;

  if (to.meta.private && !isAuthenticated) {
    next({ path: "/login" });
  } else if (to.path === "/login" && isAuthenticated) {
    next({ path: from.path });
  } else {
    next();
  }
});

export default router;
