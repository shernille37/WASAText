<script>
export default {
  name: "ErrorMsg",
  props: {
    msg: String,
  },
};
</script>

<template>
  <div class="alert alert-danger p-1" role="alert">
    {{ msg }}
  </div>
</template>

<style></style>
