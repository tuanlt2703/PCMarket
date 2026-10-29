export const ORDER_STATUSES = {
  PENDING: { label: 'Pending', color: 'secondary' },
  CONFIRMED: { label: 'Confirmed', color: 'info' },
  PROCESSING: { label: 'Processing', color: 'primary' },
  SHIPPED: { label: 'Shipped', color: 'warning' },
  DELIVERED: { label: 'Delivered', color: 'success' },
  CANCELLED: { label: 'Cancelled', color: 'danger' },
}

export const ORDER_TRANSITIONS = {
  PENDING: ['CONFIRMED', 'CANCELLED'],
  CONFIRMED: ['PROCESSING', 'CANCELLED'],
  PROCESSING: ['SHIPPED', 'CANCELLED'],
  SHIPPED: ['DELIVERED'],
  DELIVERED: [],
  CANCELLED: [],
}

// Orders in these statuses hold (reserve) stock
export const RESERVING_ORDER_STATUSES = ['PENDING', 'CONFIRMED', 'PROCESSING']

export const PAYMENT_STATUSES = {
  UNPAID: { label: 'Unpaid', color: 'secondary' },
  PAID: { label: 'Paid', color: 'success' },
  REFUNDED: { label: 'Refunded', color: 'danger' },
}

export const PAYMENT_METHODS = {
  COD: 'Cash on delivery',
  BANK_TRANSFER: 'Bank transfer',
  CARD: 'Credit / debit card',
  E_WALLET: 'E-wallet',
}

export const WARRANTY_STATUSES = {
  RECEIVED: { label: 'Received', color: 'secondary' },
  DIAGNOSING: { label: 'Diagnosing', color: 'info' },
  REPAIRING: { label: 'Repairing', color: 'primary' },
  REPLACED: { label: 'Replaced', color: 'success' },
  RETURNED: { label: 'Returned to customer', color: 'success' },
  REJECTED: { label: 'Rejected', color: 'danger' },
}

export const CUSTOMER_TIERS = {
  MEMBER: { label: 'Member', color: 'secondary' },
  SILVER: { label: 'Silver', color: 'info' },
  GOLD: { label: 'Gold', color: 'warning' },
  PLATINUM: { label: 'Platinum', color: 'primary' },
}

export const STOCK_MOVEMENT_TYPES = {
  IN: { label: 'Stock in', color: 'success' },
  OUT: { label: 'Stock out', color: 'danger' },
  ADJUST: { label: 'Adjustment', color: 'warning' },
}
