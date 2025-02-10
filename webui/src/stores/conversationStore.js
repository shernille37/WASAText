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
  conversation: {
    data: null,
    loading: true,
    error: null,
  },

  addConversationFlag: false,
  addMemberFlag: false,

  async getConversations() {
    try {
      this.conversations.loading = true;
      const res = await axios.get("/conversations", {
        headers: {
          Authorization: `Bearer ${authStore.user.data.userID}`,
        },
      });

      console.log(res.data);
      this.conversations.loading = false;
      this.conversations.data = res.data;
    } catch (error) {
      this.conversations.loading = false;
      this.conversations.error = error.toString();
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
      return resAddConversation.data;
    } catch (error) {
      this.conversations.loading = false;
      this.conversations.error = error.response.data;
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

      this.conversation.data.members.push(...res.data);

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
});
