async function loadPerformerDropdown() {
    const response = await fetch("/performers");
    const performers = await response.json();

    const select = document.getElementById("performer_id");
    select.innerHTML = "";

    performers.forEach(p => {
        const option = document.createElement("option");
        option.value = p.id;
        option.textContent = `${p.name} (${p.category})`;
        select.appendChild(option);
    });
}
function getBookingBadge(status) {
    const value = (status || "").toLowerCase();
    let badgeClass = "default";

    if (value === "active") badgeClass = "active";
    if (value === "completed") badgeClass = "completed";

    return `<span class="badge ${badgeClass}">${status || "unknown"}</span>`;
}

async function loadCurrentUser() {
    const response = await fetch("/me");
    if (!response.ok) return null;
    return await response.json();
}

async function loadBookings() {
    const currentUser = await loadCurrentUser();
    const response = await fetch("/bookings");
    const bookings = await response.json();

    const tableBody = document.getElementById("bookings-table-body");
    tableBody.innerHTML = "";

    if (!Array.isArray(bookings) || bookings.length === 0) {
        tableBody.innerHTML = `
            <tr>
                <td colspan="8">
                    <div class="empty-state">No bookings yet</div>
                </td>
            </tr>
        `;
        return;
    }

    bookings.forEach(booking => {
        const performerName = booking.performer ? booking.performer.name : booking.performer_id;

        const canDelete = currentUser && currentUser.role !== "performer";

        const row = document.createElement("tr");
        row.innerHTML = `
            <td>${booking.id}</td>
            <td>${performerName}</td>
            <td>${booking.client_name}</td>
            <td>${booking.event_date}</td>
            <td>${booking.start_time}</td>
            <td>${booking.end_time}</td>
            <td>${getBookingBadge(booking.status)}</td>
            <td>
                <div class="action-group">
                    ${canDelete ? `<button onclick="deleteBooking(${booking.id})" class="danger-btn">Delete</button>` : ``}
                </div>
            </td>
        `;

        tableBody.appendChild(row);
    });
}

async function deleteBooking(id) {
    if (!confirm("Delete this booking?")) return;

    const response = await fetch(`/bookings/${id}`, {
        method: "DELETE"
    });

    if (response.ok) {
        loadBookings();
    } else {
        const data = await response.json();
        alert(data.error || "Failed to delete booking");
    }
}

document.getElementById("booking-form").addEventListener("submit", async function (e) {
    e.preventDefault();

    const booking = {
        performer_id: parseInt(document.getElementById("performer_id").value),
        client_name: document.getElementById("client_name").value,
        event_date: document.getElementById("event_date").value,
        start_time: document.getElementById("start_time").value,
        end_time: document.getElementById("end_time").value,
        status: document.getElementById("status").value
    };

    const response = await fetch("/bookings", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(booking)
    });

    const data = await response.json();

    if (response.ok) {
        document.getElementById("booking-form").reset();
        loadBookings();
    } else {
        alert(data.error || "Failed to create booking");
    }
});

loadBookings();
loadPerformerDropdown();