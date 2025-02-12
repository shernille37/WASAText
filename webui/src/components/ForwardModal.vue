<script>
import { conversationStore } from "../stores/conversationStore";
import { messageStore } from "../stores/messageStore";
import { authStore } from "../stores/authStore";
export default {
  name: "ForwardModal",
  props: {
    message: Object,
    conversation: Object,
  },
  data() {
    return {
      apiUrl: __API_URL__,
      conversationStore,
      authStore,
      messageStore,
      selectedUser: null,
      selectedGroup: null,
    };
  },
  computed: {
    conversations() {
      return {
        data: this.conversationStore.conversations.data.filter(
          (conversation) => conversation.group
        ),
        loading: this.conversationStore.conversations.loading,
        error: this.conversationStore.conversations.error,
      };
    },
    users() {
      return {
        data: this.authStore.userList.data,
        loading: this.authStore.userList.loading,
        error: this.authStore.userList.error,
      };
    },
  },
  methods: {
    selectUser(userID) {
      // Deselect if already selected
      if (this.selectedUser === userID) this.resetFields();
      else {
        this.selectedUser = userID;
        this.selectedGroup = null;
      }
    },
    selectGroup(conversationID) {
      if (this.selectedGroup === conversationID) this.resetFields();
      else {
        this.selectedGroup = conversationID;
        this.selectedUser = null;
      }
    },
    async handleForward() {
      if (this.selectedGroup === null && this.selectedUser === null) {
        alert("Please select atleast one recipient");
        return;
      }
      // Forward the message
      await this.messageStore.forwardMessage(
        this.message.messageID,
        this.conversation.conversationID,
        this.selectedUser,
        this.selectedGroup
      );

      this.closeModal();
      this.resetFields();
    },
    isSelected(id) {
      return this.selectedUser === id || this.selectedGroup === id;
    },
    closeModal() {
      this.$emit("close-forward-modal");
    },
    hasGroupImage(conversation) {
      return !!conversation.group.groupImage;
    },
    resetFields() {
      this.selectedUser = null;
      this.selectedGroup = null;
    },
  },
  async mounted() {
    this.authStore.getUsers();
  },
};
</script>

<template>
  <div id="forward-modal" class="bg-slate rounded-3 p-3">
    <div class="mb-3 d-flex flex-column align-items-center">
      <div class="d-flex align-items-center justify-content-between w-100">
        <h4>Send To</h4>
        <i
          role="button"
          class="bi bi-x fs-5 text-danger"
          @click="closeModal"
        ></i>
      </div>
    </div>

    <div class="content d-flex flex-column gap-3 overflow-scroll">
      <!-- Existing Conversations -->
      <div
        role="button"
        :class="`d-flex justify-content-evenly align-items-center rounded-3 p-1 ${
          isSelected(conversation.conversationID) ? 'bg-info' : ''
        }`"
        v-for="conversation in conversations.data"
        :key="conversation.conversationID"
        @click="selectGroup(conversation.conversationID)"
      >
        <img
          v-if="hasGroupImage(conversation)"
          :src="`${apiUrl}${conversation.group.groupImage}`"
          alt="Profile Image"
          class="profile-image"
        />
        <i v-else class="bi bi-person-circle fs-3"></i>

        <p>
          {{ conversation.group.groupName }}
        </p>
      </div>

      <!-- Other users -->

      <div
        role="button"
        :class="`d-flex justify-content-evenly align-items-center rounded-3 p-1 ${
          isSelected(user.userID) ? 'bg-info' : ''
        }`"
        v-for="user in users.data"
        :key="user.userID"
        @click="selectUser(user.userID)"
      >
        <img
          v-if="user.image"
          :src="`${apiUrl}${user.image}`"
          alt="Profile Image"
          class="profile-image"
        />
        <i v-else class="bi bi-person-circle fs-3"></i>

        <p>
          {{ user.username }}
        </p>
      </div>
    </div>
    <button
      class="btn btn-info p-1 d-block mx-auto mt-3"
      @click="handleForward"
    >
      Forward
    </button>
  </div>
</template>

<style scoped>
#forward-modal {
  position: fixed;
  top: 50%;
  left: 60%;
  transform: translate(-50%, -50%);
  z-index: 1;
  min-height: 100px;
  max-height: 380px;

  min-width: 300px;
  max-width: 500px;
}
.content {
  max-height: 200px;
}
</style>
