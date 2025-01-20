<script>
import { authStore } from "../stores/authStore";
import EmojiPicker from "./EmojiPicker.vue";

export default {
  name: "Message",
  components: {
    EmojiPicker,
  },
  data() {
    return {
      apiUrl: import.meta.env.VITE_API_URL,
      authStore,
      emojiClick: false,
    };
  },
  props: {
    message: Object,
    conversation: Object,
  },
  computed: {
    isOwner() {
      return this.message.senderID === authStore.user.data.userID;
    },
    isRecipientOwner() {
      return this.message.replyRecipientName === authStore.user.data.username;
    },
    isReply() {
      return !!this.message.replyMessageID;
    },
    isGroup() {
      return this.conversation.type === "group";
    },
  },
  methods: {
    toggleEmojiPicker() {
      this.emojiClick = !this.emojiClick;
    },
  },
};
</script>

<template>
  <div
    :class="[
      'message-container',
      'd-flex',
      'align-items-center',
      isOwner ? 'justify-content-end' : 'justify-content-start',
    ]"
  >
    <div class="d-flex flex-column" style="max-width: 50%">
      <p v-show="isGroup && !isOwner" style="font-size: 15px">
        {{ message.senderName }}
      </p>

      <p v-show="isReply" class="text-end">
        {{ isOwner ? "You" : message.senderName }} replied to
        {{ isRecipientOwner ? "You" : message.replyRecipientName }}
      </p>

      <div
        class="bg-light rounded-2 text-black px-1 py-1 text-end"
        v-if="isReply"
      >
        <i v-if="message.hasImage" class="fs-4 bi bi-image"></i>
        <p v-else>
          {{
            message.replyMessage.length <= 20
              ? message.replyMessage
              : message.replyMessage.slice(0, 20) + "..."
          }}
        </p>
      </div>

      <div class="position-relative">
        <!-- Image -->
        <img
          class="image"
          v-if="message.hasImage"
          :src="`${apiUrl}${message.image}`"
          alt="Image"
        />
        <!-- Message -->
        <p
          v-if="message.message !== ''"
          :class="`message rounded-4 px-3 py-1 ${
            isOwner ? 'bg-primary text-white' : 'bg-light text-black'
          }`"
        >
          {{ message.message }}
        </p>

        <!-- Icons -->
        <div
          :class="`icons d-flex position-absolute ${
            isOwner ? 'end-100' : 'start-100'
          }`"
        >
          <i class="bi bi-reply text-black"></i>

          <i @click="toggleEmojiPicker" class="bi bi-emoji-smile"></i>
          <EmojiPicker
            @mouseleave="toggleEmojiPicker"
            v-if="emojiClick"
            :senderID="authStore.user.data.userID"
            :conversationID="conversation.conversationID"
          />

          <i class="bi bi-info-circle"></i>
        </div>

        <!-- Reactions -->
        <div class="position-absolute end-0 d-flex gap-1">
          <div
            class="d-flex justify-content-center align-items-center gap-1 bg-light rounded-2 p-1"
            v-for="(reaction, index) in message.reactions"
            :key="index"
          >
            <p>{{ reaction.unicode }}</p>
            <p>{{ reaction.count }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
p:not(.message) {
  font-style: italic;
  font-size: 12px;
}

.icons {
  width: 0px;
  opacity: 0;
  font-size: 20px;
  bottom: 50%;
  transform: translateY(50%);
  cursor: pointer;
  transition: width 0.5s ease, opacity 0.1s ease;
}

.message:hover ~ .icons,
.image:hover ~ .icons,
.icons:hover {
  width: 100px;
  opacity: 1;
}

.image {
  max-width: 100%;
  max-height: 500px;
  object-fit: contain;
  display: block;
}
</style>
