<script>
import { RouterView } from "vue-router";
import Navbar from "../components/Navbar.vue";

export default {
  name: "Home",
  components: {
    Navbar,
  },
};
</script>

<template>
  <Navbar />
  <RouterView />
</template>
