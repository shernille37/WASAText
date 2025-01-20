<script>
import LoadingSpinner from "../components/LoadingSpinner.vue";
import ErrorMsg from "../components/ErrorMsg.vue";
import { authStore } from "../stores/authStore";
export default {
  name: "EmojiPicker",
  components: {
    LoadingSpinner,
    ErrorMsg,
  },
  props: {
    senderID: String,
    conversationID: String,
  },
  data() {
    return {
      authStore,
      emojis: [],
      loading: true,
      error: null,
    };
  },

  async mounted() {
    try {
      this.loading = true;

      const res = await this.$axios.get("/emojis");

      this.loading = false;
      this.emojis = res.data.emojis;
    } catch (error) {
      this.error = error.response.data;
    }
  },

  methods: {
    clickEmoji(emoji) {
      console.log(emoji, this.senderID, this.conversationID);
    },
  },
};
</script>

<template>
  <div v-if="loading">
    <LoadingSpinner />
  </div>
  <div v-else-if="error">
    <ErrorMsg :msg="error" />
  </div>

  <div v-else class="d-flex" id="emoji-picker">
    <span v-for="emoji in emojis" :key="emoji" @click="clickEmoji(emoji)">{{
      emoji
    }}</span>
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
