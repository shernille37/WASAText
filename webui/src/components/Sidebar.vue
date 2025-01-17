<script>
import { conversationStore } from "../stores/conversationStore";
import Conversation from "./Conversation.vue";
import LoadingSpinner from "./LoadingSpinner.vue";
export default {
  name: "Sidebar",
  components: {
    Conversation,
    LoadingSpinner,
  },
  data() {
    return {
      conversationStore,
    };
  },
  computed: {
    conversations() {
      return {
        data: this.conversationStore.conversations.data,
        loading: this.conversationStore.conversations.loading,
        error: this.conversationStore.conversations.error,
      };
    },
  },
  mounted() {
    this.conversationStore.getConversations();
  },
};
</script>

<template>
  <div
    class="d-flex flex-column p-2 shadow-sm col-3 overflow-auto overflow-x-hidden"
  >
    <div
      class="d-flex justify-content-center justify-content-sm-between align-items-center p-2"
    >
      <h3 class="fw-bold d-none d-sm-block">Chats</h3>
      <i
        role="button"
        class="bi bi-plus-circle fs-5 hover-bg-light rounded-circle"
      ></i>
    </div>

    <div v-if="conversations.data.length == 0">
      <p class="fs-4 text-center text-uppercase">No conversations</p>
    </div>

    <LoadingSpinner v-if="conversations.loading" />

    <div v-else>
      <Conversation
        v-for="conversation in conversations.data"
        :key="conversation.conversationID"
        :conversation="conversation"
        @click="$emit('select-conversation', conversation.conversationID)"
      />
    </div>
  </div>
</template>
