<script>
import { conversationStore } from "../stores/conversationStore";
import { messageStore } from "../stores/messageStore";
import AddMemberForm from "./AddMemberForm.vue";
import LoadingSpinner from "./LoadingSpinner.vue";
import Message from "./Message.vue";
import MessageForm from "./MessageForm.vue";

export default {
  name: "ChatMessages",
  props: {
    conversationID: String,
  },
  components: {
    Message,
    MessageForm,
    LoadingSpinner,
    AddMemberForm,
  },
  data() {
    return {
      apiUrl: import.meta.env.VITE_API_URL,
      conversationStore,
      messageStore,
    };
  },
  computed: {
    conversation() {
      return {
        data: this.conversationStore.conversation.data,
        loading: this.conversationStore.conversation.loading,
        error: this.conversationStore.conversation.error,
        showAddMember: this.conversationStore.addMemberFlag,
      };
    },
    message() {
      return {
        data: this.messageStore.messages.data,
        loading: this.messageStore.messages.loading,
        error: this.messageStore.messages.error,
      };
    },
    isPrivate() {
      return !!this.conversation.data.private;
    },
    hasImage() {
      if (this.isPrivate) {
        return !!this.conversation.data.private.user.image;
      } else {
        return !!this.conversation.data.group.groupImage;
      }
    },
  },
  methods: {
    scrollToBottom() {
      this.$nextTick(() => {
        const chatContainer = this.$refs.chatContainer;
        if (chatContainer) {
          // Scroll to the bottom
          chatContainer.scrollTop = chatContainer.scrollHeight;
        }
      });
    },
    toggleAddMemberForm() {
      this.conversationStore.addMemberFlag =
        !this.conversationStore.addMemberFlag;
    },
  },

  watch: {
    conversationID: {
      handler(id) {
        this.conversationStore.getConversation(id);
        this.messageStore.getMessages(id);
        this.scrollToBottom();
      },
    },
    "message.data": {
      handler() {
        this.scrollToBottom();
      },
      deep: true,
    },
  },
};
</script>

<template>
  <div
    v-if="!conversationID"
    class="flex-grow-1 text-center my-auto fs-2 text-uppercase"
  >
    Select a conversation
  </div>

  <div
    v-else-if="conversation.loading || message.loading"
    class="mx-auto my-auto"
  >
    <LoadingSpinner />
  </div>

  <div id="heading" v-else class="d-flex flex-column flex-grow-1">
    <!-- Heading -->
    <div
      class="d-flex justify-content-between align-items-center p-2 mb-2 border-bottom"
    >
      <div class="d-flex align-items-center">
        <img
          v-if="hasImage"
          class="profile-image"
          :src="
            isPrivate
              ? `${apiUrl}${conversation.data.private.user.image}`
              : `${apiUrl}${conversation.data.group.groupImage}`
          "
          alt="Profile Picture"
        />
        <i v-else class="bi bi-person-circle fs-3"></i>
        <p class="fw-medium fs-3 ms-3">
          {{
            isPrivate
              ? conversation.data.private.user.username
              : conversation.data.group.groupName
          }}
        </p>
      </div>

      <i
        v-if="conversation.data.group"
        role="button"
        class="bi bi-person-plus-fill rounded-circle hover-bg-light fs-5 p-1 mb-1"
        @click="toggleAddMemberForm"
      ></i>
    </div>

    <div v-if="conversation.showAddMember">
      <AddMemberForm :conversationID="conversationID" />
    </div>

    <!-- Messages -->

    <div
      ref="chatContainer"
      class="chat-container d-flex flex-column flex-grow-1 gap-3 p-2 overflow-scroll"
    >
      <Message
        v-for="message in message.data"
        :key="message.messageID"
        :message="message"
        :conversation="conversation.data"
      />
    </div>

    <MessageForm :conversationID="conversationID" />
  </div>
</template>

<style scoped>
.chat-container {
  overflow-x: hidden !important;
}
</style>
