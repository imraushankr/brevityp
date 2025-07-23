import { createBrowserRouter, createRoutesFromElements, Route } from "react-router-dom"
import { MainLayout } from '../layouts';

const MainRoutes = createBrowserRouter(createRoutesFromElements(
  <>
  <Route element={<MainLayout/>}>
    <Route index element={<div>Home Page</div>}/>
    <Route path='/about' element={<div>About Page</div>}/>
    <Route path='/about' element={<div>About Page</div>}/>
    <Route path='/contact' element={<div>Contact Page</div>}/>
  </Route>
  </>
));

export default MainRoutes