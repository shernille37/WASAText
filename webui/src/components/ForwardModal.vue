<script>
import { conversationStore } from "../stores/conversationStore";
import { authStore } from "../stores/authStore";
import LoadingSpinner from "./LoadingSpinner.vue";
export default {
  name: "ForwardModal",
  props: {
    message: Object,
    conversation: Object,
  },
  components: {
    LoadingSpinner,
  },
  data() {
    return {
      apiUrl: import.meta.env.VITE_API_URL,
      conversationStore,
      authStore,
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
    users() {
      return {
        data: this.authStore.userList.data,
        loading: this.authStore.userList.loading,
        error: this.authStore.userList.error,
      };
    },
  },
  methods: {
    closeModal() {
      this.$emit("close-forward-modal");
    },
    isPrivate(conversation) {
      return !!conversation.private;
    },
    hasImage(conversation) {
      if (this.isPrivate(conversation)) {
        return !!conversation.private.user.image;
      } else {
        return !!conversation.group.groupImage;
      }
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
      <input type="text" placeholder="Search" class="form-control p-1 my-2" />
    </div>

    <div class="content d-flex flex-column gap-3 overflow-scroll">
      <!-- Existing Conversations -->
      <div
        role="button"
        class="d-flex justify-content-evenly align-items-center"
        v-for="conversation in conversations.data"
        :key="conversation.conversationID"
      >
        <img
          v-if="hasImage(conversation)"
          :src="`${apiUrl}${
            isPrivate(conversation)
              ? conversation.private.user.image
              : conversation.group.groupImage
          }`"
          alt="Profile Image"
          class="profile-image"
        />
        <i v-else class="bi bi-person-circle fs-3"></i>

        <p>
          {{
            isPrivate(conversation)
              ? conversation.private.user.username
              : conversation.group.groupName
          }}
        </p>
      </div>
    </div>
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
  max-height: 300px;

  min-width: 300px;
  max-width: 500px;
}
.content {
  max-height: 200px;
}
</style>
