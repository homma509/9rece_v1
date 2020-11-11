<template>
  <v-container>
    <v-navigation-drawer app v-model="drawer" clipped>
      <v-container>
        <v-list-item @click="drawer = !drawer">
          <v-list-item-icon>
            <v-icon>mdi-menu</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title class="title grey--text text--darken-2"
              >9rece</v-list-item-title
            >
          </v-list-item-content>
        </v-list-item>

        <v-divider></v-divider>

        <v-list nav dense>
          <template v-for="nav_list in nav_lists">
            <v-list-item
              v-if="!nav_list.lists"
              :to="nav_list.link"
              :key="nav_list.name"
              @click="menu_close"
            >
              <v-list-item-icon>
                <v-icon>{{ nav_list.icon }}</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>{{ nav_list.name }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-group
              v-else
              no-action
              :prepend-icon="nav_list.icon"
              :key="nav_list.name"
              v-model="nav_list.active"
            >
              <template v-slot:activator>
                <v-list-item-content>
                  <v-list-item-title>{{ nav_list.name }}</v-list-item-title>
                </v-list-item-content>
              </template>
              <v-list-item
                v-for="list in nav_list.lists"
                :key="list.name"
                :to="list.link"
              >
                <v-list-item-title>{{ list.name }}</v-list-item-title>
              </v-list-item>
            </v-list-group>
          </template>
        </v-list>
      </v-container>
    </v-navigation-drawer>
    <v-app-bar app dark clipped-left dense>
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <!-- <v-btn to="/home" text="true"> -->
      <v-toolbar-title>9rece</v-toolbar-title>
      <!-- </v-btn> -->
      <v-spacer></v-spacer>
      <v-toolbar-items>
        <v-menu offset-y>
          <template v-slot:activator="{ on }">
            <v-btn v-on="on" text>
              <v-icon>mdi-account</v-icon>
              <v-icon>mdi-menu-down</v-icon>
            </v-btn>
          </template>
          <v-list dense>
            <v-subheader>アカウント: k-homma</v-subheader>
            <v-divider></v-divider>
            <v-list-item
              v-for="support in supports"
              :key="support.name"
              :to="support.link"
            >
              <v-list-item-icon>
                <v-icon>{{ support.icon }}</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>{{ support.name }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-divider></v-divider>
            <v-list-item to="/signout">
              <v-list-item-icon>
                <v-icon>mdi-logout</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>サインアウト</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-toolbar-items>
    </v-app-bar>
    <v-content>
      <router-view />
    </v-content>
    <v-footer dark app>9rece</v-footer>
  </v-container>
</template>

<script>
export default {
  name: "Top",

  components: {},

  methods: {
    menu_close() {
      this.nav_lists.forEach((nav_list) => {
        nav_list.active = false;
      });
    },
  },

  data: () => ({
    drawer: null,
    nav_lists: [
      {
        name: "ホーム",
        icon: "mdi-home",
        link: "/home",
      },
      {
        name: "データ入力",
        icon: "mdi-database-import",
        active: false,
        link: "",
        lists: [
          { name: "レセプト", link: "/import/uke" },
          { name: "EFファイル", link: "/import/ef" },
          { name: "施設", link: "/import/facility" },
        ],
      },
      { name: "Dashboard", icon: "mdi-view-dashboard", link: "/dashboard" },
      {
        name: "レポート",
        icon: "mdi-file-document",
        active: false,
        link: "",
        lists: [
          { name: "患者別明細リスト", link: "/report/list" },
          { name: "入院基本料別サマリ", link: "/report/summary" },
        ],
      },
      {
        name: "分析",
        icon: "mdi-chart-areaspline-variant",
        active: false,
        link: "",
        lists: [
          { name: "年間推移", link: "/analystic/graph" },
          { name: "ピボット", link: "/analystic/pivot" },
        ],
      },
    ],
    supports: [
      {
        name: "マイアカウント",
        icon: "mdi-account-cog",
        link: "/account",
      },
      {
        name: "設定",
        icon: "mdi-cog",
        link: "/settings",
      },
    ],
  }),
};
</script>
