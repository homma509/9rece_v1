import Vue from "vue";
import Vuetify from "vuetify/lib";
import ja from "vuetify/es5/locale/ja.js";

Vue.use(Vuetify);

export default new Vuetify({
  lang: {
    locales: { ja },
    current: "ja",
  },
});
