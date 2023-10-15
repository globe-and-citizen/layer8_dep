import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import CallbackView from "../views/CallbackView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
    },
    // {
    //   path: "/login",
    //   name: "login",
    //   component: () => import("../views/LoginView.vue"),
    // },
    {
      path: '/oauth2/callback',
      name: 'oauth2-callback',
      component: CallbackView,
    },
    // {
    //   path: '/poems/:id',
    //   name: 'poem',
    //   component: () => import('../views/PoemView.vue')
    // }
  ],
});

export default router;