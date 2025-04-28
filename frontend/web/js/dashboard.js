document.addEventListener("DOMContentLoaded", () => {
  const API_URL = "/api"
  let currentUserRole = ""

  // Проверяем авторизацию
  checkAuth()

  // Обработчики навигации
  document.querySelectorAll(".nav-link").forEach((link) => {
    link.addEventListener("click", function (e) {
      e.preventDefault()

      // Убираем активный класс у всех ссылок
      document.querySelectorAll(".nav-link").forEach((l) => l.classList.remove("active"))

      // Добавляем активный класс текущей ссылке
      this.classList.add("active")

      // Получаем ID страницы для отображения
      const pageId = this.getAttribute("data-page")

      // Обновляем заголовок страницы
      document.getElementById("pageTitle").textContent = this.textContent.trim()

      // Скрываем все страницы
      document.querySelectorAll(".page-container").forEach((page) => page.classList.add("d-none"))

      // Показываем выбранную страницу
      document.getElementById(`${pageId}Page`).classList.remove("d-none")

      // Загружаем данные для страницы
      if (pageId === "dashboard") {
        loadDashboardData()
      } else if (pageId === "products") {
        loadProducts()
      } else if (pageId === "users") {
        loadUsers()
      } else if (pageId === "profile") {
        loadProfile()
      }
    })
  })

  // Обработчик выхода
  document.getElementById("logoutBtn").addEventListener("click", () => {
    localStorage.removeItem("accessToken")
    localStorage.removeItem("refreshToken")
    localStorage.removeItem("tokenType")
    window.location.href = "/web/index.html"
  })

  // Обработчик формы изменения пароля
  document.getElementById("changePasswordForm").addEventListener("submit", async (e) => {
    e.preventDefault()

    const oldPassword = document.getElementById("oldPassword").value
    const newPassword = document.getElementById("newPassword").value

    try {
      document.getElementById("passwordError").classList.add("d-none")
      document.getElementById("passwordSuccess").classList.add("d-none")

      const response = await fetchWithAuth(`${API_URL}/profile/password`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ oldPassword, newPassword }),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Ошибка изменения пароля")
      }

      document.getElementById("passwordSuccess").classList.remove("d-none")
      document.getElementById("changePasswordForm").reset()
    } catch (error) {
      document.getElementById("passwordError").textContent = error.message
      document.getElementById("passwordError").classList.remove("d-none")
    }
  })

  // Обработчик кнопки добавления товара
  document.getElementById("addProductBtn").addEventListener("click", () => {
    document.getElementById("productForm").reset()
    document.getElementById("productId").value = ""
    document.getElementById("productModalLabel").textContent = "Добавить товар"

    const productModalElement = document.getElementById("productModal")
    const productModal = new bootstrap.Modal(productModalElement)
    productModal.show()
  })

  // Обработчик сохранения товара
  document.getElementById("saveProductBtn").addEventListener("click", async () => {
    const productId = document.getElementById("productId").value
    const name = document.getElementById("productName").value
    const description = document.getElementById("productDescription").value
    const price = Number.parseFloat(document.getElementById("productPrice").value)
    const stock = Number.parseInt(document.getElementById("productStock").value)

    if (!name || isNaN(price) || isNaN(stock)) {
      alert("Пожалуйста, заполните все обязательные поля")
      return
    }

    const productData = {
      name,
      description,
      price,
      stock,
    }

    try {
      let response

      if (productId) {
        // Обновление существующего товара
        response = await fetchWithAuth(`${API_URL}/products/${productId}`, {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(productData),
        })
      } else {
        // Создание нового товара
        response = await fetchWithAuth(`${API_URL}/products`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(productData),
        })
      }

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Ошибка сохранения товара")
      }

      // Закрываем модальное окно
      const productModalElement = document.getElementById("productModal")
      const productModal = bootstrap.Modal.getInstance(productModalElement)
      productModal.hide()

      // Обновляем список товаров
      loadProducts()
    } catch (error) {
      alert(`Ошибка: ${error.message}`)
    }
  })

  // Обработчик сохранения роли пользователя
  document.getElementById("saveUserRoleBtn").addEventListener("click", async () => {
    const userId = document.getElementById("userRoleUserId").value
    const role = document.getElementById("userRoleSelect").value

    try {
      const response = await fetchWithAuth(`${API_URL}/admin/assign-role`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ userId, role }),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Ошибка изменения роли")
      }

      // Закрываем модальное окно
      const userRoleModalElement = document.getElementById("userRoleModal")
      const userRoleModal = bootstrap.Modal.getInstance(userRoleModalElement)
      userRoleModal.hide()

      // Обновляем список пользователей
      loadUsers()
    } catch (error) {
      alert(`Ошибка: ${error.message}`)
    }
  })

  // Функция проверки авторизации
  async function checkAuth() {
    const accessToken = localStorage.getItem("accessToken")

    if (!accessToken) {
      window.location.href = "/web/index.html"
      return
    }

    try {
      const response = await fetchWithAuth(`${API_URL}/profile`)

      if (!response.ok) {
        throw new Error("Ошибка авторизации")
      }

      const userData = await response.json()
      currentUserRole = userData.role

      // Отображаем роль пользователя
      document.getElementById("userRole").textContent = `Роль: ${getRoleDisplayName(currentUserRole)}`

      // Настраиваем интерфейс в зависимости от роли
      setupUIForRole(currentUserRole)

      // Загружаем данные для дашборда
      loadDashboardData()
    } catch (error) {
      console.error("Ошибка проверки авторизации:", error)

      // Пробуем обновить токен
      const refreshed = await refreshToken()

      if (!refreshed) {
        // Если не удалось обновить токен, перенаправляем на страницу входа
        localStorage.removeItem("accessToken")
        localStorage.removeItem("refreshToken")
        window.location.href = "/web/index.html"
      } else {
        // Если токен обновлен, повторяем проверку
        checkAuth()
      }
    }
  }

  // Функция настройки интерфейса в зависимости от роли
  function setupUIForRole(role) {
    // Показываем/скрываем элементы для администраторов
    const adminElements = document.querySelectorAll(".admin-only")
    if (role === "director") {
      adminElements.forEach((el) => el.classList.remove("d-none"))
    } else {
      adminElements.forEach((el) => el.classList.add("d-none"))
    }

    // Показываем/скрываем элементы для менеджеров
    const managerElements = document.querySelectorAll(".manager-only")
    if (role === "manager" || role === "director") {
      managerElements.forEach((el) => el.classList.remove("d-none"))
    } else {
      managerElements.forEach((el) => el.classList.add("d-none"))
    }
  }

  // Функция загрузки данных для дашборда
  async function loadDashboardData() {
    try {
      // Загружаем список товаров для подсчета статистики
      const productsResponse = await fetchWithAuth(`${API_URL}/products`)

      if (!productsResponse.ok) {
        throw new Error("Ошибка загрузки товаров")
      }

      const products = await productsResponse.json()

      // Подсчитываем статистику
      const productsCount = products.length
      const totalValue = products.reduce((sum, product) => sum + product.price * product.stock, 0)

      // Обновляем информацию на дашборде
      document.getElementById("productsCount").textContent = productsCount
      document.getElementById("totalValue").textContent = `${totalValue.toFixed(2)} ₽`

      // Если пользователь - директор, загружаем информацию о пользователях
      if (currentUserRole === "director") {
        const usersResponse = await fetchWithAuth(`${API_URL}/admin/users`)

        if (!usersResponse.ok) {
          throw new Error("Ошибка загрузки пользователей")
        }

        const users = await usersResponse.json()
        document.getElementById("usersCount").textContent = users.length
      }
    } catch (error) {
      console.error("Ошибка загрузки данных для дашборда:", error)
    }
  }

  // Функция загрузки списка товаров
  async function loadProducts() {
    try {
      const response = await fetchWithAuth(`${API_URL}/products`)

      if (!response.ok) {
        throw new Error("Ошибка загрузки товаров")
      }

      const products = await response.json()
      const tableBody = document.getElementById("productsTableBody")
      tableBody.innerHTML = ""

      products.forEach((product) => {
        const row = document.createElement("tr")

        row.innerHTML = `
                    <td>${product.name}</td>
                    <td>${product.description || "-"}</td>
                    <td>${product.price.toFixed(2)} ₽</td>
                    <td>${product.stock}</td>
                    <td><span class="badge ${product.status === "active" ? "bg-success" : "bg-danger"}">${product.status}</span></td>
                    <td class="manager-only ${currentUserRole === "manager" || currentUserRole === "director" ? "" : "d-none"}">
                        <button class="btn btn-sm btn-primary edit-product" data-id="${product.id}">
                            <i class="bi bi-pencil"></i>
                        </button>
                        <button class="btn btn-sm btn-danger delete-product" data-id="${product.id}">
                            <i class="bi bi-trash"></i>
                        </button>
                    </td>
                `

        tableBody.appendChild(row)
      })

      // Добавляем обработчики для кнопок редактирования и удаления
      document.querySelectorAll(".edit-product").forEach((button) => {
        button.addEventListener("click", function () {
          const productId = this.getAttribute("data-id")
          editProduct(productId)
        })
      })

      document.querySelectorAll(".delete-product").forEach((button) => {
        button.addEventListener("click", function () {
          const productId = this.getAttribute("data-id")
          deleteProduct(productId)
        })
      })
    } catch (error) {
      console.error("Ошибка загрузки товаров:", error)
    }
  }

  // Функция загрузки списка пользователей
  async function loadUsers() {
    if (currentUserRole !== "director") {
      return
    }

    try {
      const response = await fetchWithAuth(`${API_URL}/admin/users`)

      if (!response.ok) {
        throw new Error("Ошибка загрузки пользователей")
      }

      const users = await response.json()
      const tableBody = document.getElementById("usersTableBody")
      tableBody.innerHTML = ""

      users.forEach((user) => {
        const row = document.createElement("tr")

        row.innerHTML = `
                    <td>${user.username}</td>
                    <td>${user.email}</td>
                    <td>${getRoleDisplayName(user.role)}</td>
                    <td><span class="badge ${getStatusBadgeClass(user.status)}">${user.status}</span></td>
                    <td>
                        <button class="btn btn-sm btn-primary change-role" data-id="${user.id}" data-role="${user.role}">
                            <i class="bi bi-person-gear"></i> Роль
                        </button>
                        ${
                          user.status === "active"
                            ? `<button class="btn btn-sm btn-warning ban-user" data-id="${user.id}">
                                <i class="bi bi-slash-circle"></i> Блок
                            </button>`
                            : `<button class="btn btn-sm btn-success unban-user" data-id="${user.id}">
                                <i class="bi bi-check-circle"></i> Разблок
                            </button>`
                        }
                        ${
                          user.status === "deleted"
                            ? `<button class="btn btn-sm btn-info restore-user" data-id="${user.id}">
                                <i class="bi bi-arrow-counterclockwise"></i> Восстановить
                            </button>`
                            : `<button class="btn btn-sm btn-danger delete-user" data-id="${user.id}">
                                <i class="bi bi-trash"></i> Удалить
                            </button>`
                        }
                    </td>
                `

        tableBody.appendChild(row)
      })

      // Добавляем обработчики для кнопок управления пользователями
      document.querySelectorAll(".change-role").forEach((button) => {
        button.addEventListener("click", function () {
          const userId = this.getAttribute("data-id")
          const currentRole = this.getAttribute("data-role")
          changeUserRole(userId, currentRole)
        })
      })

      document.querySelectorAll(".ban-user").forEach((button) => {
        button.addEventListener("click", function () {
          const userId = this.getAttribute("data-id")
          banUser(userId)
        })
      })

      document.querySelectorAll(".unban-user").forEach((button) => {
        button.addEventListener("click", function () {
          const userId = this.getAttribute("data-id")
          unbanUser(userId)
        })
      })

      document.querySelectorAll(".delete-user").forEach((button) => {
        button.addEventListener("click", function () {
          const userId = this.getAttribute("data-id")
          deleteUser(userId)
        })
      })

      document.querySelectorAll(".restore-user").forEach((button) => {
        button.addEventListener("click", function () {
          const userId = this.getAttribute("data-id")
          restoreUser(userId)
        })
      })
    } catch (error) {
      console.error("Ошибка загрузки пользователей:", error)
    }
  }

  // Функция загрузки профиля пользователя
  async function loadProfile() {
    try {
      const response = await fetchWithAuth(`${API_URL}/profile`)

      if (!response.ok) {
        throw new Error("Ошибка загрузки профиля")
      }

      const profile = await response.json()

      document.getElementById("profileUsername").value = profile.username
      document.getElementById("profileEmail").value = profile.email
      document.getElementById("profileRole").value = getRoleDisplayName(profile.role)
    } catch (error) {
      console.error("Ошибка загрузки профиля:", error)
    }
  }

  // Функция редактирования товара
  async function editProduct(productId) {
    try {
      const response = await fetchWithAuth(`${API_URL}/products/${productId}`)

      if (!response.ok) {
        throw new Error("Ошибка загрузки информации о товаре")
      }

      const product = await response.json()

      document.getElementById("productId").value = product.id
      document.getElementById("productName").value = product.name
      document.getElementById("productDescription").value = product.description || ""
      document.getElementById("productPrice").value = product.price
      document.getElementById("productStock").value = product.stock

      document.getElementById("productModalLabel").textContent = "Редактировать товар"

      const productModalElement = document.getElementById("productModal")
      const productModal = new bootstrap.Modal(productModalElement)
      productModal.show()
    } catch (error) {
      console.error("Ошибка редактирования товара:", error)
      alert(`Ошибка: ${error.message}`)
    }
  }

  // Функция удаления товара
  async function deleteProduct(productId) {
    if (!confirm("Вы уверены, что хотите удалить этот товар?")) {
      return
    }

    try {
      const response = await fetchWithAuth(`${API_URL}/products/${productId}`, {
        method: "DELETE",
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Ошибка удаления товара")
      }

      // Обновляем список товаров
      loadProducts()
    } catch (error) {
      console.error("Ошибка удаления товара:", error)
      alert(`Ошибка: ${error.message}`)
    }
  }

  // Функция изменения роли пользователя
  function changeUserRole(userId, currentRole) {
    document.getElementById("userRoleUserId").value = userId
    document.getElementById("userRoleSelect").value = currentRole

    const userRoleModalElement = document.getElementById("userRoleModal")
    const userRoleModal = new bootstrap.Modal(userRoleModalElement)
    userRoleModal.show()
  }

  // Функция блокировки пользователя
  async function banUser(userId) {
    if (!confirm("Вы уверены, что хотите заблокировать этого пользователя?")) {
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
        throw new Error(data.message || "Ошибка блокировки пользователя")
      }

      // Обновляем список пользователей
      loadUsers()
    } catch (error) {
      console.error("Ошибка блокировки пользователя:", error)
      alert(`Ошибка: ${error.message}`)
    }
  }

  // Функция разблокировки пользователя
  async function unbanUser(userId) {
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
        throw new Error(data.message || "Ошибка разблокировки пользователя")
      }

      // Обновляем список пользователей
      loadUsers()
    } catch (error) {
      console.error("Ошибка разблокировки пользователя:", error)
      alert(`Ошибка: ${error.message}`)
    }
  }

  // Функция удаления пользователя
  async function deleteUser(userId) {
    if (!confirm("Вы уверены, что хотите удалить этого пользователя?")) {
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
        throw new Error(data.message || "Ошибка удаления пользователя")
      }

      // Обновляем список пользователей
      loadUsers()
    } catch (error) {
      console.error("Ошибка удаления пользователя:", error)
      alert(`Ошибка: ${error.message}`)
    }
  }

  // Функция восстановления пользователя
  async function restoreUser(userId) {
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
        throw new Error(data.message || "Ошибка восстановления пользователя")
      }

      // Обновляем список пользователей
      loadUsers()
    } catch (error) {
      console.error("Ошибка восстановления пользователя:", error)
      alert(`Ошибка: ${error.message}`)
    }
  }

  // Функция для выполнения запросов с авторизацией
  async function fetchWithAuth(url, options = {}) {
    const accessToken = localStorage.getItem("accessToken")
    const tokenType = localStorage.getItem("tokenType") || "Bearer"

    if (!options.headers) {
      options.headers = {}
    }

    options.headers["Authorization"] = `${tokenType} ${accessToken}`

    try {
      const response = await fetch(url, options)

      // Если получили 401, пробуем обновить токен
      if (response.status === 401) {
        const refreshed = await refreshToken()

        if (refreshed) {
          // Если токен обновлен, повторяем запрос
          const newAccessToken = localStorage.getItem("accessToken")
          options.headers["Authorization"] = `${tokenType} ${newAccessToken}`
          return fetch(url, options)
        }
      }

      return response
    } catch (error) {
      console.error("Ошибка запроса:", error)
      throw error
    }
  }

  // Функция обновления токена
  async function refreshToken() {
    const refreshToken = localStorage.getItem("refreshToken")

    if (!refreshToken) {
      return false
    }

    try {
      const response = await fetch(`${API_URL}/refresh`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ refresh_token: refreshToken }),
      })

      if (!response.ok) {
        throw new Error("Не удалось обновить токен")
      }

      const data = await response.json()

      localStorage.setItem("accessToken", data.access_token)
      localStorage.setItem("refreshToken", data.refresh_token)

      return true
    } catch (error) {
      console.error("Ошибка обновления токена:", error)
      return false
    }
  }

  // Вспомогательные функции
  function getRoleDisplayName(role) {
    switch (role) {
      case "viewer":
        return "Просмотр"
      case "manager":
        return "Менеджер"
      case "director":
        return "Директор"
      default:
        return role
    }
  }

  function getStatusBadgeClass(status) {
    switch (status) {
      case "active":
        return "bg-success"
      case "banned":
        return "bg-danger"
      case "deleted":
        return "bg-secondary"
      default:
        return "bg-primary"
    }
  }
})
