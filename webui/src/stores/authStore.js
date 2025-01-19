import { reactive } from "vue";
import axios from "../services/axios";

export const authStore = reactive({
  user: {
    data: null,
    loading: false,
    error: null,
  },

  userList: {
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

    this.userList.data = null;
    this.userList.loading = false;
    this.userList.error = null;
  },

  async getUsers() {
    try {
      this.userList.loading = true;

      const res = await axios.get("/users", {
        headers: {
          Authorization: `Bearer ${this.user.data.userID}`,
        },
      });

      this.userList.loading = false;
      this.userList.data = res.data;
      console.log(res.data);
    } catch (error) {
      this.userList.loading = false;
      this.userList.error = error.toString();
    }
  },
});
