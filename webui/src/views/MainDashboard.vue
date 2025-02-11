<script>
import Sidebar from "../components/Sidebar.vue";
import ChatMessages from "../components/ChatMessages.vue";
import NewChatMessages from "../components/NewChatMessages.vue";
import { conversationStore } from "../stores/conversationStore";

export default {
  name: "MainDashboard",
  components: {
    Sidebar,
    ChatMessages,
    NewChatMessages,
  },
  data() {
    return {
      conversationStore,
      localSelectedConversation: null,
    };
  },
  computed: {
    addConversationFlag() {
      return this.conversationStore.addConversationFlag;
    },
    selectedConversation() {
      return this.conversationStore.selectedConversation;
    },
  },
  watch: {
    selectedConversation: {
      handler(newSelectedConversation) {
        this.conversationStore.addConversationFlag = false;
        this.$nextTick(
          () => (this.localSelectedConversation = newSelectedConversation)
        );
      },
      deep: true,
    },
  },
  methods: {
    selectConversation(conversationID) {
      this.conversationStore.addConversationFlag = false;
      this.$nextTick(() => (this.localSelectedConversation = conversationID));
    },
  },
};
</script>

<template>
  <main class="d-flex">
    <Sidebar @select-conversation="selectConversation" />
    <NewChatMessages
      @add-conversation="selectConversation"
      v-if="addConversationFlag"
    />
    <ChatMessages v-else :conversationID="localSelectedConversation" />
  </main>
</template>

<style scoped>
main {
  height: calc(100vh - var(--navbar-height));
}
</style>
