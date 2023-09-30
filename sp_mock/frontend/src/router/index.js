import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";

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
    // {
    //   path: '/oauth2/callback',
    //   name: 'oauth2-callback',
    //   component: () => import('../views/CallbackView.vue')
    // },
    // {
    //   path: '/poems/:id',
    //   name: 'poem',
    //   component: () => import('../views/PoemView.vue')
    // }
  ],
});

export default router;