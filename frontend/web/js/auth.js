document.addEventListener("DOMContentLoaded", () => {
  const API_URL = "/api"

  // Форма входа
  const loginForm = document.getElementById("loginForm")
  const loginError = document.getElementById("loginError")

  // Форма регистрации
  const registerForm = document.getElementById("registerForm")
  const registerError = document.getElementById("registerError")
  const registerSuccess = document.getElementById("registerSuccess")

  // Обработчик формы входа
  loginForm.addEventListener("submit", async (e) => {
    e.preventDefault()

    const email = document.getElementById("loginEmail").value
    const password = document.getElementById("loginPassword").value

    try {
      loginError.classList.add("d-none")

      const response = await fetch(`${API_URL}/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.message || "Ошибка входа")
      }

      // Сохраняем токены в localStorage
      localStorage.setItem("accessToken", data.access_token)
      localStorage.setItem("refreshToken", data.refresh_token)
      localStorage.setItem("tokenType", data.token_type)

      // Перенаправляем на страницу профиля или дашборда
      window.location.href = "/web/dashboard.html"
    } catch (error) {
      loginError.textContent = error.message
      loginError.classList.remove("d-none")
    }
  })

  // Обработчик формы регистрации
  registerForm.addEventListener("submit", async (e) => {
    e.preventDefault()

    const username = document.getElementById("registerUsername").value
    const email = document.getElementById("registerEmail").value
    const password = document.getElementById("registerPassword").value

    try {
      registerError.classList.add("d-none")
      registerSuccess.classList.add("d-none")

      const response = await fetch(`${API_URL}/register`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, email, password }),
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.message || "Ошибка регистрации")
      }

      // Показываем сообщение об успешной регистрации
      registerSuccess.classList.remove("d-none")
      registerForm.reset()

      // Переключаемся на вкладку входа через 2 секунды
      setTimeout(() => {
        document.getElementById("login-tab").click()
      }, 2000)
    } catch (error) {
      registerError.textContent = error.message
      registerError.classList.remove("d-none")
    }
  })

  // Функция для обновления токена
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

  // Проверяем, есть ли токен, и если есть, перенаправляем на дашборд
  const accessToken = localStorage.getItem("accessToken")
  if (accessToken) {
    window.location.href = "/web/dashboard.html"
  }
})
