import { useState, useEffect } from 'react'
import { Card, Table, InputNumber, Button, Tag, Space, message, Typography, Row, Col, Statistic, Descriptions, List, Modal } from 'antd'
import { DollarOutlined, EyeOutlined, SyncOutlined } from '@ant-design/icons'
import { listWageSummaries, getMonthlySummary, getRealtimeAccumulate, getWorkerDetails, settleMonth } from '../api'
import dayjs from 'dayjs'

const { Title } = Typography
const { Search } = InputNumber

export default function WageQuery() {
  const [summaries, setSummaries] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [realtimeAmt, setRealtimeAmt] = useState(null)
  const [detailModalOpen, setDetailModalOpen] = useState(false)
  const [detailList, setDetailList] = useState([])
  const [selectedWorker, setSelectedWorker] = useState(null)
  const [selectedMonth, setSelectedMonth] = useState(dayjs().format('YYYY-MM'))

  const userStr = localStorage.getItem('user')
  const user = userStr ? JSON.parse(userStr) : {}

  const fetchSummaries = async (p = page) => {
    setLoading(true)
    try {
      const res = await listWageSummaries({ page: p, pageSize: 20 })
      setSummaries(res.data?.list || [])
      setTotal(res.data?.total || 0)
    } catch (err) {
      message.error('获取工资汇总失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchRealtime = async () => {
    if (!user.userId) return
    try {
      const month = selectedMonth || dayjs().format('YYYY-MM')
      const res = await getRealtimeAccumulate(user.userId, month)
      setRealtimeAmt(res.data?.accumulateAmount || 0)
    } catch {
      setRealtimeAmt(null)
    }
  }

  useEffect(() => { fetchSummaries(); fetchRealtime() }, [])

  const handleSettle = async () => {
    const month = selectedMonth || dayjs().format('YYYY-MM')
    try {
      await settleMonth(month)
      message.success(`${month} 月度结算完成`)
      fetchSummaries()
    } catch (err) {
      message.error(err.message || '结算失败')
    }
  }

  const handleViewDetails = async (workerId, workerName) => {
    setSelectedWorker({ id: workerId, name: workerName })
    try {
      const month = selectedMonth || dayjs().format('YYYY-MM')
      const startDate = `${month}-01`
      const endDate = `${month}-31`
      const res = await getWorkerDetails({ workerId, startDate, endDate })
      setDetailList(res.data || [])
      setDetailModalOpen(true)
    } catch {
      message.error('获取明细失败')
    }
  }

  const columns = [
    { title: '工人', dataIndex: 'worker', width: 120,
      render: (w) => w?.realName || '-' },
    { title: '月份', dataIndex: 'summaryMonth', width: 100 },
    { title: '合格数', dataIndex: 'totalQtyGood', width: 80, align: 'right' },
    { title: '不良数', dataIndex: 'totalQtyDefect', width: 80, align: 'right' },
    { title: '计件总额', dataIndex: 'grossAmount', width: 110, align: 'right',
      render: (v) => `¥${v}` },
    { title: '不良扣款', dataIndex: 'defectAmount', width: 100, align: 'right',
      render: (v) => v > 0 ? <span style={{ color: '#cf1322' }}>-¥{v}</span> : '¥0' },
    { title: '班组分配', dataIndex: 'allocationAmt', width: 100, align: 'right',
      render: (v) => v > 0 ? <span style={{ color: '#52c41a' }}>+¥{v}</span> : '¥0' },
    { title: '应发合计', dataIndex: 'netAmount', width: 120, align: 'right',
      render: (v) => <span style={{ fontWeight: 'bold', color: '#cf1322', fontSize: 15 }}>¥{v}</span> },
    { title: '状态', dataIndex: 'calcStatus', width: 80, align: 'center',
      render: (v) => v === 1 ? <Tag color="green">已结算</Tag> : <Tag color="orange">未结算</Tag> },
    {
      title: '操作', width: 100,
      render: (_, r) => (
        <Button size="small" icon={<EyeOutlined />}
          onClick={() => handleViewDetails(r.workerId, r.worker?.realName)}>
          明细
        </Button>
      ),
    },
  ]

  const detailColumns = [
    { title: '日期', dataIndex: 'wageDate', width: 110 },
    { title: '类型', dataIndex: 'detailType', width: 90,
      render: (v) => v === 1 ? <Tag color="blue">直接计件</Tag> : v === 2 ? <Tag color="green">班组分配</Tag> : <Tag>调整</Tag> },
    { title: '合格数', dataIndex: 'qtyGood', width: 80, align: 'right' },
    { title: '不良数', dataIndex: 'qtyDefect', width: 80, align: 'right' },
    { title: '单价', dataIndex: 'unitPrice', width: 90, align: 'right',
      render: (v) => v > 0 ? `¥${v}` : '-' },
    { title: '金额', dataIndex: 'amount', width: 110, align: 'right',
      render: (v) => <span style={{ fontWeight: 'bold', color: v >= 0 ? '#cf1322' : '#999' }}>
        {v >= 0 ? '+' : ''}¥{v}
      </span> },
    { title: '备注', dataIndex: 'remark', ellipsis: true },
  ]

  return (
    <div>
      <Title level={4}>工资查询</Title>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={8}>
          <Card>
            <Statistic title="当月实时累计计件额"
              value={realtimeAmt ?? '-'} prefix="¥"
              valueStyle={{ color: '#cf1322', fontSize: 28 }}
              suffix={realtimeAmt !== null ? '' : ' (查询中)'} />
            <Button size="small" icon={<SyncOutlined />} onClick={fetchRealtime}
              style={{ marginTop: 8 }}>刷新</Button>
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title="查询月份" value={selectedMonth} />
            <InputNumber stringMode style={{ width: 160, marginTop: 8 }}
              placeholder="YYYY-MM" onChange={(v) => setSelectedMonth(v)} />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title="工资汇总数" value={total} prefix={<DollarOutlined />} />
            <Button type="primary" style={{ marginTop: 8 }}
              onClick={handleSettle}>月度结算</Button>
          </Card>
        </Col>
      </Row>

      <Card>
        <Table rowKey="id" dataSource={summaries} columns={columns}
          loading={loading} size="middle" scroll={{ x: 1100 }}
          pagination={{
            current: page, total, pageSize: 20,
            onChange: (p) => { setPage(p); fetchSummaries(p) },
          }} />
      </Card>

      <Modal title={selectedWorker ? `${selectedWorker.name} - 工资明细` : '工资明细'}
        open={detailModalOpen} onCancel={() => setDetailModalOpen(false)}
        footer={null} width={800}>
        <Table rowKey="id" dataSource={detailList} columns={detailColumns}
          size="small" pagination={false} />
      </Modal>
    </div>
  )
}
