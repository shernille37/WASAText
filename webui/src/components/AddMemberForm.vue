<script>
import { authStore } from "../stores/authStore";
import { conversationStore } from "../stores/conversationStore";
import ErrorMsg from "./ErrorMsg.vue";
import LoadingSpinner from "./LoadingSpinner.vue";

export default {
  name: "AddMemberForm",
  components: {
    LoadingSpinner,
    ErrorMsg,
  },
  data() {
    return {
      authStore,
      conversationStore,
      apiUrl: import.meta.env.API_URL,
      suggestedMembers: [],
      selectedMembers: [],
    };
  },
  computed: {
    users() {
      return {
        data: this.authStore.userList.data,
        loading: this.authStore.userList.loading,
        error: this.authStore.userList.error,
      };
    },

    conversation() {
      return {
        data: this.conversationStore.conversation.data,
        loading: this.conversationStore.conversation.loading,
        error: this.conversationStore.conversation.error,
      };
    },

    show() {
      return this.conversationStore.addMemberFlag;
    },
  },
  methods: {
    async getUsers() {
      await this.authStore.getUsers();

      // Just display the members that are NOT already selected (NOT in selectedMembers)
      this.suggestedMembers = this.users.data.filter(
        (user) =>
          !this.conversation.data.members.some((u) => u.userID === user.userID)
      );
    },
    handleAddMember(user) {
      this.selectedMembers.push(user);

      // Remove from suggested members that are already selected
      this.suggestedMembers = this.suggestedMembers.filter(
        (u) => u.userID !== user.userID
      );
    },
    handleRemoveSelected(user) {
      // Remove the user from the selected Members
      this.selectedMembers = this.selectedMembers.filter(
        (u) => u.userID !== user.userID
      );
      // Re-add the user to the suggested members
      this.suggestedMembers.push(user);
    },
    async handleSubmit() {
      if (!this.selectedMembers.length) {
        alert("Please select atleast 1 Member");
      }

      const data = {
        members: this.selectedMembers.map((user) => user.userID),
      };

      await conversationStore.addMembersToGroup(data);
    },
    close() {
      this.conversationStore.addMemberFlag =
        !this.conversationStore.addMemberFlag;
    },
  },
  mounted() {
    this.getUsers();
  },
};
</script>

<template>
  <div v-if="conversation.loading">
    <LoadingSpinner />
  </div>

  <div v-else-if="conversation.error">
    <ErrorMsg msg="error" />
  </div>

  <div
    v-else-if="show"
    class="h-100 w-100 d-flex justify-content-center align-items-center"
  >
    <div class="position-relative p-2 bg-info w-50 rounded-2">
      <p class="fw-semibold text-center">Add People</p>
      <p
        v-if="!this.suggestedMembers.length && !this.selectedMembers.length"
        class="text-center"
      >
        All users in the system is already part of this group
      </p>
      <i
        role="button"
        class="position-absolute bi bi-x fs-3 end-0 top-0"
        @click="close"
      ></i>
      <form @submit.prevent="handleSubmit">
        <input
          type="text"
          name="addMember"
          id="addMember"
          class="form-control p-1 mt-3 mb-3"
          autocomplete="off"
          placeholder="Add People"
        />
      </form>

      <div class="d-flex gap-3">
        <div
          role="button"
          id="selected-members"
          class="p-2 hover-bg-light rounded-2 text-center"
          v-for="user in selectedMembers"
          :key="user.userID"
          @click="handleRemoveSelected(user)"
        >
          <img v-if="user.image" :src="`${apiUrl}${user.image}`" alt="" />
          <i class="bi bi-person-circle fs-3"></i>
          <p>{{ user.username }}</p>
        </div>
      </div>

      <div
        role="button"
        class="d-flex gap-3 p-2 hover-bg-light rounded-2"
        v-for="user in suggestedMembers"
        :key="user.userID"
        @click="handleAddMember(user)"
      >
        <img v-if="user.image" :src="`${apiUrl}${user.image}`" alt="" />
        <i v-else class="bi bi-person-circle fs-3"></i>
        <p>{{ user.username }}</p>
      </div>

      <div>
        <button
          class="btn btn-primary w-50 mx-auto d-block"
          @click="handleSubmit"
        >
          Add
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
#selected-members p {
  font-size: 12px;
}
</style>
