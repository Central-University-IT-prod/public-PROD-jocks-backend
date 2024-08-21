import { FunctionComponent } from "react";
import Discount from "./Discount";

interface Props {}

const Advertising: FunctionComponent<Props> = () => {
  return (
    <article className="grid gap-12">
      <h1 className="font-size text-3xl text-center">
        Скачайте наше приложение, и получите <br />
        <span className="font-black text-4xl">купон</span>:
      </h1>
      <Discount />
      <button className="w-full py-3 bg-primary text-white font-medium rounded-3xl shadow-md" onClick={() => {window.location.href = 'https://github.com/Central-University-IT-prod/PROD-jocks-client-mobile/releases/download/apkv2/app-release.apk'}}>
        Скачать
      </button>
    </article>
  );
};

export default Advertising;
