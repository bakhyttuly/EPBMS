function getStatusBadge(status) {
    const value = (status || "").toLowerCase();
    let badgeClass = "default";

    if (value === "available") badgeClass = "available";
    if (value === "busy") badgeClass = "busy";

    return `<span class="badge ${badgeClass}">${status || "unknown"}</span>`;
}

async function loadPerformers() {
    const response = await fetch("/performers");
    const performers = await response.json();

    const tableBody = document.getElementById("performers-table-body");
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

    performers.forEach(performer => {
        const row = document.createElement("tr");

        row.innerHTML = `
            <td>${performer.id}</td>
            <td>${performer.name}</td>
            <td>${performer.category}</td>
            <td>${performer.price}</td>
            <td>${performer.description || ""}</td>
            <td>${getStatusBadge(performer.availability_status)}</td>
            <td>
                <div class="action-group">
                    <button onclick="deletePerformer(${performer.id})" class="danger-btn">Delete</button>
                </div>
            </td>
        `;

        tableBody.appendChild(row);
    });
}

async function deletePerformer(id) {
    if (!confirm("Delete this performer?")) return;

    const response = await fetch(`/performers/${id}`, {
        method: "DELETE"
    });

    if (response.ok) {
        loadPerformers();
    } else {
        alert("Failed to delete performer");
    }
}

document.getElementById("performer-form").addEventListener("submit", async function (e) {
    e.preventDefault();

    const performer = {
        name: document.getElementById("name").value,
        category: document.getElementById("category").value,
        price: parseFloat(document.getElementById("price").value),
        description: document.getElementById("description").value,
        availability_status: document.getElementById("availability_status").value
    };

    const response = await fetch("/performers", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(performer)
    });

    if (response.ok) {
        document.getElementById("performer-form").reset();
        loadPerformers();
    } else {
        alert("Failed to add performer");
    }
});

loadPerformers();