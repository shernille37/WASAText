<script>
export default {
  name: "ErrorMsg",
  props: {
    msg: String,
  },
};
</script>

<template>
  <p class="alert alert-danger p-1" role="alert">
    {{ msg }}
  </p>
</template>

<style></style>
