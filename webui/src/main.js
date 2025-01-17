import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import axios from "./services/axios.js";

import "./assets/main.css";
import "bootstrap-icons/font/bootstrap-icons.css";

const app = createApp(App);
app.config.globalProperties.$axios = axios;
app.use(router);
app.mount("#app");
