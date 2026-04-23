function getCalendarBadge(status) {
    const value = (status || "").toLowerCase();
    let badgeClass = "default";

    if (value === "active") badgeClass = "active";
    if (value === "completed") badgeClass = "completed";

    return `<span class="badge ${badgeClass}">${status || "unknown"}</span>`;
}

document.getElementById("calendar-form").addEventListener("submit", async function (e) {
    e.preventDefault();

    const date = document.getElementById("calendar-date").value;
    const response = await fetch(`/bookings/by-date?date=${date}`);
    const bookings = await response.json();

    const tableBody = document.getElementById("calendar-table-body");
    tableBody.innerHTML = "";

    if (!Array.isArray(bookings) || bookings.length === 0) {
        tableBody.innerHTML = `
            <tr>
                <td colspan="7">
                    <div class="empty-state">No bookings for this date</div>
                </td>
            </tr>
        `;
        return;
    }

    bookings.forEach(booking => {
        const performerName = booking.performer ? booking.performer.name : booking.performer_id;

        const row = document.createElement("tr");
        row.innerHTML = `
            <td>${booking.id}</td>
            <td>${performerName}</td>
            <td>${booking.client_name}</td>
            <td>${booking.event_date}</td>
            <td>${booking.start_time}</td>
            <td>${booking.end_time}</td>
            <td>${getCalendarBadge(booking.status)}</td>
        `;
        tableBody.appendChild(row);
    });
});