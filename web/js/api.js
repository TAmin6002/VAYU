const API_BASE = "/api/v1";

/**
 * Wrapper around fetch that attaches the JWT (admin or student, if
 * present) and unwraps the { data, error } envelope the backend returns.
 * Browsing events is public and works with no token at all.
 */
async function apiRequest(path, { method = "GET", body } = {}) {
  const headers = { "Content-Type": "application/json" };
  const token = localStorage.getItem("token");
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${API_BASE}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  let payload = {};
  try {
    payload = await res.json();
  } catch (_) {
    // no body
  }

  // An expired/invalid token on a protected call: drop the session and
  // send the user to the login page.
  if (res.status === 401 && !path.startsWith("/auth/") && auth.isLoggedIn()) {
    auth.clearSession();
    window.location.href = "/login.html";
    throw new Error("نشست شما منقضی شده است. لطفاً دوباره وارد شوید.");
  }

  if (!res.ok) {
    throw new Error(payload.error || `request failed with status ${res.status}`);
  }

  return payload.data;
}

const api = {
  login: (data) => apiRequest("/auth/login", { method: "POST", body: data }),
  signup: (data) => apiRequest("/auth/signup", { method: "POST", body: data }),
  logout: () => apiRequest("/auth/logout", { method: "POST" }),

  listEvents: () => apiRequest("/events"),
  getEvent: (id) => apiRequest(`/events/${id}`),
  createEvent: (data) => apiRequest("/events", { method: "POST", body: data }),
  updateEvent: (id, data) => apiRequest(`/events/${id}`, { method: "PUT", body: data }),
  deleteEvent: (id) => apiRequest(`/events/${id}`, { method: "DELETE" }),

  registerForEvent: (id, data) => apiRequest(`/events/${id}/register`, { method: "POST", body: data }),
  cancelRegistration: (eventId) => apiRequest(`/events/${eventId}/cancel`, { method: "POST" }),
  myRegistrations: () => apiRequest("/student/registrations"),
  listParticipants: (id) => apiRequest(`/events/${id}/participants`),
};
