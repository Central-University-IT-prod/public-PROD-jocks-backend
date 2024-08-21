import { FunctionComponent } from "react";

interface Props {}

const Discount: FunctionComponent<Props> = () => {
  return (
    <article className="px-6 pt-6 pb-3 grid gap-3 border border-gray-200 rounded-2xl">
        <h3 className="text-base">Следующий заказ</h3>
        <h2 className="text-base font-normal">Скидка 100%</h2>
      <hr className="h-px bg-black border-dashed"></hr>
      <div className="flex flex-col text-right">
        <span className="text-xs font-medium">Купон</span>
        <span className="text-sm">777</span>
      </div>
    </article>
  );
};

export default Discount;
