import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import vuetify from "./plugins/vuetify";

// Amplify読み込み
import Amplify, * as AmplifyModules from "aws-amplify";
import { AmplifyPlugin } from "aws-amplify-vue";
import awsmobile from "./aws-exports";

Amplify.configure(awsmobile);
Vue.use(AmplifyPlugin, AmplifyModules);

Vue.config.productionTip = false;

let messageResource = {
  ja: {
    "Create Account": "",
    "Have an account? ": "",
    "Sign in": "",
    Username: "ユーザー名",
    "Enter your Username": "ユーザー名もしくはemailアドレス",
    Password: "パスワード",
    "Enter your password": "パスワード",
    "Forget your password? ": "パスワードを忘れた方",
    "Reset password": "パスワードリセット",
    "No account? ": "会員登録",
    "Create account": "アカウント作成",
    "Sign in to your account": "サインイン",
    "Reset your password": "パスワードリセット",
    "Sign In": "サインイン",
    "Send Code": "リセット",
    "Back to Sign In": "戻る",
    "Enter new password": "新しいパスワードの設定",
    "New Password": "新しいパスワード",
    Submit: "設定",
  },
};
AmplifyModules.I18n.putVocabularies(messageResource);

new Vue({
  router,
  store,
  vuetify,
  render: (h) => h(App),
}).$mount("#app");
