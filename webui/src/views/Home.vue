<script>
import Navbar from "../components/Navbar.vue";
import Sidebar from "../components/Sidebar.vue";
import ChatMessages from "../components/ChatMessages.vue";
import NewChatMessages from "../components/NewChatMessages.vue";
import { conversationStore } from "../stores/conversationStore";

export default {
  name: "Home",
  components: {
    Navbar,
    Sidebar,
    ChatMessages,
    NewChatMessages,
  },
  data() {
    return {
      conversationStore,
      selectedConversation: null,
      addConversation: false,
    };
  },
  methods: {
    selectConveration(conversationID) {
      this.addConversation = false;
      // To avoid race conditions
      this.$nextTick(() => (this.selectedConversation = conversationID));
    },
    toggleAddConversation(addConversation) {
      this.addConversation = addConversation;
    },
  },
};
</script>

<template>
  <Navbar />
  <main class="d-flex">
    <Sidebar
      @toggle-add-conversation="toggleAddConversation"
      @select-conversation="selectConveration"
    />
    <NewChatMessages v-if="addConversation" />
    <ChatMessages v-else :conversationID="selectedConversation" />
  </main>
</template>

<style scoped>
main {
  height: calc(100vh - var(--navbar-height));
}
</style>
