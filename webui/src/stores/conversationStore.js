import { reactive } from "vue";
import axios from "../services/axios";

import { authStore } from "../stores/authStore";

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
      this.conversation.error = error.toString();
    }
  },
});
