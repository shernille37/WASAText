<script>
import { conversationStore } from "../stores/conversationStore";
import Conversation from "./Conversation.vue";
import LoadingSpinner from "./LoadingSpinner.vue";
import NewConversation from "./NewConversation.vue";

export default {
  name: "Sidebar",
  components: {
    Conversation,
    LoadingSpinner,
    NewConversation,
  },
  data() {
    return {
      conversationStore,
      addConversation: false,
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
  methods: {
    toggleAddConversation() {
      this.addConversation = !this.addConversation;
      this.$emit("toggle-add-conversation", this.addConversation);
    },
    handleSelectConversation(conversationID) {
      this.addConversation = false;
      this.$emit("select-conversation", conversationID);
    },
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
        @click="toggleAddConversation"
      ></i>
    </div>

    <NewConversation v-if="addConversation" />

    <div v-if="conversations.data.length == 0">
      <p class="fs-4 text-center text-uppercase">No conversations</p>
    </div>

    <LoadingSpinner v-if="conversations.loading" />

    <div v-else>
      <Conversation
        v-for="conversation in conversations.data"
        :key="conversation.conversationID"
        :conversation="conversation"
        @click="handleSelectConversation(conversation.conversationID)"
      />
    </div>
  </div>
</template>
