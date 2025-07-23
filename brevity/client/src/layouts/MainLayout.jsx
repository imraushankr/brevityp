import { Outlet } from 'react-router-dom'
import { Footer, Header } from '../components'

const MainLayout = () => {
  return (
    <div>
      <Header/>
      <div>
        <Outlet/>
      </div>
      <Footer/>
    </div>
  )
}

export default MainLayout