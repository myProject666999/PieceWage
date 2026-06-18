import { useState, useEffect } from 'react'
import { Table, Button, Modal, Form, InputNumber, Select, Input, Tag, Space, message, Card, Typography, Row, Col, Descriptions } from 'antd'
import { PlusOutlined, EyeOutlined, StopOutlined } from '@ant-design/icons'
import { listReports, createReport, getReport, voidReport, listProcessSteps, listTeams, listStepsByProduct, getEffectivePrice } from '../api'
import dayjs from 'dayjs'

const { Title } = Typography

const gradeMap = { STD: '标准', PRE: '特级', ADV: '高级' }

export default function ReportManagement() {
  const [reports, setReports] = useState([])
  const [steps, setSteps] = useState([])
  const [teams, setTeams] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [detailOpen, setDetailOpen] = useState(false)
  const [detail, setDetail] = useState(null)
  const [previewPrice, setPreviewPrice] = useState(null)
  const [form] = Form.useForm()

  const fetchReports = async (p = page) => {
    setLoading(true)
    try {
      const res = await listReports({ page: p, pageSize: 20 })
      setReports(res.data?.list || [])
      setTotal(res.data?.total || 0)
    } catch (err) {
      message.error('获取报工列表失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchOptions = async () => {
    try {
      const [stepsRes, teamsRes] = await Promise.all([
        listProcessSteps({ page: 1, pageSize: 200 }),
        listTeams(),
      ])
      setSteps(stepsRes.data?.list || [])
      setTeams(teamsRes.data || [])
    } catch {}
  }

  useEffect(() => { fetchReports(); fetchOptions() }, [])

  const handlePreviewPrice = async () => {
    const processId = form.getFieldValue('processId')
    const gradeLevel = form.getFieldValue('gradeLevel') || 'STD'
    const reportDate = form.getFieldValue('reportDate')
    if (!processId || !reportDate) {
      message.warning('请先选择工序和报工日期')
      return
    }
    try {
      const res = await getEffectivePrice({ processId, gradeLevel, date: reportDate })
      setPreviewPrice(res.data)
      message.success(`报工日期${reportDate}生效单价: ¥${res.data.unitPrice} (版本${res.data.versionNo})`)
    } catch {
      message.error('该日期暂无生效单价，无法报工')
      setPreviewPrice(null)
    }
  }

  const handleCreate = async (values) => {
    try {
      await createReport(values)
      message.success('报工成功！系统已按报工当时生效单价计算计件收入')
      setModalOpen(false)
      form.resetFields()
      setPreviewPrice(null)
      fetchReports()
    } catch (err) {
      message.error(err.message || '报工失败')
    }
  }

  const handleViewDetail = async (id) => {
    try {
      const res = await getReport(id)
      setDetail(res.data)
      setDetailOpen(true)
    } catch {
      message.error('获取详情失败')
    }
  }

  const handleVoid = async (id) => {
    Modal.confirm({
      title: '确认作废',
      content: '作废后将冲销该笔计件收入，确定要作废吗？',
      okType: 'danger',
      onOk: async () => {
        try {
          await voidReport(id)
          message.success('已作废')
          fetchReports()
        } catch (err) {
          message.error(err.message || '作废失败')
        }
      },
    })
  }

  const columns = [
    { title: '报工单号', dataIndex: 'reportNo', width: 160 },
    {
      title: '工人', dataIndex: 'worker', width: 100,
      render: (w) => w?.realName || '-',
    },
    {
      title: '工序', dataIndex: 'process', width: 150,
      render: (p) => p ? `${p.processName}` : '-',
    },
    {
      title: '单价快照', dataIndex: 'unitPrice', width: 100, align: 'right',
      render: (v) => <span style={{ color: '#1677ff' }}>¥{v}</span>,
    },
    {
      title: '等级', dataIndex: 'gradeLevel', width: 70, align: 'center',
      render: (v) => <Tag>{gradeMap[v] || v}</Tag>,
    },
    { title: '报工日期', dataIndex: 'reportDate', width: 110 },
    { title: '合格数', dataIndex: 'qtyGood', width: 80, align: 'right' },
    { title: '不良数', dataIndex: 'qtyDefect', width: 80, align: 'right' },
    { title: '总额', dataIndex: 'grossAmount', width: 100, align: 'right',
      render: (v) => `¥${v}` },
    { title: '扣款', dataIndex: 'defectAmount', width: 90, align: 'right',
      render: (v) => v > 0 ? <span style={{ color: '#cf1322' }}>-¥{v}</span> : '¥0' },
    {
      title: '实发', dataIndex: 'netAmount', width: 110, align: 'right',
      render: (v) => <span style={{ fontWeight: 'bold', color: '#cf1322', fontSize: 14 }}>¥{v}</span>,
    },
    {
      title: '状态', dataIndex: 'status', width: 80, align: 'center',
      render: (v) => v === 1 ? <Tag color="green">生效</Tag> : <Tag color="red">已作废</Tag>,
    },
    {
      title: '操作', width: 120, fixed: 'right',
      render: (_, r) => (
        <Space>
          <Button size="small" icon={<EyeOutlined />} onClick={() => handleViewDetail(r.id)} />
          {r.status === 1 && (
            <Button size="small" danger icon={<StopOutlined />} onClick={() => handleVoid(r.id)} />
          )}
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Title level={4}>报工管理</Title>
      <Card>
        <Space style={{ marginBottom: 16 }}>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
            新建报工
          </Button>
        </Space>
        <Table rowKey="id" dataSource={reports} columns={columns}
          loading={loading} scroll={{ x: 1500 }}
          pagination={{
            current: page, total, pageSize: 20,
            onChange: (p) => { setPage(p); fetchReports(p) },
          }} size="middle" />
      </Card>

      <Modal title="新建报工" open={modalOpen} onCancel={() => { setModalOpen(false); setPreviewPrice(null) }}
        onOk={() => form.submit()} width={600}>
        <Form form={form} layout="vertical" onFinish={handleCreate}
          initialValues={{ gradeLevel: 'STD', qtyGood: 0, qtyDefect: 0, unitDefect: 0.5 }}>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="processId" label="工序" rules={[{ required: true }]}>
                <Select placeholder="选择工序" showSearch optionFilterProp="label"
                  options={steps.map(s => ({
                    value: s.id,
                    label: `${s.processName} ${s.isShared === 1 ? '[共享]' : ''}`,
                  }))} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="gradeLevel" label="等级" rules={[{ required: true }]}>
                <Select options={[
                  { value: 'STD', label: '标准(STD)' },
                  { value: 'PRE', label: '特级(PRE)' },
                  { value: 'ADV', label: '高级(ADV)' },
                ]} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="reportDate" label="报工日期" rules={[{ required: true }]}>
                <Input placeholder="YYYY-MM-DD" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item>
                <Button onClick={handlePreviewPrice} style={{ marginTop: 30 }}>
                  查询生效单价
                </Button>
                {previewPrice && (
                  <Tag color="green" style={{ marginLeft: 8 }}>
                    ¥{previewPrice.unitPrice}/件 (V{previewPrice.versionNo})
                  </Tag>
                )}
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={8}>
              <Form.Item name="qtyGood" label="合格数量" rules={[{ required: true }]}>
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="qtyDefect" label="不良/返工数量">
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="unitDefect" label="不良扣款比例">
                <InputNumber min={0} max={1} step={0.1} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="workOrderNo" label="关联工单号">
            <Input placeholder="选填" />
          </Form.Item>
          <Form.Item name="remark" label="备注">
            <Input.TextArea rows={2} />
          </Form.Item>
          <div style={{ padding: 12, background: '#f6f8fa', borderRadius: 6, fontSize: 13, color: '#666' }}>
            <strong>计件规则说明：</strong>系统将自动查询报工日期当时生效的单价（非最新单价），
            按合格数量×单价计算总额，不良品按数量×单价×扣款比例扣除，得出实发金额。
          </div>
        </Form>
      </Modal>

      <Modal title="报工单详情" open={detailOpen} onCancel={() => setDetailOpen(false)} footer={null} width={600}>
        {detail && (
          <Descriptions column={2} bordered size="small">
            <Descriptions.Item label="报工单号">{detail.reportNo}</Descriptions.Item>
            <Descriptions.Item label="工人">{detail.worker?.realName}</Descriptions.Item>
            <Descriptions.Item label="工序">{detail.process?.processName}</Descriptions.Item>
            <Descriptions.Item label="产品">{detail.process?.product?.productName}</Descriptions.Item>
            <Descriptions.Item label="单价快照">¥{detail.unitPrice}</Descriptions.Item>
            <Descriptions.Item label="单价版本">
              {detail.price ? `V${detail.price.versionNo}` : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="等级">{gradeMap[detail.gradeLevel] || detail.gradeLevel}</Descriptions.Item>
            <Descriptions.Item label="报工日期">{detail.reportDate}</Descriptions.Item>
            <Descriptions.Item label="合格数量">{detail.qtyGood}</Descriptions.Item>
            <Descriptions.Item label="不良数量">{detail.qtyDefect}</Descriptions.Item>
            <Descriptions.Item label="计件总额">¥{detail.grossAmount}</Descriptions.Item>
            <Descriptions.Item label="不良扣款">¥{detail.defectAmount}</Descriptions.Item>
            <Descriptions.Item label="实发金额" span={2}>
              <span style={{ fontWeight: 'bold', color: '#cf1322', fontSize: 18 }}>¥{detail.netAmount}</span>
            </Descriptions.Item>
            <Descriptions.Item label="状态">
              {detail.status === 1 ? <Tag color="green">生效</Tag> : <Tag color="red">已作废</Tag>}
            </Descriptions.Item>
            <Descriptions.Item label="关联工单">{detail.workOrderNo || '-'}</Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </div>
  )
}
