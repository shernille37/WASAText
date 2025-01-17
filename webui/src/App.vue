<script setup>
import { RouterView } from "vue-router";
import Navbar from "./components/Navbar.vue";
import Sidebar from "./components/Sidebar.vue";
import ChatMessages from "./components/ChatMessages.vue";
</script>
<script>
export default {
  components: {
    Navbar,
    Sidebar,
    ChatMessages,
  },
};
</script>

<template>
  <RouterView />
</template>
