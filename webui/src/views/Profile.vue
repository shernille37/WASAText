<script>
import { authStore } from "../stores/authStore";
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
      authStore,
      apiUrl: import.meta.env.VITE_API_URL,
      username: "",
      image: null,
      fileToUpload: null,
    };
  },
  computed: {
    user() {
      return {
        data: authStore.user.data,
        loading: authStore.user.loading,
        error: authStore.user.error,
      };
    },
  },
  methods: {
    openFileUpload() {
      this.$refs.profileImage.click();
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
      if (this.fileToUpload || this.username)
        await this.authStore.updateProfile(this.fileToUpload, this.username);

      if (!this.user.error) this.$router.push("/");
    },
  },
  watch: {
    "user.error": {
      handler() {
        setTimeout(() => {
          this.authStore.user.error = null;
        }, 3000);
      },
      deep: true,
    },
  },
};
</script>

<template>
  <div class="profile d-flex justify-content-center align-items-center">
    <div
      v-if="user.data"
      class="d-flex flex-column align-items-center justify-content-evenly rounded-3"
    >
      <ErrorMsg :msg="user.error" />
      <p class="fs-5">Edit Profile Details</p>

      <div
        v-if="user.data.image || image"
        class="h-50 w-50 d-flex flex-column align-items-center justify-content-evenly"
      >
        <img v-if="image" :src="image" alt="Profile Preview" />
        <img v-else :src="`${apiUrl}${user.data.image}`" alt="Profile" />
      </div>
      <i v-else class="profile-icon bi bi-person-circle"></i>

      <i role="button" class="bi bi-pencil-square" @click="openFileUpload"></i>

      <input
        @change="handleUploadProfile"
        type="file"
        name="image"
        id="image"
        class="d-none"
        ref="profileImage"
      />

      <input
        v-model="username"
        type="text"
        :placeholder="`${user.data.username}`"
        class="form-control p-1 w-50"
      />

      <button
        class="btn btn-info p-2"
        :disabled="!image && !username"
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
