import { createSlice } from "@reduxjs/toolkit";

interface Data {
  receivedData: {
    checl_id: number;
    business_id: number;
    base_review_fields: Field[];
    products_review_fields: Field[];
  };
  reviewesData: {
    base_review_fields: { [key: number]: number };
    products_review_fields: { [key: number]: number };
  };
}
interface Field {
  id: number;
  name: string;
  type: string;
}

const DataSlice = createSlice({
  name: "habits",
  initialState: {
    receivedData: {},
    reviewesData: {
      base_review_fields: {},
      products_review_fields: {},
    },
  } as Data,

  reducers: {
    setData: (state, action) => {
      const data = action.payload;
      state["receivedData"] = data;
    },
    setRating: (state, action) => {
      const {
        isProduct,
        id,
        value,
      }: { isProduct?: boolean; id: number; value: number } = action.payload;
      if (isProduct) {
        state["reviewesData"]["products_review_fields"][id] = value;
      } else {
        state["reviewesData"]["base_review_fields"][id] = value;
      }
    },
  },
});

export const DataReducer = DataSlice.reducer;
export const DataActions = DataSlice.actions;
