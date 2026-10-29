import { getDb, clone, delay } from './mockDb'
import { reservedQuantity } from './productApi'

const DAY = 24 * 60 * 60 * 1000

export const dashboardApi = {
  async summary({ days = 30 } = {}) {
    await delay()
    const db = getDb()
    const since = Date.now() - days * DAY
    const recent = db.orders.filter((o) => Date.parse(o.created_at) >= since)
    const valid = recent.filter((o) => o.status !== 'CANCELLED')

    // Daily revenue / order count series for the chart
    const series = []
    for (let i = days - 1; i >= 0; i--) {
      const date = new Date(Date.now() - i * DAY).toISOString().slice(0, 10)
      const dayOrders = valid.filter((o) => o.created_at.slice(0, 10) === date)
      series.push({ date, revenue: dayOrders.reduce((s, o) => s + o.total, 0), orders: dayOrders.length })
    }

    const typeName = (id) => db.lookups.component_types.find((t) => t.id === id)?.name ?? 'Other'
    const salesByType = {}
    for (const o of valid) {
      for (const it of o.items) {
        const m = db.models.find((x) => x.id === it.model_id)
        const key = m ? typeName(m.component_type_id) : 'Other'
        salesByType[key] = (salesByType[key] ?? 0) + it.quantity * it.unit_price
      }
    }

    const lowStock = db.inventory
      .map((inv) => {
        const available = inv.quantity_on_hand - reservedQuantity(db, inv.model_id)
        const model = db.models.find((m) => m.id === inv.model_id)
        return { ...inv, model_name: model?.model_name, available, is_active: model?.is_active }
      })
      .filter((x) => x.is_active && x.available <= x.reorder_level)
      .sort((a, b) => a.available - b.available)

    const warrantyOpen = db.warranties.filter((w) => ['RECEIVED', 'DIAGNOSING', 'REPAIRING'].includes(w.status))

    return clone({
      newOrders: recent.length,
      pendingOrders: db.orders.filter((o) => o.status === 'PENDING').length,
      revenue: valid.reduce((s, o) => s + o.total, 0),
      customers: db.customers.length,
      newCustomers: db.customers.filter((c) => Date.parse(c.created_at) >= since).length,
      productCount: db.models.length,
      lowStockCount: lowStock.length,
      openWarranties: warrantyOpen.length,
      series,
      salesByType,
      lowStock: lowStock.slice(0, 6),
      recentOrders: [...db.orders]
        .sort((a, b) => b.created_at.localeCompare(a.created_at))
        .slice(0, 7)
        .map((o) => ({ ...o, customer_name: db.customers.find((c) => c.id === o.customer_id)?.full_name })),
      warrantyByStatus: db.warranties.reduce((acc, w) => ((acc[w.status] = (acc[w.status] ?? 0) + 1), acc), {}),
    })
  },
}
