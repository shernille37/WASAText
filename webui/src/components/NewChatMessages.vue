<script>
import MessageForm from "./MessageForm.vue";
import LoadingSpinner from "./LoadingSpinner.vue";
import ErrorMsg from "./ErrorMsg.vue";
import { authStore } from "../stores/authStore";
import { conversationStore } from "../stores/conversationStore";

export default {
  name: "NewChatMessages",
  components: {
    MessageForm,
    LoadingSpinner,
    ErrorMsg,
  },
  data() {
    return {
      apiUrl: __API_URL__,
      authStore,
      conversationStore,
      conversationType: null,
      groupName: null,
      groupImage: null,
      groupImageToUpload: null,
      searchMembers: "",
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
    conversations() {
      return {
        data: this.conversationStore.conversations.data,
        loading: this.conversationStore.conversations.loading,
        error: this.conversationStore.conversations.error,
      };
    },
  },
  watch: {
    "conversations.error": {
      handler() {
        setTimeout(() => {
          this.conversationStore.conversations.error = null;
        }, 3000);
      },
      deep: true,
    },
  },
  methods: {
    openGroupImageUpload() {
      this.$refs.group_image.click();
    },
    closeGroupImageUpload() {
      this.groupImage = null;
      this.$refs.group_image.value = "";
    },
    handleGroupImageUpload(e) {
      const file = e.target.files[0];

      if (file && file.type.startsWith("image/")) {
        this.groupImage = URL.createObjectURL(file);
        this.groupImageToUpload = file;
      } else {
        alert("Please upload a valid image file");
      }
    },
    async getUsers() {
      await this.authStore.getUsers();
      // Just display the members that are NOT already selected (NOT in selectedMembers)
      this.suggestedMembers = this.users.data.filter(
        (user) => !this.selectedMembers.some((u) => u.userID === user.userID)
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
    handleAddMember(user) {
      if (this.conversationType === "group") {
        this.selectedMembers.push(user);

        // Remove from suggested members that are already selected
        this.suggestedMembers = this.suggestedMembers.filter(
          (u) => u.userID !== user.userID
        );
      } else {
        this.selectedMembers = [];
        this.selectedMembers.push(user);
        this.suggestedMembers = [];
      }
    },
    async handleAddConversation(data) {
      if (!this.conversationType) {
        alert("Please select a conversation type");
        return;
      }

      if (this.conversationType === "private") {
        if (!this.selectedMembers.length) {
          alert("Please pick a user");
          return;
        }
      } else {
        if (!this.groupName || data.message == "") {
          alert("Please insert a Groupname and a message");
          return;
        }

        if (this.selectedMembers.length + 1 <= 2) {
          alert("Members should be more than 2");
          return;
        }
      }

      data = {
        ...data,
        conversationType: this.conversationType,
        groupName: this.groupName,
        groupImage: this.groupImageToUpload,
        members: this.selectedMembers.map((member) => member.userID),
      };

      const res = await this.conversationStore.addConversation(data);
      this.resetFields();
      if (res) this.$emit("add-conversation", res.conversationID);
    },

    resetFields() {
      this.groupName = null;
      this.groupImage = null;
      this.groupImageToUpload = null;
      this.searchMembers = "";
      this.suggestedMembers = [];
      this.selectedMembers = [];
    },
  },
};
</script>

<template>
  <div class="p-2 d-flex flex-column flex-grow-1 bg-light">
    <div
      class="d-flex justify-content-start align-items-center gap-3 p-2 mb-2 border-bottom"
    >
      <button
        type="button"
        class="btn btn-primary p-1 rounded-2 text-white"
        @click="
          () => {
            conversationType = 'private';
            resetFields();
          }
        "
      >
        Private
      </button>
      <button
        type="button"
        class="btn btn-secondary p-1 rounded-2 text-white"
        @click="
          () => {
            conversationType = 'group';
            resetFields();
          }
        "
      >
        Group
      </button>
    </div>

    <!-- Private Input -->
    <div v-if="conversationType === 'private'" class="flex-grow-1">
      <input
        type="text"
        class="form-control p-1 mb-2"
        id="userMember"
        v-model="searchMembers"
        :placeholder="`${
          this.conversationType === 'private' && this.selectedMembers.length
            ? this.selectedMembers[0].username
            : 'Username'
        }`"
        @focus="getUsers"
      />

      <div v-if="users.loading">
        <LoadingSpinner />
      </div>

      <div v-else class="col-3 d-flex flex-column">
        <div
          role="button"
          class="d-flex gap-3 p-2 hover-bg-light rounded-2"
          v-for="user in suggestedMembers"
          :key="user.userID"
          @click="handleAddMember(user)"
        >
          <img
            v-if="user.image"
            :src="`${apiUrl}${user.image}`"
            alt=""
            class="profile-image"
          />
          <i v-else class="bi bi-person-circle fs-3"></i>
          <p>{{ user.username }}</p>
        </div>
      </div>
    </div>

    <!-- Group Input -->
    <div v-else-if="conversationType === 'group'" class="flex-grow-1">
      <button
        type="button"
        class="btn btn-info p-1 rounded-2 text-white mb-2"
        @click="openGroupImageUpload"
      >
        Add Group Image
      </button>
      <input
        @change="handleGroupImageUpload"
        type="file"
        name="image"
        id="group-image"
        class="d-none"
        ref="group_image"
      />
      <!-- Group Image Preview -->
      <div v-if="groupImage" class="ms-4">
        <i
          role="button"
          class="position-absolute bi bi-x fs-5 text-danger bg-white rounded-circle"
          @click="closeGroupImageUpload"
        ></i>
        <img :src="groupImage" alt="GroupImage" width="150" />
      </div>

      <!-- Group Name -->
      <input
        type="text"
        class="form-control p-1 mb-2"
        v-model="groupName"
        id="groupname"
        placeholder="Groupname"
      />

      <!-- Group Members -->
      <input
        type="text"
        class="form-control p-1 mb-3"
        id="members"
        v-model="searchMembers"
        placeholder="Pick members"
        @focus="getUsers"
      />
      <div v-if="users.loading">
        <LoadingSpinner />
      </div>
      <div v-else class="row">
        <!-- Suggested Members -->
        <div class="col-3 d-flex flex-column">
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
        </div>

        <!-- Selected Members -->
        <div v-if="selectedMembers.length" class="col">
          <p class="text-center fw-bold">Selected Members</p>
          <div
            role="button"
            class="d-flex justify-content-center gap-3 p-2 hover-bg-light rounded-2"
            v-for="user in selectedMembers"
            :key="user.userID"
            @click="handleRemoveSelected(user)"
          >
            <img v-if="user.image" :src="`${apiUrl}${user.image}`" alt="" />
            <i v-else class="bi bi-person-circle fs-3"></i>
            <p>{{ user.username }}</p>
          </div>
        </div>
      </div>
    </div>

    <div v-if="conversations.error">
      <ErrorMsg :msg="conversations.error" />
    </div>

    <MessageForm
      @add-conversation="handleAddConversation"
      :newConversation="true"
    />
  </div>
</template>
