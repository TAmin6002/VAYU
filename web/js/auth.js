// Two account types: "admin" (manages events) and "student" (registers
// for events). The login response tells us which one we got; we keep the
// token, the role and the public profile in localStorage.
const auth = {
  saveSession(token, role, user) {
    localStorage.setItem("token", token);
    localStorage.setItem("role", role);
    localStorage.setItem("user", JSON.stringify(user));
  },
  clearSession() {
    localStorage.removeItem("token");
    localStorage.removeItem("role");
    localStorage.removeItem("user");
  },
  getUser() {
    try {
      return JSON.parse(localStorage.getItem("user"));
    } catch (_) {
      return null;
    }
  },
  role() {
    return localStorage.getItem("role");
  },
  isLoggedIn() {
    return !!localStorage.getItem("token") && !!this.role();
  },
  isAdmin() {
    return this.isLoggedIn() && this.role() === "admin";
  },
  isStudent() {
    return this.isLoggedIn() && this.role() === "student";
  },
  async logout() {
    try {
      await api.logout();
    } catch (_) {
      // Token might already be expired/invalid — doesn't matter, we're
      // clearing it locally either way.
    }
    this.clearSession();
    window.location.href = "/login.html";
  },
  // Page guards: send anyone without the right role to the login page.
  requireAdmin() {
    if (!this.isAdmin()) {
      window.location.href = "/login.html";
      return false;
    }
    return true;
  },
  requireStudent() {
    if (!this.isStudent()) {
      window.location.href = "/login.html";
      return false;
    }
    return true;
  },
};

function renderNav() {
  const nav = document.getElementById("nav");
  if (!nav) return;

  let links = `<a href="/events.html">رویدادها</a>`;

  if (auth.isAdmin()) {
    links += `<a href="/groupAdmin_dashboard.html">پنل ادمین</a>
              <button id="logout-btn" class="nav-logout">خروج</button>`;
  } else if (auth.isStudent()) {
    links += `<a href="/myevents.html">ثبت‌نام‌های من</a>
              <button id="logout-btn" class="nav-logout">خروج</button>`;
  } else {
    links += `<a href="/login.html">ورود</a><a href="/signup.html">ثبت‌نام دانشجو</a>`;
  }

  nav.innerHTML = links;

  const logoutBtn = document.getElementById("logout-btn");
  if (logoutBtn) logoutBtn.addEventListener("click", () => auth.logout());
}

document.addEventListener("DOMContentLoaded", renderNav);
