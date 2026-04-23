async function loadBookings() {
    try {
        const response = await authenticatedFetch("/api/v1/bookings");
        if (!response) return;

        const result = await response.json();
        const bookings = result.data;

        const tableBody = document.getElementById("bookings-table-body");
        if (!tableBody) return;
        
        tableBody.innerHTML = "";

        if (!Array.isArray(bookings) || bookings.length === 0) {
            tableBody.innerHTML = `<tr><td colspan="8"><div class="empty-state">No bookings found</div></td></tr>`;
            return;
        }

        const user = JSON.parse(localStorage.getItem("user") || "{}");

        bookings.forEach(booking => {
            const row = document.createElement("tr");

            let actions = "";
            if (user.role === "admin" && booking.status === "pending") {
                actions = `
                    <button onclick="updateStatus(${booking.id}, 'confirmed')" class="success-btn">Approve</button>
                    <button onclick="updateStatus(${booking.id}, 'rejected')" class="danger-btn">Reject</button>
                `;
            } else if (user.role === "admin" && booking.status === "confirmed") {
                actions = `<button onclick="updateStatus(${booking.id}, 'completed')" class="primary-btn">Complete</button>`;
            } else if (user.role === "admin") {
                actions = `<button onclick="deleteBooking(${booking.id})" class="danger-btn">Delete</button>`;
            }

            row.innerHTML = `
                <td>${booking.id}</td>
                <td>${booking.performer ? booking.performer.name : 'N/A'}</td>
                <td>${booking.event_date}</td>
                <td>${booking.start_time} - ${booking.end_time}</td>
                <td><span class="badge ${booking.status}">${booking.status}</span></td>
                <td>${booking.notes || ""}</td>
                <td>
                    <div class="action-group">
                        ${actions}
                    </div>
                </td>
            `;

            tableBody.appendChild(row);
        });
    } catch (error) {
        console.error("Error loading bookings:", error);
    }
}

async function updateStatus(id, status) {
    try {
        const response = await authenticatedFetch(`/api/v1/admin/bookings/${id}/status`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ status })
        });

        if (response && response.ok) {
            loadBookings();
        } else {
            const result = await response.json();
            alert("Failed to update status: " + (result.error || "Unknown error"));
        }
    } catch (error) {
        console.error("Error updating status:", error);
    }
}

async function deleteBooking(id) {
    if (!confirm("Delete this booking?")) return;

    const response = await authenticatedFetch(`/api/v1/admin/bookings/${id}`, {
        method: "DELETE"
    });

    if (response && response.ok) {
        loadBookings();
    } else {
        const result = await response.json();
        alert("Failed to delete booking: " + (result.error || "Unknown error"));
    }
}

loadBookings();
