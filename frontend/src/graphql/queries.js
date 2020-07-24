import { API, graphqlOperation } from "aws-amplify";

export const GET_RECES = {
  async getReces(facilityID, yearMonth) {
    const variables = {
      filter: {
        ID: {
          eq: `${facilityID}:${yearMonth}`,
        },
      },
    };
    const listReces = `
      query listReces($filter: TableReceFilterInput) {
        listReces(filter: $filter) {
          items {
              FacilityID
              CaredOn
              ClientID
              BuildingID
              PlanClass
              Value
          }
        }
      }
    `;
    const reces = await API.graphql(graphqlOperation(listReces, variables));
    return reces.data.listReces.items;
  },
};

export const GET_FACILITIES = {
  async getFacilities() {
    const variables = {
      filter: {
        DataType: {
          eq: "FacilityInfo",
        },
      },
    };
    const listReces = `
      query listReces($filter: TableReceFilterInput) {
        listReces(filter: $filter) {
          items {
              FacilityID
              FacilityName
          }
        }
      }
    `;
    const reces = await API.graphql(graphqlOperation(listReces, variables));
    return reces.data.listReces.items.sort((a, b) => {
      if (a.FacilityID < b.FacilityID) return -1;
      if (a.FacilityID > b.FacilityID) return 1;
      return 0;
    });
  },
};
