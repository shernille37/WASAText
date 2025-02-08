<script>
import LoadingSpinner from "../components/LoadingSpinner.vue";
import ErrorMsg from "../components/ErrorMsg.vue";
import { reactionStore } from "../stores/reactionStore";
import { authStore } from "../stores/authStore";
export default {
  name: "EmojiPicker",
  components: {
    LoadingSpinner,
    ErrorMsg,
  },
  props: {
    messageID: String,
    senderID: String,
    conversationID: String,
  },
  emits: ["toggle-emoji-picker"],
  data() {
    return {
      reactionStore,
      authStore,
    };
  },

  computed: {
    owner() {
      return this.authStore.user.data.userID;
    },
    emojis() {
      return {
        data: this.reactionStore.emojis.data,
        loading: this.reactionStore.emojis.loading,
        error: this.reactionStore.emojis.error,
      };
    },
    reactions() {
      return {
        data: this.reactionStore.reactions.data,
        loading: this.reactionStore.reactions.loading,
        error: this.reactionStore.reactions.error,
      };
    },
  },

  mounted() {
    this.reactionStore.getEmojis();
    this.reactionStore.getReactions(this.conversationID, this.messageID);
  },

  methods: {
    async clickEmoji(emoji) {
      this.$emit("toggle-emoji-picker");

      // Check if the auth user has already reacted to the message
      const userReaction = this.reactions.data.find(
        (reaction) => reaction.reactor.userID === this.owner
      );

      // Add reaction
      if (!userReaction) {
        this.reactionStore.addReaction(
          emoji,
          this.messageID,
          this.conversationID
        );
      } else {
        // Delete the reaction and replace
        await this.reactionStore.deleteReaction(
          userReaction.unicode,
          userReaction.reactionID,
          this.messageID,
          this.conversationID
        );

        // Add the picked emoji if the picked emoji is different with the existing emoji
        if (emoji !== userReaction.unicode)
          await this.reactionStore.addReaction(
            emoji,
            this.messageID,
            this.conversationID
          );
      }
    },
  },
};
</script>

<template>
  <div v-if="emojis.loading">
    <LoadingSpinner />
  </div>
  <div v-else-if="emojis.error">
    <ErrorMsg :msg="error" />
  </div>

  <div v-else class="d-flex" id="emoji-picker">
    <span
      v-for="emoji in emojis.data"
      :key="emoji"
      @click="clickEmoji(emoji)"
      >{{ emoji }}</span
    >
  </div>
</template>

<style scoped>
#emoji-picker {
  width: auto;
  position: absolute;
  border: 1px solid #ccc;
  border-radius: 5px;
  font-size: 20px;
  background-color: white;
}

#emoji-picker span {
  cursor: pointer;
  margin-right: 5px !important;
}

#emoji-picker span:hover {
  background-color: var(--color-slate);
}
</style>
