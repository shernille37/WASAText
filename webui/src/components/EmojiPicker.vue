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
    clickEmoji(emoji) {
      console.log(emoji, this.senderID, this.conversationID);

      // Check if the auth user has already reacted to the message
      const hasReacted = this.reactions.data.some(
        (reaction) => reaction.reactor.userID === this.owner
      );

      // Add reaction
      if (!hasReacted) {
        this.reactionStore.addReaction(
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
