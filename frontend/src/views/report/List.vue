<template>
  <v-container>
    <v-row>
      <v-col cols="12" sm="6" md="4">
        <v-select
          :items="facilities"
          item-text="FacilityName"
          item-value="FacilityID"
          v-model="facility"
          @change="loadReces"
          label="施設"
          prepend-icon="mdi-hospital-building"
        ></v-select>
      </v-col>
      <v-col cols="12" sm="6" md="4">
        <v-menu
          v-model="menu"
          :close-on-content-click="false"
          :nudge-right="40"
          transition="scale-transition"
          offset-y
          min-width="290px"
        >
          <template v-slot:activator="{ on, attrs }">
            <v-text-field
              v-model="yearMonth"
              label="年月"
              prepend-icon="mdi-calendar"
              readonly
              v-bind="attrs"
              v-on="on"
            ></v-text-field>
          </template>
          <v-date-picker
            dark
            @change="loadReces"
            v-model="yearMonth"
            @input="menu = false"
            type="month"
          ></v-date-picker>
        </v-menu>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-data-table
          :headers="headers"
          :footer-props="{
            'items-per-page-options': [5, 10, 20, 50, 100, -1],
            showFirstLastPage: true,
          }"
          :items="points"
          :loading="loading"
        ></v-data-table>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { GET_RECES, GET_FACILITIES } from "../../graphql/queries";

export default {
  data: () => ({
    loading: false,
    facility: null,
    yearMonth: new Date().toISOString().substr(0, 7),
    menu: false,
    headers: [
      { text: "施設コード", value: "FacilityID" },
      { text: "実施年月日", value: "CaredOn" },
      { text: "データ識別番号", value: "ClientID" },
      { text: "病棟コード", value: "BuildingID" },
      { text: "入院基本料区分", value: "PlanClass" },
      { text: "日当点", value: "Value" }
    ],
    points: [],
    facilities: []
  }),

  async mounted() {
    await GET_FACILITIES.getFacilities()
      .then(data => {
        this.facilities = data;
        console.log(data);
      })
      .catch(err => {
        console.log(err);
      });
  },

  methods: {
    async loadReces() {
      this.loading = true;
      await GET_RECES.getReces(this.facility, this.yearMonth)
        .then(data => {
          this.points = data;
          this.loading = false;
        })
        .catch(err => {
          console.log(err);
          this.loading = false;
        });
    }
  }
};
</script>
