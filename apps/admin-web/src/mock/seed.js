import { LOOKUP_SEED } from './lookupSeed'
import { buildProductSeed } from './productSeed'
import { LOOKUP_TABLES } from '@/config/lookupTables'
import { getSchema, emptySpec } from '@/config/componentSchemas'

const DAY = 24 * 60 * 60 * 1000

// Deterministic PRNG so the seed is stable between resets
function mulberry32(seed) {
  return () => {
    seed |= 0
    seed = (seed + 0x6d2b79f5) | 0
    let t = Math.imul(seed ^ (seed >>> 15), 1 | seed)
    t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
}

const CUSTOMERS = [
  ['Nguyễn Văn An', 'an.nguyen@gmail.com', '0901234567', '12 Lê Lợi, Quận 1, TP.HCM'],
  ['Trần Thị Bình', 'binh.tran@gmail.com', '0912345678', '45 Nguyễn Huệ, Quận 1, TP.HCM'],
  ['Lê Hoàng Cường', 'cuong.le@outlook.com', '0923456789', '88 Cầu Giấy, Hà Nội'],
  ['Phạm Minh Đức', 'duc.pham@gmail.com', '0934567890', '23 Trần Phú, Đà Nẵng'],
  ['Hoàng Thu Hà', 'ha.hoang@yahoo.com', '0945678901', '7 Phan Chu Trinh, Huế'],
  ['Vũ Quốc Huy', 'huy.vu@gmail.com', '0956789012', '102 Võ Văn Tần, Quận 3, TP.HCM'],
  ['Đặng Ngọc Khánh', 'khanh.dang@gmail.com', '0967890123', '56 Láng Hạ, Hà Nội'],
  ['Bùi Thanh Long', 'long.bui@gmail.com', '0978901234', '9 Nguyễn Văn Linh, Cần Thơ'],
  ['Đỗ Mai Linh', 'linh.do@gmail.com', '0989012345', '31 Hai Bà Trưng, Hải Phòng'],
  ['Ngô Gia Minh', 'minh.ngo@gmail.com', '0990123456', '64 Điện Biên Phủ, Bình Thạnh, TP.HCM'],
  ['Dương Hải Nam', 'nam.duong@gmail.com', '0381234567', '17 Kim Mã, Hà Nội'],
  ['Lý Phương Oanh', 'oanh.ly@gmail.com', '0392345678', '5 Lý Thường Kiệt, Nha Trang'],
  ['Trịnh Công Phúc', 'phuc.trinh@gmail.com', '0703456789', '210 Cách Mạng Tháng 8, Quận 10, TP.HCM'],
  ['Mai Anh Quân', 'quan.mai@gmail.com', '0764567890', '41 Bạch Đằng, Đà Nẵng'],
  ['Phan Thị Thảo', 'thao.phan@gmail.com', '0775678901', '78 Nguyễn Trãi, Thanh Xuân, Hà Nội'],
]

const TIER_BY_POINTS = (p) => (p >= 5000 ? 'PLATINUM' : p >= 2000 ? 'GOLD' : p >= 500 ? 'SILVER' : 'MEMBER')

export function createSeed() {
  const rand = mulberry32(20261028)
  const pick = (arr) => arr[Math.floor(rand() * arr.length)]
  const now = Date.now()
  const iso = (t) => new Date(t).toISOString()
  const created = iso(now - 120 * DAY)

  // ---- lookups ----
  const lookups = {}
  for (const { table } of LOOKUP_TABLES) {
    lookups[table] = (LOOKUP_SEED[table] ?? []).map(([code, name, description], i) => ({
      id: i + 1,
      code,
      name,
      description: description ?? null,
      is_active: true,
      created_at: created,
      updated_at: created,
    }))
  }
  const L = (table, code) => {
    const row = lookups[table].find((r) => r.code === code)
    if (!row) throw new Error(`Seed lookup not found: ${table}.${code}`)
    return row.id
  }

  // ---- models + inventory ----
  const models = []
  const inventory = []
  const stockMovements = []
  buildProductSeed(L).forEach((p, i) => {
    const id = i + 1
    const schema = getSchema(p.type)
    models.push({
      id,
      component_type_id: L('component_types', p.type),
      manufacturer_id: L('manufacturers', p.manufacturer),
      model_name: p.model_name,
      description: p.description,
      release_date: p.release_date,
      is_active: true,
      created_at: created,
      updated_at: created,
      spec: { ...emptySpec(schema), ...p.spec },
    })
    inventory.push({
      model_id: id,
      sku: p.sku,
      selling_price: p.price,
      unit_cost: p.cost,
      quantity_on_hand: p.stock,
      reorder_level: 5,
      location: pick(['WH-HCM-A1', 'WH-HCM-B2', 'WH-HN-A1']),
      updated_at: created,
    })
    stockMovements.push({
      id: stockMovements.length + 1,
      model_id: id,
      type: 'IN',
      quantity: p.stock + 10,
      reason: 'Initial stock import',
      reference: 'PO-INIT',
      created_at: created,
    })
  })

  // ---- customers ----
  const customers = CUSTOMERS.map(([full_name, email, phone, address], i) => ({
    id: i + 1,
    code: `CUS${String(i + 1).padStart(5, '0')}`,
    full_name,
    email,
    phone,
    address,
    loyalty_points: 0,
    tier: 'MEMBER',
    is_active: i !== 11,
    created_at: iso(now - (90 - i * 3) * DAY),
    updated_at: iso(now - (90 - i * 3) * DAY),
  }))

  // ---- orders ----
  const orders = []
  const statusFlow = ['PENDING', 'CONFIRMED', 'PROCESSING', 'SHIPPED', 'DELIVERED']
  for (let i = 0; i < 48; i++) {
    const daysAgo = Math.floor(rand() * 60)
    const createdAt = now - daysAgo * DAY - Math.floor(rand() * DAY)
    const customer = pick(customers)
    const itemCount = 1 + Math.floor(rand() * 3)
    const items = []
    for (let k = 0; k < itemCount; k++) {
      const m = pick(models)
      if (items.some((it) => it.model_id === m.id)) continue
      const inv = inventory.find((x) => x.model_id === m.id)
      items.push({
        model_id: m.id,
        sku: inv.sku,
        product_name: m.model_name,
        quantity: 1 + Math.floor(rand() * 2),
        unit_price: inv.selling_price,
      })
    }
    const subtotal = items.reduce((s, it) => s + it.quantity * it.unit_price, 0)
    const shipping_fee = subtotal > 5_000_000 ? 0 : 30_000

    // Older orders are more likely to be finished
    let status
    if (rand() < 0.08) status = 'CANCELLED'
    else if (daysAgo > 10) status = 'DELIVERED'
    else status = statusFlow[Math.min(4, Math.floor(rand() * 5))]

    const history = [{ status: 'PENDING', note: 'Order placed', at: iso(createdAt) }]
    if (status !== 'CANCELLED') {
      const idx = statusFlow.indexOf(status)
      for (let s = 1; s <= idx; s++) {
        history.push({ status: statusFlow[s], note: null, at: iso(createdAt + s * 0.2 * DAY) })
      }
    } else {
      history.push({ status: 'CANCELLED', note: 'Customer requested cancellation', at: iso(createdAt + 0.1 * DAY) })
    }

    const payment_method = pick(['COD', 'BANK_TRANSFER', 'CARD', 'E_WALLET'])
    let payment_status = 'UNPAID'
    if (status === 'DELIVERED' || (payment_method !== 'COD' && status !== 'CANCELLED')) payment_status = 'PAID'

    orders.push({
      id: i + 1,
      order_number: '',
      customer_id: customer.id,
      items,
      subtotal,
      discount: 0,
      shipping_fee,
      total: subtotal + shipping_fee,
      status,
      payment_status,
      payment_method,
      shipping_address: customer.address,
      note: null,
      status_history: history,
      created_at: iso(createdAt),
      updated_at: history.at(-1).at,
    })
  }
  orders.sort((a, b) => a.created_at.localeCompare(b.created_at))
  orders.forEach((o, i) => {
    o.id = i + 1
    o.order_number = `ORD-${o.created_at.slice(0, 10).replaceAll('-', '')}-${String(i + 1).padStart(4, '0')}`
  })

  for (const c of customers) {
    const spent = orders
      .filter((o) => o.customer_id === c.id && o.status === 'DELIVERED')
      .reduce((s, o) => s + o.total, 0)
    c.loyalty_points = Math.floor(spent / 10_000)
    c.tier = TIER_BY_POINTS(c.loyalty_points)
  }

  // ---- warranty claims ----
  const warrantyStatuses = ['RECEIVED', 'DIAGNOSING', 'REPAIRING', 'REPLACED', 'RETURNED', 'REJECTED']
  const issues = [
    'No display output after 2 weeks of use.',
    'Coil whine and random shutdowns under load.',
    'Fan makes grinding noise.',
    'Drive not detected in BIOS.',
    'Memory errors in MemTest86.',
    'System fails to POST.',
  ]
  const delivered = orders.filter((o) => o.status === 'DELIVERED')
  const warranties = []
  for (let i = 0; i < 9 && delivered.length; i++) {
    const order = delivered[Math.floor(rand() * delivered.length)]
    const item = pick(order.items)
    const status = warrantyStatuses[i % warrantyStatuses.length]
    const createdAt = Math.min(now, Math.max(Date.parse(order.updated_at) + DAY, now - Math.floor(rand() * 10) * DAY))
    warranties.push({
      id: i + 1,
      claim_number: `WAR-${String(i + 1).padStart(5, '0')}`,
      order_id: order.id,
      customer_id: order.customer_id,
      model_id: item.model_id,
      serial_number: `SN${Math.floor(rand() * 1e10).toString().padStart(10, '0')}`,
      purchase_date: order.created_at.slice(0, 10),
      warranty_expiry: new Date(Date.parse(order.created_at) + 3 * 365 * DAY).toISOString().slice(0, 10),
      issue_description: pick(issues),
      status,
      resolution: ['REPLACED', 'RETURNED', 'REJECTED'].includes(status)
        ? status === 'REJECTED'
          ? 'Physical damage found — not covered by warranty.'
          : 'Issue resolved and unit returned to customer.'
        : null,
      created_at: iso(createdAt),
      updated_at: iso(createdAt),
    })
  }

  return {
    version: 1,
    lookups,
    models,
    inventory,
    stockMovements,
    customers,
    orders,
    warranties,
  }
}
