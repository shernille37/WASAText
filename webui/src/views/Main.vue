<script>
import Sidebar from "../components/Sidebar.vue";
import ChatMessages from "../components/ChatMessages.vue";
import NewChatMessages from "../components/NewChatMessages.vue";
import { conversationStore } from "../stores/conversationStore";

export default {
  name: "Main",
  components: {
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
  computed: {
    addConversationFlag() {
      return this.conversationStore.addConversationFlag;
    },
  },
  methods: {
    selectConveration(conversationID) {
      this.conversationStore.addConversationFlag = false;

      this.$nextTick(() => (this.selectedConversation = conversationID));
    },
  },
};
</script>

<template>
  <main class="d-flex">
    <Sidebar @select-conversation="selectConveration" />
    <NewChatMessages
      v-if="addConversationFlag"
      @add-conversation="selectConveration"
    />
    <ChatMessages v-else :conversationID="selectedConversation" />
  </main>
</template>

<style scoped>
main {
  height: calc(100vh - var(--navbar-height));
}
</style>
