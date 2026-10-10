function requireGroupAdmin() {
  return auth.requireAdmin();
}

function initialsOf(name) {
  const parts = (name || "").trim().split(" ").filter(Boolean);
  const letters = parts.slice(0, 2).map((p) => p[0]);
  return letters.join("") || "?";
}

// ---- groupAdmin_dashboard.html ----
async function initAdminDashboard() {
  if (!requireGroupAdmin()) return;

  const params = new URLSearchParams(window.location.search);
  const participantsFor = params.get("participants");

  if (participantsFor) {
    await renderParticipants(participantsFor);
    return;
  }

  const container = document.getElementById("admin-events-list");
  const statStrip = document.getElementById("stat-strip");
  container.innerHTML = "در حال بارگذاری...";

  try {
    const events = await api.listEvents();
    if (!events || events.length === 0) {
      if (statStrip) statStrip.style.display = "none";
      container.className = "";
      container.innerHTML = emptyStateHTML(
        'هنوز رویدادی ثبت نکرده‌اید. <a href="/groupAdmin_EventForm.html">ایجاد رویداد جدید</a>'
      );
      return;
    }

    if (statStrip) {
      const totalCapacity = events.reduce((sum, e) => sum + e.capacity, 0);
      const totalRegistered = events.reduce((sum, e) => sum + e.registered_count, 0);
      statStrip.style.display = "grid";
      statStrip.innerHTML = `
        <div class="stat-pill"><span class="num">${events.length}</span><span class="label">رویداد فعال</span></div>
        <div class="stat-pill"><span class="num">${totalRegistered}</span><span class="label">مجموع ثبت‌نام‌ها</span></div>
        <div class="stat-pill"><span class="num">${totalCapacity}</span><span class="label">مجموع ظرفیت</span></div>
      `;
    }

    container.className = "events-grid";
    container.innerHTML = events
      .map((e) => {
        const { day, month } = formatDateBadge(e.event_time);
        const pct = capacityPct(e);
        return `
      <div class="event-card" style="cursor:default">
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
          <span class="capacity-label ${e.is_full ? "full" : ""}">${e.registered_count} / ${e.capacity} ثبت‌نام</span>
          <div class="event-actions-inline">
            <a href="/groupAdmin_EventForm.html?id=${e.id}"><button class="ghost">ویرایش</button></a>
            <a href="/groupAdmin_dashboard.html?participants=${e.id}"><button class="ghost">شرکت‌کنندگان</button></a>
            <button class="danger" data-delete-id="${e.id}">حذف</button>
          </div>
        </div>
      </div>`;
      })
      .join("");

    container.querySelectorAll("[data-delete-id]").forEach((btn) => {
      btn.addEventListener("click", async () => {
        if (!confirm("از حذف این رویداد مطمئن هستید؟")) return;
        try {
          await api.deleteEvent(btn.dataset.deleteId);
          initAdminDashboard();
        } catch (err) {
          alert(err.message);
        }
      });
    });
  } catch (err) {
    container.innerHTML = `<p class="error">${err.message}</p>`;
  }
}

async function renderParticipants(eventId) {
  const container = document.getElementById("admin-events-list");
  const statStrip = document.getElementById("stat-strip");
  if (statStrip) statStrip.style.display = "none";
  container.className = "";
  container.innerHTML = "در حال بارگذاری...";

  try {
    const [event, participants] = await Promise.all([
      api.getEvent(eventId),
      api.listParticipants(eventId),
    ]);

    const rows = (participants || [])
      .map(
        (p) => `
      <tr>
        <td>
          <div class="person-cell">
            <span class="avatar-circle">${initialsOf(p.full_name)}</span>
            ${p.full_name}
          </div>
        </td>
        <td>${p.email}</td>
        <td>${p.student_number || "—"}</td>
        <td>${p.status === "registered" ? "ثبت‌نام شده" : "لغو شده"}</td>
        <td>${formatDateTime(p.registered_at)}</td>
      </tr>`
      )
      .join("");

    container.innerHTML = `
      <div class="page-toolbar" style="margin-bottom:16px;">
        <div>
          <h3 style="margin-bottom:4px;">شرکت‌کنندگان رویداد: ${event.title}</h3>
          <a class="kicker" style="margin:0;" href="/groupAdmin_dashboard.html">&larr; بازگشت به داشبورد</a>
        </div>
      </div>
      <table>
        <thead>
          <tr><th>نام</th><th>ایمیل</th><th>شماره دانشجویی</th><th>وضعیت</th><th>تاریخ ثبت‌نام</th></tr>
        </thead>
        <tbody>${rows || '<tr><td colspan="5">هنوز کسی ثبت‌نام نکرده است.</td></tr>'}</tbody>
      </table>
    `;
  } catch (err) {
    container.innerHTML = `<p class="error">${err.message}</p>`;
  }
}

// ---- groupAdmin_EventForm.html ----
async function initEventForm() {
  if (!requireGroupAdmin()) return;

  const params = new URLSearchParams(window.location.search);
  const id = params.get("id");
  const form = document.getElementById("event-form");
  const msg = document.getElementById("form-msg");
  const heading = document.getElementById("form-heading");

  if (id) {
    heading.textContent = "ویرایش رویداد";
    try {
      const event = await api.getEvent(id);
      form.title.value = event.title;
      form.description.value = event.description;
      form.location.value = event.location;
      form.capacity.value = event.capacity;
      form.event_time.value = new Date(event.event_time).toISOString().slice(0, 16);
    } catch (err) {
      msg.className = "error";
      msg.textContent = err.message;
    }
  } else {
    heading.textContent = "ایجاد رویداد جدید";
  }

  form.addEventListener("submit", async (ev) => {
    ev.preventDefault();
    msg.textContent = "";

    const payload = {
      title: form.title.value,
      description: form.description.value,
      location: form.location.value,
      capacity: parseInt(form.capacity.value, 10),
      event_time: new Date(form.event_time.value).toISOString(),
    };

    try {
      if (id) {
        await api.updateEvent(id, payload);
      } else {
        await api.createEvent(payload);
      }
      window.location.href = "/groupAdmin_dashboard.html";
    } catch (err) {
      msg.className = "error";
      msg.textContent = err.message;
    }
  });
}
