import { createApp } from 'vue';
import { createRouter, createWebHistory } from "vue-router";

import App from "./App.vue"
import Index from "./components/Index.vue";

const routes = [
  { path: "/", name:"index", component: Index },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

createApp(App).use(router).mount("#app");