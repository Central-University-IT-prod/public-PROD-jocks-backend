import Rating from "../components/Rating";
import { Link } from "react-router-dom";
import { RootState } from "../redux/store";
import {
  selectProductFields,
} from "../redux/selectors";
import { useSelector } from "react-redux";
const ProductRating = () => {
  const Data = useSelector((state: RootState) => selectProductFields(state));
  return (
    <section className="min-h-screen p-6 grid gap-3">
      <div className="flex flex-col gap-6">
        {Data?.map(({ name, id }) => (
          <Rating title={name} id={id} key={id} isProduct={true}/>
        ))}
      </div>
      <Link
        to="/sending"
        className=" w-full py-3 bg-primary text-white font-medium rounded-3xl shadow-md mt-auto text-center"
      >
        Отправить
      </Link>
    </section>
  );
};

export default ProductRating;
