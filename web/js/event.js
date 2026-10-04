function formatDateTime(iso) {
  const d = new Date(iso);
  return d.toLocaleString("fa-IR", { dateStyle: "medium", timeStyle: "short" });
}

const FA_MONTHS = [
  "فروردین", "اردیبهشت", "خرداد", "تیر", "مرداد", "شهریور",
  "مهر", "آبان", "آذر", "دی", "بهمن", "اسفند",
];

function formatDateBadge(iso) {
  const parts = new Intl.DateTimeFormat("fa-IR-u-nu-latn", {
    day: "numeric",
    month: "numeric",
  }).formatToParts(new Date(iso));
  const day = parts.find((p) => p.type === "day")?.value || "";
  const monthIdx = parseInt(parts.find((p) => p.type === "month")?.value || "1", 10) - 1;
  return { day, month: FA_MONTHS[monthIdx] || "" };
}

const ICON_CALENDAR = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="3"/><path d="M16 2v4M8 2v4M3 10h18"/></svg>`;
const ICON_PIN = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/><circle cx="12" cy="10" r="3"/></svg>`;
const ICON_USERS = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"/></svg>`;
const ICON_EMPTY = `<svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="3"/><path d="M16 2v4M8 2v4M3 10h18M9 16l2 2 4-4"/></svg>`;

function capacityPct(e) {
  if (!e.capacity) return 0;
  return Math.min(100, Math.round((e.registered_count / e.capacity) * 100));
}

function eventCardHTML(e) {
  const { day, month } = formatDateBadge(e.event_time);
  const pct = capacityPct(e);
  return `
    <a class="event-card" href="/event_detail.html?id=${e.id}">
      <div class="event-date">
        <span class="day">${day}</span>
        <span class="month">${month}</span>
      </div>
      <div class="event-body">
        <h3>${e.title}</h3>
        <div class="event-meta-row">
          <span>${ICON_CALENDAR} ${formatDateTime(e.event_time)}</span>
          <span>${ICON_PIN} ${e.location}</span>
        </div>
        <div class="capacity-track">
          <div class="capacity-fill ${e.is_full ? "full" : ""}" style="width:${pct}%"></div>
        </div>
        <span class="capacity-label ${e.is_full ? "full" : ""}">
          ${e.registered_count} / ${e.capacity} ثبت‌نام${e.is_full ? " — تکمیل ظرفیت" : ""}
        </span>
      </div>
    </a>`;
}

function emptyStateHTML(text) {
  return `<div class="empty-state">${ICON_EMPTY}<p>${text}</p></div>`;
}

// ---- events.html ----
async function initEventsList() {
  const container = document.getElementById("events-list");
  container.innerHTML = "در حال بارگذاری...";

  try {
    const events = await api.listEvents();
    if (!events || events.length === 0) {
      container.className = "";
      container.innerHTML = emptyStateHTML("در حال حاضر رویدادی ثبت نشده است.");
      return;
    }
    container.className = "events-grid";
    container.innerHTML = events.map(eventCardHTML).join("");
  } catch (err) {
    container.innerHTML = `<p class="error">${err.message}</p>`;
  }
}

// ---- event_detail.html ----
function escapeHTML(v) {
  return String(v ?? "").replace(/[&<>"']/g, (c) => ({
    "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;",
  }[c]));
}

// Registration form — pre-filled from the logged-in student's account.
function registerFormHTML(user) {
  return `
    <form id="register-form" class="card" style="margin-top:0">
      <h4 style="margin-top:0">ثبت‌نام در این رویداد</h4>
      <input type="text" name="full_name" placeholder="نام و نام خانوادگی" value="${escapeHTML(user?.full_name)}" required />
      <input type="email" name="email" placeholder="ایمیل" value="${escapeHTML(user?.email)}" required />
      <input type="text" name="student_number" placeholder="شماره دانشجویی" value="${escapeHTML(user?.student_number)}" inputmode="numeric" autocomplete="off" pattern="[0-9۰-۹٠-٩]{5,20}" title="شماره دانشجویی باید فقط شامل ارقام باشد (۵ تا ۲۰ رقم)" required />
      <button type="submit">ثبت‌نام</button>
      <p id="register-msg"></p>
    </form>`;
}

function registeredCardHTML() {
  return `
    <div class="card" style="margin-top:0">
      <h4 style="margin-top:0">شما در این رویداد ثبت‌نام کرده‌اید ✓</h4>
      <p>اگر دیگر نمی‌توانید شرکت کنید، می‌توانید ثبت‌نام خود را لغو کنید.</p>
      <button id="cancel-btn" class="danger">انصراف از ثبت‌نام</button>
      <p id="cancel-msg"></p>
    </div>`;
}

function loginPromptHTML() {
  return `
    <div class="card" style="margin-top:0">
      <h4 style="margin-top:0">ثبت‌نام در این رویداد</h4>
      <p>برای ثبت‌نام در رویداد ابتدا وارد حساب دانشجویی خود شوید.</p>
      <a href="/login.html"><button>ورود</button></a>
      <a href="/signup.html"><button class="ghost">ساخت حساب دانشجو</button></a>
    </div>`;
}

async function initEventDetail() {
  const params = new URLSearchParams(window.location.search);
  const id = params.get("id");
  const container = document.getElementById("event-detail");

  try {
    const event = await api.getEvent(id);
    container.innerHTML = `
      <div class="card">
        <h3>${escapeHTML(event.title)}</h3>
        <div class="detail-meta">
          <span>${ICON_CALENDAR} ${formatDateTime(event.event_time)}</span>
          <span>${ICON_PIN} ${escapeHTML(event.location)}</span>
          <span>${ICON_USERS} ظرفیت ${event.capacity} نفر</span>
        </div>
        <p>${escapeHTML(event.description)}</p>
      </div>
    `;

    await renderEventActions(id, event);
  } catch (err) {
    container.innerHTML = `<p class="error">${err.message}</p>`;
  }
}

// Admin  -> management buttons only.
// Guest  -> prompt to log in / sign up.
// Student-> registration form, or (if already registered) a cancel button.
async function renderEventActions(eventId, event) {
  const actions = document.getElementById("event-actions");

  if (auth.isAdmin()) {
    actions.innerHTML = `<a href="/groupAdmin_EventForm.html?id=${event.id}"><button class="ghost">ویرایش رویداد</button></a>
              <a href="/groupAdmin_dashboard.html?participants=${event.id}"><button class="ghost">مشاهده شرکت‌کنندگان</button></a>`;
    return;
  }

  if (!auth.isStudent()) {
    actions.innerHTML = loginPromptHTML();
    return;
  }

  let isRegistered = false;
  try {
    const regs = await api.myRegistrations();
    isRegistered = (regs || []).some((r) => r.event_id === eventId);
  } catch (err) {
    actions.innerHTML = `<p class="error">${err.message}</p>`;
    return;
  }

  if (isRegistered) {
    actions.innerHTML = registeredCardHTML();
    const btn = document.getElementById("cancel-btn");
    btn.addEventListener("click", async () => {
      if (!confirm("از انصراف مطمئن هستید؟")) return;
      const msg = document.getElementById("cancel-msg");
      msg.textContent = "";
      try {
        await api.cancelRegistration(eventId);
        await renderEventActions(eventId, event);
      } catch (err) {
        msg.className = "error";
        msg.textContent = err.message;
      }
    });
    return;
  }

  actions.innerHTML = registerFormHTML(auth.getUser());
  const form = document.getElementById("register-form");
  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const msg = document.getElementById("register-msg");
    msg.textContent = "";
    try {
      await api.registerForEvent(eventId, {
        full_name: form.full_name.value,
        email: form.email.value,
        student_number: form.student_number.value.trim(),
      });
      await renderEventActions(eventId, event);
    } catch (err) {
      msg.className = "error";
      msg.textContent = err.message;
    }
  });
}

// ---- myevents.html ----
function myEventCardHTML(e) {
  const { day, month } = formatDateBadge(e.event_time);
  return `
    <div class="event-card" style="cursor:default">
      <div class="event-date">
        <span class="day">${day}</span>
        <span class="month">${month}</span>
      </div>
      <div class="event-body">
        <h3>${escapeHTML(e.title)}</h3>
        <div class="event-meta-row">
          <span>${ICON_CALENDAR} ${formatDateTime(e.event_time)}</span>
          <span>${ICON_PIN} ${escapeHTML(e.location)}</span>
        </div>
        <div class="event-actions-inline">
          <a href="/event_detail.html?id=${e.id}"><button class="ghost">مشاهده رویداد</button></a>
          <button class="danger" data-event-id="${e.id}">انصراف</button>
        </div>
      </div>
    </div>`;
}

async function initMyEvents() {
  if (!auth.requireStudent()) return;

  const container = document.getElementById("my-events-list");
  container.innerHTML = "در حال بارگذاری...";

  try {
    const regs = (await api.myRegistrations()) || [];

    if (regs.length === 0) {
      container.className = "";
      container.innerHTML = emptyStateHTML("شما هنوز در رویدادی ثبت‌نام نکرده‌اید.");
      return;
    }

    const events = await Promise.all(regs.map((r) => api.getEvent(r.event_id).catch(() => null)));

    container.className = "events-grid";
    container.innerHTML = events.filter(Boolean).map(myEventCardHTML).join("");

    container.querySelectorAll("[data-event-id]").forEach((btn) => {
      btn.addEventListener("click", async () => {
        if (!confirm("از انصراف مطمئن هستید؟")) return;
        try {
          await api.cancelRegistration(btn.dataset.eventId);
          initMyEvents();
        } catch (err) {
          alert(err.message);
        }
      });
    });
  } catch (err) {
    container.innerHTML = `<p class="error">${err.message}</p>`;
  }
}
