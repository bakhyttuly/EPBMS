async function loadPerformers() {
    try {
        const response = await authenticatedFetch("/api/v1/performers");
        if (!response) return;
        
        const result = await response.json();
        const performers = result.data;

        const tableBody = document.getElementById("performers-table-body");
        if (!tableBody) return;
        
        tableBody.innerHTML = "";

        if (!Array.isArray(performers) || performers.length === 0) {
            tableBody.innerHTML = `
                <tr>
                    <td colspan="7">
                        <div class="empty-state">No performers yet</div>
                    </td>
                </tr>
            `;
            return;
        }

        const user = JSON.parse(localStorage.getItem("user") || "{}");

        performers.forEach(performer => {
            const row = document.createElement("tr");

            let actions = "";
            if (user.role === "admin") {
                actions = `<button onclick="deletePerformer(${performer.id})" class="danger-btn">Delete</button>`;
            } else if (user.role === "client") {
                actions = `<button onclick="openBookingModal(${performer.id}, '${performer.name}')" class="primary-btn">Book</button>`;
            }

            row.innerHTML = `
                <td>${performer.id}</td>
                <td>${performer.name}</td>
                <td>${performer.category}</td>
                <td>$${performer.price}</td>
                <td>${performer.description || ""}</td>
                <td><span class="badge available">Available</span></td>
                <td>
                    <div class="action-group">
                        ${actions}
                    </div>
                </td>
            `;

            tableBody.appendChild(row);
        });
    } catch (error) {
        console.error("Error loading performers:", error);
    }
}

async function deletePerformer(id) {
    if (!confirm("Delete this performer?")) return;

    const response = await authenticatedFetch(`/api/v1/performers/${id}`, {
        method: "DELETE"
    });

    if (response && response.ok) {
        loadPerformers();
    } else {
        alert("Failed to delete performer");
    }
}

async function openBookingModal(id, name) {
    const date = prompt(`Enter date for ${name} (YYYY-MM-DD):`, "2025-12-01");
    if (!date) return;
    const start = prompt("Enter start time (HH:MM):", "10:00");
    if (!start) return;
    const end = prompt("Enter end time (HH:MM):", "12:00");
    if (!end) return;

    const response = await authenticatedFetch("/api/v1/bookings", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            performer_id: id,
            event_date: date,
            start_time: start,
            end_time: end
        })
    });

    if (response && response.ok) {
        alert("Booking request sent! Waiting for admin approval.");
    } else {
        const result = await response.json();
        alert("Booking failed: " + (result.error || "Unknown error"));
    }
}

const performerForm = document.getElementById("performer-form");
if (performerForm) {
    performerForm.addEventListener("submit", async function (e) {
        e.preventDefault();

        const performer = {
            name: document.getElementById("name").value,
            category: document.getElementById("category").value,
            price: parseFloat(document.getElementById("price").value),
            description: document.getElementById("description").value
        };

        const response = await authenticatedFetch("/api/v1/performers", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(performer)
        });

        if (response && response.ok) {
            performerForm.reset();
            loadPerformers();
        } else {
            alert("Failed to add performer");
        }
    });
}

loadPerformers();
