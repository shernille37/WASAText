<script>
import { reactionStore } from "../stores/reactionStore";
import LoadingSpinner from "./LoadingSpinner.vue";
export default {
  name: "ReactionModal",
  props: {
    message: Object,
    conversation: Object,
  },
  components: {
    LoadingSpinner,
  },
  data() {
    return {
      apiUrl: __API_URL__,
      reactionStore,
    };
  },
  computed: {
    reactions() {
      return {
        data: this.reactionStore.reactions.data,
        loading: this.reactionStore.reactions.loading,
        error: this.reactionStore.reactions.error,
      };
    },
  },
  methods: {
    closeModal() {
      this.$emit("close-modal");
    },
  },
  async mounted() {
    await this.reactionStore.getReactions(
      this.conversation.conversationID,
      this.message.messageID
    );
  },
};
</script>

<template>
  <div id="reaction-modal" class="bg-slate rounded-3 p-3">
    <div class="mb-3 d-flex justify-content-between align-items-center">
      <h4>Reactions</h4>
      <i role="button" class="bi bi-x fs-5 text-danger" @click="closeModal"></i>
    </div>
    <LoadingSpinner v-if="reactions.loading" />
    <div v-else class="d-flex flex-column gap-3">
      <div
        v-for="reaction in reactions.data"
        :key="reaction.reactionID"
        class="d-flex justify-content-evenly align-items-center"
      >
        <img
          v-if="reaction.reactor.image"
          :src="`${apiUrl}${reaction.reactor.image}`"
          alt="Profile Image"
          class="profile-image"
        />
        <i v-else class="bi bi-person-circle fs-3"></i>

        <p>{{ reaction.reactor.username }}</p>
        <p>{{ reaction.unicode }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
#reaction-modal {
  position: fixed;
  top: 50%;
  left: 60%;
  transform: translate(-50%, -50%);
  z-index: 1;
  min-height: 100px;
  max-height: 300px;

  min-width: 300px;
  max-width: 500px;
}
</style>
