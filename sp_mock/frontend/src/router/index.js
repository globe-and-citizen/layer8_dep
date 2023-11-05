import { createRouter, createWebHistory } from "vue-router";
import Home from "../views/Home.vue";
import Callback from "../views/CallBack.vue";


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: Home,
    },
    {
      path: "/stress-test",
      name: "stress-test",
      component: () => import("../views/StressTest.vue"),
    },
    {
      path: '/oauth2/callback',
      name: 'oauth2-callback',
      component: Callback,
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