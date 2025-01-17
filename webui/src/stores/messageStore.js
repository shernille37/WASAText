import { reactive } from "vue";
import axios from "../services/axios";

import { authStore } from "../stores/authStore";

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

      console.log(res.data);
    } catch (error) {
      this.messages.loading = false;
      this.messages.error = error.toString();
    }
  },

  async sendMessage(data) {
    this.messages.sendMessageLoading = true;
    try {
      // Check if there's an image
      if (data.image) {
        const formData = new FormData();
        formData.append("image", data.image);
        // Upload the image
        const resUpload = await axios.post("/upload", formData, {
          headers: {
            "Content-Type": "multipart/form-data",
            Authorization: `Bearer ${authStore.user.data.userID}`,
          },
        });

        // Send the message
        const resSendMessage = await axios.post(
          `/conversations/${data.conversationID}/messages`,
          {
            message: data.message,
            image: resUpload.data.image,
          },
          {
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${authStore.user.data.userID}`,
            },
          }
        );

        this.messages.data.push(resSendMessage.data);
      } else {
        const resSendMessage = await axios.post(
          `/conversations/${data.conversationID}/messages`,
          {
            message: data.message,
          },
          {
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${authStore.user.data.userID}`,
            },
          }
        );

        this.messages.data.push(resSendMessage.data);
      }

      this.messages.sendMessageLoading = false;
    } catch (error) {
      console.error(error.message);
      this.messages.sendMessageLoading = false;
      this.messages.error = error.toString();
    }
  },
});
