import { FunctionComponent } from "react";
import { useDispatch, useSelector } from "react-redux";
import { DataActions } from "../redux";
import { RootState } from "../redux/store";
import { selectBaseReviews, selectProductReviews } from "../redux/selectors";

interface Props {
  title: string;
  id: number;
  isProduct?: boolean;
}

interface StarProps {
  rating: number;
  value: number;
  id: number;
  isProduct?: boolean;
}

const Rating: FunctionComponent<Props> = ({ title, id, isProduct }) => {
  const rating = useSelector((state: RootState) => {
    if (isProduct) {
      return selectProductReviews(state);
    } else {
      return selectBaseReviews(state);
    }
  });
  return (
    <div>
      <h6 className="text-3xl font-bold mb-3">{title}</h6>
      <div className="flex gap-2">
        {[...Array(5)].map((star, index) => {
          return (
            <Star
              key={index}
              value={index + 1}
              rating={rating[id] ? rating[id] : 0}
              id={id}
              isProduct={isProduct}
            />
          );
        })}
      </div>
    </div>
  );
};

const Star: FunctionComponent<StarProps> = ({
  value,
  rating,
  id,
  isProduct,
}) => {
  const dispatch = useDispatch();
  return (
    <div
      className="star"
      onClick={() => dispatch(DataActions.setRating({ id, value, isProduct }))}
    >
      <svg
        data-rating={value}
        fill={value <= rating ? "#FFEB39" : "#FFFFFF"}
        height={32}
        viewBox="0 0 34 32"
        width={32}
      >
        <path
          d="M16.0848 1.56886C16.4592 0.810381 17.5408 0.810381 17.9152 1.56886L21.8314 9.5026C22.1284 10.1044 22.7024 10.5218 23.3665 10.6189L32.1259 11.8992C32.9627 12.0215 33.2962 13.0501 32.6904 13.6402L26.3546 19.8112C25.8731 20.2802 25.6534 20.9561 25.767 21.6185L27.262 30.3355C27.405 31.1692 26.5298 31.805 25.7811 31.4113L17.9501 27.293C17.3553 26.9803 16.6447 26.9803 16.0499 27.293L8.21893 31.4113C7.47019 31.805 6.59497 31.1692 6.73798 30.3355L8.23303 21.6185C8.34665 20.9561 8.12689 20.2802 7.64542 19.8112L1.3096 13.6402C0.703783 13.0501 1.0373 12.0215 1.8741 11.8992L10.6335 10.6189C11.2976 10.5218 11.8716 10.1044 12.1686 9.50261L16.0848 1.56886Z"
          stroke="#393939"
          strokeOpacity="0.5"
          strokeWidth="0.5"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
      </svg>
    </div>
  );
};

export default Rating;

