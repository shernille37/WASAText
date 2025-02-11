<script>
import { conversationStore } from "../stores/conversationStore";
import { messageStore } from "../stores/messageStore";
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
      messageStore,
      conversationType: "all",
    };
  },
  computed: {
    conversations() {
      const data =
        this.conversationType === "all"
          ? this.conversationStore.conversations.data
          : this.conversationType === "private"
          ? this.conversationStore.privateConversations.data
          : this.conversationStore.groupConversations.data;

      const loading =
        this.conversationType === "all"
          ? this.conversationStore.conversations.loading
          : this.conversationType === "private"
          ? this.conversationStore.privateConversations.loading
          : this.conversationStore.groupConversations.loading;

      const error =
        this.conversationType === "all"
          ? this.conversationStore.conversations.error
          : this.conversationType === "private"
          ? this.conversationStore.privateConversations.error
          : this.conversationStore.groupConversations.error;

      return {
        data,
        loading,
        error,
      };
    },
  },
  mounted() {
    this.conversationStore.getConversations();
  },
  methods: {
    toggleAddConversation() {
      this.conversationStore.addConversationFlag =
        !this.conversationStore.addConversationFlag;
    },
    handleSelectConversation(conversationID) {
      this.conversationStore.selectedConversation = conversationID;
      this.conversationStore.addMemberFlag = false;
      this.messageStore.replyMessage = null;
    },
    async toggleConversationType(type) {
      this.conversationType = type;
      if (type === "all") await this.conversationStore.getConversations();
      else if (type === "private")
        await this.conversationStore.getPrivateConversations();
      else await this.conversationStore.getGroupConversations();
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

      <div class="btn-group" role="group" aria-label="Basic example">
        <button
          type="button"
          class="btn btn-info p-1"
          @click="toggleConversationType('all')"
        >
          All
        </button>
        <button
          type="button"
          class="btn btn-primary"
          @click="toggleConversationType('private')"
        >
          Private
        </button>
        <button
          type="button"
          class="btn btn-secondary"
          @click="toggleConversationType('group')"
        >
          Group
        </button>
      </div>

      <i
        role="button"
        class="bi bi-plus-circle fs-5 hover-bg-light rounded-circle"
        @click="toggleAddConversation"
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
        @click="handleSelectConversation(conversation.conversationID)"
      />
    </div>
  </div>
</template>
