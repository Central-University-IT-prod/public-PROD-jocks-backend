import Advertising from "../components/Advertising";
import { Link } from "react-router-dom";


const Start = () => {

  return (
    <section className="min-h-screen flex flex-col items-center justify-center gap-12 p-6">
      <Advertising />
      <Link to='business' className="w-full py-3 bg-gray-100 text-gray-200 font-medium rounded-3xl shadow-md text-center">
        Оставить отзыв без награды
      </Link>
    </section>
  );
};

export default Start;
