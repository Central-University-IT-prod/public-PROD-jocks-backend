import { RootState } from './store';

const selectDataModule = (state: RootState): RootState['data'] => state.data;

export const selectData = (state: RootState) => selectDataModule(state);

export const selectReceivedData = (state: RootState) => selectData(state)['receivedData'];
export const selectBaseFields = (state: RootState) => selectReceivedData(state)['base_review_fields'];
export const selectProductFields = (state: RootState) => selectReceivedData(state)['products_review_fields'];

export const selectReviewsData =  (state: RootState) => selectData(state)['reviewesData'];
export const selectBaseReviews = (state: RootState) => selectReviewsData(state)['base_review_fields'];
export const selectProductReviews = (state: RootState) => selectReviewsData(state)['products_review_fields'];

