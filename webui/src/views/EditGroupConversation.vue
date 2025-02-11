<script>
import { conversationStore } from "../stores/conversationStore";
import LoadingSpinner from "../components/LoadingSpinner.vue";
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  name: "Profile",
  components: {
    LoadingSpinner,
    ErrorMsg,
  },
  data() {
    return {
      apiUrl: __API_URL__,
      conversationStore,
      groupName: "",
      image: null,
      fileToUpload: null,
    };
  },
  computed: {
    conversation() {
      return {
        data: conversationStore.conversation.data,
        loading: conversationStore.conversation.loading,
        error: conversationStore.conversation.error,
      };
    },
  },
  methods: {
    openFileUpload() {
      this.$refs.groupImage.click();
    },
    handleUploadProfile(e) {
      const file = e.target.files[0];
      if (file && file.type.startsWith("image/")) {
        this.image = URL.createObjectURL(file);
        this.fileToUpload = file;
      } else {
        alert("Please upload a valid image file");
      }
    },
    async handleSubmit() {
      if (this.fileToUpload || this.groupName) {
        await this.conversationStore.editGroupConversation(
          this.fileToUpload,
          this.groupName
        );
      }
      if (!this.conversation.error) this.$router.push("/");
    },
  },
  watch: {
    "conversation.error": {
      handler() {
        setTimeout(() => {
          this.conversationStore.conversation.error = null;
        }, 3000);
      },
      deep: true,
    },
  },
  mounted() {
    this.conversationStore.getConversation(this.$route.params.conversationID);
  },
};
</script>

<template>
  <div v-if="conversation.loading">
    <LoadingSpinner />
  </div>
  <div v-else class="profile d-flex justify-content-center align-items-center">
    <div
      v-if="conversation.data.group"
      class="d-flex flex-column align-items-center justify-content-evenly rounded-3"
    >
      <ErrorMsg :msg="conversation.error" />
      <p class="fs-5">Edit Group Details</p>

      <div
        v-if="conversation.data.group.groupImage || image"
        class="position-relative h-50 w-50 d-flex flex-column align-items-center justify-content-evenly"
      >
        <img v-if="image" :src="image" alt="Profile Preview" />
        <img
          v-else
          :src="`${apiUrl}${conversation.data.group.groupImage}`"
          alt="Profile"
        />
      </div>
      <i v-else class="profile-icon bi bi-person-circle"></i>
      <i role="button" class="bi bi-pencil-square" @click="openFileUpload"></i>

      <input
        @change="handleUploadProfile"
        type="file"
        name="image"
        id="image"
        class="d-none"
        ref="groupImage"
      />

      <input
        v-model="groupName"
        type="text"
        :placeholder="`${conversation.data.group.groupName}`"
        class="form-control p-1 w-50"
      />

      <button
        class="btn btn-info p-2"
        :disabled="!image && !groupName"
        @click="handleSubmit"
      >
        Save
      </button>
    </div>
  </div>
</template>

<style scoped>
.profile {
  height: calc(100vh - var(--navbar-height));
}
.profile div {
  background-color: var(--color-slate);
  height: 500px;
  width: 500px;
}
.profile-icon {
  font-size: 5rem;
  transform: translateX(-20px);
}
img {
  width: 200px;
  height: 200px;
  border-radius: 999px;
}
p {
  font-size: 30px;
}
</style>
