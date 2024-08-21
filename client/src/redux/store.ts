import { configureStore } from '@reduxjs/toolkit';
import { DataReducer } from './index';

export const store = configureStore({
  reducer: {
    data: DataReducer,
  }
});

export type storeDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;