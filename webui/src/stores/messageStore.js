import { reactive } from "vue";
import axios from "../services/axios";

import { authStore } from "../stores/authStore";
import { uploadImage } from "../utils/upload";

export const messageStore = reactive({
  messages: {
    data: [],
    loading: true,
    error: null,
    sendMessageLoading: false,
  },

  async getMessages(conversationID) {
    try {
      this.messages.loading = true;

      const res = await axios.get(`/conversations/${conversationID}/messages`, {
        headers: {
          Authorization: `Bearer ${authStore.user.data.userID}`,
        },
      });

      this.messages.loading = false;
      this.messages.data = res.data;
    } catch (error) {
      this.messages.loading = false;
      this.messages.error = error.response.data;
    }
  },

  async sendMessage(data) {
    this.messages.sendMessageLoading = true;
    let resUploadImage = null;
    try {
      // Check if there's an image
      if (data.image) {
        resUploadImage = await uploadImage(data.image);
      }

      const resSendMessage = await axios.post(
        `/conversations/${data.conversationID}/messages`,
        {
          message: data.message,
          image: resUploadImage,
        },
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${authStore.user.data.userID}`,
          },
        }
      );

      this.messages.data.push(resSendMessage.data);
      this.messages.sendMessageLoading = false;
    } catch (error) {
      console.error(error.response.data);
      this.messages.sendMessageLoading = false;
      this.messages.error = error.response.data;
    }
  },
});
