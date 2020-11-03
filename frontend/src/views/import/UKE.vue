<template>
  <v-container>
    <v-row>
      <v-col cols="10" offset="1">
        <v-file-input
          chips
          multiple
          accept="text/plain,.uke"
          label="UKEファイルを選択してください"
          v-model="file"
          @change="onChangeSelected"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col offset="1">
        <v-btn
          :loading="loading"
          :disabled="file.length == 0"
          color="grey"
          class="ma-2 white--text"
          @click="putFiles"
        >
          アップロード
          <v-icon right dark>mdi-cloud-upload</v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-row v-if="message">
      <v-col cols="10" offset="1">
        <v-alert dense :type="type">
          {{ message }}
        </v-alert>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { Storage } from "aws-amplify";

export default {
  data: () => ({
    file: [],
    loading: false,
    message: null,
    type: null
  }),

  methods: {
    onChangeSelected() {
      this.message = null;
    },
    putFile(file) {
      return new Promise((resolve, reject) => {
        Storage.put("uke/" + file.name, file, {
          level: "private",
          contentType: "text/plain"
        })
          .then(result => {
            console.log(result);
            resolve();
          })
          .catch(error => {
            console.log(error);
            reject(error);
          });
      });
    },
    async putFiles() {
      this.message = null;
      if (!this.file) {
        return;
      }
      this.loading = true;

      try {
        for (let f of this.file) {
          await this.putFile(f);
        }
        this.message = "ファイルをアップロードしました.";
        this.type = "info";
        this.file = [];
      } catch (error) {
        this.message = "ファイルのアップロードに失敗しました. " + error;
        this.type = "error";
      }

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
