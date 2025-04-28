import { API_URL } from "./constants"

// Function to refresh the access token
export async function refreshToken(): Promise<boolean> {
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
      throw new Error("Failed to refresh token")
    }

    const data = await response.json()

    localStorage.setItem("accessToken", data.access_token)
    localStorage.setItem("refreshToken", data.refresh_token)

    return true
  } catch (error) {
    console.error("Error refreshing token:", error)
    return false
  }
}

// Function to make authenticated API requests
export async function fetchWithAuth(url: string, options: RequestInit = {}): Promise<Response> {
  const accessToken = localStorage.getItem("accessToken")
  const tokenType = localStorage.getItem("tokenType") || "Bearer"

  if (!options.headers) {
    options.headers = {}
  }
  // Add authorization header
  ;(options.headers as Record<string, string>)["Authorization"] = `${tokenType} ${accessToken}`

  try {
    const response = await fetch(url, options)

    // If unauthorized, try to refresh the token
    if (response.status === 401) {
      const refreshed = await refreshToken()

      if (refreshed) {
        // If token refreshed successfully, retry the request
        const newAccessToken = localStorage.getItem("accessToken")
        if (newAccessToken) {
          ;(options.headers as Record<string, string>)["Authorization"] = `${tokenType} ${newAccessToken}`
          return fetch(url, options)
        }
      }
    }

    return response
  } catch (error) {
    console.error("Error making authenticated request:", error)
    throw error
  }
}
