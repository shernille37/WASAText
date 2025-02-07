import { reactive } from "vue";
import axios from "../services/axios";

import { authStore } from "../stores/authStore";
import { messageStore } from "./messageStore";

export const reactionStore = reactive({
  emojis: {
    data: [],
    loading: true,
    error: null,
  },
  reactions: {
    data: [],
    loading: true,
    error: null,
  },
  async getEmojis() {
    try {
      this.emojis.loading = true;

      const res = await axios.get("/emojis");

      this.emojis.loading = false;
      this.emojis.data = res.data.emojis;
    } catch (error) {
      this.emojis.error = error.response.data;
    }
  },
  async getReactions(conversationID, messageID) {
    try {
      const res = await axios.get(
        `/conversations/${conversationID}/messages/${messageID}/reactions`,
        {
          headers: {
            Authorization: `Bearer ${authStore.user.data.userID}`,
          },
        }
      );

      this.reactions.data = res.data;
    } catch (error) {
      this.reactions.error = error.response.data;
    }
  },

  async addReaction(emoji, messageID, conversationID) {
    try {
      const res = await axios.post(
        `/conversations/${conversationID}/messages/${messageID}/reactions`,
        {
          unicode: emoji,
        },
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${authStore.user.data.userID}`,
          },
        }
      );

      // Update reaction UI
      messageStore.messages.data.map((message) => {
        if (message.messageID === messageID) {
          const reaction = message.reactions.find((r) => r.unicode === emoji);

          // If the reaction is new
          if (!reaction) {
            message.reactions.push(res.data);
          } else {
            // Increment count
            reaction.count += 1;
          }
        }
      });
    } catch (error) {
      this.reactions.error = error.response.data;
    }
  },

  async deleteReaction(emoji, reactionID, messageID, conversationID) {
    try {
      await axios.delete(
        `/conversations/${conversationID}/messages/${messageID}/reactions/${reactionID}`,
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${authStore.user.data.userID}`,
          },
        }
      );

      // Update reaction UI
      messageStore.messages.data.map((message) => {
        if (message.messageID === messageID) {
          const reaction = message.reactions.find((r) => r.unicode === emoji);

          // If the count decreases to 0, remove the emoji in the UI
          if (reaction.count - 1 == 0) {
            message.reactions = message.reactions.filter(
              (r) => r.unicode !== emoji
            );
          } else {
            // Decrement the count
            reaction.count -= 1;
          }
        }
      });
    } catch (error) {
      this.reactions.error = error.response.data;
    }
  },
});
