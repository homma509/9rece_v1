import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    user: null,
  },
  // getters: {
  //   user: (state) => state.user,
  // },
  mutations: {
    // ユーザー情報保存
    setUser(state, user) {
      state.user = user;
    },
  },
});

export default store;
