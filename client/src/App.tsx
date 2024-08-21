import { BrowserRouter, Route, Routes } from "react-router-dom";
import Start from "./screens/Start";
import OverallRating from "./screens/OverallRating";
import ProductRating from "./screens/ProductRating";
import End from "./screens/End";
import Screen404 from "./screens/404";
import Loader from "./screens/Loader";
import LoaderSend from "./screens/LoaderSend";

function App() {
  return (
    <BrowserRouter>
      <main className="sm:max-w-96 sm:mx-auto">
        <Routes>
          <Route path="/check/:id" element={<Loader />} />
          <Route path="/" element={<Start />} />
          <Route path="/business" element={<OverallRating />} />
          <Route path="/products" element={<ProductRating />} />
          <Route path="/sending" element={<LoaderSend />} />
          <Route path="/good-job" element={<End />} />
          <Route path="*" element={<Screen404 />} />
        </Routes>
      </main>
    </BrowserRouter>
  );
}

export default App;
