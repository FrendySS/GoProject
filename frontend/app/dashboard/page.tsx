"use client"

import type React from "react"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Alert, AlertDescription } from "@/components/ui/alert"
import {
  LayoutDashboard,
  Package,
  Users,
  User,
  LogOut,
  Plus,
  Pencil,
  Trash,
  Ban,
  CheckCircle,
  RotateCcw,
} from "lucide-react"
import { API_URL } from "@/lib/constants"
import { fetchWithAuth } from "@/lib/api-helpers"
import type { Product, User as UserType, DashboardStats } from "@/lib/types"

export default function DashboardPage() {
  const router = useRouter()
  const [activeTab, setActiveTab] = useState("dashboard")
  const [userRole, setUserRole] = useState("")
  const [stats, setStats] = useState<DashboardStats>({
    productsCount: 0,
    totalValue: 0,
    usersCount: 0,
  })

  // Products state
  const [products, setProducts] = useState<Product[]>([])
  const [productDialogOpen, setProductDialogOpen] = useState(false)
  const [currentProduct, setCurrentProduct] = useState<Partial<Product>>({})
  const [productError, setProductError] = useState("")

  // Users state
  const [users, setUsers] = useState<UserType[]>([])
  const [userRoleDialogOpen, setUserRoleDialogOpen] = useState(false)
  const [selectedUserId, setSelectedUserId] = useState("")
  const [selectedUserRole, setSelectedUserRole] = useState("")

  // Profile state
  const [profile, setProfile] = useState<UserType | null>(null)
  const [oldPassword, setOldPassword] = useState("")
  const [newPassword, setNewPassword] = useState("")
  const [passwordError, setPasswordError] = useState("")
  const [passwordSuccess, setPasswordSuccess] = useState(false)

  // Check authentication on component mount
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const response = await fetchWithAuth(`${API_URL}/profile`)

        if (!response.ok) {
          throw new Error("Authentication failed")
        }

        const userData = await response.json()
        setUserRole(userData.role)

        // Load initial data
        loadDashboardData()
      } catch (error) {
        // Redirect to login page if authentication fails
        localStorage.removeItem("accessToken")
        localStorage.removeItem("refreshToken")
        router.push("/login")
      }
    }

    checkAuth()
  }, [router])

  // Load dashboard data
  const loadDashboardData = async () => {
    try {
      // Load products for statistics
      const productsResponse = await fetchWithAuth(`${API_URL}/products`)

      if (!productsResponse.ok) {
        throw new Error("Failed to load products")
      }

      const productsData = await productsResponse.json()

      // Calculate statistics
      const productsCount = productsData.length
      const totalValue = productsData.reduce((sum: number, product: Product) => sum + product.price * product.stock, 0)

      setStats({
        ...stats,
        productsCount,
        totalValue,
      })

      // If user is director, load users count
      if (userRole === "director") {
        const usersResponse = await fetchWithAuth(`${API_URL}/admin/users`)

        if (usersResponse.ok) {
          const usersData = await usersResponse.json()
          setStats({
            ...stats,
            productsCount,
            totalValue,
            usersCount: usersData.length,
          })
        }
      }
    } catch (error) {
      console.error("Error loading dashboard data:", error)
    }
  }

  // Load products
  const loadProducts = async () => {
    try {
      const response = await fetchWithAuth(`${API_URL}/products`)

      if (!response.ok) {
        throw new Error("Failed to load products")
      }

      const data = await response.json()
      setProducts(data)
    } catch (error) {
      console.error("Error loading products:", error)
    }
  }

  // Load users
  const loadUsers = async () => {
    if (userRole !== "director") return

    try {
      const response = await fetchWithAuth(`${API_URL}/admin/users`)

      if (!response.ok) {
        throw new Error("Failed to load users")
      }

      const data = await response.json()
      setUsers(data)
    } catch (error) {
      console.error("Error loading users:", error)
    }
  }

  // Load profile
  const loadProfile = async () => {
    try {
      const response = await fetchWithAuth(`${API_URL}/profile`)

      if (!response.ok) {
        throw new Error("Failed to load profile")
      }

      const data = await response.json()
      setProfile(data)
    } catch (error) {
      console.error("Error loading profile:", error)
    }
  }

  // Handle tab change
  const handleTabChange = (value: string) => {
    setActiveTab(value)

    // Load data based on selected tab
    if (value === "dashboard") {
      loadDashboardData()
    } else if (value === "products") {
      loadProducts()
    } else if (value === "users") {
      loadUsers()
    } else if (value === "profile") {
      loadProfile()
    }
  }

  // Handle logout
  const handleLogout = () => {
    localStorage.removeItem("accessToken")
    localStorage.removeItem("refreshToken")
    localStorage.removeItem("tokenType")
    router.push("/login")
  }

  // Handle product form submission
  const handleProductSubmit = async () => {
    try {
      setProductError("")

      if (!currentProduct.name || !currentProduct.price || currentProduct.stock === undefined) {
        setProductError("Please fill all required fields")
        return
      }

      let response

      if (currentProduct.id) {
        // Update existing product
        response = await fetchWithAuth(`${API_URL}/products/${currentProduct.id}`, {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(currentProduct),
        })
      } else {
        // Create new product
        response = await fetchWithAuth(`${API_URL}/products`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(currentProduct),
        })
      }

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Failed to save product")
      }

      // Close dialog and reload products
      setProductDialogOpen(false)
      loadProducts()

      // Also refresh dashboard data
      if (activeTab === "dashboard") {
        loadDashboardData()
      }
    } catch (error) {
      setProductError(error instanceof Error ? error.message : "An unknown error occurred")
    }
  }

  // Handle product edit
  const handleEditProduct = async (productId: string) => {
    try {
      const response = await fetchWithAuth(`${API_URL}/products/${productId}`)

      if (!response.ok) {
        throw new Error("Failed to load product details")
      }

      const product = await response.json()
      setCurrentProduct(product)
      setProductDialogOpen(true)
    } catch (error) {
      console.error("Error editing product:", error)
    }
  }

  // Handle product delete
  const handleDeleteProduct = async (productId: string) => {
    if (!confirm("Are you sure you want to delete this product?")) {
      return
    }

    try {
      const response = await fetchWithAuth(`${API_URL}/products/${productId}`, {
        method: "DELETE",
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Failed to delete product")
      }

      // Reload products
      loadProducts()

      // Also refresh dashboard data
      if (activeTab === "dashboard") {
        loadDashboardData()
      }
    } catch (error) {
      console.error("Error deleting product:", error)
      alert(error instanceof Error ? error.message : "An unknown error occurred")
    }
  }

  // Handle user role change
  const handleUserRoleChange = async () => {
    try {
      const response = await fetchWithAuth(`${API_URL}/admin/assign-role`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ userId: selectedUserId, role: selectedUserRole }),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Failed to change user role")
      }

      // Close dialog and reload users
      setUserRoleDialogOpen(false)
      loadUsers()
    } catch (error) {
      console.error("Error changing user role:", error)
      alert(error instanceof Error ? error.message : "An unknown error occurred")
    }
  }

  // Handle user ban
  const handleBanUser = async (userId: string) => {
    if (!confirm("Are you sure you want to ban this user?")) {
      return
    }

    try {
      const response = await fetchWithAuth(`${API_URL}/admin/ban-user`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ userId }),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Failed to ban user")
      }

      // Reload users
      loadUsers()
    } catch (error) {
      console.error("Error banning user:", error)
      alert(error instanceof Error ? error.message : "An unknown error occurred")
    }
  }

  // Handle user unban
  const handleUnbanUser = async (userId: string) => {
    try {
      const response = await fetchWithAuth(`${API_URL}/admin/unban-user`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ userId }),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Failed to unban user")
      }

      // Reload users
      loadUsers()
    } catch (error) {
      console.error("Error unbanning user:", error)
      alert(error instanceof Error ? error.message : "An unknown error occurred")
    }
  }

  // Handle user delete
  const handleDeleteUser = async (userId: string) => {
    if (!confirm("Are you sure you want to delete this user?")) {
      return
    }

    try {
      const response = await fetchWithAuth(`${API_URL}/admin/delete-user`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ userId }),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Failed to delete user")
      }

      // Reload users
      loadUsers()
    } catch (error) {
      console.error("Error deleting user:", error)
      alert(error instanceof Error ? error.message : "An unknown error occurred")
    }
  }

  // Handle user restore
  const handleRestoreUser = async (userId: string) => {
    try {
      const response = await fetchWithAuth(`${API_URL}/admin/restore-user`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ userId }),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Failed to restore user")
      }

      // Reload users
      loadUsers()
    } catch (error) {
      console.error("Error restoring user:", error)
      alert(error instanceof Error ? error.message : "An unknown error occurred")
    }
  }

  // Handle password change
  const handlePasswordChange = async (e: React.FormEvent) => {
    e.preventDefault()
    setPasswordError("")
    setPasswordSuccess(false)

    try {
      const response = await fetchWithAuth(`${API_URL}/profile/password`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ oldPassword, newPassword }),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Failed to change password")
      }

      // Show success message and reset form
      setPasswordSuccess(true)
      setOldPassword("")
      setNewPassword("")
    } catch (error) {
      setPasswordError(error instanceof Error ? error.message : "An unknown error occurred")
    }
  }

  // Helper function to get role display name
  const getRoleDisplayName = (role: string) => {
    switch (role) {
      case "viewer":
        return "Viewer"
      case "manager":
        return "Manager"
      case "director":
        return "Director"
      default:
        return role
    }
  }

  // Helper function to get status badge variant
  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case "active":
        return "success"
      case "banned":
        return "destructive"
      case "deleted":
        return "secondary"
      default:
        return "default"
    }
  }

  return (
      <div className="flex min-h-screen bg-background">
        {/* Sidebar */}
        <div className="hidden md:flex flex-col w-64 bg-sidebar border-r">
          <div className="p-6">
            <h2 className="text-2xl font-bold">MarketEase</h2>
            <p className="text-sm text-muted-foreground">Warehouse Management</p>
          </div>
          <div className="flex-1 px-4 space-y-2">
            <Button
                variant={activeTab === "dashboard" ? "secondary" : "ghost"}
                className="w-full justify-start"
                onClick={() => handleTabChange("dashboard")}
            >
              <LayoutDashboard className="mr-2 h-4 w-4" />
              Dashboard
            </Button>
            <Button
                variant={activeTab === "products" ? "secondary" : "ghost"}
                className="w-full justify-start"
                onClick={() => handleTabChange("products")}
            >
              <Package className="mr-2 h-4 w-4" />
              Products
            </Button>
            {userRole === "director" && (
                <Button
                    variant={activeTab === "users" ? "secondary" : "ghost"}
                    className="w-full justify-start"
                    onClick={() => handleTabChange("users")}
                >
                  <Users className="mr-2 h-4 w-4" />
                  Users
                </Button>
            )}
            <Button
                variant={activeTab === "profile" ? "secondary" : "ghost"}
                className="w-full justify-start"
                onClick={() => handleTabChange("profile")}
            >
              <User className="mr-2 h-4 w-4" />
              Profile
            </Button>
          </div>
          <div className="p-4 border-t">
            <Button variant="destructive" className="w-full" onClick={handleLogout}>
              <LogOut className="mr-2 h-4 w-4" />
              Logout
            </Button>
          </div>
        </div>

        {/* Mobile navigation */}
        <div className="md:hidden fixed bottom-0 left-0 right-0 bg-background border-t z-10">
          <div className="flex justify-around p-2">
            <Button
                variant="ghost"
                size="sm"
                className={activeTab === "dashboard" ? "bg-accent" : ""}
                onClick={() => handleTabChange("dashboard")}
            >
              <LayoutDashboard className="h-5 w-5" />
            </Button>
            <Button
                variant="ghost"
                size="sm"
                className={activeTab === "products" ? "bg-accent" : ""}
                onClick={() => handleTabChange("products")}
            >
              <Package className="h-5 w-5" />
            </Button>
            {userRole === "director" && (
                <Button
                    variant="ghost"
                    size="sm"
                    className={activeTab === "users" ? "bg-accent" : ""}
                    onClick={() => handleTabChange("users")}
                >
                  <Users className="h-5 w-5" />
                </Button>
            )}
            <Button
                variant="ghost"
                size="sm"
                className={activeTab === "profile" ? "bg-accent" : ""}
                onClick={() => handleTabChange("profile")}
            >
              <User className="h-5 w-5" />
            </Button>
            <Button variant="ghost" size="sm" onClick={handleLogout}>
              <LogOut className="h-5 w-5 text-destructive" />
            </Button>
          </div>
        </div>

        {/* Main content */}
        <div className="flex-1 flex flex-col min-h-screen pb-16 md:pb-0">
          <header className="border-b p-4 flex justify-between items-center">
            <h1 className="text-xl font-bold">{activeTab.charAt(0).toUpperCase() + activeTab.slice(1)}</h1>
            <Badge variant="outline">{getRoleDisplayName(userRole)}</Badge>
          </header>

          <main className="flex-1 p-4 md:p-6 overflow-auto">
            {/* Dashboard Tab */}
            {activeTab === "dashboard" && (
                <div className="grid gap-4 md:grid-cols-3">
                  <Card>
                    <CardHeader className="pb-2">
                      <CardTitle className="text-sm font-medium">Total Products</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <div className="text-2xl font-bold">{stats.productsCount}</div>
                      <p className="text-xs text-muted-foreground">Products in inventory</p>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardHeader className="pb-2">
                      <CardTitle className="text-sm font-medium">Total Value</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <div className="text-2xl font-bold">{stats.totalValue.toFixed(2)} ₽</div>
                      <p className="text-xs text-muted-foreground">Total inventory value</p>
                    </CardContent>
                  </Card>
                  {userRole === "director" && (
                      <Card>
                        <CardHeader className="pb-2">
                          <CardTitle className="text-sm font-medium">Active Users</CardTitle>
                        </CardHeader>
                        <CardContent>
                          <div className="text-2xl font-bold">{stats.usersCount}</div>
                          <p className="text-xs text-muted-foreground">Registered users</p>
                        </CardContent>
                      </Card>
                  )}
                </div>
            )}

            {/* Products Tab */}
            {activeTab === "products" && (
                <div className="space-y-4">
                  <div className="flex justify-between items-center">
                    <h2 className="text-lg font-semibold">Product List</h2>
                    {(userRole === "manager" || userRole === "director") && (
                        <Dialog open={productDialogOpen} onOpenChange={setProductDialogOpen}>
                          <DialogTrigger asChild>
                            <Button onClick={() => setCurrentProduct({})} size="sm">
                              <Plus className="mr-2 h-4 w-4" />
                              Add Product
                            </Button>
                          </DialogTrigger>
                          <DialogContent>
                            <DialogHeader>
                              <DialogTitle>{currentProduct.id ? "Edit Product" : "Add New Product"}</DialogTitle>
                              <DialogDescription>Fill in the product details below</DialogDescription>
                            </DialogHeader>
                            {productError && (
                                <Alert variant="destructive">
                                  <AlertDescription>{productError}</AlertDescription>
                                </Alert>
                            )}
                            <div className="grid gap-4 py-4">
                              <div className="grid gap-2">
                                <Label htmlFor="name">Name</Label>
                                <Input
                                    id="name"
                                    value={currentProduct.name || ""}
                                    onChange={(e) => setCurrentProduct({ ...currentProduct, name: e.target.value })}
                                    required
                                />
                              </div>
                              <div className="grid gap-2">
                                <Label htmlFor="description">Description</Label>
                                <Textarea
                                    id="description"
                                    value={currentProduct.description || ""}
                                    onChange={(e) => setCurrentProduct({ ...currentProduct, description: e.target.value })}
                                    rows={3}
                                />
                              </div>
                              <div className="grid grid-cols-2 gap-4">
                                <div className="grid gap-2">
                                  <Label htmlFor="price">Price</Label>
                                  <Input
                                      id="price"
                                      type="number"
                                      min="0"
                                      step="0.01"
                                      value={currentProduct.price || ""}
                                      onChange={(e) =>
                                          setCurrentProduct({ ...currentProduct, price: Number.parseFloat(e.target.value) })
                                      }
                                      required
                                  />
                                </div>
                                <div className="grid gap-2">
                                  <Label htmlFor="stock">Stock</Label>
                                  <Input
                                      id="stock"
                                      type="number"
                                      min="0"
                                      step="1"
                                      value={currentProduct.stock || ""}
                                      onChange={(e) =>
                                          setCurrentProduct({ ...currentProduct, stock: Number.parseInt(e.target.value) })
                                      }
                                      required
                                  />
                                </div>
                              </div>
                            </div>
                            <DialogFooter>
                              <Button variant="outline" onClick={() => setProductDialogOpen(false)}>
                                Cancel
                              </Button>
                              <Button onClick={handleProductSubmit}>Save</Button>
                            </DialogFooter>
                          </DialogContent>
                        </Dialog>
                    )}
                  </div>

                  <div className="border rounded-md">
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>Name</TableHead>
                          <TableHead className="hidden md:table-cell">Description</TableHead>
                          <TableHead>Price</TableHead>
                          <TableHead>Stock</TableHead>
                          <TableHead>Status</TableHead>
                          {(userRole === "manager" || userRole === "director") && (
                              <TableHead className="text-right">Actions</TableHead>
                          )}
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {products.length === 0 ? (
                            <TableRow>
                              <TableCell colSpan={6} className="text-center py-4">
                                No products found
                              </TableCell>
                            </TableRow>
                        ) : (
                            products.map((product) => (
                                <TableRow key={product.id}>
                                  <TableCell className="font-medium">{product.name}</TableCell>
                                  <TableCell className="hidden md:table-cell">{product.description || "-"}</TableCell>
                                  <TableCell>{product.price.toFixed(2)} ₽</TableCell>
                                  <TableCell>{product.stock}</TableCell>
                                  <TableCell>
                                    <Badge variant={product.status === "active" ? "success" : "destructive"}>
                                      {product.status}
                                    </Badge>
                                  </TableCell>
                                  {(userRole === "manager" || userRole === "director") && (
                                      <TableCell className="text-right">
                                        <Button variant="ghost" size="icon" onClick={() => handleEditProduct(product.id)}>
                                          <Pencil className="h-4 w-4" />
                                          <span className="sr-only">Edit</span>
                                        </Button>
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            className="text-destructive"
                                            onClick={() => handleDeleteProduct(product.id)}
                                        >
                                          <Trash className="h-4 w-4" />
                                          <span className="sr-only">Delete</span>
                                        </Button>
                                      </TableCell>
                                  )}
                                </TableRow>
                            ))
                        )}
                      </TableBody>
                    </Table>
                  </div>
                </div>
            )}

            {/* Users Tab */}
            {activeTab === "users" && userRole === "director" && (
                <div className="space-y-4">
                  <h2 className="text-lg font-semibold">User Management</h2>

                  <div className="border rounded-md">
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>Username</TableHead>
                          <TableHead>Email</TableHead>
                          <TableHead>Role</TableHead>
                          <TableHead>Status</TableHead>
                          <TableHead className="text-right">Actions</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {users.length === 0 ? (
                            <TableRow>
                              <TableCell colSpan={5} className="text-center py-4">
                                No users found
                              </TableCell>
                            </TableRow>
                        ) : (
                            users.map((user) => (
                                <TableRow key={user.id}>
                                  <TableCell className="font-medium">{user.username}</TableCell>
                                  <TableCell>{user.email}</TableCell>
                                  <TableCell>{getRoleDisplayName(user.role)}</TableCell>
                                  <TableCell>
                                    <Badge variant={getStatusBadgeVariant(user.status)}>{user.status}</Badge>
                                  </TableCell>
                                  <TableCell className="text-right space-x-1">
                                    <Dialog
                                        open={userRoleDialogOpen && selectedUserId === user.id}
                                        onOpenChange={(open) => {
                                          setUserRoleDialogOpen(open)
                                          if (!open) setSelectedUserId("")
                                        }}
                                    >
                                      <DialogTrigger asChild>
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            onClick={() => {
                                              setSelectedUserId(user.id)
                                              setSelectedUserRole(user.role)
                                            }}
                                        >
                                          Role
                                        </Button>
                                      </DialogTrigger>
                                      <DialogContent>
                                        <DialogHeader>
                                          <DialogTitle>Change User Role</DialogTitle>
                                          <DialogDescription>Select a new role for this user</DialogDescription>
                                        </DialogHeader>
                                        <div className="grid gap-4 py-4">
                                          <div className="grid gap-2">
                                            <Label htmlFor="role">Role</Label>
                                            <select
                                                id="role"
                                                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                                                value={selectedUserRole}
                                                onChange={(e) => setSelectedUserRole(e.target.value)}
                                            >
                                              <option value="viewer">Viewer</option>
                                              <option value="manager">Manager</option>
                                              <option value="director">Director</option>
                                            </select>
                                          </div>
                                        </div>
                                        <DialogFooter>
                                          <Button variant="outline" onClick={() => setUserRoleDialogOpen(false)}>
                                            Cancel
                                          </Button>
                                          <Button onClick={handleUserRoleChange}>Save</Button>
                                        </DialogFooter>
                                      </DialogContent>
                                    </Dialog>

                                    {user.status === "active" ? (
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            className="text-amber-500"
                                            onClick={() => handleBanUser(user.id)}
                                        >
                                          <Ban className="h-4 w-4 mr-1" />
                                          Ban
                                        </Button>
                                    ) : user.status === "banned" ? (
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            className="text-green-500"
                                            onClick={() => handleUnbanUser(user.id)}
                                        >
                                          <CheckCircle className="h-4 w-4 mr-1" />
                                          Unban
                                        </Button>
                                    ) : null}

                                    {user.status === "deleted" ? (
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            className="text-blue-500"
                                            onClick={() => handleRestoreUser(user.id)}
                                        >
                                          <RotateCcw className="h-4 w-4 mr-1" />
                                          Restore
                                        </Button>
                                    ) : (
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            className="text-destructive"
                                            onClick={() => handleDeleteUser(user.id)}
                                        >
                                          <Trash className="h-4 w-4 mr-1" />
                                          Delete
                                        </Button>
                                    )}
                                  </TableCell>
                                </TableRow>
                            ))
                        )}
                      </TableBody>
                    </Table>
                  </div>
                </div>
            )}

            {/* Profile Tab */}
            {activeTab === "profile" && (
                <div className="grid gap-6 md:grid-cols-2">
                  <Card>
                    <CardHeader>
                      <CardTitle>Profile Information</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <div className="space-y-4">
                        <div className="space-y-2">
                          <Label>Username</Label>
                          <Input value={profile?.username || ""} readOnly />
                        </div>
                        <div className="space-y-2">
                          <Label>Email</Label>
                          <Input value={profile?.email || ""} readOnly />
                        </div>
                        <div className="space-y-2">
                          <Label>Role</Label>
                          <Input value={getRoleDisplayName(profile?.role || "")} readOnly />
                        </div>
                      </div>
                    </CardContent>
                  </Card>

                  <Card>
                    <CardHeader>
                      <CardTitle>Change Password</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <form onSubmit={handlePasswordChange} className="space-y-4">
                        {passwordError && (
                            <Alert variant="destructive">
                              <AlertDescription>{passwordError}</AlertDescription>
                            </Alert>
                        )}
                        {passwordSuccess && (
                            <Alert className="bg-green-50 text-green-800 border-green-200">
                              <AlertDescription>Password changed successfully!</AlertDescription>
                            </Alert>
                        )}
                        <div className="space-y-2">
                          <Label htmlFor="oldPassword">Current Password</Label>
                          <Input
                              id="oldPassword"
                              type="password"
                              value={oldPassword}
                              onChange={(e) => setOldPassword(e.target.value)}
                              required
                          />
                        </div>
                        <div className="space-y-2">
                          <Label htmlFor="newPassword">New Password</Label>
                          <Input
                              id="newPassword"
                              type="password"
                              value={newPassword}
                              onChange={(e) => setNewPassword(e.target.value)}
                              required
                              minLength={6}
                          />
                          <p className="text-xs text-muted-foreground">Password must be at least 6 characters long</p>
                        </div>
                        <Button type="submit" className="w-full">
                          Change Password
                        </Button>
                      </form>
                    </CardContent>
                  </Card>
                </div>
            )}
          </main>
        </div>
      </div>
  )
}