<script>
export default {
  name: "Conversation",
  data() {
    return {
      apiUrl: import.meta.env.VITE_API_URL,
    };
  },
  props: {
    conversation: Object,
  },
  computed: {
    isPrivate() {
      return !!this.conversation.private;
    },
    hasImage() {
      if (this.isPrivate) {
        return !!this.conversation.private.user.image;
      } else {
        return !!this.conversation.group.groupImage;
      }
    },
  },
};
</script>

<template>
  <div
    role="button"
    class="d-flex p-2 justify-content-center justify-content-md-between align-items-center rounded-3 my-2 hover-bg-light"
  >
    <img
      v-if="hasImage"
      :src="`${apiUrl}${
        isPrivate
          ? conversation.private.user.image
          : conversation.group.groupImage
      }`"
      alt="Profile Picture"
      class="profile-image"
    />
    <i v-else class="bi bi-person-circle fs-3"></i>

    <div class="d-none d-md-flex">
      <div>
        <p class="fw-semibold">
          {{
            isPrivate
              ? conversation.private.user.username
              : conversation.group.groupName
          }}
        </p>

        <i v-if="conversation.latestMessage.image" class="fs-5 bi bi-image"></i>
        <p v-else class="fs-6">
          {{
            conversation.latestMessage.message.length <= 20
              ? conversation.latestMessage.message
              : conversation.latestMessage.message.slice(0, 20) + "..."
          }}
        </p>

        <div class="info d-flex justify-content-between">
          <p
            :class="[
              'p-1',
              'rounded-2',
              'text-white',
              isPrivate ? 'bg-primary' : 'bg-secondary',
            ]"
          >
            {{ conversation.type }}
          </p>
          <p class="ms-5">{{ conversation.latestMessage.timestamp }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.info {
  font-size: 10px;
}
</style>
