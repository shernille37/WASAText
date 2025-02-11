import { reactive } from "vue";
import axios from "../services/axios";
import { uploadImage } from "../utils/upload";

import { conversationStore } from "./conversationStore";
import { messageStore } from "./messageStore";

export const authStore = reactive({
  user: {
    data: null,
    loading: false,
    error: null,
    success: false,
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

    conversationStore.resetFields();
    messageStore.resetFields();
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
    } catch (error) {
      this.userList.loading = false;
      this.userList.error = error.response.data;
    }
  },

  async updateProfile(imageFile, username) {
    let promises = [];
    let resUploadImage = null;
    try {
      if (imageFile) {
        resUploadImage = await uploadImage(imageFile);
        const updateProfileImage = axios.put(
          `/users/${this.user.data.userID}/image`,
          {
            image: resUploadImage,
          },
          {
            headers: {
              Authorization: `Bearer ${this.user.data.userID}`,
            },
          }
        );

        promises.push(updateProfileImage);
      }
      if (username) {
        const updateUsername = axios.put(
          `/users/${this.user.data.userID}/username`,
          {
            username: username,
          },
          {
            headers: {
              Authorization: `Bearer ${this.user.data.userID}`,
            },
          }
        );

        promises.push(updateUsername);
      }

      await Promise.all(promises);
      if (imageFile) this.user.data.image = resUploadImage;
      if (username) this.user.data.username = username;
    } catch (error) {
      this.user.error = error.response.data;
    }
  },
});
