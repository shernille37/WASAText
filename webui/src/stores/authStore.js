import { reactive } from "vue";
import axios from "../services/axios";

export const authStore = reactive({
  user: {
    data: null,
    loading: false,
    error: null,
  },

  async login(username) {
    try {
      this.user.loading = true;

      const res = await axios.post(
        "/login",
        { username },
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      );

      this.user.loading = false;
      this.user.data = res.data;
    } catch (error) {
      this.user.error = error.toString();
    }
  },

  logout() {
    this.user.data = null;
    this.user.loading = false;
    this.user.error = null;
  },
});
