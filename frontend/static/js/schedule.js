async function loadMySchedule() {
    try {
        const response = await authenticatedFetch("/api/v1/bookings");
        if (!response) return;

        const result = await response.json();
        const bookings = result.data;

        const tableBody = document.getElementById("schedule-table-body");
        if (!tableBody) return;
        
        tableBody.innerHTML = "";

        if (!Array.isArray(bookings) || bookings.length === 0) {
            tableBody.innerHTML = `<tr><td colspan="6"><div class="empty-state">No confirmed bookings in your schedule</div></td></tr>`;
            return;
        }

        bookings.forEach(booking => {
            const row = document.createElement("tr");
            row.innerHTML = `
                <td>${booking.id}</td>
                <td>${booking.event_date}</td>
                <td>${booking.start_time}</td>
                <td>${booking.end_time}</td>
                <td><span class="badge confirmed">Confirmed</span></td>
                <td>${booking.notes || ""}</td>
            `;
            tableBody.appendChild(row);
        });
    } catch (error) {
        console.error("Error loading schedule:", error);
    }
}

loadMySchedule();
