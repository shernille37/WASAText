<script>
import { messageStore } from "../stores/messageStore";
import LoadingSpinner from "./LoadingSpinner.vue";
export default {
  name: "ReadersModal",
  props: {
    message: Object,
    conversation: Object,
  },
  components: {
    LoadingSpinner,
  },
  data() {
    return {
      apiUrl: __API_URL__,
      messageStore,
    };
  },
  computed: {
    readers() {
      return {
        data: this.messageStore.readers.data,
        loading: this.messageStore.readers.loading,
        error: this.messageStore.readers.error,
      };
    },
  },
  methods: {
    closeModal() {
      this.$emit("close-reader-modal");
    },
  },
  async mounted() {
    await this.messageStore.getReaders(
      this.conversation.conversationID,
      this.message.messageID
    );
  },
};
</script>

<template>
  <div id="reader-modal" class="bg-slate rounded-3 p-3">
    <div class="mb-3 d-flex justify-content-between align-items-center">
      <h4>Readers</h4>
      <i role="button" class="bi bi-x fs-5 text-danger" @click="closeModal"></i>
    </div>
    <LoadingSpinner v-if="readers.loading" />

    <div v-else class="d-flex flex-column gap-3">
      <div
        v-for="reader in readers.data"
        :key="reader.user.userID"
        class="d-flex justify-content-evenly align-items-center"
      >
        <img
          v-if="reader.user.image"
          :src="`${apiUrl}${reader.user.image}`"
          alt="Profile Image"
          class="profile-image"
        />
        <i v-else class="bi bi-person-circle fs-3"></i>

        <p>{{ reader.user.username }}</p>
        <p class="fs-7">{{ reader.timestamp.toLocaleString() }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
#reader-modal {
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
</style>
