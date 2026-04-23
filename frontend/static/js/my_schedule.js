function getCalendarBadge(status) {
    const value = (status || "").toLowerCase();
    let badgeClass = "default";

    if (value === "active") badgeClass = "active";
    if (value === "completed") badgeClass = "completed";

    return `<span class="badge ${badgeClass}">${status || "unknown"}</span>`;
}

async function loadMySchedule() {
    const response = await fetch("/bookings");
    const bookings = await response.json();

    const tableBody = document.getElementById("my-schedule-body");
    tableBody.innerHTML = "";

    if (!Array.isArray(bookings) || bookings.length === 0) {
        tableBody.innerHTML = `
            <tr>
                <td colspan="6">
                    <div class="empty-state">No bookings yet</div>
                </td>
            </tr>
        `;
        return;
    }

    bookings.forEach(booking => {
        const row = document.createElement("tr");
        row.innerHTML = `
            <td>${booking.id}</td>
            <td>${booking.client_name}</td>
            <td>${booking.event_date}</td>
            <td>${booking.start_time}</td>
            <td>${booking.end_time}</td>
            <td>${getCalendarBadge(booking.status)}</td>
        `;
        tableBody.appendChild(row);
    });
}

loadMySchedule();