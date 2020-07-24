import Vue from "vue";
import VueRouter from "vue-router";
import Signin from "@/views/Signin";
import Top from "@/views/Top";
import Home from "@/views/Home";
import EF from "@/views/import/EF";
import Facility from "@/views/import/Facility";
import List from "@/views/report/List";
import Graph from "@/views/analystic/Graph";
import store from "@/store";

// Amplify読み込み
import { Auth } from "aws-amplify";
import { AmplifyEventBus } from "aws-amplify-vue";

Vue.use(VueRouter);

const signout = (to, from, next) => {
  return Auth.signOut()
    .then(() => {
      next("/signin");
    })
    .catch((error) => {
      console.log(error);
    });
};

let user;

function getUser() {
  return Auth.currentAuthenticatedUser()
    .then((data) => {
      if (data && data.signInUserSession) {
        store.commit("setUser", data);
        return data;
      }
    })
    .catch(() => {
      store.commit("setUser", null);
      return null;
    });
}

// ログイン状態管理
AmplifyEventBus.$on("authState", async (state) => {
  if (state === "signedOut") {
    user = null;
    store.commit("setUser", null);
    router.push({ path: "/signin" });
  } else if (state === "signedIn") {
    user = await getUser();
    router.push({ path: "/home" });
  }
});

const routes = [
  {
    path: "/signin",
    name: "signin",
    component: Signin,
  },
  {
    path: "/signout",
    beforeEnter: signout,
  },
  {
    path: "/",
    name: "top",
    component: Top,
    meta: { requiresAuth: true },
    children: [
      {
        path: "/home",
        name: "home",
        component: Home,
      },
      {
        path: "import/ef",
        name: "ef",
        component: EF,
      },
      {
        path: "import/facility",
        name: "facility",
        component: Facility,
      },
      {
        path: "report/list",
        name: "list",
        component: List,
      },
      {
        path: "analystic/graph",
        name: "graph",
        component: Graph,
      },
    ],
  },
  {
    path: "*",
    redirect: "/home",
  },
];

// TODO setting BASE_URL
const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

router.beforeResolve(async (to, from, next) => {
  if (to.matched.some((record) => record.meta.requiresAuth)) {
    user = await getUser();
    if (!user) {
      return next({
        path: "/signin",
      });
    }
    return next();
  }
  return next();
});

export default router;
