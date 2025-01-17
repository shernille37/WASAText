<script>
export default {
  name: "EmojiPicker",
};
</script>

<template>
  <div id="emoji-picker" class="d-flex">
    <span>💙</span>
    <span>👍</span>
    <span>👎</span>
    <span>😂</span>
    <!-- Add more emojis here -->
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
