<script>
import { conversationStore } from "../stores/conversationStore";
import LoadingSpinner from "./LoadingSpinner.vue";
export default {
  name: "MembersModal",
  props: {
    conversationID: String,
  },
  components: {
    LoadingSpinner,
  },
  data() {
    return {
      apiUrl: __API_URL__,
      conversationStore,
    };
  },
  computed: {
    members() {
      return {
        data: this.conversationStore.conversationMembers.data,
        loading: this.conversationStore.conversationMembers.loading,
        error: this.conversationStore.conversationMembers.error,
      };
    },
  },
  methods: {
    closeModal() {
      this.conversationStore.membersListFlag =
        !this.conversationStore.membersListFlag;
    },
  },
  async mounted() {
    await this.conversationStore.getMembers(this.conversationID);
  },
};
</script>

<template>
  <div id="members-modal" class="bg-slate rounded-3 p-3">
    <div class="mb-3 d-flex justify-content-between align-items-center">
      <h4>Group Members</h4>
      <i role="button" class="bi bi-x fs-5 text-danger" @click="closeModal"></i>
    </div>
    <LoadingSpinner v-if="members.loading" />
    <div v-else class="d-flex flex-column gap-3">
      <div
        v-for="member in members.data"
        :key="member.userID"
        class="d-flex justify-content-evenly align-items-center"
      >
        <img
          v-if="member.image"
          :src="`${apiUrl}${member.image}`"
          alt="Profile Image"
          class="profile-image"
        />
        <i v-else class="bi bi-person-circle fs-3"></i>

        <p>{{ member.username }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
#members-modal {
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
