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
      console.log("Reactions: ", res.data);
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

      // Update reaction count
      let messageStoreMessages = messageStore.messages.data;

      const message = messageStoreMessages.find(
        (message) => message.messageID === messageID
      );
      const reaction = message.reactions.find((r) => r.unicode === emoji);

      if (!reaction) {
        message.reactions.push(res.data);
      } else {
        reaction.count += 1;
      }

      messageStore.messages.data = messageStoreMessages;
      console.log(res.data);
    } catch (error) {
      this.reactions.error = error.response.data;
    }
  },
});
