import { createRoot } from "react-dom/client";
import "./index.css";
import MainRoutes from "./routers/MainRoutes.jsx";
import { RouterProvider } from "react-router-dom";

createRoot(document.getElementById("root")).render(<RouterProvider router={MainRoutes}/>);