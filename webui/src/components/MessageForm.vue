<script>
import { messageStore } from "../stores/messageStore";
export default {
  name: "MessageForm",
  props: {
    conversationID: String,
    newConversation: Boolean,
  },
  data() {
    return {
      messageStore,
      message: "",
      image: null,
      fileToUpload: null,
    };
  },
  methods: {
    openFile() {
      this.$refs.file.click();
    },
    closeFile() {
      this.image = null;
      this.$refs.file.value = "";
    },
    handleSubmit() {
      if (this.message == "" && !this.fileToUpload) {
        alert("Please insert a message or an image");
        return;
      }

      const data = {
        conversationID: this.conversationID,
        message: this.message,
        image: this.fileToUpload,
      };

      if (this.newConversation) {
        this.$emit("add-conversation", data);
      } else messageStore.sendMessage(data);

      this.message = "";
      this.image = null;
      this.fileToUpload = null;
    },
    handleFileUpload(e) {
      const file = e.target.files[0];

      if (file && file.type.startsWith("image/")) {
        this.image = URL.createObjectURL(file);
        this.fileToUpload = file;
      } else {
        alert("Please upload a valid image file");
      }
    },
  },
};
</script>

<template>
  <div class="d-flex align-items-end p-3 w-100">
    <i
      role="button"
      class="bi bi-plus-circle fs-5 hover-bg-light rounded-circle"
      @click="openFile"
    >
    </i>

    <div class="flex-grow-1">
      <!-- Image Preview -->
      <div v-show="image" class="ms-4">
        <i
          role="button"
          class="position-absolute bi bi-x fs-5 text-danger bg-white rounded-circle"
          @click="closeFile"
        ></i>
        <img :src="image" alt="Image" width="150" />
      </div>
      <!-- Reply Preview -->
      <!-- <div>
        <p class="fs-6 ms-4 bg-light p-2">Reply Preview</p>
      </div> -->
      <form
        @submit.prevent="handleSubmit"
        class="d-flex flex-grow-1 ms-3 p-1 rounded-4"
      >
        <input
          @change="handleFileUpload"
          type="file"
          name="image"
          id="image"
          class="d-none"
          ref="file"
        />
        <input
          type="text"
          name="message"
          id="message"
          v-model="message"
          :placeholder="`${
            messageStore.messages.sendMessageLoading
              ? 'Sending'
              : 'Type a message...'
          }`"
          class="form-control p-1"
          autocomplete="off"
        />

        <button type="submit" class="d-none">Hidden Submit</button>
      </form>
    </div>
  </div>
</template>
