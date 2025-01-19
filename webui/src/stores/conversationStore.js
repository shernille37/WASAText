import { reactive } from "vue";
import axios from "../services/axios";

import { authStore } from "../stores/authStore";
import { uploadImage } from "../utils/upload";

export const conversationStore = reactive({
  conversations: {
    data: [],
    loading: true,
    error: null,
  },
  conversation: {
    data: null,
    loading: true,
    error: null,
  },
  async getConversations() {
    try {
      this.conversations.loading = true;
      const res = await axios.get("/conversations", {
        headers: {
          Authorization: `Bearer ${authStore.user.data.userID}`,
        },
      });
      this.conversations.loading = false;
      this.conversations.data = res.data;
    } catch (error) {
      this.conversations.loading = false;
      this.conversations.error = error.toString();
    }
  },

  async getConversation(id) {
    try {
      this.conversation.loading = true;
      const res = await axios.get(`/conversations/${id}`, {
        headers: {
          Authorization: `Bearer ${authStore.user.data.userID}`,
        },
      });

      this.conversation.loading = false;
      this.conversation.data = res.data;
    } catch (error) {
      this.conversation.loading = false;
      this.conversation.error = error.response.data;
    }
  },

  async addConversation(data) {
    this.conversation.loading = true;

    let resImageUpload = null;
    let resGroupImageUpload = null;

    try {
      if (data.conversationType === "group" && data.groupImage) {
        resGroupImageUpload = await uploadImage(data.groupImage);
      }

      if (data.image) {
        resImageUpload = await uploadImage(data.image);
      }

      // Add Conversation
      if (data.conversationType === "group") {
        const resAddConversation = await axios.post(
          `/group-conversations`,
          {
            groupName: data.groupName,
            groupImage: resGroupImageUpload,
            message: data.message,
            image: resImageUpload,
            members: data.members,
          },
          {
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${authStore.user.data.userID}`,
            },
          }
        );

        this.conversations.loading = false;

        this.conversations.data.push(resAddConversation.data);
      } else {
        const resAddConversation = await axios.post(
          `/private-conversations`,
          {
            receiverID: data.members[0],
            message: data.message,
            image: resImageUpload,
          },
          {
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${authStore.user.data.userID}`,
            },
          }
        );

        this.conversations.loading = false;

        this.conversations.data.push(resAddConversation.data);
      }
    } catch (error) {
      this.conversations.loading = false;
      this.conversations.error = error.response.data;
    }
  },
});
