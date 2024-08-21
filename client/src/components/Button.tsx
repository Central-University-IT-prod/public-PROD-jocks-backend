import { FunctionComponent } from "react";

interface Props {
  text: string;
  addStyle?: string;
}

const Button: FunctionComponent<Props> = ({ text, addStyle }) => {
  return (
    <button className={`w-full py-3 bg-primary text-white font-medium rounded-3xl shadow-md ${addStyle}`}>
      {text}
    </button>
  );
};

export default Button;
