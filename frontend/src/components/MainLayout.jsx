import { Layout, Menu, Dropdown, Avatar, Typography, message } from 'antd'
import {
  DashboardOutlined,
  DollarOutlined,
  ToolOutlined,
  FileTextOutlined,
  TeamOutlined,
  UserOutlined,
  LogoutOutlined,
} from '@ant-design/icons'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { useState } from 'react'

const { Header, Sider, Content } = Layout
const { Text } = Typography

const menuItems = [
  { key: '/', icon: <DashboardOutlined />, label: '工作台' },
  { key: '/prices', icon: <DollarOutlined />, label: '单价管理' },
  { key: '/reports', icon: <FileTextOutlined />, label: '报工管理' },
  { key: '/wage', icon: <DollarOutlined />, label: '工资查询' },
  { key: '/process', icon: <ToolOutlined />, label: '工序产品' },
  { key: '/teams', icon: <TeamOutlined />, label: '班组管理' },
]

export default function MainLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const [collapsed, setCollapsed] = useState(false)

  const userStr = localStorage.getItem('user')
  const user = userStr ? JSON.parse(userStr) : {}

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    message.success('已退出登录')
    navigate('/login')
  }

  const dropdownItems = {
    items: [
      { key: 'info', icon: <UserOutlined />, label: `${user.realName || user.username}` },
      { key: 'role', label: `角色: ${user.role === 9 ? '管理员' : user.role === 2 ? '核算员' : '工人'}`, disabled: true },
      { type: 'divider' },
      { key: 'logout', icon: <LogoutOutlined />, label: '退出登录', danger: true },
    ],
    onClick: ({ key }) => { if (key === 'logout') handleLogout() },
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider collapsible collapsed={collapsed} onCollapse={setCollapsed}
        style={{ background: '#001529' }}>
        <div style={{ height: 48, margin: 12, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <Text style={{ color: '#fff', fontSize: collapsed ? 14 : 16, fontWeight: 'bold', whiteSpace: 'nowrap' }}>
            {collapsed ? '计件' : '计件薪资核算'}
          </Text>
        </div>
        <Menu theme="dark" mode="inline" selectedKeys={[location.pathname]}
          items={menuItems} onClick={({ key }) => navigate(key)} />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', display: 'flex',
          justifyContent: 'flex-end', alignItems: 'center', boxShadow: '0 1px 4px rgba(0,0,0,0.08)' }}>
          <Dropdown menu={dropdownItems} placement="bottomRight">
            <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 8 }}>
              <Avatar icon={<UserOutlined />} style={{ background: '#1677ff' }} />
              <span>{user.realName || user.username}</span>
            </div>
          </Dropdown>
        </Header>
        <Content style={{ margin: 16, padding: 20, background: '#fff', borderRadius: 8,
          minHeight: 'calc(100vh - 112px)', overflow: 'auto' }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
