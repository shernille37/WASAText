import { reactive } from "vue";
import axios from "../services/axios";

import { authStore } from "../stores/authStore";
import { uploadImage } from "../utils/upload";

export const conversationStore = reactive({
  conversations: {
    data: [],
    loading: true,
    error: null,
  },
  privateConversations: {
    data: [],
    loading: true,
    error: null,
  },
  groupConversations: {
    data: [],
    loading: true,
    error: null,
    leaveGroupConversationError: null,
  },
  conversation: {
    data: null,
    loading: true,
    error: null,
  },
  conversationMembers: {
    data: [],
    loading: true,
    error: null,
  },

  selectedConversation: false,
  addConversationFlag: false,
  addMemberFlag: false,
  membersListFlag: false,

  resetFields() {
    this.selectedConversation = null;
    this.addConversationFlag = false;
    this.addMemberFlag = false;
  },

  async getConversations() {
    try {
      this.conversations.loading = true;
      const res = await axios.get("/conversations", {
        headers: {
          Authorization: `Bearer ${authStore.user.data.userID}`,
        },
      });

      // Update messages to delivered!
      await this.updateMessageToDelivered();

      this.conversations.loading = false;
      this.conversations.data = res.data;
    } catch (error) {
      this.conversations.loading = false;
      this.conversations.error = error.response.data;
    }
  },

  async getConversation(id) {
    try {
      this.conversation.loading = true;
      const res = await axios.get(`/conversations/${id}`, {
        headers: {
          Authorization: `Bearer ${authStore.user.data.userID}`,
        },
      });

      this.conversation.loading = false;
      this.conversation.data = res.data;
    } catch (error) {
      this.conversation.loading = false;
      this.conversation.error = error.response.data;
    }
  },

  async getPrivateConversations() {
    try {
      this.privateConversations.loading = true;
      const res = await axios.get("/private-conversations", {
        headers: {
          Authorization: `Bearer ${authStore.user.data.userID}`,
        },
      });

      // Update messages to delivered!
      await this.updateMessageToDelivered();

      this.privateConversations.loading = false;
      this.privateConversations.data = res.data;
    } catch (error) {
      this.privateConversations.loading = false;
      this.privateConversations.error = error.response.data;
    }
  },

  async getGroupConversations() {
    try {
      this.groupConversations.loading = true;
      const res = await axios.get("/group-conversations", {
        headers: {
          Authorization: `Bearer ${authStore.user.data.userID}`,
        },
      });

      // Update messages to delivered!
      await this.updateMessageToDelivered();

      this.groupConversations.loading = false;
      this.groupConversations.data = res.data;
    } catch (error) {
      this.groupConversations.loading = false;
      this.groupConversations.error = error.response.data;
    }
  },

  async addConversation(data) {
    this.conversation.loading = true;

    let resImageUpload = null;
    let resGroupImageUpload = null;

    try {
      if (data.conversationType === "group" && data.groupImage) {
        resGroupImageUpload = await uploadImage(data.groupImage);
      }

      if (data.image) {
        resImageUpload = await uploadImage(data.image);
      }

      const groupData = {
        groupName: data.groupName,
        groupImage: resGroupImageUpload,
        message: data.message,
        image: resImageUpload,
        members: data.members,
      };

      const privateData = {
        receiverID: data.members[0],
        message: data.message,
        image: resImageUpload,
      };

      const resAddConversation = await axios.post(
        `${
          data.conversationType === "group"
            ? "/group-conversations"
            : "/private-conversations"
        }`,
        data.conversationType === "group" ? groupData : privateData,
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${authStore.user.data.userID}`,
          },
        }
      );

      this.conversations.loading = false;
      this.conversations.data.push(resAddConversation.data);

      if (data.conversationType === "group")
        this.groupConversations.data.push(resAddConversation.data);
      else this.privateConversations.data.push(resAddConversation.data);

      return resAddConversation.data;
    } catch (error) {
      this.conversations.loading = false;
      this.conversations.error = error.response.data;
    }
  },

  async getMembers(conversationID) {
    try {
      this.conversationMembers.loading = true;
      const resMembers = await axios.get(
        `/group-conversations/${conversationID}/members`,
        {
          headers: {
            Authorization: `Bearer ${authStore.user.data.userID}`,
          },
        }
      );

      this.conversationMembers.loading = false;
      this.conversationMembers.data = resMembers.data;
    } catch (error) {
      this.conversationMembers.loading = false;
      this.conversationMembers.error = error.response.data;
    }
  },

  async addMembersToGroup(data) {
    try {
      this.conversation.loading = true;

      const res = await axios.post(
        `/group-conversations/${this.conversation.data.conversationID}/members`,
        data,
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${authStore.user.data.userID}`,
          },
        }
      );

      this.conversationMembers.data.push(...res.data);

      this.conversation.loading = false;
      this.addMemberFlag = false;
    } catch (error) {
      this.conversation.loading = false;
      this.conversation.error = error.response.data;
    }
  },

  async editGroupConversation(groupImage, groupName) {
    let promises = [];
    let resUploadImage = null;
    try {
      if (groupImage) {
        resUploadImage = await uploadImage(groupImage);
        const updateGroupImage = axios.put(
          `/group-conversations/${this.conversation.data.conversationID}/photo`,
          {
            groupImage: resUploadImage,
          },
          {
            headers: {
              Authorization: `Bearer ${authStore.user.data.userID}`,
            },
          }
        );

        promises.push(updateGroupImage);
      }
      if (groupName) {
        const updateGroupName = axios.put(
          `/group-conversations/${this.conversation.data.conversationID}/name`,
          {
            groupName: groupName,
          },
          {
            headers: {
              Authorization: `Bearer ${authStore.user.data.userID}`,
            },
          }
        );

        promises.push(updateGroupName);
      }

      await Promise.all(promises);
    } catch (error) {
      this.conversation.error = error.response.data;
    }
  },

  async leaveGroupConversation(conversationID) {
    try {
      await axios.delete(
        `/group-conversations/${this.conversation.data.conversationID}/members`,
        {
          headers: {
            Authorization: `Bearer ${authStore.user.data.userID}`,
          },
        }
      );

      this.conversations.data = this.conversations.data.filter(
        (conversation) => conversation.conversationID !== conversationID
      );
      this.groupConversations.data = this.groupConversations.data.filter(
        (conversation) => conversation.conversationID !== conversationID
      );
      this.selectedConversation = null;
    } catch (error) {
      this.groupConversations.leaveGroupConversationError = error.response.data;
    }
  },

  async updateMessageToDelivered() {
    for (const conversation of this.conversations.data) {
      try {
        await axios.put(
          `/conversations/${conversation.conversationID}/messages/deliver`,
          {},
          {
            headers: {
              Authorization: `Bearer ${authStore.user.data.userID}`,
            },
          }
        );
      } catch (error) {
        throw new Error(error);
      }
    }
  },

  async updateMessageToRead() {
    for (const conversation of this.conversations.data) {
      try {
        await axios.put(
          `/conversations/${conversation.conversationID}/messages/read`,
          {},
          {
            headers: {
              Authorization: `Bearer ${authStore.user.data.userID}`,
            },
          }
        );
      } catch (error) {
        throw new Error(error);
      }
    }
  },
});
