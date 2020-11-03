<template>
  <v-container>
    <v-row>
      <v-col cols="10" offset="1">
        <v-file-input
          accept=".txt"
          label="施設ファイルを選択してください"
          v-model="file"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col offset="1">
        <v-btn
          :loading="loading"
          :disabled="loading"
          color="grey"
          class="ma-2 white--text"
          @click="upload"
        >
          アップロード
          <v-icon right dark>mdi-cloud-upload</v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-row v-if="errorMessage">
      <v-col cols="10" offset="1">
        <v-alert dense type="warning">
          {{ errorMessage }}
        </v-alert>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
// TODO add csv download
// TODO Upload成功時に数秒間メッセージ表示
import { Storage } from "aws-amplify";

export default {
  data: () => ({
    file: null,
    loading: false,
    errorMessage: null
  }),

  methods: {
    async upload() {
      this.errorMessage = null;
      if (!this.file) {
        return;
      }
      this.loading = true;

      await Storage.put("facility/" + this.file.name, this.file, {
        level: "public"
      })
        .then(result => {
          console.log(result);
          this.file = null;
        })
        .catch(err => {
          console.log(err);
          this.errorMessage = err;
        });

      this.loading = false;
    }
  }
};
</script>

<style>
@-moz-keyframes loader {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(360deg);
  }
}
@-webkit-keyframes loader {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(360deg);
  }
}
@-o-keyframes loader {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(360deg);
  }
}
@keyframes loader {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
