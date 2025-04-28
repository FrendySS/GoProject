export interface Product {
  id: string
  name: string
  description: string
  price: number
  stock: number
  status: string
  createdAt: string
  updatedAt: string
  deletedAt?: string
}

export interface User {
  id: string
  username: string
  email: string
  role: string
  status: string
  createdAt: string
  updatedAt: string
  deletedAt?: string
}

export interface DashboardStats {
  productsCount: number
  totalValue: number
  usersCount: number
}
