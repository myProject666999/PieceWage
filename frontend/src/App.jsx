import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import MainLayout from './components/MainLayout'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import PriceManagement from './pages/PriceManagement'
import ReportManagement from './pages/ReportManagement'
import WageQuery from './pages/WageQuery'
import ProcessProduct from './pages/ProcessProduct'
import TeamManagement from './pages/TeamManagement'

function PrivateRoute({ children }) {
  const token = localStorage.getItem('token')
  return token ? children : <Navigate to="/login" replace />
}

function App() {
  return (
    <ConfigProvider locale={zhCN} theme={{
      token: { colorPrimary: '#1677ff', borderRadius: 6 },
    }}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/" element={<PrivateRoute><MainLayout /></PrivateRoute>}>
            <Route index element={<Dashboard />} />
            <Route path="prices" element={<PriceManagement />} />
            <Route path="reports" element={<ReportManagement />} />
            <Route path="wage" element={<WageQuery />} />
            <Route path="process" element={<ProcessProduct />} />
            <Route path="teams" element={<TeamManagement />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  )
}

export default App
