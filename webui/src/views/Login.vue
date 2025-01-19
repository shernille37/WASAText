<script>
import { authStore } from "../stores/authStore";
import LoadingSpinner from "../components/LoadingSpinner.vue";
import Copyright from "../components/Copyright.vue";
export default {
  name: "Login",
  data() {
    return {
      authStore,
      username: "",
    };
  },
  components: {
    Copyright,
    LoadingSpinner,
  },
  methods: {
    async handleSubmit() {
      if (this.username != "") await this.authStore.login(this.username);
      else alert("Please enter a username");

      this.$router.push("/");
    },
  },
};
</script>

<template>
  <div
    class="vh-100 d-flex align-items-center justify-content-center"
    v-if="authStore.user.loading"
  >
    <LoadingSpinner />
  </div>

  <div v-else class="d-flex vh-100 justify-content-center align-items-center">
    <div
      class="d-flex flex-column justify-content-between p-4 rounded-3 shadow-lg col-8 col-md-4"
      style="min-height: 250px"
    >
      <p class="text-center fs-2">WASAText</p>
      <form
        class="d-flex flex-column flex-grow-1 justify-content-evenly"
        @submit.prevent="handleSubmit"
      >
        <div>
          <input
            v-model="username"
            type="text"
            class="form-control p-2 mb-4"
            id="username"
            placeholder="Enter your username"
          />
        </div>
        <button
          type="submit"
          class="bg-primary p-2 rounded-3 text-white border-0 w-100 fw-semibold"
        >
          Log in
        </button>
      </form>
    </div>
  </div>
  <Copyright />
</template>
