import { FunctionComponent } from "react";

interface Props {}

const Screen404: FunctionComponent<Props> = (props) => {
  return (
    <section className="min-h-screen flex flex-col items-center justify-center p-6">
      <h1 className="font-size text-3xl text-center">404 страница не найдена :(</h1>
    </section>
  );
};

export default Screen404;
