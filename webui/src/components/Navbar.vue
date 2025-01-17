<script>
import ErrorMsg from "../components/ErrorMsg.vue";
import { authStore } from "../stores/authStore";

export default {
  name: "Navbar",
  data() {
    return {
      apiUrl: import.meta.env.VITE_API_URL,
      authStore,
    };
  },
  components: {
    ErrorMsg,
  },
  computed: {
    user() {
      return this.authStore.user.data;
    },
  },
  methods: {
    handleLogout() {
      this.authStore.logout();
      this.$router.push("/login");
    },
  },
};
</script>

<template>
  <nav
    class="navbar bg-white py-1 px-4 position-sticky top-0 d-flex justify-content-between align-items-center shadow-sm"
  >
    <div>
      <p class="fs-2 fw-semibold navbar-brand">WASAText</p>
    </div>

    <div class="d-flex align-items-center gap-2">
      <img
        v-if="user.image"
        :src="`${apiUrl}${user.image}`"
        alt="Profile Picture"
        class="profile-image"
      />
      <i v-else class="bi bi-person-circle fs-3"></i>

      <p class="ms-2">{{ user.username }}</p>
      <div class="dropdown">
        <i
          data-bs-toggle="dropdown"
          aria-expanded="false"
          role="button"
          class="bi bi-chevron-down ms-2 p-2 hover-bg-light"
        ></i>
        <ul class="dropdown-menu dropdown-menu-end mt-3">
          <li>
            <p role="button" class="dropdown-item p-1 rounded-1">Profile</p>
          </li>
          <li>
            <p
              @click="handleLogout"
              role="button"
              class="dropdown-item p-1 rounded-1"
            >
              Logout
            </p>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<style scoped>
nav {
  height: var(--navbar-height);
}

i {
  border-radius: 999px;
}
</style>
