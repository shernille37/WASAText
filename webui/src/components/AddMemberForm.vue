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
  props: {
    conversationID: String,
  },
  emits: ["refresh-conversation"],
  data() {
    return {
      apiUrl: __API_URL__,
      authStore,
      conversationStore,
      suggestedMembers: [],
      selectedMembers: [],
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

    users() {
      return {
        data: this.authStore.userList.data,
        loading: this.authStore.userList.loading,
        error: this.authStore.userList.error,
      };
    },

    show() {
      return this.conversationStore.addMemberFlag;
    },
  },
  methods: {
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
  watch: {
    "conversation.error": {
      handler() {
        setTimeout(() => {
          this.conversationStore.conversation.error = null;
        }, 3000);
      },
      deep: true,
    },
  },
  async mounted() {
    await this.authStore.getUsers();
    await this.conversationStore.getMembers(this.conversationID);
    // Just display the members that are NOT currently members
    this.suggestedMembers = this.users.data.filter(
      (user) => !this.members.data.some((u) => u.userID === user.userID)
    );
  },
};
</script>

<template>
  <div id="addMemberContainer">
    <div v-if="members.loading">
      <LoadingSpinner />
    </div>

    <div v-else-if="members.error">
      <ErrorMsg :msg="members.error" />
    </div>

    <div class="d-flex justify-content-center align-items-center">
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
            <img
              v-if="user.image"
              :src="`${apiUrl}${user.image}`"
              alt="Profile Image"
              class="profile-image"
            />
            <i v-else class="bi bi-person-circle fs-3"></i>
            <p>{{ user.username }}</p>
          </div>
        </div>

        <div id="suggestedMembersContainer" class="overflow-scroll">
          <div
            role="button"
            class="d-flex gap-3 p-2 hover-bg-light rounded-2 align-items-center"
            v-for="user in suggestedMembers"
            :key="user.userID"
            @click="handleAddMember(user)"
          >
            <img
              v-if="user.image"
              :src="`${apiUrl}${user.image}`"
              alt="Profile Picture"
              class="profile-image"
            />
            <i v-else class="bi bi-person-circle fs-3"></i>
            <p>{{ user.username }}</p>
          </div>
        </div>

        <div>
          <button
            class="btn btn-primary w-50 mx-auto d-block mt-2"
            @click="handleSubmit"
          >
            Add
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
#selected-members p {
  font-size: 12px;
}

#addMemberContainer {
  position: fixed;
  top: 50%;
  left: 60%;

  transform: translate(-50%, -50%);
  z-index: 1;
  min-height: 100px;

  min-width: 700px;
  max-width: 500px;
}

#suggestedMembersContainer {
  max-height: 200px;
}
</style>
