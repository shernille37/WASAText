<script>
import { conversationStore } from "../stores/conversationStore";
import { messageStore } from "../stores/messageStore";
import AddMemberForm from "./AddMemberForm.vue";
import LoadingSpinner from "./LoadingSpinner.vue";
import MembersModal from "./MembersModal.vue";
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
    MembersModal,
  },
  data() {
    return {
      apiUrl: __API_URL__,
      pollingInterval: null,
      POLLING_DELAY: 5000,
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
        showMemberList: this.conversationStore.membersListFlag,
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
    async refresh() {
      if (this.conversationID) {
        this.conversationStore.getConversation(this.conversationID);
        this.messageStore.getMessages(this.conversationID);
        this.scrollToBottom();
      }
    },
    startPolling() {
      this.pollingInterval = setInterval(this.refresh, this.POLLING_DELAY);
    },
    stopPolling() {
      if (this.pollingInterval) {
        clearInterval(this.pollingInterval);
        this.pollingInterval = null;
      }
    },
    scrollToBottom() {
      this.$nextTick(() => {
        const chatContainer = this.$refs.chatContainer;
        if (chatContainer) {
          // Scroll to the bottom
          chatContainer.scrollTop = chatContainer.scrollHeight;
        }
      });
    },
    async handleLeaveConversation() {
      if (confirm("Are you sure?"))
        await this.conversationStore.leaveGroupConversation(
          this.conversationID
        );
    },
    toggleAddMemberForm() {
      this.conversationStore.membersListFlag = false;
      this.conversationStore.addMemberFlag =
        !this.conversationStore.addMemberFlag;
    },
    toggleMembersList() {
      this.conversationStore.addMemberFlag = false;
      this.conversationStore.membersListFlag =
        !this.conversationStore.membersListFlag;
    },
    goToEditGroupConversation() {
      if (this.conversation.data.group)
        this.$router.push(
          `/group-conversations/${this.conversation.data.conversationID}/edit`
        );
    },
  },

  watch: {
    conversationID: {
      async handler(newID) {
        await this.conversationStore.getConversation(this.conversationID);
        await this.messageStore.getMessages(this.conversationID);
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

  async mounted() {
    if (this.conversationID) {
      await this.conversationStore.getConversation(this.conversationID);
      await this.messageStore.getMessages(this.conversationID);
      this.scrollToBottom();
    }

    // this.startPolling();
  },
  beforeUnmount() {
    this.stopPolling();
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
        <i
          v-if="conversation.data.group"
          role="button"
          class="bi bi-pencil-square ms-2"
          @click="goToEditGroupConversation"
        ></i>
      </div>

      <div class="d-flex gap-2">
        <i
          v-if="conversation.data.group"
          role="button"
          class="bi bi-info-circle rounded-circle hover-bg-light fs-5 p-1 mb-1"
          @click="toggleMembersList"
        ></i>
        <i
          v-if="conversation.data.group"
          role="button"
          class="bi bi-person-plus-fill rounded-circle hover-bg-light fs-5 p-1 mb-1"
          @click="toggleAddMemberForm"
        ></i>
        <i
          role="button"
          v-if="conversation.data.group"
          class="bi bi-box-arrow-right rounded-circle hover-bg-light fs-5 p-1 mb-1"
          @click="handleLeaveConversation"
        ></i>
      </div>
    </div>

    <AddMemberForm
      v-if="conversation.showAddMember"
      :conversationID="conversationID"
    />

    <MembersModal
      v-if="conversation.showMemberList"
      :conversationID="conversationID"
    />

    <!-- Messages -->
    <div
      ref="chatContainer"
      class="chat-container position-relative d-flex flex-column flex-grow-1 gap-5 p-2 overflow-scroll"
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
