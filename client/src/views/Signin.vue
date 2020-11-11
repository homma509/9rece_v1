<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6">
        <!--ログインコンポーネント-->
        <amplify-authenticator></amplify-authenticator>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: "signin",
};
</script>
