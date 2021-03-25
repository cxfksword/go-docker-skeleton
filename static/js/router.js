const router = new VueRouter({
  routes: [
    { path: '/', component: httpVueLoader('/public/view/index.vue'), name: '首页'},
  ]
})

const app = new Vue({
  router
}).$mount('#app')
