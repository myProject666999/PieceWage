import request from '../utils/request'

export const login = (data) => request.post('/auth/login', data)

export const getProfile = () => request.get('/auth/profile')

export const listProducts = (params) => request.get('/products', { params })

export const createProduct = (data) => request.post('/products', data)

export const listProcessSteps = (params) => request.get('/process-steps', { params })

export const createProcessStep = (data) => request.post('/process-steps', data)

export const listStepsByProduct = (productId) => request.get(`/products/${productId}/steps`)

export const listProcessPrices = (params) => request.get('/process-prices', { params })

export const createProcessPrice = (data) => request.post('/process-prices', data)

export const getEffectivePrice = (params) => request.get('/process-prices/effective', { params })

export const listPricesByProcess = (processId) => request.get(`/process-prices/process/${processId}`)

export const createReport = (data) => request.post('/reports', data)

export const getReport = (id) => request.get(`/reports/${id}`)

export const listReports = (params) => request.get('/reports', { params })

export const voidReport = (id) => request.put(`/reports/${id}/void`)

export const getMonthlySummary = (workerId, month) =>
  request.get(`/wage/summary/${workerId}/${month}`)

export const listWageSummaries = (params) => request.get('/wage/summaries', { params })

export const getWorkerDetails = (params) => request.get('/wage/details', { params })

export const getWorkerDailyDetails = (workerId, date) =>
  request.get(`/wage/daily/${workerId}/${date}`)

export const getRealtimeAccumulate = (workerId, month) =>
  request.get(`/wage/realtime/${workerId}/${month}`)

export const settleMonth = (month) => request.post(`/wage/settle/${month}`)

export const listTeams = () => request.get('/teams')

export const createTeam = (data) => request.post('/teams', data)

export const getTeamMembers = (teamId) => request.get(`/teams/${teamId}/members`)

export const getAllocation = (reportId) => request.get(`/allocations/${reportId}`)

export const listUsers = (params) => request.get('/users', { params })
