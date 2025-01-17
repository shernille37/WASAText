<script>
import Navbar from "../components/Navbar.vue";
import Sidebar from "../components/Sidebar.vue";
import ChatMessages from "../components/ChatMessages.vue";

export default {
  name: "Home",
  components: {
    Navbar,
    Sidebar,
    ChatMessages,
  },
  data() {
    return {
      selectedConversation: null,
    };
  },
  methods: {
    selectConveration(conversationID) {
      this.selectedConversation = conversationID;
    },
  },
};
</script>

<template>
  <Navbar />
  <main class="d-flex">
    <Sidebar @select-conversation="selectConveration" />
    <ChatMessages :conversationID="selectedConversation" />
  </main>
</template>

<style scoped>
main {
  height: calc(100vh - var(--navbar-height));
}
</style>
