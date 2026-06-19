import { Card, Row, Col, Statistic, Typography, Tag } from 'antd'
import {
  DollarOutlined, FileTextOutlined, TeamOutlined, ToolOutlined,
} from '@ant-design/icons'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { listReports, listWageSummaries } from '../api'
import dayjs from 'dayjs'

const { Title } = Typography

export default function Dashboard() {
  const [stats, setStats] = useState({ reportCount: 0, summaryCount: 0 })
  const navigate = useNavigate()

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [reportRes, summaryRes] = await Promise.all([
          listReports({ page: 1, pageSize: 1 }),
          listWageSummaries({ page: 1, pageSize: 1 }),
        ])
        setStats({
          reportCount: reportRes.data?.total || 0,
          summaryCount: summaryRes.data?.total || 0,
        })
      } catch {}
    }
    fetchStats()
  }, [])

  const userStr = localStorage.getItem('user')
  const user = userStr ? JSON.parse(userStr) : {}
  const roleLabel = user.role === 9 ? '管理员' : user.role === 2 ? '核算员' : '工人'
  const currentMonth = dayjs().format('YYYY-MM')

  return (
    <div>
      <Title level={4}>工作台</Title>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic title="当前用户" value={user.realName || user.username}
              prefix={<TeamOutlined />} />
            <Tag color="blue" style={{ marginTop: 8 }}>{roleLabel}</Tag>
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="当前月份" value={currentMonth}
              prefix={<FileTextOutlined />} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="报工单总数" value={stats.reportCount}
              prefix={<FileTextOutlined />} valueStyle={{ color: '#1677ff' }} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="工资汇总数" value={stats.summaryCount}
              prefix={<DollarOutlined />} valueStyle={{ color: '#cf1322' }} />
          </Card>
        </Col>
      </Row>

      <Card title="快捷操作">
        <Row gutter={16}>
          <Col span={8}>
            <Card hoverable style={{ textAlign: 'center' }}
              onClick={() => navigate('/reports')}>
              <FileTextOutlined style={{ fontSize: 32, color: '#1677ff' }} />
              <p style={{ marginTop: 8 }}>报工管理</p>
            </Card>
          </Col>
          <Col span={8}>
            <Card hoverable style={{ textAlign: 'center' }}
              onClick={() => navigate('/wage')}>
              <DollarOutlined style={{ fontSize: 32, color: '#cf1322' }} />
              <p style={{ marginTop: 8 }}>工资查询</p>
            </Card>
          </Col>
          <Col span={8}>
            <Card hoverable style={{ textAlign: 'center' }}
              onClick={() => navigate('/prices')}>
              <ToolOutlined style={{ fontSize: 32, color: '#52c41a' }} />
              <p style={{ marginTop: 8 }}>单价管理</p>
            </Card>
          </Col>
        </Row>
      </Card>
    </div>
  )
}
