<template>
  <v-container>
    <v-row>
      <v-col cols="12" sm="4" md="3">
        <v-select
          :items="facilities"
          v-model="facility"
          @change="load"
          label="施設"
          prepend-icon="mdi-hospital-building"
        ></v-select>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" sm="10" md="10" lg="10"><SaleGraph /></v-col>
    </v-row>
  </v-container>
</template>

<script>
import rece from "../../graphql/queries";
import SaleGraph from "@/components/SaleGraph.vue";

export default {
  components: {
    SaleGraph,
  },
  data: () => ({
    loading: false,
    facility: null,
    points: [],
    facilities: [
      { text: "世田谷病院", value: "101" },
      { text: "中野病院", value: "201" },
      { text: "杉並病院", value: "301" },
      { text: "武蔵野病院", value: "401" },
      { text: "三鷹病院", value: "501" },
    ],
  }),

  methods: {
    async load() {
      this.loading = true;
      await rece
        .getReces(this.facility, this.yearMonth)
        .then((data) => {
          this.points = data;
          this.loading = false;
        })
        .catch((err) => {
          console.log(err);
          this.loading = false;
        });
    },
  },
};
</script>
