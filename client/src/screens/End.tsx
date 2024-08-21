import { FunctionComponent } from "react";
import Advertising from "../components/Advertising";

interface Props {}

const End: FunctionComponent<Props> = (props) => {
  return (
    <section className="min-h-screen flex flex-col items-center justify-center gap-12 p-6">
      <h2 className="font-size text-3xl text-center">Отзыв отправлен</h2>
      <Advertising />
    </section>
  );
};

export default End;
