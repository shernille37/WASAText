import { reactive } from "vue";
import axios from "../services/axios";

import { authStore } from "../stores/authStore";
import { conversationStore } from "./conversationStore";
import { uploadImage } from "../utils/upload";

export const messageStore = reactive({
  messages: {
    data: [],
    loading: true,
    error: null,
    sendMessageLoading: false,
  },

  replyMessage: null,

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
          replyMessageID: this.replyMessage
            ? this.replyMessage.messageID
            : null,
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
      this.replyMessage = null;

      // Update latest message
      conversationStore.conversations.data.map((conversation) => {
        if (
          conversation.conversationID === resSendMessage.data.conversationID
        ) {
          const { image, message, timestamp } = resSendMessage.data;
          conversation.latestMessage = {
            image,
            message,
            timestamp,
          };
        }
      });
    } catch (error) {
      console.error(error.response.data);
      this.messages.sendMessageLoading = false;
      this.messages.error = error.response.data;
    }
  },

  async deleteMessage(conversationID, messageID) {
    try {
      await axios.delete(
        `/conversations/${conversationID}/messages/${messageID}`,
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${authStore.user.data.userID}`,
          },
        }
      );

      // Update UI
      this.messages.data = this.messages.data.filter(
        (message) => message.messageID !== messageID
      );

      // Update latest message
      conversationStore.conversations.data.map((conversation) => {
        if (conversation.conversationID === conversationID) {
          const latestMessage =
            this.messages.data[this.messages.data.length - 1];
          if (latestMessage) {
            const { image, message, timestamp } = latestMessage;
            conversation.latestMessage = {
              image,
              message,
              timestamp,
            };
          } else {
            conversation.latestMessage = null;
          }
        }
      });
    } catch (error) {
      this.messages.error = error.response.data;
    }
  },
});
